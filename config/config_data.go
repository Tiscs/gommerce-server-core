package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// LoadRootConfig loads RootConfig from file.
// The path of the file is specified by environment variable GOMMERCE_CONFIG_PATH,
// If the environment variable is not set, it defaults to "./config/app-deploy.yaml".
func LoadRootConfig() (RootConfig, error) {
	path, ok := os.LookupEnv("GOMMERCE_CONFIG_PATH")
	if !ok {
		path = "./config/app-deploy.yaml"
	}
	txt, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &rootConfig{}
	if err := yaml.Unmarshal(txt, cfg); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}

type rootConfig struct {
	Server   *serverConfig
	IdWorker *idWorkerConfig `yaml:"id-worker"`
	Logging  *loggingConfig
	Trace    *traceConfig
	Metric   *metricConfig
	Token    *tokenConfig
}

func (c *rootConfig) GetServerConfig() ServerConfig {
	if c.Server == nil {
		c.Server = &serverConfig{}
	}
	return c.Server
}

func (c *rootConfig) GetIdWorkerConfig() IdWorkerConfig {
	if c.IdWorker == nil {
		c.IdWorker = &idWorkerConfig{}
	}
	return c.IdWorker
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

func (c *rootConfig) GetTokenConfig() TokenConfig {
	if c.Token == nil {
		c.Token = &tokenConfig{}
	}
	return c.Token
}

type serverConfig struct {
	Debug   *bool
	Name    *string
	Version *string
	HTTP    *serverHTTPConfig
	DB      *serverDBConfig
	Redis   *serverRedisConfig
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
		return "0.0.0"
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

func (c *serverConfig) GetNATSConfig() ServerNATSConfig {
	if c.NATS == nil {
		c.NATS = &serverNATSConfig{}
	}
	return c.NATS
}

type serverHTTPConfig struct {
	Addr *string
}

func (c *serverHTTPConfig) GetAddr() string {
	if c.Addr == nil {
		return ":50050"
	} else {
		return *c.Addr
	}
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
	InitAddress []string `yaml:"init-address"`
	SelectDB    *int     `yaml:"select-db"`
}

func (c *serverRedisConfig) GetInitAddress() []string {
	if c.InitAddress == nil {
		return []string{"127.0.0.1:6379"}
	} else {
		return c.InitAddress
	}
}

func (c *serverRedisConfig) GetSelectDB() int {
	if c.SelectDB == nil {
		return 0
	} else {
		return *c.SelectDB
	}
}

type serverNATSConfig struct {
	URL *string
}

func (c *serverNATSConfig) GetURL() string {
	if c.URL == nil {
		return "nats://127.0.0.1:4222"
	} else {
		return *c.URL
	}
}

type idWorkerConfig struct {
	IdEpoch       *int64  `yaml:"id-epoch"`
	ClusterId     *int64  `yaml:"cluster-id"`
	WorkerId      *int64  `yaml:"worker-id"`
	WorkerSeqKey  *string `yaml:"worker-seq-key"`
	ClusterIdBits *int32  `yaml:"cluster-id-bits"`
	WorkerIdBits  *int32  `yaml:"worker-id-bits"`
	SequenceBits  *int32  `yaml:"sequence-bits"`
}

func (c *idWorkerConfig) GetIdEpoch() int64 {
	if c.IdEpoch == nil {
		return int64(1640995200000) // Defaults to: 2022-01-01T00:00:00Z
	} else {
		return *c.IdEpoch
	}
}

func (c *idWorkerConfig) GetClusterId() int64 {
	if c.ClusterId == nil {
		return 0
	} else {
		return *c.ClusterId
	}
}

func (c *idWorkerConfig) GetWorkerId() int64 {
	if c.WorkerId == nil {
		return 0
	} else {
		return *c.WorkerId
	}
}

func (c *idWorkerConfig) GetWorkerSeqKey() string {
	if c.WorkerSeqKey == nil {
		return ""
	} else {
		return *c.WorkerSeqKey
	}
}

func (c *idWorkerConfig) GetClusterIdBits() int32 {
	if c.ClusterIdBits == nil {
		return 5
	} else {
		return *c.ClusterIdBits
	}
}

func (c *idWorkerConfig) GetWorkerIdBits() int32 {
	if c.WorkerIdBits == nil {
		return 5
	} else {
		return *c.WorkerIdBits
	}
}

func (c *idWorkerConfig) GetSequenceBits() int32 {
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

type tokenConfig struct {
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

func (c *tokenConfig) GetStore() string {
	if c.Store == nil {
		return "redis"
	} else {
		return *c.Store
	}
}

func (c *tokenConfig) GetBucket() string {
	if c.Bucket == nil {
		return "tokens"
	} else {
		return *c.Bucket
	}
}

func (c *tokenConfig) GetAccessTokenTTL() time.Duration {
	if c.AccessTokenTTL == nil {
		return 2 * 24 * time.Hour // default to 2 days
	} else {
		return *c.AccessTokenTTL
	}
}

func (c *tokenConfig) GetRefreshTokenTTL() time.Duration {
	if c.RefreshTokenTTL == nil {
		return 7 * 24 * time.Hour // default to 7 days
	} else {
		return *c.RefreshTokenTTL
	}
}

func (c *tokenConfig) GetIssuer() string {
	if c.Issuer == nil {
		return "unnamed-issuer"
	} else {
		return *c.Issuer
	}
}

func (c *tokenConfig) GetAudience() string {
	if c.Audience == nil {
		return "unnamed-audience"
	} else {
		return *c.Audience
	}
}

func (c *tokenConfig) GetSigningMethod() string {
	if c.SigningMethod == nil {
		return "RS256"
	} else {
		return *c.SigningMethod
	}
}

func (c *tokenConfig) GetPublicKey() []byte {
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

func (c *tokenConfig) GetPrivateKey() []byte {
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
