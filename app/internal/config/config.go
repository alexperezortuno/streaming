package config

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

type cliFlags struct {
	envFile  string
	port     string
	dbURL    string
	logLevel string
}

func Load() *Config {
	var flags cliFlags

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.Usage = func() {
		fs.PrintDefaults()
	}
	fs.StringVar(&flags.envFile, "env-file", "", "Path to .env file")
	fs.StringVar(&flags.port, "port", "", "Server port")
	fs.StringVar(&flags.dbURL, "database-url", "", "PostgreSQL connection string")
	fs.StringVar(&flags.logLevel, "log-level", "", "Log level (debug, info, warn, error)")
	if err := fs.Parse(os.Args[1:]); err != nil {
		if err == flag.ErrHelp {
			os.Exit(0)
		}
	}

	loadEnvFiles(flags.envFile)

	cfg := &Config{
		Port:             firstNonEmpty(flags.port, getEnv("PORT", "3000")),
		DatabaseURL:      firstNonEmpty(flags.dbURL, getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/streaming?sslmode=disable")),
		JWTSecret:        getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiration:    getDurationEnv("JWT_EXPIRATION", 72*time.Hour),
		LogLevel:         firstNonEmpty(flags.logLevel, getEnv("LOG_LEVEL", "info")),
		CORSOrigins:      getSliceEnv("CORS_ORIGINS", []string{"*"}),
		MaxUploadSize:    getInt64Env("MAX_UPLOAD_SIZE", 500<<20),
		TranscodeWorkers: int(getInt64Env("TRANSCODE_WORKERS", 2)),
	}

	cfg.MediaPath = resolvePath(getEnv("MEDIA_PATH", "./media"))

	return cfg
}

func loadEnvFiles(flagEnvFile string) {
	if flagEnvFile != "" {
		tryLoadEnv(flagEnvFile)
		return
	}

	candidates := []string{
		".env",
		"bin/.env",
		"../.env",
		"../bin/.env",
	}

	loaded := make(map[string]bool)
	for _, path := range candidates {
		abs, err := filepath.Abs(path)
		if err != nil || loaded[abs] {
			continue
		}
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			continue
		}
		if err := godotenv.Load(abs); err != nil {
			slog.Warn("loading env file", "path", abs, "error", err)
			continue
		}
		loaded[abs] = true
		slog.Debug("loaded env file", "path", abs)
	}
}

func tryLoadEnv(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return
	}
	if _, err := os.Stat(abs); os.IsNotExist(err) {
		return
	}
	if err := godotenv.Load(abs); err != nil {
		slog.Warn("loading env file", "path", abs, "error", err)
		return
	}
	slog.Debug("loaded env file", "path", abs)
}

func resolvePath(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
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
