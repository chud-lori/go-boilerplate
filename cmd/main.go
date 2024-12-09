package main

import (
	//"context"
	"fmt"
	"os"

	//"log"
	//"math/rand"
	"boilerplate/adapters/controllers"
	"boilerplate/adapters/repositories"
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

	fmt.Println("App running on port ", os.Getenv("APP_PORT"))

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
