package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_DevEnvironment(t *testing.T) {
	// Change to project root to access configs directory
	originalDir, _ := os.Getwd()
	os.Chdir("../../") // Go to project root
	defer os.Chdir(originalDir)

	// Load dev configuration
	cfg, err := LoadConfig("dev")

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify server config loaded from dev.yaml
	assert.NotEmpty(t, cfg.Server.GraphQLPort)
	assert.NotEmpty(t, cfg.Server.RESTPort)

	// Verify database config
	assert.NotEmpty(t, cfg.Database.Host)
	assert.Greater(t, cfg.Database.Port, 0)
	assert.NotEmpty(t, cfg.Database.User)
	assert.NotEmpty(t, cfg.Database.DBName)
	assert.NotEmpty(t, cfg.Database.SSLMode)

	// Verify logging config
	assert.NotEmpty(t, cfg.Logging.Level)
}

func TestLoadConfig_ProdEnvironment(t *testing.T) {
	originalDir, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalDir)

	// Load prod configuration
	cfg, err := LoadConfig("prod")

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify different values from dev
	assert.NotEmpty(t, cfg.Server.GraphQLPort)
	assert.NotEmpty(t, cfg.Database.Host)
}

func TestLoadConfig_TestEnvironment(t *testing.T) {
	originalDir, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalDir)

	// Load test configuration
	cfg, err := LoadConfig("test")

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify config loaded
	assert.NotEmpty(t, cfg.Server.GraphQLPort)
	assert.NotEmpty(t, cfg.Database.Host)
}

func TestLoadConfig_DefaultsToDevWhenEmpty(t *testing.T) {
	originalDir, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalDir)

	// Load config with empty environment (should default to dev)
	cfg, err := LoadConfig("")

	// Assert success - should load dev config
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Should have loaded valid configuration
	assert.NotEmpty(t, cfg.Server.GraphQLPort)
	assert.NotEmpty(t, cfg.Database.Host)
}

func TestLoadConfig_NonExistentEnvironment(t *testing.T) {
	// Try to load non-existent environment
	cfg, err := LoadConfig("nonexistent")

	// Should not error (will use env vars and defaults)
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestLoadConfig_EnvironmentVariableOverride(t *testing.T) {
	originalDir, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalDir)

	// Set environment variables to override config
	customPort := "9999"
	customHost := "custom-host"

	os.Setenv("GINAPI_SERVER_GRAPHQL_PORT", customPort)
	os.Setenv("GINAPI_DATABASE_HOST", customHost)

	// Cleanup after test
	defer func() {
		os.Unsetenv("GINAPI_SERVER_GRAPHQL_PORT")
		os.Unsetenv("GINAPI_DATABASE_HOST")
	}()

	// Load config
	cfg, err := LoadConfig("dev")

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify environment variables override YAML values
	assert.Equal(t, customPort, cfg.Server.GraphQLPort)
	assert.Equal(t, customHost, cfg.Database.Host)
}

func TestLoadConfig_MultipleEnvVarOverrides(t *testing.T) {
	originalDir, _ := os.Getwd()
	os.Chdir("../../")
	defer os.Chdir(originalDir)

	// Set multiple environment variables
	os.Setenv("GINAPI_SERVER_GRAPHQL_PORT", "8888")
	os.Setenv("GINAPI_SERVER_REST_PORT", "7777")
	os.Setenv("GINAPI_DATABASE_HOST", "test-db")
	os.Setenv("GINAPI_DATABASE_PORT", "5555")
	os.Setenv("GINAPI_DATABASE_USER", "testuser")
	os.Setenv("GINAPI_DATABASE_PASSWORD", "testpass")
	os.Setenv("GINAPI_DATABASE_DBNAME", "testdb")
	os.Setenv("GINAPI_DATABASE_SSLMODE", "require")
	os.Setenv("GINAPI_LOGGING_LEVEL", "debug")
	os.Setenv("GINAPI_LOGGING_PRETTY", "true")

	// Cleanup after test
	defer func() {
		os.Unsetenv("GINAPI_SERVER_GRAPHQL_PORT")
		os.Unsetenv("GINAPI_SERVER_REST_PORT")
		os.Unsetenv("GINAPI_DATABASE_HOST")
		os.Unsetenv("GINAPI_DATABASE_PORT")
		os.Unsetenv("GINAPI_DATABASE_USER")
		os.Unsetenv("GINAPI_DATABASE_PASSWORD")
		os.Unsetenv("GINAPI_DATABASE_DBNAME")
		os.Unsetenv("GINAPI_DATABASE_SSLMODE")
		os.Unsetenv("GINAPI_LOGGING_LEVEL")
		os.Unsetenv("GINAPI_LOGGING_PRETTY")
	}()

	// Load config
	cfg, err := LoadConfig("dev")

	// Assert success
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify all overrides
	assert.Equal(t, "8888", cfg.Server.GraphQLPort)
	assert.Equal(t, "7777", cfg.Server.RESTPort)
	assert.Equal(t, "test-db", cfg.Database.Host)
	assert.Equal(t, 5555, cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.DBName)
	assert.Equal(t, "require", cfg.Database.SSLMode)
	assert.Equal(t, "debug", cfg.Logging.Level)
	assert.True(t, cfg.Logging.Pretty)
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary directory for test config
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid.yaml")

	// Write invalid YAML
	err := os.WriteFile(configPath, []byte("invalid: yaml: content: :::"), 0644)
	require.NoError(t, err)

	// Create a temporary config directory
	configDir := filepath.Join(tempDir, "configs")
	err = os.Mkdir(configDir, 0755)
	require.NoError(t, err)

	// Copy invalid YAML to configs directory
	invalidPath := filepath.Join(configDir, "invalid.yaml")
	err = os.WriteFile(invalidPath, []byte("invalid: yaml: content: :::"), 0644)
	require.NoError(t, err)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tempDir)

	// Try to load invalid config
	cfg, err := LoadConfig("invalid")

	// Should error due to invalid YAML syntax
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "error reading config file")
}

func TestDatabaseConfig_DSN(t *testing.T) {
	testCases := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "Development DSN",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "password",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=password dbname=testdb sslmode=disable",
		},
		{
			name: "Production DSN with SSL",
			config: DatabaseConfig{
				Host:     "prod-db.example.com",
				Port:     5432,
				User:     "produser",
				Password: "securepass",
				DBName:   "proddb",
				SSLMode:  "require",
			},
			expected: "host=prod-db.example.com port=5432 user=produser password=securepass dbname=proddb sslmode=require",
		},
		{
			name: "Custom port",
			config: DatabaseConfig{
				Host:     "custom-host",
				Port:     5433,
				User:     "user",
				Password: "pass",
				DBName:   "db",
				SSLMode:  "verify-full",
			},
			expected: "host=custom-host port=5433 user=user password=pass dbname=db sslmode=verify-full",
		},
		{
			name: "Empty password",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password= dbname=testdb sslmode=disable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dsn := tc.config.DSN()
			assert.Equal(t, tc.expected, dsn)
		})
	}
}

func TestDatabaseConfig_DSN_Components(t *testing.T) {
	config := DatabaseConfig{
		Host:     "myhost",
		Port:     1234,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "mydb",
		SSLMode:  "require",
	}

	dsn := config.DSN()

	// Verify all components are in the DSN
	assert.Contains(t, dsn, "host=myhost")
	assert.Contains(t, dsn, "port=1234")
	assert.Contains(t, dsn, "user=myuser")
	assert.Contains(t, dsn, "password=mypassword")
	assert.Contains(t, dsn, "dbname=mydb")
	assert.Contains(t, dsn, "sslmode=require")
}
