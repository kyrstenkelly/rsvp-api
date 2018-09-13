package utils

import (
	"encoding/json"
	"github.com/ansel1/merry"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/kyrstenkelly/rsvp-api/db"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// ErrorResponse type
type ErrorResponse struct {
	Message string
}

// WrapHandler Extract attributes of errors and write them to ResponseWriter
func WrapHandler(handler func(request *http.Request, vars map[string]string) ([]byte, int, error)) http.HandlerFunc {
	f := func(writer http.ResponseWriter, request *http.Request) {
		buf, statusCode, err := handler(request, mux.Vars(request))

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

// RunWithTransaction runs a database access call within a transaction
func RunWithTransaction(call func(*pg.Tx) (interface{}, error)) (interface{}, error) {
	conn := db.GetDBConn()
	tx, err := conn.Begin()
	if err != nil {
		return nil, SQLError
	}

	result, err := call(tx)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return nil, rollbackErr
		}
		log.Debug("Transaction rollback successful")
		return nil, err
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return nil, commitErr
	}

	return result, nil
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

// GetIDFromVars parses the ID from vars to int64
func GetIDFromVars(vars map[string]string) int64 {
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		log.Error("Could not parse id as int")
	}
	return id
}
