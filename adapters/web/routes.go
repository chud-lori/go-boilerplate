package web

import (
	"net/http"

	"github.com/chud-lori/go-boilerplate/domain/ports"
)

func UserRouter(controller ports.UserController, serve *http.ServeMux) {
	serve.HandleFunc("POST /user", controller.Create)
	serve.HandleFunc("PUT /user/{userId}", controller.Update)
	serve.HandleFunc("DELETE /user/{userId}", controller.Delete)
	serve.HandleFunc("GET /user/{userId}", controller.FindById)
	serve.HandleFunc("GET /user", controller.FindAll)
}

func AuthRouter(controller ports.AuthController, serve *http.ServeMux) {
	serve.HandleFunc("POST /signin", controller.SignIn)
	serve.HandleFunc("POST /signup", controller.SignUp)
}
