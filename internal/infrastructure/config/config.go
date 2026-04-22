package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort          string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSLMode           string
	JWTSecret           string
	JWTInternalSecret   string
	WalletBaseURL       string
}

func Load() *Config {
	_ = godotenv.Load()
	return &Config{
		ServerPort:        getEnv("SERVER_PORT", "3002"),
		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            mustGetEnv("DB_USER"),
		DBPassword:        mustGetEnv("DB_PASSWORD"),
		DBName:            mustGetEnv("DB_NAME"),
		DBSSLMode:         getEnv("DB_SSLMODE", "disable"),
		JWTSecret:         mustGetEnv("ILIACHALLENGE"),
		JWTInternalSecret: mustGetEnv("ILIACHALLENGE_INTERNAL"),
		WalletBaseURL:     mustGetEnv("WALLET_BASE_URL"),
	}
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic(fmt.Sprintf("required env var %q is not set", key))
	}
	return v
}
