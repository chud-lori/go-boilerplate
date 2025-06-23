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
	"github.com/chud-lori/go-boilerplate/pkg/auth"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Boilerplate API
// @version 1.0
// @description A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project includes best practices for clean architecture, modularity, testing, and observability.

// @BasePath /api

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name X-API-KEY

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer {token}"
func main() {
	// ========== Environment Setup ==========

	if err := godotenv.Load(); err != nil {
		log.Fatal("Failed to load environment variables")
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	baseLogger := logger.NewLogger(cfg.LogLevel)

	// ========== Infrastructure ==========

	db, err := datastore.NewDatabase(cfg.DatabaseURL, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	cache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to cache server:", err)
	}
	defer cache.Close()

	ctxTimeout := 60 * time.Second

	// ========== Domain Service Dependencies ==========

	encryptor := &auth.BcryptEncryptor{}
	tokenManager := &auth.JWTManager{
		SecretKey:  cfg.JwtSecret,
		Expiration: 24 * time.Hour,
	}

	// ========== Repositories ==========

	userRepo := &repositories.UserRepositoryPostgre{DB: db}

	// ========== Services ==========

	authService := &services.AuthServiceImpl{
		DB:             db,
		UserRepository: userRepo,
		Encryptor:      encryptor,
		TokenManager:   tokenManager,
		CtxTimeout:     ctxTimeout,
	}

	userService := &services.UserServiceImpl{
		DB:             db,
		UserRepository: userRepo,
		Encryptor:      encryptor,
		Cache:          cache,
		CtxTimeout:     ctxTimeout,
	}

	// ========== Controllers ==========

	authController := &controllers.AuthController{
		AuthService: authService,
	}

	userController := &controllers.UserController{
		UserService: userService,
	}

	// ========== Routers ==========
	// router := http.NewServeMux()

	// // Public routes (no JWT middleware)
	// if cfg.AppEnv != "production" {
	// 	router.Handle("/docs/", httpSwagger.WrapHandler)
	// }
	// web.AuthRouter(authController, router)

	// // Protected routes (with JWT middleware)
	// protectedRouter := http.NewServeMux()
	// web.UserRouter(userController, protectedRouter)

	// // Apply JWT middleware only to protected routes
	// var protectedHandler http.Handler = protectedRouter
	// protectedHandler = middleware.JWTMiddleware(protectedHandler, tokenManager, baseLogger)

	// // Mount the protected handler to the main router
	// router.Handle("/api/user/", http.StripPrefix("/api/user", protectedHandler))
	// router.Handle("/api/user", protectedHandler)

	// ========== Routers ==========
	router := http.NewServeMux()

	// Documentation (no /api prefix)
	if cfg.AppEnv != "production" {
		router.Handle("/docs/", httpSwagger.WrapHandler)
	}

	// Create a single API router that will contain all API routes
	apiRouter := http.NewServeMux()

	// Auth routes (public)
	web.AuthRouter(authController, apiRouter)

	// User routes (protected)
	userRouter := http.NewServeMux()
	web.UserRouter(userController, userRouter)
	var protectedUserHandler http.Handler = userRouter
	protectedUserHandler = middleware.JWTMiddleware(protectedUserHandler, tokenManager, baseLogger)

	// Mount user routes to API router
	apiRouter.Handle("/user", protectedUserHandler)
	apiRouter.Handle("/user/", protectedUserHandler)

	// When you add more domains, just add them to apiRouter:
	// productRouter := http.NewServeMux()
	// web.ProductRouter(productController, productRouter)
	// apiRouter.Handle("/product", productRouter)
	// apiRouter.Handle("/product/", productRouter)

	// orderRouter := http.NewServeMux()
	// web.OrderRouter(orderController, orderRouter)
	// var protectedOrderHandler http.Handler = orderRouter
	// protectedOrderHandler = middleware.JWTMiddleware(protectedOrderHandler, tokenManager, baseLogger)
	// apiRouter.Handle("/order", protectedOrderHandler)
	// apiRouter.Handle("/order/", protectedOrderHandler)

	// Single mount point for all API routes
	router.Handle("/api/", http.StripPrefix("/api", apiRouter))

	// ========== Global Middleware Chain ==========

	var handler http.Handler = router
	handler = middleware.LogTrafficMiddleware(handler, baseLogger)
	handler = middleware.APIKeyMiddleware(handler, cfg.APIKey, baseLogger)
	// handler = middleware.JWTMiddleware(handler, tokenManager, baseLogger)

	// ========== HTTP Server Setup ==========

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: handler,
	}

	// Start server in goroutine
	go func() {
		utils.Banner(cfg)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// ========== Graceful Shutdown ==========

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
