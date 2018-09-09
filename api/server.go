package api

import (
	"github.com/gorilla/mux"
	"github.com/kyrstenkelly/rsvp-api/handlers"
	"log"
	"net/http"
)

const swaggerPath = "/api/v1/swagger"

// Serve sets up handlers and serves at the given port
func Serve(port int64) {
	router := mux.NewRouter()

	router.HandleFunc("/guests", handlers.GetGuests).Methods("GET")
	router.HandleFunc("/guests/{id}", handlers.GetGuest).Methods("GET")
	router.HandleFunc("/guests/{id}", handlers.CreateGuest).Methods("POST")
	router.HandleFunc("/guests/{id}", handlers.DeleteGuest).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
