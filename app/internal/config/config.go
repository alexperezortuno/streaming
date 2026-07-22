package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	DatabaseURL      string
	MediaPath        string
	JWTSecret        string
	JWTExpiration    time.Duration
	LogLevel         string
	CORSOrigins      []string
	MaxUploadSize    int64
	TranscodeWorkers int
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "3000"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/streaming?sslmode=disable"),
		MediaPath:        getEnv("MEDIA_PATH", "./media"),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration:    getDurationEnv("JWT_EXPIRATION", 72*time.Hour),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		CORSOrigins:      getSliceEnv("CORS_ORIGINS", []string{"*"}),
		MaxUploadSize:    getInt64Env("MAX_UPLOAD_SIZE", 500<<20), // 500MB
		TranscodeWorkers: int(getInt64Env("TRANSCODE_WORKERS", 2)),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		d, err := time.ParseDuration(val)
		if err == nil {
			return d
		}
	}
	return fallback
}

func getInt64Env(key string, fallback int64) int64 {
	if val := os.Getenv(key); val != "" {
		n, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return n
		}
	}
	return fallback
}

func getSliceEnv(key string, fallback []string) []string {
	if val := os.Getenv(key); val != "" {
		return []string{val}
	}
	return fallback
}
