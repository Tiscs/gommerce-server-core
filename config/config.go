package config

import (
	"log/slog"
	"time"

	"go.uber.org/dig"
)

type extractSectionsResult struct {
	dig.Out

	ServerConfig      ServerConfig
	ServerHTTPConfig  ServerHTTPConfig
	ServerDBConfig    ServerDBConfig
	ServerRedisConfig ServerRedisConfig
	ServerNATSConfig  ServerNATSConfig
	IdWorkerConfig    IdWorkerConfig
	LoggingConfig     LoggingConfig
	TraceConfig       TraceConfig
	MetricConfig      MetricConfig
	TokenConfig       TokenConfig
}

// ExtractSections extracts sections from RootConfig.
// It is used for dependency injection.
func ExtractSections(cfg RootConfig) extractSectionsResult {
	result := extractSectionsResult{
		ServerConfig:      cfg.GetServerConfig(),
		ServerHTTPConfig:  cfg.GetServerConfig().GetHTTPConfig(),
		ServerDBConfig:    cfg.GetServerConfig().GetDBConfig(),
		ServerRedisConfig: cfg.GetServerConfig().GetRedisConfig(),
		ServerNATSConfig:  cfg.GetServerConfig().GetNATSConfig(),
		IdWorkerConfig:    cfg.GetIdWorkerConfig(),
		LoggingConfig:     cfg.GetLoggingConfig(),
		TraceConfig:       cfg.GetTraceConfig(),
		MetricConfig:      cfg.GetMetricConfig(),
		TokenConfig:       cfg.GetTokenConfig(),
	}
	return result
}

type RootConfig interface {
	GetServerConfig() ServerConfig
	GetIdWorkerConfig() IdWorkerConfig
	GetLoggingConfig() LoggingConfig
	GetTraceConfig() TraceConfig
	GetMetricConfig() MetricConfig
	GetTokenConfig() TokenConfig
}

type ServerConfig interface {
	GetDebug() bool
	GetName() string
	GetVersion() string
	GetInstanceId() string
	GetHTTPConfig() ServerHTTPConfig
	GetDBConfig() ServerDBConfig
	GetRedisConfig() ServerRedisConfig
	GetNATSConfig() ServerNATSConfig
}

type ServerHTTPConfig interface {
	GetAddr() string
}

type ServerDBConfig interface {
	GetDriver() string
	GetSource() string
}

type ServerRedisConfig interface {
	GetInitAddress() []string
	GetSelectDB() int
}

type ServerNATSConfig interface {
	GetURL() string
}

type IdWorkerConfig interface {
	GetIdEpoch() int64
	GetClusterId() int64
	GetWorkerId() int64
	GetWorkerSeqKey() string
	GetClusterIdBits() int32
	GetWorkerIdBits() int32
	GetSequenceBits() int32
}

type LoggingConfig interface {
	GetZapLogger() LoggingZapLoggerConfig
	GetSlogLogger() LoggingSlogLoggerConfig
}

type LoggingZapLoggerConfig interface {
	GetPreset() string
}

type LoggingSlogLoggerConfig interface {
	GetHandler() string
	GetAddSource() bool
	GetLeveler() slog.Leveler
}

type TraceConfig interface {
	GetExporterConfig() TraceExporterConfig
}

type TraceExporterConfig interface {
	GetProtocol() string
	GetEndpoint() string
	GetInsecure() bool
}

type MetricConfig interface {
	GetExporterConfig() MetricExporterConfig
}

type MetricExporterConfig interface {
	GetProtocol() string
	GetEndpoint() string
	GetInsecure() bool
}

type TokenConfig interface {
	GetStore() string
	GetBucket() string
	GetAccessTokenTTL() time.Duration
	GetRefreshTokenTTL() time.Duration
	GetIssuer() string
	GetAudience() string
	GetSigningMethod() string
	GetPublicKey() []byte
	GetPrivateKey() []byte
}
