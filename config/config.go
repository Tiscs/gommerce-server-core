package config

import (
	"log/slog"
	"time"

	"go.uber.org/fx"
)

type extractSectionsResult struct {
	fx.Out

	ServerConfig      ServerConfig
	ServerHTTPConfig  ServerHTTPConfig
	ServerDBConfig    ServerDBConfig
	ServerRedisConfig ServerRedisConfig
	ServerMinIOConfig ServerMinIOConfig
	ServerNATSConfig  ServerNATSConfig
	SnowflakeConfig   SnowflakeConfig
	LoggingConfig     LoggingConfig
	TraceConfig       TraceConfig
	MetricConfig      MetricConfig
	SecureConfig      SecureConfig
	SecureTokenConfig SecureTokenConfig
}

// ExtractSections extracts sections from RootConfig.
// It is used for dependency injection.
func ExtractSections(cfg RootConfig) extractSectionsResult {
	result := extractSectionsResult{
		ServerConfig:      cfg.GetServerConfig(),
		ServerHTTPConfig:  cfg.GetServerConfig().GetHTTPConfig(),
		ServerDBConfig:    cfg.GetServerConfig().GetDBConfig(),
		ServerRedisConfig: cfg.GetServerConfig().GetRedisConfig(),
		ServerMinIOConfig: cfg.GetServerConfig().GetMinIOConfig(),
		ServerNATSConfig:  cfg.GetServerConfig().GetNATSConfig(),
		SnowflakeConfig:   cfg.GetSnowflakeConfig(),
		LoggingConfig:     cfg.GetLoggingConfig(),
		TraceConfig:       cfg.GetTraceConfig(),
		MetricConfig:      cfg.GetMetricConfig(),
		SecureConfig:      cfg.GetSecureConfig(),
		SecureTokenConfig: cfg.GetSecureConfig().GetToken(),
	}
	return result
}

type RootConfig interface {
	GetServerConfig() ServerConfig
	GetSnowflakeConfig() SnowflakeConfig
	GetLoggingConfig() LoggingConfig
	GetTraceConfig() TraceConfig
	GetMetricConfig() MetricConfig
	GetSecureConfig() SecureConfig
}

type ServerConfig interface {
	GetDebug() bool
	GetName() string
	GetVersion() string
	GetInstanceId() string
	GetHTTPConfig() ServerHTTPConfig
	GetDBConfig() ServerDBConfig
	GetRedisConfig() ServerRedisConfig
	GetMinIOConfig() ServerMinIOConfig
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
	GetInitAddr() string
	GetSelectDB() int
}

type ServerMinIOConfig interface {
	GetEndpoint() string
	GetAccessKey() string
	GetSecretKey() string
	GetUseSSL() bool
}

type ServerNATSConfig interface {
	GetSeedURL() string
	GetNoEcho() bool
}

type SnowflakeConfig interface {
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

type SecureConfig interface {
	GetToken() SecureTokenConfig
}

type SecureTokenConfig interface {
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
