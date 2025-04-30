package config

import (
	"os"
)

type PostgresConfig struct {
	DSN string
}

type Config struct {
	MigrateDIR string
	DB         PostgresConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NewConfig() *Config {
	return &Config{
		MigrateDIR: getEnv("MIGRATE_DIR", "./migrate"),
		DB: PostgresConfig{
			DSN: getEnv("DB_DSN", "postgres://postgres:postgres@db:5432/scraper?sslmode=disable"),
		},
	}
}
