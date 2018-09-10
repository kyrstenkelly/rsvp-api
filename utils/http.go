package utils

import (
	"encoding/json"
	"github.com/ansel1/merry"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// ErrorResponse type
type ErrorResponse struct {
	Message string
}

// WrapHandler Extract attributes of errors and write them to ResponseWriter
func WrapHandler(handler func(request *http.Request) ([]byte, int, error)) func(writer http.ResponseWriter, request *http.Request) {
	f := func(writer http.ResponseWriter, request *http.Request) {
		buf, statusCode, err := handler(request)

		if err != nil {
			log.Error(err)

			message := merry.Message(err)
			responseBody := &ErrorResponse{
				Message: message,
			}

			buf, _ = json.Marshal(&responseBody)
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(statusCode)
		_, err = writer.Write(buf)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to write response writer")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	}
	return f
}

// SerializeResponse will marshal a response object into json
func SerializeResponse(obj interface{}, status int) ([]byte, int, error) {
	if obj == nil {
		return nil, status, nil
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		return nil, http.StatusInternalServerError, JSONMarshalingError
	}
	return buf, status, nil
}
