package config

import "os"

type Config struct {
	Port   string
	DB_url string
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable: " + name)
	}
	return v
}

func NewConfig() *Config {
	return &Config{Port: getenv("PORT"), DB_url: getenv("DATABASE_URL")}
}
