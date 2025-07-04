package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APIKey        string
	DatabaseURL   string
	ServerPort    int
	AppEnv        string
	CtxTimeout    int
	LogLevel      string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
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

	cfg.MailServer = os.Getenv("MAIL_GRPC_SERVER")

	// Version
	cfg.Version = "1.0.0"
	cfg.JwtSecret = os.Getenv("JWT_SECRET")

	return cfg, nil
}

type MailConfig struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

var Mail MailConfig

func LoadMailConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found or error reading it")
	}

	port, err := strconv.Atoi(os.Getenv("GRPC_MAIL_PORT"))
	if err != nil {
		port = 2525 // default fallback
	}

	Mail = MailConfig{
		Host: os.Getenv("GRPC_MAIL_HOST"),
		Port: port,
		User: os.Getenv("GRPC_MAIL_USER"),
		Pass: os.Getenv("GRPC_MAIL_PASS"),
		From: os.Getenv("GRPC_MAIL_FROM"),
	}
}
