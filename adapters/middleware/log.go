package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt" // Import fmt for Sprintf
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

// func (lrw *loggingTraffic) Flush() {
// 	if flusher, ok := lrw.ResponseWriter.(http.Flusher); ok {
// 		flusher.Flush()
// 	}
// 	// If the underlying ResponseWriter is not a Flusher, we can't flush,
// 	// but we don't want to panic or error here, just don't flush.
// 	// This might happen with certain HTTP servers or proxy setups.
// }

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

		var requestBodyLog interface{}
		var reqBodyBuffer bytes.Buffer // Buffer to hold a copy of the request body

		// Only process body for methods that typically carry one
		if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			contentType := r.Header.Get("Content-Type")

			// Always create a TeeReader to duplicate the body into reqBodyBuffer.
			// This ensures the body content is copied before being consumed by any subsequent reads
			// (e.g., by r.ParseMultipartForm, io.ReadAll, or the next handler).
			r.Body = io.NopCloser(io.TeeReader(r.Body, &reqBodyBuffer))

			// --- Content-Type Specific Body Logging Logic ---
			if strings.HasPrefix(contentType, "multipart/form-data") {
				// For multipart, we need to parse it to get specific text fields (like 'file_name', 'file_type').
				// This also means we specifically do NOT log the actual file content.
				// r.ParseMultipartForm reads from r.Body (our TeeReader), consuming it.
				err := r.ParseMultipartForm(32 << 20) // 32 MB limit for memory for non-file fields. Files are spooled to disk.
				if err != nil {
					newLogger.WithError(err).Warn("Failed to parse multipart/form-data for logging")
					requestBodyLog = "multipart_parse_failed"
				} else {
					// Initialize a map to hold only the non-file fields (text fields)
					parsedMultipartLog := make(map[string]interface{})

					// Log ONLY regular form fields (like 'file_name', 'file_type') from r.Form.
					// r.Form is populated by ParseMultipartForm.
					for key, values := range r.Form {
						if _, sensitive := sensitivePayloadKeys[strings.ToLower(key)]; sensitive {
							parsedMultipartLog[key] = "[SENSITIVE]"
						} else {
							if len(values) == 1 {
								parsedMultipartLog[key] = values[0]
							} else {
								parsedMultipartLog[key] = values // Log all values if multiple
							}
						}
					}

					// Add a note if files were present but skipped for logging
					if r.MultipartForm != nil && len(r.MultipartForm.File) > 0 {
						parsedMultipartLog["_files_present_skipped_logging"] = true
					}
					requestBodyLog = parsedMultipartLog
				}

			} else if strings.HasPrefix(contentType, "application/") &&
				!strings.Contains(contentType, "json") && // Explicitly allow application/json
				!strings.Contains(contentType, "x-www-form-urlencoded") { // Explicitly allow application/x-www-form-urlencoded
				// This condition catches typical binary application types (e.g., application/pdf,
				// application/octet-stream, application/zip, application/protobuf, etc.)
				// The body content has already been copied to reqBodyBuffer by TeeReader,
				// so we don't need to explicitly read from r.Body again here.
				newLogger.Debugf("Skipping request body logging for binary Content-Type: %s", contentType)
				requestBodyLog = fmt.Sprintf("body_skipped_for_type_%s", strings.ReplaceAll(contentType, "/", "_"))

			} else {
				// This handles:
				// - application/json (explicitly allowed above)
				// - application/x-www-form-urlencoded (explicitly allowed above)
				// - text/* types (e.g., text/plain, text/html)
				// The full body content is available in reqBodyBuffer.
				reqBodyBytes := reqBodyBuffer.Bytes()

				if strings.Contains(contentType, "application/json") {
					var jsonBody map[string]interface{}
					if err := json.Unmarshal(reqBodyBytes, &jsonBody); err == nil {
						// Apply sensitive key filtering for JSON
						for key := range sensitivePayloadKeys {
							delete(jsonBody, key)
						}
						requestBodyLog = jsonBody
					} else {
						newLogger.WithError(err).Debug("Request body is not valid JSON, logging as raw string.")
						requestBodyLog = string(reqBodyBytes) // Log raw if not valid JSON
					}
				} else {
					// Log other non-JSON textual bodies as raw string
					requestBodyLog = string(reqBodyBytes)
				}
			}
		}

		// Crucial: Restore the request body for the next handler.
		// This must be done regardless of whether the body was logged or skipped for logging.
		// The `reqBodyBuffer` always contains the full original request body.
		r.Body = io.NopCloser(bytes.NewBuffer(reqBodyBuffer.Bytes()))

		lrw := newLoggingTraffic(w)
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
			logFields["request_body"] = requestBodyLog
		}

		newLogger.WithFields(logFields).Info("Processed request")
	})
}
