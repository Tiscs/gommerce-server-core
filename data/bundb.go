package data

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/choral-io/gommerce-server-core/config"
	"github.com/choral-io/gommerce-server-core/logging"
	_ "github.com/shopspring/decimal"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mssqldialect"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bunotel"
	"github.com/uptrace/bun/schema"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type pagetoken interface {
	GetPage() int32
	GetSize() int32
}

// WithPaging is a bun.SelectQuery modifier that adds paging to the query.
func WithPaging(p pagetoken) func(*bun.SelectQuery) *bun.SelectQuery {
	var size = 10
	if p.GetSize() > 0 {
		size = int(p.GetSize())
	}
	var page = 1
	if p.GetPage() > 1 {
		page = int(p.GetPage())
	}
	return func(query *bun.SelectQuery) *bun.SelectQuery {
		return query.Offset((page - 1) * size).Limit(size)
	}
}

type globalQueryHook struct {
	logger logging.Logger
}

func (h *globalQueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h *globalQueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	duration := time.Since(event.StartTime)
	if event.Err != nil {
		h.logger.Error(ctx, "", slog.String("query", event.Query), slog.String("operation", event.Operation()), slog.Duration("duration", duration), slog.String("error", event.Err.Error()))
	} else {
		h.logger.Debug(ctx, "", slog.String("query", event.Query), slog.String("operation", event.Operation()), slog.Duration("duration", duration))
	}
}

// NewBunDB creates a new bun.IDB instance with metrics, tracing and logging.
func NewBunDB(cfg config.ServerDBConfig, logger logging.Logger, tp trace.TracerProvider, mp metric.MeterProvider) (bun.IDB, error) {
	var dialect schema.Dialect
	switch cfg.GetDriver() {
	case "pg", "pgsql":
		dialect = pgdialect.New()
	case "mysql":
		dialect = mysqldialect.New()
	case "mssql":
		dialect = mssqldialect.New()
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.GetDriver())
	}
	sdb, err := sql.Open(cfg.GetDriver(), cfg.GetSource())
	if err != nil {
		return nil, err
	}
	bdb := bun.NewDB(sdb, dialect, bun.WithDiscardUnknownColumns())
	bdb.AddQueryHook(bunotel.NewQueryHook(bunotel.WithTracerProvider(tp), bunotel.WithMeterProvider(mp)))
	bdb.AddQueryHook(&globalQueryHook{logger: logger})
	bdb.RegisterModel()
	return bdb, nil
}
