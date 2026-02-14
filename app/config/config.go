package config

import (
	"os"
	"strings"
)

// Config holds application configuration loaded from environment or file.
type Config struct {
	HTTPPort       string
	WSPath         string
	AllowedOrigins []string
}

// Load reads configuration and returns Config. Source can be extended (env, file).
func Load() (*Config, error) {
	origins := []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173", "http://127.0.0.1:3000"}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		origins = strings.Split(v, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
	}
	return &Config{
		HTTPPort:       "8080",
		WSPath:         "/ws",
		AllowedOrigins: origins,
	}, nil
}
