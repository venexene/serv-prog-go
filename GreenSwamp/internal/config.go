package internal

import "os"

type Config struct {
	Port      string
	TplDir    string
	StaticDir string
}

func LoadConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		TplDir:    getEnv("TEMPLATE_DIR", "templates"),
		StaticDir: getEnv("STATIC_DIR", "static"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}