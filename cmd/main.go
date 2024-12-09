package main

import (
	//"context"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	//"log"
	//"math/rand"
	"boilerplate/adapters/controllers"
	"boilerplate/adapters/repositories"
	"boilerplate/adapters/utils"
	"boilerplate/adapters/web"
	"boilerplate/domain/services"
	"boilerplate/infrastructure"
	"boilerplate/pkg/logger"
	"net/http"

	"github.com/joho/godotenv"
)

//type Middleware func(http.HandlerFunc) http.HandlerFunc

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey != "secret-api-key" {
			//logger, _ := r.Context().Value("logger").(*logrus.Entry)
			//logger.Warn("Unauth bnruhhh")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Failed load keys")
	}

	postgredb := infrastructure.NewPostgreDB()
	defer postgredb.Close()

	userRepository, _ := repositories.NewUserRepositoryPostgre(postgredb)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	router := http.NewServeMux()

	web.UserRouter(userController, router)

	var handler http.Handler = router
	handler = logger.LogTrafficMiddleware(handler)
	handler = APIKeyMiddleware(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("APP_PORT")),
		Handler: handler,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", os.Getenv("APP_PORT"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	wait := utils.GracefullShutdown(context.Background(), 5*time.Second, map[string]utils.Operation{
		"database": func(ctx context.Context) error {
			return postgredb.Close()
		},
		"http-server": func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})

	<-wait
}
