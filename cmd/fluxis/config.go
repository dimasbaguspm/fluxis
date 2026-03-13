package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	authConfig "github.com/dimasbaguspm/fluxis/internal/auth/service"
	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/cors"
	"github.com/dimasbaguspm/fluxis/pkg/postgres"
	ratelimit "github.com/dimasbaguspm/fluxis/pkg/rate-limit"
)

type Config struct {
	Env       string
	DB        postgres.Config
	Server    ServerConfig
	Auth      authConfig.Config
	DataCache cache.Config
	RateLimit ratelimit.Config
	CORS      cors.Config
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (c ServerConfig) addr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

func LoadEnv() *Config {
	slog.Info("[Config]: Attempting to load few environment variables")

	cfg := &Config{
		Env: getEnv("ENV", "development"),
		Server: ServerConfig{
			Host:         getEnv("HOST", "0.0.0.0"),
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDuration("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		DB: postgres.Config{
			Primary:  mustEnv("DATABASE_URL"),
			MaxConns: getInt("DB_MAX_CONNS", 25),
			MinConns: getInt("DB_MIN_CONNS", 5),
		},
		Auth: authConfig.Config{
			AccessTokenSecret:  mustEnv("JWT_ACCESS_SECRET"),
			RefreshTokenSecret: mustEnv("JWT_REFRESH_SECRET"),
			AccessTokenExpiry:  getDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			BcryptCost:         getInt("BCRYPT_COST", 12),
		},
		DataCache: cache.Config{
			DefaultTTL: getDuration("CACHE_DEFAULT_TTL", 15*time.Minute),
			HMACKey:    mustEnv("CACHE_HMAC_KEY"),
		},
		RateLimit: ratelimit.Config{
			MaxRequests: getInt("RATE_LIMIT_MAX_REQUESTS", 100),
			Window:      getDuration("RATE_LIMIT_WINDOW", 1*time.Minute),
		},
		CORS: cors.Config{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
			AllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
			AllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization"),
			AllowedMaxAge:  getInt("CORS_MAX_AGE", 3600),
		},
	}

	slog.Info(fmt.Sprintf("[Config]: Environment %s is established", cfg.Env))
	return cfg
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("[Config]: Required environment variable %q is not set", key))
	}
	return v
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	slog.Info(fmt.Sprintf("[Config]: Env %s is missing, using '%s' as a fallback", key, fallback))
	return fallback
}

func getInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("[Config]: Env var %q must be an integer, got %q", key, v))
	}
	return n
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		panic(fmt.Sprintf("[Config]: Env var %q must be a duration (e.g. 15m, 7h), got %q", key, v))
	}
	return d
}
