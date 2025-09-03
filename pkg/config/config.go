package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App         App
	Gateway     GatewayConfig
	UserService UserServiceConfig
	Redis       RedisConfig
	Database    DatabaseConfig
}

type App struct {
	Env             string
	WebDomain       string
	AuthSessionTtl  time.Duration
	AuthSessionsTtl time.Duration
	AuthExetendTtl  time.Duration
}

type GatewayConfig struct {
	Port string
}

type UserServiceConfig struct {
	Address string
	Timeout time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Loading config")

	return &Config{
		App: App{
			Env:             getEnv("APP_ENV", "dev"),
			WebDomain:       getEnv("APP_WEB_DOMAIN", "localhost"),
			AuthSessionTtl:  getEnvDuration("APP_AUTH_SESSION_TTL", 24*time.Hour),
			AuthSessionsTtl: getEnvDuration("APP_AUTH_SESSIONS_TTL", 7*24*time.Hour),
			AuthExetendTtl:  getEnvDuration("APP_AUTH_EXTEND_TTL", 30*time.Minute),
		},
		Gateway: GatewayConfig{
			Port: getEnv("GATEWAY_PORT", "8080"),
		},
		UserService: UserServiceConfig{
			Address: getEnv("USER_SERVICE_ADDR", "localhost:50051"),
			Timeout: getEnvDuration("USER_SERVICE_TIMEOUT", 10*time.Second),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "booking"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 20),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 1*time.Hour),
			ConnMaxIdleTime: getEnvDuration("DB_CONN_MAX_IDLE_TIME", 15*time.Minute),
		},
	}
}

// getEnv returns environment variable or fallback value
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// getEnvInt returns environment variable as int or fallback value
func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

// getEnvDuration returns environment variable as time.Duration or fallback value
func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
