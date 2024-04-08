package server

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"strings"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
	"github.com/choral-io/gommerce-server-core/secure"
	"github.com/choral-io/gommerce-server-core/validator"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// ServerServiceRegister is implemented by servers that register grpc services.
type ServerServiceRegister interface {
	RegisterServerService(grpc.ServiceRegistrar)
}

// ServerServiceRegisterFunc is a function that registers grpc services.
type ServerServiceRegisterFunc func(grpc.ServiceRegistrar)

// GatewayClientRegister is implemented by servers that register grpc gateway clients.
type GatewayClientRegister interface {
	RegisterGatewayClient(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

// GatewayClientRegisterFunc is a function that registers grpc gateway clients.
type GatewayClientRegisterFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

// GRPCHandler is an implementation of http.Handler for gRPC.
type GRPCHandler struct {
	h2cHandler http.Handler

	srvOptions []grpc.ServerOption      // grpc server options
	gcdOptions []grpc.DialOption        // grpc client dial options
	gtwOptions []runtime.ServeMuxOption // grpc gateway options

	unaryInts  []grpc.UnaryServerInterceptor  // grpc unary interceptors
	streamInts []grpc.StreamServerInterceptor // grpc stream interceptors

	srvServers []ServerServiceRegisterFunc // grpc server services
	gtwClients []GatewayClientRegisterFunc // grpc gateway clients

	useHealthz bool // whether to use healthz endpoint
}

// GRPCHandlerOption is an option for GRPCHandler, used to configure it.
type GRPCHandlerOption func(*GRPCHandler) error

// NewGRPCHandler returns a new GRPCHandler with the given config and options.
func NewGRPCHandler(cfg config.ServerHTTPConfig, opts ...GRPCHandlerOption) (*GRPCHandler, error) {
	h := &GRPCHandler{
		srvOptions: []grpc.ServerOption{},
		gtwOptions: []runtime.ServeMuxOption{},
		unaryInts:  []grpc.UnaryServerInterceptor{},
		streamInts: []grpc.StreamServerInterceptor{},
		gcdOptions: []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	}

	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, cfg.GetAddr(), h.gcdOptions...)
	if err != nil {
		return nil, err
	}

	h.srvOptions = append(h.srvOptions,
		grpc.ChainUnaryInterceptor(h.unaryInts...),
		grpc.ChainStreamInterceptor(h.streamInts...),
	)

	if h.useHealthz {
		h.gtwOptions = append(h.gtwOptions, runtime.WithHealthzEndpoint(health.NewHealthClient(conn)))
	}

	grpcServer := grpc.NewServer(h.srvOptions...)
	gatewayMux := runtime.NewServeMux(h.gtwOptions...)

	for _, srv := range h.srvServers {
		srv(grpcServer)
	}

	for _, crf := range h.gtwClients {
		if err := crf(ctx, gatewayMux, conn); err != nil {
			return nil, err
		}
	}

	reflection.Register(grpcServer)
	gtwHandler := cors.Default().Handler(gatewayMux)

	h.h2cHandler = h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			gtwHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	return h, nil
}

// ServeHTTP implements http.Handler.
func (h *GRPCHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h2cHandler.ServeHTTP(w, r)
}

// WithOTELStatsHandler returns a GRPCHandlerOption that use an opentelemetry stats handler for grpc server.
func WithOTELStatsHandler(tp trace.TracerProvider, mp metric.MeterProvider) GRPCHandlerOption {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	return func(h *GRPCHandler) error {
		// otelgrpc.UnaryServerInterceptor and otelgrpc.StreamServerInterceptor are deprecated,
		// Use otelgrpc.NewServerHandler instead
		h.srvOptions = append(h.srvOptions, grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(propagator),
		)))
		h.gcdOptions = append(h.gcdOptions, grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithTracerProvider(tp),
			otelgrpc.WithMeterProvider(mp),
			otelgrpc.WithPropagators(propagator),
		)))
		return nil
	}
}

// WithStaticFileHandler returns a GRPCHandlerOption that adds a static file handler to grpc gateway.
func WithStaticFileHandler(pattern string, sfs fs.FS) GRPCHandlerOption {
	hfs := http.FileServer(http.FS(sfs))
	return func(h *GRPCHandler) error {
		h.gtwOptions = append(h.gtwOptions, func(mux *runtime.ServeMux) {
			err := mux.HandlePath("GET", pattern, func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
				fi, err := fs.Stat(sfs, strings.TrimLeft(r.URL.Path, "/"))
				if errors.Is(err, fs.ErrNotExist) || fi.IsDir() {
					http.ServeFileFS(w, r, sfs, "index.html")
				} else if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					hfs.ServeHTTP(w, r)
				}
			})
			if err != nil {
				panic(err)
			}
		})
		return nil
	}
}

// WithServeMuxRoutes returns a GRPCHandlerOption that adds routes to grpc gateway.
func WithServeMuxRoutes(routes ...ServerMuxRoute) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		h.gtwOptions = append(h.gtwOptions, func(mux *runtime.ServeMux) {
			for _, route := range routes {
				for _, method := range route.Methods {
					err := mux.HandlePath(method, route.Pattern, route.Handler)
					if err != nil {
						panic(err)
					}
				}
			}
		})
		return nil
	}
}

// WithUnaryInterceptors returns a GRPCHandlerOption that adds the given unary interceptors to grpc handler.
func WithUnaryInterceptors(ints ...grpc.UnaryServerInterceptor) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		h.unaryInts = append(h.unaryInts, ints...)
		return nil
	}
}

// WithStreamInterceptors returns a GRPCHandlerOption that adds the given stream interceptors to grpc handler.
func WithStreamInterceptors(ints ...grpc.StreamServerInterceptor) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		h.streamInts = append(h.streamInts, ints...)
		return nil
	}
}

// WithLoggingInterceptor returns a GRPCHandlerOption that adds a logging interceptor to grpc handler.
func WithLoggingInterceptor(logger logging.Logger) GRPCHandlerOption {
	grpclog := logging.NewGRPCLogger(logger)
	return func(h *GRPCHandler) error {
		h.unaryInts = append(h.unaryInts, grpclog.UnaryServerInterceptor())
		h.streamInts = append(h.streamInts, grpclog.StreamServerInterceptor())
		return nil
	}
}

// WithRecoveryInterceptor returns a GRPCHandlerOption that adds a recovery interceptor to grpc handler.
// The given recovery handler will be called when a panic occurs.
func WithRecoveryInterceptor(f recovery.RecoveryHandlerFuncContext) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		h.unaryInts = append(h.unaryInts, recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(f)))
		h.streamInts = append(h.streamInts, recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(f)))
		return nil
	}
}

// WithValidatorInterceptor returns a GRPCHandlerOption that adds a validator interceptor to grpc handler.
func WithValidatorInterceptor() GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		h.unaryInts = append(h.unaryInts, validator.UnaryServerInterceptor())
		h.streamInts = append(h.streamInts, validator.StreamServerInterceptor())
		return nil
	}
}

// WithSecureInterceptor returns a GRPCHandlerOption that adds a secure interceptor to grpc handler.
// The given matcher will be used to determine which methods should be secured.
func WithSecureInterceptor(auth *secure.ServerAuthorizer, matcher selector.Matcher) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		if matcher == nil {
			h.unaryInts = append(h.unaryInts, auth.UnaryServerInterceptor())
			h.streamInts = append(h.streamInts, auth.StreamServerInterceptor())
		} else {
			h.unaryInts = append(h.unaryInts, selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(), matcher))
			h.streamInts = append(h.streamInts, selector.StreamServerInterceptor(auth.StreamServerInterceptor(), matcher))
		}
		return nil
	}
}

// WithRegistrations returns a GRPCHandlerOption that registers grpc servers and gateway clients.
// The given registrations must implement ServerServiceRegister or GatewayClientRegister.
func WithRegistrations(regs ...any) GRPCHandlerOption {
	return func(h *GRPCHandler) error {
		for _, reg := range regs {
			if _, ok := reg.(health.HealthServer); ok {
				h.useHealthz = true
			}
			registered := false
			if r, ok := reg.(ServerServiceRegister); ok {
				registered = true
				h.srvServers = append(h.srvServers, r.RegisterServerService)
			}
			if r, ok := reg.(ServerServiceRegisterFunc); ok {
				registered = true
				h.srvServers = append(h.srvServers, r)
			}
			if r, ok := reg.(GatewayClientRegister); ok {
				registered = true
				h.gtwClients = append(h.gtwClients, r.RegisterGatewayClient)
			}
			if r, ok := reg.(GatewayClientRegisterFunc); ok {
				registered = true
				h.gtwClients = append(h.gtwClients, r)
			}
			if !registered {
				return errors.New("registration must implement ServerServiceRegister or GatewayClientRegister")
			}
		}
		return nil
	}
}
