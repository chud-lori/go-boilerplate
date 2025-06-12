package helper

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/chud-lori/go-boilerplate/pkg/errors"
	"github.com/chud-lori/go-boilerplate/pkg/logger"
	"github.com/sirupsen/logrus"
)

func GetPayload(request *http.Request, result interface{}) error {
	logger, ok := request.Context().Value(logger.LoggerContextKey).(*logrus.Entry)
	if !ok {
		// Fallback: If logger not found, use a basic logger or panic,
		// depending on how critical logging is at this point.
		// For simplicity, let's use a standard logger if not found,
		// though in a production app, you might want to ensure it's always there.
		// Or, if GetPayload is always called after LogTrafficMiddleware, you might panic.
		// For now, let's just log a warning.
		logrus.Warn("GetPayload: Request logger not found in context. Using base logrus.")
		logger = logrus.NewEntry(logrus.StandardLogger()) // Use a default entry if not found
	}

	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)

	if err != nil {
		logger.WithError(err).Error("Failed to decode request payload")
		return appErrors.NewBadRequestError("User not found", err)
	}
	return nil
}

func WriteResponse(writer http.ResponseWriter, response interface{}, httpCode int64) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(int(httpCode))
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)

	if err != nil {
		panic(err)
	}
}
