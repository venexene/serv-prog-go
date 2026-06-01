package conf

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Cfg struct {
	Port     string
	DBPath   string
	JWTKey   string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	LogLevel logrus.Level
}

func LoadCfg() (*Cfg, error) {
	_ = godotenv.Load()

	cfg := &Cfg{
		Port:   envOr("PORT", "8080"),
		DBPath: envOr("DB_PATH", "data/gravity.db"),
		JWTKey: envOr("JWT_SECRET", "super-secret-key-change-in-production"),
	}

	access, err := time.ParseDuration(envOr("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		return nil, err
	}
	cfg.JWTAccessTTL = access

	refresh, err := time.ParseDuration(envOr("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		return nil, err
	}
	cfg.JWTRefreshTTL = refresh

	lvl, err := logrus.ParseLevel(envOr("LOG_LEVEL", "debug"))
	if err != nil {
		lvl = logrus.InfoLevel
	}
	cfg.LogLevel = lvl

	return cfg, nil
}

func envOr(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}