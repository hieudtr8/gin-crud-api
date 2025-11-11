package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	GraphQLPort string `mapstructure:"graphql_port"` // GraphQL API server port
	RESTPort    string `mapstructure:"rest_port"`    // Legacy REST API server port
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`      // Database host
	Port     int    `mapstructure:"port"`      // Database port
	User     string `mapstructure:"user"`      // Database username
	Password string `mapstructure:"password"`  // Database password
	DBName   string `mapstructure:"dbname"`    // Database name
	SSLMode  string `mapstructure:"sslmode"`   // SSL mode (disable, require, verify-full)
	MaxConns int    `mapstructure:"max_conns"` // Maximum connection pool size
	MinConns int    `mapstructure:"min_conns"` // Minimum idle connections
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`  // Log level: debug, info, warn, error
	Pretty bool   `mapstructure:"pretty"` // Pretty console output vs JSON
}

// Config is the top-level configuration structure
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`   // Server configuration
	Database DatabaseConfig `mapstructure:"database"` // Database configuration
	Logging  LoggingConfig  `mapstructure:"logging"`  // Logging configuration
}

// LoadConfig loads configuration from YAML file and environment variables
// The env parameter determines which config file to load (dev, prod, test)
// Environment variables with GINAPI_ prefix override YAML values
func LoadConfig(env string) (*Config, error) {
	// Default to development environment
	if env == "" {
		env = "dev"
	}

	v := viper.New()

	// Configure config file settings
	v.SetConfigName(env)        // e.g., dev, prod, test
	v.SetConfigType("yaml")     // YAML format
	v.AddConfigPath("./configs") // Look in configs directory
	v.AddConfigPath("../configs") // For tests run from subdirectories
	v.AddConfigPath(".")         // Fallback to current directory

	// Configure environment variable settings
	v.SetEnvPrefix("GINAPI")     // Environment variables must start with GINAPI_
	v.AutomaticEnv()             // Automatically read environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Convert dots to underscores

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is acceptable - we'll use defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Log that we're using environment variables only
		fmt.Printf("Config file not found for environment '%s', using environment variables and defaults\n", env)
	} else {
		fmt.Printf("Loaded configuration from: %s\n", v.ConfigFileUsed())
	}

	// Unmarshal configuration into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &cfg, nil
}

// DSN returns the PostgreSQL connection string
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}