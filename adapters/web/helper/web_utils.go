package helper

import (
	"encoding/json"
	"net/http"
)

func GetPayload(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)

	if err != nil {
		panic(err)
	}
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
