package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int
	MinConns int
}

type Config struct {
	Database  DatabaseConfig
	Port      string
	LogLevel  string // "debug", "info", "warn", "error"
	LogPretty bool   // Pretty logging for development
}

func Load() (*Config, error) {
	port, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	maxConns, _ := strconv.Atoi(getEnv("DB_MAX_CONNS", "25"))
	minConns, _ := strconv.Atoi(getEnv("DB_MIN_CONNS", "5"))

	// Parse log pretty (default: true for development)
	logPretty := getEnv("LOG_PRETTY", "true") == "true"

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     port,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "gin_crud_api"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			MaxConns: maxConns,
			MinConns: minConns,
		},
		Port:      getEnv("SERVER_PORT", "8080"),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogPretty: logPretty,
	}, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}