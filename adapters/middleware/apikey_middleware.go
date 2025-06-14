package middleware

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func APIKeyMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	mwLogger := logger.WithFields(logrus.Fields{
		"layer": "middleware",
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("x-api-key")
		if apiKey != os.Getenv("API_KEY") {
			mwLogger.Error("Invalid API KEY")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
