package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func APIKeyMiddleware(next http.Handler, logger *logrus.Logger) http.Handler {
	mwLogger := logger.WithFields(logrus.Fields{
		"layer": "middleware",
	})
	// Skip API key check for Swagger docs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/docs/") {
			next.ServeHTTP(w, r)
			return
		}
		apiKey := r.Header.Get("X-API-KEY")

		if apiKey != os.Getenv("API_KEY") {
			mwLogger.Error("Invalid API KEY")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
