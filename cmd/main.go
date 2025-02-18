package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chud-lori/go-boilerplate/adapters/controllers"
	"github.com/chud-lori/go-boilerplate/adapters/repositories"
	"github.com/chud-lori/go-boilerplate/adapters/utils"
	"github.com/chud-lori/go-boilerplate/adapters/web"
	"github.com/chud-lori/go-boilerplate/domain/services"
	"github.com/chud-lori/go-boilerplate/infrastructure/datastore"
	"github.com/chud-lori/go-boilerplate/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed load keys")
	}

	baseLogger := logger.NewLogger()

	db, err := datastore.NewDatabase(os.Getenv("DB_URL"), baseLogger)
	if err != nil {
		baseLogger.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	ctxTimeout := time.Duration(60) * time.Second

	userRepository, _ := repositories.NewUserRepositoryPostgre(db)
	userService := services.NewUserService(db, userRepository, ctxTimeout)
	userController := controllers.NewUserController(userService)

	router := http.NewServeMux()

	web.UserRouter(userController, router)

	var handler http.Handler = router
	handler = logger.LogTrafficMiddleware(handler, baseLogger)
	handler = utils.APIKeyMiddleware(handler, baseLogger)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: handler,
	}

	// Run server in a goroutine
	go func() {
        // log.Printf("Server is running on port %s", os.Getenv("APP_PORT"))
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
	})

	<-wait
}
