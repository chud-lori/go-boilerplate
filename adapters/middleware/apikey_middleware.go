package middleware

import (
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

func APIKeyMiddleware(next http.Handler, apiKey string, logger *logrus.Logger) http.Handler {
	mwLogger := logger.WithFields(logrus.Fields{
		"layer": "middleware",
	})
	// Skip API key check for Swagger docs
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/docs/") {
			next.ServeHTTP(w, r)
			return
		}
		reqApiKey := r.Header.Get("X-API-KEY")

		if reqApiKey != apiKey {
			mwLogger.Error("Invalid API KEY")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
