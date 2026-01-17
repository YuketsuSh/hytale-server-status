package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Hytale     HytaleConfig     `yaml:"hytale"`
	Cache      CacheConfig      `yaml:"cache"`
	Logging    LoggingConfig    `yaml:"logging"`
	Security   SecurityConfig   `yaml:"security"`
	Monitoring MonitoringConfig `yaml:"monitoring"`
}

type ServerConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`
}

type HytaleConfig struct {
	DefaultPort    int           `yaml:"default_port" env:"HYTALE_PORT" env-default:"5520"`
	Timeout        time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"10s"`
	MaxConnections int           `yaml:"max_connections" env:"MAX_CONNECTIONS" env-default:"100"`
	UserAgent      string        `yaml:"user_agent" env:"USER_AGENT" env-default:"HytaleStatusDaemon/1.0"`
}

type CacheConfig struct {
	TTL             time.Duration `yaml:"ttl" env:"CACHE_TTL" env-default:"30s"`
	MaxEntries      int           `yaml:"max_entries" env:"CACHE_MAX_ENTRIES" env-default:"1000"`
	CleanupInterval time.Duration `yaml:"cleanup_interval" env:"CACHE_CLEANUP_INTERVAL" env-default:"60s"`
}

type LoggingConfig struct {
	Level  string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`
	Output string `yaml:"output" env:"LOG_OUTPUT" env-default:"stdout"`
}

type SecurityConfig struct {
	RateLimit      int      `yaml:"rate_limit" env:"RATE_LIMIT" env-default:"100"`
	TrustedProxies []string `yaml:"trusted_proxies" env:"TRUSTED_PROXIES" env-separator:","`
	CORSOrigins    []string `yaml:"cors_origins" env:"CORS_ORIGINS" env-separator:"," env-default:"*"`
}

type MonitoringConfig struct {
	EnableMetrics bool `yaml:"enable_metrics" env:"ENABLE_METRICS" env-default:"true"`
	MetricsPort   int  `yaml:"metrics_port" env:"METRICS_PORT" env-default:"9090"`
	HealthCheck   bool `yaml:"health_check" env:"HEALTH_CHECK" env-default:"true"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Hytale.Timeout == 0 {
		cfg.Hytale.Timeout = 10 * time.Second
	}
	if cfg.Cache.TTL == 0 {
		cfg.Cache.TTL = 30 * time.Second
	}
	if cfg.Cache.CleanupInterval == 0 {
		cfg.Cache.CleanupInterval = 60 * time.Second
	}

	return &cfg, nil
}
