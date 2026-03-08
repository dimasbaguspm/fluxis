package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	authConfig "github.com/dimasbaguspm/fluxis/internal/auth/service"
	"github.com/dimasbaguspm/fluxis/pkg/postgres"
	"github.com/dimasbaguspm/fluxis/pkg/redis"
)

type Config struct {
	Env    string
	DB     postgres.Config
	Redis  redis.Config
	Server ServerConfig
	Auth   authConfig.Config
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
		Redis: redis.Config{
			Addr:     getEnv("REDIS_ADDR", "redis:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Auth: authConfig.Config{
			AccessTokenSecret:  mustEnv("JWT_ACCESS_SECRET"),
			RefreshTokenSecret: mustEnv("JWT_REFRESH_SECRET"),
			AccessTokenExpiry:  getDuration("JWT_ACCESS_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour),
			BcryptCost:         getInt("BCRYPT_COST", 12),
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
