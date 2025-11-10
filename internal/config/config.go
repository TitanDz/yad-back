package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string
	Env      string
	Database DatabaseConfig
	JWT      JWTConfig
	Datadog  DatadogConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	MaxConn  int
}

type JWTConfig struct {
	Secret          string
	AccessDuration  int // minutes
	RefreshDuration int // hours
}

type DatadogConfig struct {
	Enabled     bool
	APIKey      string
	AppKey      string
	Environment string
	Service     string
}

func LoadConfig() *Config {
	// Set defaults
	viper.SetDefault("port", "8080")
	viper.SetDefault("env", "development")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "3306")
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.name", "yad_db")
	viper.SetDefault("database.maxconn", 10)

	// Read from .env
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	// Read from environment variables
	viper.AutomaticEnv()

	cfg := &Config{
		Port: getEnv("PORT", "8080"),
		Env:  getEnv("ENV", "development"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "yad_db"),
			MaxConn:  10,
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "dev-secret-key-change-in-production"),
			AccessDuration:  15,     // 15 minutes
			RefreshDuration: 7 * 24, // 7 days
		},
		Datadog: DatadogConfig{
			Enabled:     getEnv("DATADOG_ENABLED", "false") == "true",
			APIKey:      getEnv("DATADOG_API_KEY", ""),
			AppKey:      getEnv("DATADOG_APP_KEY", ""),
			Environment: getEnv("DATADOG_ENV", "development"),
			Service:     getEnv("DATADOG_SERVICE", "yad-backend"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
