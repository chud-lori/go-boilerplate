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
	"github.com/chud-lori/go-boilerplate/infrastructure/grpc_clients"
	"github.com/chud-lori/go-boilerplate/infrastructure/api_clients"
	"github.com/chud-lori/go-boilerplate/internal/utils"
	"github.com/chud-lori/go-boilerplate/pkg/auth"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	db, err := datastore.NewPostgreDatabase(cfg.DatabaseURL, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	cache, err := cache.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB, baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to cache server:", err)
	}
	defer cache.Close()

	mailGrpcConn, err := grpc.NewClient(cfg.MailServer, grpc.WithTransportCredentials(insecure.NewCredentials())) // Use WithTransportCredentials for production
	if err != nil {
		log.Fatalf("did not connect to mail gRPC service: %v", err)
	}
	defer mailGrpcConn.Close()

	// To use the API-based mail client instead of gRPC, uncomment the following line and comment out the gRPC one above:
	// mailClient := api_clients.NewApiMailClient("http://localhost:8081/send-mail") // Replace with your real API endpoint
	mailClient := grpc_clients.NewGrpcMailClient(mailGrpcConn)

	ctxTimeout := 60 * time.Second

	// ========== Domain Service Dependencies ==========

	encryptor := &auth.BcryptEncryptor{}
	tokenManager := &auth.JWTManager{
		SecretKey:  cfg.JwtSecret,
		Expiration: 24 * time.Hour,
	}
	mailService := &services.MailServiceImpl{
		MailClient: mailClient,
	}

	// ========== Repositories ==========

	userRepo := &repositories.UserRepositoryPostgre{}
	postRepo := &repositories.PostRepositoryPostgre{}

	// ========== Services ==========

	authService := &services.AuthServiceImpl{
		DB:             db,
		UserRepository: userRepo,
		MailService:    mailService,
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

	postService := &services.PostServiceImpl{
		DB:             db,
		PostRepository: postRepo,
		UserRepository: userRepo,
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

	postController := &controllers.PostController{
		PostService: postService,
	}

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

	// Post routes (public + protected)
	web.PostRouter(postController, apiRouter, tokenManager, baseLogger)

	// User routes (protected)
	userRouter := http.NewServeMux()
	web.UserRouter(userController, userRouter)
	var protectedUserHandler http.Handler = userRouter
	protectedUserHandler = middleware.JWTMiddleware(protectedUserHandler, tokenManager, baseLogger)

	// Mount user routes to API router
	apiRouter.Handle("/user", protectedUserHandler)
	apiRouter.Handle("/user/", protectedUserHandler)

	// Single mount point for all API routes
	router.Handle("/api/", http.StripPrefix("/api", apiRouter))

	// ========== Global Middleware Chain ==========

	var handler http.Handler = router
	handler = middleware.LogTrafficMiddleware(handler, baseLogger)
	handler = middleware.APIKeyMiddleware(handler, cfg.APIKey, baseLogger)

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
