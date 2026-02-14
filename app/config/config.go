package config

// Config holds application configuration loaded from environment or file.
type Config struct {
	HTTPPort string
	WSPath   string
}

// Load reads configuration and returns Config. Source can be extended (env, file).
func Load() (*Config, error) {
	return &Config{
		HTTPPort: "8080",
		WSPath:   "/ws",
	}, nil
}
