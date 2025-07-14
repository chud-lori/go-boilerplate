package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/chud-lori/go-boilerplate/adapters/web/dto"
	"github.com/chud-lori/go-boilerplate/adapters/web/helper"
	"github.com/sirupsen/logrus"
)

func RecoveryMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logger.Infof("PANIC: %v\n%s", rvr, debug.Stack())

				helper.WriteResponse(w, dto.WebResponse{
					Message: "An unexpected error occurred",
					Status:  0,
					Data:    nil,
				}, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
