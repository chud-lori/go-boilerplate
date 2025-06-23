package ports

import (
	"net/http"
)

type AuthController interface {
	SignIn(writer http.ResponseWriter, request *http.Request)
	SignUp(writer http.ResponseWriter, request *http.Request)
}
