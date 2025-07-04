package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type loggingTraffic struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingTraffic(w http.ResponseWriter) *loggingTraffic {
	return &loggingTraffic{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingTraffic) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

var sensitivePayloadKeys = map[string]struct{}{
	"password":        {},
	"credit_card_num": {},
	"cvv":             {},
	// Add more sensitive keys here
}

func LogTrafficMiddleware(next http.Handler, baseLogger *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		newLogger := baseLogger.WithField("RequestID", requestID)
		ctx := context.WithValue(r.Context(), "logger", newLogger)
		r = r.WithContext(ctx)

		// TODO: if showing source in log
		// baseLogger.SetReportCaller(true)
		//_, file, line, ok := runtime.Caller(1)
		//source := "unknown"
		//if ok {
		//    source = fmt.Sprintf("%s:%d", file, line)
		//}

		// --- Capture and Process Request Body ---
		var requestBodyLog interface{}
		var err error
		var reqBodyBytes []byte

		// Only read the body if it's a method that typically carries a body
		if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			reqBodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				newLogger.WithError(err).Warn("Failed to read request body")
			} else {
				// Restore the body for the next handler
				r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))

				// Attempt to unmarshal as JSON and exclude sensitive fields
				if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
					var jsonBody map[string]interface{}
					if err := json.Unmarshal(reqBodyBytes, &jsonBody); err == nil {
						// Exclude sensitive fields
						for key := range sensitivePayloadKeys {
							delete(jsonBody, key) // This is the change: delete the key
						}
						requestBodyLog = jsonBody
					} else {
						newLogger.WithError(err).Debug("Request body is not valid JSON, logging as raw string.")
						requestBodyLog = string(reqBodyBytes)
					}
				} else {
					// Log non-JSON bodies as raw string
					requestBodyLog = string(reqBodyBytes)
				}
			}
		}

		lrw := newLoggingTraffic(w)
		// call the next handler
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		logFields := logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": duration.String(),
			"status":   lrw.statusCode,
			"query":    r.URL.RawQuery,
		}

		if requestBodyLog != nil {
			logFields["request_body"] = requestBodyLog // Add processed request body
		}

		newLogger.WithFields(logFields).Info("Processed request")
	})
}
