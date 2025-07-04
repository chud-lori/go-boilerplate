package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	APIKey        string
	DatabaseURL   string
	ServerPort    int
	AppEnv        string
	CtxTimeout    int
	LogLevel      string
	RedisAddr     string // New: Redis server address
	RedisPassword string // New: Redis password (can be empty)
	RedisDB       int    // New: Redis DB number
	Version       string
	JwtSecret     string
	MailServer    string
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{}

	cfg.APIKey = os.Getenv("API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("API_KEY environment variable not set")
	}

	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}
	cfg.ServerPort = port

	cfg.AppEnv = os.Getenv("APP_ENV")

	ctxTimoutStr := os.Getenv("CTX_TIMEOUT")
	ctxTimeout, err := strconv.Atoi(ctxTimoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CTX_TIMEOUT: %w", err)
	}
	cfg.CtxTimeout = ctxTimeout

	cfg.LogLevel = os.Getenv("LOG_LEVEL")

	// --- Redis Configuration ---
	cfg.RedisAddr = os.Getenv("REDIS_ADDR")
	if cfg.RedisAddr == "" {
		cfg.RedisAddr = "localhost:6379" // Default Redis address
	}

	cfg.RedisPassword = os.Getenv("REDIS_PASSWORD") // Can be empty

	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr == "" {
		redisDBStr = "0" // Default Redis DB
	}
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
	}
	cfg.RedisDB = redisDB

	cfg.MailServer = os.Getenv("MAIL_GRPC_SERVER") // Can be empty

	// Version
	cfg.Version = "1.0.0"
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	return cfg, nil
}
