package web

import (
	"net/http"

	"github.com/chud-lori/go-boilerplate/adapters/controllers"
	"github.com/chud-lori/go-boilerplate/adapters/middleware"
	"github.com/chud-lori/go-boilerplate/domain/ports"
	"github.com/sirupsen/logrus"
)

func UserRouter(controller *controllers.UserController, serve *http.ServeMux) {
	serve.HandleFunc("POST /user", controller.Create)
	serve.HandleFunc("PUT /user/{userId}", controller.Update)
	serve.HandleFunc("DELETE /user/{userId}", controller.Delete)
	serve.HandleFunc("GET /user/{userId}", controller.FindById)
	serve.HandleFunc("GET /user", controller.FindAll)
}

func AuthRouter(controller *controllers.AuthController, serve *http.ServeMux) {
	serve.HandleFunc("POST /signin", controller.SignIn)
	serve.HandleFunc("POST /signup", controller.SignUp)
}

func PostRouter(controller *controllers.PostController, serve *http.ServeMux, tokenManager ports.TokenManager, logger *logrus.Logger) {
	// Protected endpoints
	createHandler := middleware.JWTMiddleware(http.HandlerFunc(controller.Create), tokenManager, logger)
	serve.Handle("POST /post", createHandler)

	updateHandler := middleware.JWTMiddleware(http.HandlerFunc(controller.Update), tokenManager, logger)
	serve.Handle("PUT /post/{postId}", updateHandler)

	deleteHandler := middleware.JWTMiddleware(http.HandlerFunc(controller.Delete), tokenManager, logger)
	serve.Handle("DELETE /post/{postId}", deleteHandler)

	uploadHandler := middleware.JWTMiddleware(http.HandlerFunc(controller.UploadAttachment), tokenManager, logger)
	serve.Handle("POST /post/{postId}/upload", uploadHandler)

	// Public endpoints
	serve.HandleFunc("GET /post/{postId}", controller.GetById)
	serve.HandleFunc("GET /post", controller.GetAll)
	serve.HandleFunc("GET /uploads/{uploadId}/events", controller.UploadStatusSSE)
}
