package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Storage  StorageConfig  `yaml:"storage"`
	Auth     AuthConfig     `yaml:"auth"`
	Features FeaturesConfig `yaml:"features"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host        string   `yaml:"host"`
	Port        string   `yaml:"port"`
	CORSOrigins []string `yaml:"cors_origins"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Path string `yaml:"path"`
}

// StorageConfig represents storage configuration
type StorageConfig struct {
	OutputDir      string `yaml:"output_dir"`
	MaxScanAgeDays int    `yaml:"max_scan_age_days"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	JWTSecret   string `yaml:"jwt_secret"`
	TokenExpiry string `yaml:"token_expiry"`
}

// FeaturesConfig represents feature flags
type FeaturesConfig struct {
	WebSocketEnabled bool `yaml:"websocket_enabled"`
	AutoCleanup      bool `yaml:"auto_cleanup"`
}

// Load loads configuration from file or environment variables
func Load() (*Config, error) {
	// Default configuration
	cfg := &Config{
		Server: ServerConfig{
			Host:        getEnv("SERVER_HOST", "0.0.0.0"),
			Port:        getEnv("SERVER_PORT", "8080"),
			CORSOrigins: getEnvSlice("CORS_ORIGINS", []string{"http://localhost:5173"}),
		},
		Database: DatabaseConfig{
			Path: getEnv("DATABASE_PATH", "./data/scans.db"),
		},
		Storage: StorageConfig{
			OutputDir:      getEnv("OUTPUT_DIR", "./output"),
			MaxScanAgeDays: getEnvInt("MAX_SCAN_AGE_DAYS", 90),
		},
		Auth: AuthConfig{
			JWTSecret:   getEnv("JWT_SECRET", "change-this-secret-in-production"),
			TokenExpiry: getEnv("TOKEN_EXPIRY", "24h"),
		},
		Features: FeaturesConfig{
			WebSocketEnabled: getEnvBool("WEBSOCKET_ENABLED", true),
			AutoCleanup:      getEnvBool("AUTO_CLEANUP", true),
		},
	}

	// Try to load from config file
	configPath := getEnv("CONFIG_PATH", "config.yaml")
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return cfg, nil
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var result int
		fmt.Sscanf(value, "%d", &result)
		return result
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1"
	}
	return defaultValue
}
