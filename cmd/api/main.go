package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/controllers"
	"github.com/chud-lori/go-boilerplate/adapters/middleware"
	"github.com/chud-lori/go-boilerplate/adapters/repositories"
	"github.com/chud-lori/go-boilerplate/adapters/web"
	"github.com/chud-lori/go-boilerplate/config"
	_ "github.com/chud-lori/go-boilerplate/docs"
	"github.com/chud-lori/go-boilerplate/domain/services"
	"github.com/chud-lori/go-boilerplate/infrastructure/cache"
	"github.com/chud-lori/go-boilerplate/infrastructure/datastore"
	"github.com/chud-lori/go-boilerplate/internal/utils"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Boilerplate API
// @version 1.0
// @description A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project includes best practices for clean architecture, modularity, testing, and observability.

// @BasePath /api

// @securityDefinitions.apiKey ApiKeyAuth
// @type apiKey
// @in header
// @name X-API-KEY
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed load keys")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	baseLogger := logger.NewLogger(cfg.LogLevel)

	db, err := datastore.NewDatabase(cfg.DatabaseURL, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	cache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to cache server: ", err)
	}
	defer cache.Close()

	ctxTimeout := time.Duration(60) * time.Second

	userRepository := &repositories.UserRepositoryPostgre{DB: db}
	userService := &services.UserServiceImpl{
		DB:             db,
		UserRepository: userRepository,
		Cache:          cache,
		CtxTimeout:     ctxTimeout,
	}
	userController := &controllers.UserController{
		UserService: userService,
	}

	router := http.NewServeMux()
	if cfg.AppEnv != "production" {
		router.Handle("/docs/", httpSwagger.WrapHandler)
	}

	web.UserRouter(userController, router)

	var handler http.Handler = router
	handler = middleware.LogTrafficMiddleware(handler, baseLogger)
	handler = middleware.APIKeyMiddleware(handler, cfg.APIKey, baseLogger)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: handler,
	}

	// Run server in a goroutine
	go func() {
		utils.Banner(cfg)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	wait := utils.GracefullShutdown(context.Background(), 5*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return db.Close()
		},
		"http-server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
		"cache": func(ctx context.Context) error {
			return cache.Close()
		},
	})

	<-wait
}
