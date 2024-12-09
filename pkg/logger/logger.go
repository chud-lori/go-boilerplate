package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)



var Log *logrus.Logger

func Init(ctx context.Context) {
    Log = logrus.New()
    Log.SetFormatter(&logrus.JSONFormatter{})
    Log.SetLevel(logrus.InfoLevel)
    requestID, _ := ctx.Value("RequestID").(string)
    Log = Log.WithField("RequestID", requestID).Logger
}

type loggingTraffic struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingTraffic(w http.ResponseWriter) *loggingTraffic {
	return &loggingTraffic{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (lrw *loggingTraffic) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogTrafficMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = uuid.New().String()
        }

        baseLogger := logrus.New()
        baseLogger.SetFormatter(&logrus.JSONFormatter{})

        logger := baseLogger.WithField("RequestID", requestID)

        ctx := context.WithValue(r.Context(), "logger", logger)

        r = r.WithContext(ctx)

		lrw := NewLoggingTraffic(w)

		// call the next handler
        next.ServeHTTP(lrw, r)

        //_, file, line, ok := runtime.Caller(1)
        //source := "unknown"
        //if ok {
        //    source = fmt.Sprintf("%s:%d", file, line)
        //}

		duration := time.Since(start)

        logger.WithFields(logrus.Fields{
            //"timestamp": time.Now().Format(time.RFC3339),
			//"requestID": requestID,
			"method":    r.Method,
			"path":      r.URL.Path,
			"duration":  duration.String(),
			"status":    lrw.statusCode,
        }).Info("Processed request")

	})
}

