package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
)

type rootConfig struct {
	Server    *serverConfig
	Snowflake *snowflakeConfig
	Logging   *loggingConfig
	Trace     *traceConfig
	Metric    *metricConfig
	Secure    *secureConfig
}

func (c *rootConfig) GetServerConfig() ServerConfig {
	if c.Server == nil {
		c.Server = &serverConfig{}
	}
	return c.Server
}

func (c *rootConfig) GetSnowflakeConfig() SnowflakeConfig {
	if c.Snowflake == nil {
		c.Snowflake = &snowflakeConfig{}
	}
	return c.Snowflake
}

func (c *rootConfig) GetLoggingConfig() LoggingConfig {
	if c.Logging == nil {
		c.Logging = &loggingConfig{}
	}
	return c.Logging
}

func (c *rootConfig) GetTraceConfig() TraceConfig {
	if c.Trace == nil {
		c.Trace = &traceConfig{}
	}
	return c.Trace
}

func (c *rootConfig) GetMetricConfig() MetricConfig {
	if c.Metric == nil {
		c.Metric = &metricConfig{}
	}
	return c.Metric
}

func (c *rootConfig) GetSecureConfig() SecureConfig {
	if c.Secure == nil {
		c.Secure = &secureConfig{}
	}
	return c.Secure
}

type serverConfig struct {
	Debug   *bool
	Name    *string
	Version *string
	HTTP    *serverHTTPConfig
	DB      *serverDBConfig
	Redis   *serverRedisConfig
	MinIO   *serverMinIOConfig
	NATS    *serverNATSConfig
}

func (c *serverConfig) GetDebug() bool {
	if c.Debug == nil {
		return false
	} else {
		return *c.Debug
	}
}

func (c *serverConfig) GetName() string {
	if c.Name == nil {
		return "<UNKNOWN>"
	} else {
		return *c.Name
	}
}

func (c *serverConfig) GetVersion() string {
	if c.Version == nil {
		return "0.0.1"
	} else {
		return *c.Version
	}
}

func (c *serverConfig) GetInstanceId() string {
	if name, ok := os.LookupEnv("SERVER_INSTANCE_ID"); ok {
		return name
	} else if hostname, err := os.Hostname(); err != nil {
		return hostname
	} else {
		return uuid.NewString()
	}
}

func (c *serverConfig) GetHTTPConfig() ServerHTTPConfig {
	if c.HTTP == nil {
		c.HTTP = &serverHTTPConfig{}
	}
	return c.HTTP
}

func (c *serverConfig) GetDBConfig() ServerDBConfig {
	if c.DB == nil {
		c.DB = &serverDBConfig{}
	}
	return c.DB
}

func (c *serverConfig) GetRedisConfig() ServerRedisConfig {
	if c.Redis == nil {
		c.Redis = &serverRedisConfig{}
	}
	return c.Redis
}

func (c *serverConfig) GetMinIOConfig() ServerMinIOConfig {
	if c.MinIO == nil {
		c.MinIO = &serverMinIOConfig{}
	}
	return c.MinIO
}

func (c *serverConfig) GetNATSConfig() ServerNATSConfig {
	if c.NATS == nil {
		c.NATS = &serverNATSConfig{}
	}
	return c.NATS
}

type serverHTTPConfig struct {
	Addr *string
	Cors *serverHTTPCorsConfig
}

func (c *serverHTTPConfig) GetAddr() string {
	if c.Addr == nil {
		return ":5050"
	} else {
		return *c.Addr
	}
}

func (c *serverHTTPConfig) GetCors() cors.Options {
	return c.Cors.corsOptions()
}

type serverHTTPCorsConfig struct {
	AllowedOrigins       *[]string `yaml:"allowed-origins"`
	AllowedMethods       *[]string `yaml:"allowed-methods"`
	AllowedHeaders       *[]string `yaml:"allowed-headers"`
	ExposedHeaders       *[]string `yaml:"exposed-headers"`
	MaxAge               *int      `yaml:"max-age"`
	AllowCredentials     *bool     `yaml:"allow-credentials"`
	AllowPrivateNetwork  *bool     `yaml:"allow-private-network"`
	OptionsPassthrough   *bool     `yaml:"options-passthrough"`
	OptionsSuccessStatus *int      `yaml:"options-success-status"`
}

func (c *serverHTTPCorsConfig) corsOptions() cors.Options {
	opts := cors.Options{
		AllowedOrigins:       []string{"*"},
		AllowedMethods:       []string{"HEAD", "GET", "POST"},
		AllowedHeaders:       []string{"Authorization", "Content-Type", "Content-Length"},
		ExposedHeaders:       []string{"Content-Type", "Content-Length"},
		MaxAge:               5,
		AllowCredentials:     false,
		AllowPrivateNetwork:  false,
		OptionsPassthrough:   false,
		OptionsSuccessStatus: 204,
	}
	if c == nil {
		return opts
	}
	if c.AllowedOrigins != nil {
		opts.AllowedOrigins = *c.AllowedOrigins
	}
	if c.AllowedMethods != nil {
		opts.AllowedMethods = *c.AllowedMethods
	}
	if c.AllowedHeaders != nil {
		opts.AllowedHeaders = *c.AllowedHeaders
	}
	if c.ExposedHeaders != nil {
		opts.ExposedHeaders = *c.ExposedHeaders
	}
	if c.MaxAge != nil {
		opts.MaxAge = *c.MaxAge
		if opts.MaxAge == 0 {
			opts.MaxAge = 5
		} else if opts.MaxAge == -1 {
			opts.MaxAge = 0
		}
	}
	if c.AllowCredentials != nil {
		opts.AllowCredentials = *c.AllowCredentials
	}
	if c.AllowPrivateNetwork != nil {
		opts.AllowPrivateNetwork = *c.AllowPrivateNetwork
	}
	if c.OptionsPassthrough != nil {
		opts.OptionsPassthrough = *c.OptionsPassthrough
	}
	if c.OptionsSuccessStatus != nil {
		opts.OptionsSuccessStatus = *c.OptionsSuccessStatus
		if opts.OptionsSuccessStatus == 0 {
			opts.OptionsSuccessStatus = 204
		}
	}
	return opts
}

type serverDBConfig struct {
	Driver *string
	Source *string
}

func (c *serverDBConfig) GetDriver() string {
	if c.Driver == nil {
		panic("DB driver is not set")
	} else {
		return *c.Driver
	}
}

func (c *serverDBConfig) GetSource() string {
	if c.Source == nil {
		panic("DB source is not set")
	} else {
		return *c.Source
	}
}

type serverRedisConfig struct {
	InitAddr *string `yaml:"init-addr"`
	SelectDB *int    `yaml:"select-db"`
}

func (c *serverRedisConfig) GetInitAddr() string {
	if c.InitAddr == nil {
		return "127.0.0.1:6379"
	} else {
		return *c.InitAddr
	}
}

func (c *serverRedisConfig) GetSelectDB() int {
	if c.SelectDB == nil {
		return 0
	} else {
		return *c.SelectDB
	}
}

type serverMinIOConfig struct {
	Endpoint  *string
	AccessKey *string `yaml:"access-key"`
	SecretKey *string `yaml:"secret-key"`
	UseSSL    *bool   `yaml:"use-ssl"`
}

func (c *serverMinIOConfig) GetEndpoint() string {
	if c.Endpoint == nil {
		panic("MinIO endpoint is not set")
	} else {
		return *c.Endpoint
	}
}

func (c *serverMinIOConfig) GetAccessKey() string {
	if c.AccessKey == nil {
		return ""
	} else {
		return *c.AccessKey
	}
}

func (c *serverMinIOConfig) GetSecretKey() string {
	if c.SecretKey == nil {
		return ""
	} else {
		return *c.SecretKey
	}
}

func (c *serverMinIOConfig) GetUseSSL() bool {
	if c.UseSSL == nil {
		return true
	} else {
		return *c.UseSSL
	}
}

type serverNATSConfig struct {
	SeedURL *string `yaml:"seed-url"`
	NoEcho  *bool   `yaml:"no-echo"`
}

func (c *serverNATSConfig) GetSeedURL() string {
	if c.SeedURL == nil {
		return "nats://127.0.0.1:4222"
	} else {
		return *c.SeedURL
	}
}

func (c *serverNATSConfig) GetNoEcho() bool {
	if c.NoEcho == nil {
		return false
	} else {
		return *c.NoEcho
	}
}

type snowflakeConfig struct {
	IdEpoch       *int64  `yaml:"id-epoch"`
	ClusterId     *int64  `yaml:"cluster-id"`
	WorkerId      *int64  `yaml:"worker-id"`
	WorkerSeqKey  *string `yaml:"worker-seq-key"`
	ClusterIdBits *int32  `yaml:"cluster-id-bits"`
	WorkerIdBits  *int32  `yaml:"worker-id-bits"`
	SequenceBits  *int32  `yaml:"sequence-bits"`
}

func (c *snowflakeConfig) GetIdEpoch() int64 {
	if c.IdEpoch == nil {
		return int64(1640995200000) // Defaults to: 2022-01-01T00:00:00Z
	} else {
		return *c.IdEpoch
	}
}

func (c *snowflakeConfig) GetClusterId() int64 {
	if c.ClusterId == nil {
		return 0
	} else {
		return *c.ClusterId
	}
}

func (c *snowflakeConfig) GetWorkerId() int64 {
	if c.WorkerId == nil {
		return 0
	} else {
		return *c.WorkerId
	}
}

func (c *snowflakeConfig) GetWorkerSeqKey() string {
	if c.WorkerSeqKey == nil {
		return ""
	} else {
		return *c.WorkerSeqKey
	}
}

func (c *snowflakeConfig) GetClusterIdBits() int32 {
	if c.ClusterIdBits == nil {
		return 5
	} else {
		return *c.ClusterIdBits
	}
}

func (c *snowflakeConfig) GetWorkerIdBits() int32 {
	if c.WorkerIdBits == nil {
		return 5
	} else {
		return *c.WorkerIdBits
	}
}

func (c *snowflakeConfig) GetSequenceBits() int32 {
	if c.SequenceBits == nil {
		return 12
	} else {
		return *c.SequenceBits
	}
}

type loggingConfig struct {
	SlogLogger *loggingSlogLoggerConfig `yaml:"slog-logger"`
	ZapLogger  *loggingZapLoggerConfig  `yaml:"zap-logger"`
}

func (c *loggingConfig) GetSlogLogger() LoggingSlogLoggerConfig {
	if c.SlogLogger == nil {
		return nil
	}
	return c.SlogLogger
}

func (c *loggingConfig) GetZapLogger() LoggingZapLoggerConfig {
	if c.ZapLogger == nil {
		return nil
	}
	return c.ZapLogger
}

type loggingSlogLoggerConfig struct {
	Handler   *string
	AddSource *bool
	Leveler   *slog.Level
}

func (c *loggingSlogLoggerConfig) GetHandler() string {
	if c.Handler == nil {
		return "text"
	} else {
		return *c.Handler
	}
}

func (c *loggingSlogLoggerConfig) GetAddSource() bool {
	if c.AddSource == nil {
		return true
	} else {
		return *c.AddSource
	}
}

func (c *loggingSlogLoggerConfig) GetLeveler() slog.Leveler {
	if c.Leveler == nil {
		return slog.LevelDebug
	} else {
		return *c.Leveler
	}
}

type loggingZapLoggerConfig struct {
	Preset *string
}

func (c *loggingZapLoggerConfig) GetPreset() string {
	if c.Preset == nil {
		return "development"
	} else {
		return *c.Preset
	}
}

type traceConfig struct {
	Exporter *traceExporterConfig
}

func (c *traceConfig) GetExporterConfig() TraceExporterConfig {
	return c.Exporter
}

type traceExporterConfig struct {
	Protocol *string
	Endpoint *string
	Insecure *bool
}

func (c *traceExporterConfig) GetProtocol() string {
	if c.Protocol == nil {
		return "noop"
	} else {
		return *c.Protocol
	}
}

func (c *traceExporterConfig) GetEndpoint() string {
	if c.Endpoint == nil {
		return ""
	} else {
		return *c.Endpoint
	}
}

func (c *traceExporterConfig) GetInsecure() bool {
	if c.Insecure == nil {
		return false
	} else {
		return *c.Insecure
	}
}

type metricConfig struct {
	Exporter *metricExporterConfig
}

func (c *metricConfig) GetExporterConfig() MetricExporterConfig {
	return c.Exporter
}

type metricExporterConfig struct {
	Protocol *string
	Endpoint *string
	Insecure *bool
}

func (c *metricExporterConfig) GetProtocol() string {
	if c.Protocol == nil {
		return "noop"
	} else {
		return *c.Protocol
	}
}

func (c *metricExporterConfig) GetEndpoint() string {
	if c.Endpoint == nil {
		return ""
	} else {
		return *c.Endpoint
	}
}

func (c *metricExporterConfig) GetInsecure() bool {
	if c.Insecure == nil {
		return false
	} else {
		return *c.Insecure
	}
}

type secureConfig struct {
	Token *secureTokenConfig
}

func (c *secureConfig) GetToken() SecureTokenConfig {
	if c.Token == nil {
		c.Token = &secureTokenConfig{}
	}
	return c.Token
}

type secureTokenConfig struct {
	Store           *string
	Bucket          *string
	AccessTokenTTL  *time.Duration `yaml:"access-token-ttl"`
	RefreshTokenTTL *time.Duration `yaml:"refresh-token-ttl"`
	Issuer          *string
	Audience        *string
	SigningMethod   *string `yaml:"signing-method"`
	PublicKeyFile   *string `yaml:"public-key-file"`
	PrivateKeyFile  *string `yaml:"private-key-file"`
	PublicKeyValue  *string `yaml:"public-key-value"`
	PrivateKeyValue *string `yaml:"private-key-value"`
}

func (c *secureTokenConfig) GetStore() string {
	if c.Store == nil {
		return "redis"
	} else {
		return *c.Store
	}
}

func (c *secureTokenConfig) GetBucket() string {
	if c.Bucket == nil {
		return "tokens"
	} else {
		return *c.Bucket
	}
}

func (c *secureTokenConfig) GetAccessTokenTTL() time.Duration {
	if c.AccessTokenTTL == nil {
		return 2 * 24 * time.Hour // default to 2 days
	} else {
		return *c.AccessTokenTTL
	}
}

func (c *secureTokenConfig) GetRefreshTokenTTL() time.Duration {
	if c.RefreshTokenTTL == nil {
		return 7 * 24 * time.Hour // default to 7 days
	} else {
		return *c.RefreshTokenTTL
	}
}

func (c *secureTokenConfig) GetIssuer() string {
	if c.Issuer == nil {
		return "unnamed-issuer"
	} else {
		return *c.Issuer
	}
}

func (c *secureTokenConfig) GetAudience() string {
	if c.Audience == nil {
		return "unnamed-audience"
	} else {
		return *c.Audience
	}
}

func (c *secureTokenConfig) GetSigningMethod() string {
	if c.SigningMethod == nil {
		return "RS256"
	} else {
		return *c.SigningMethod
	}
}

func (c *secureTokenConfig) GetPublicKey() []byte {
	if c.PublicKeyValue != nil {
		return []byte(*c.PublicKeyValue)
	}
	if c.PublicKeyFile != nil {
		body, err := os.ReadFile(*c.PublicKeyFile)
		if err != nil {
			panic(err)
		}
		return body
	}
	return nil
}

func (c *secureTokenConfig) GetPrivateKey() []byte {
	if c.PrivateKeyValue != nil {
		return []byte(*c.PrivateKeyValue)
	}
	if c.PrivateKeyFile != nil {
		body, err := os.ReadFile(*c.PrivateKeyFile)
		if err != nil {
			panic(err)
		}
		return body
	}
	return nil
}
