package api

import (
	"github.com/gorilla/mux"
	"github.com/kyrstenkelly/rsvp-api/db/access"
	"github.com/kyrstenkelly/rsvp-api/handlers"
	"github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const swaggerPath = "/api/v1/swagger"

// Serve sets up handlers and serves at the given port
func Serve(port int64) {
	router := mux.NewRouter()

	addressDAO := access.NewAddressesDAO()
	addressHandler := handlers.NewAddressesHandler(addressDAO)

	guestsDAO := access.NewGuestsDAO()
	guestHandler := handlers.NewGuestsHandler(guestsDAO)

	router.HandleFunc("/addresses", utils.WrapHandler(addressHandler.GetAddressesHandler)).Methods("GET")
	router.HandleFunc("/addresses", utils.WrapHandler(addressHandler.CreateAddressHandler)).Methods("POST")
	router.HandleFunc("/addresses/{id}", utils.WrapHandler(addressHandler.GetAddressHandler)).Methods("GET")
	router.HandleFunc("/addresses/{id}", utils.WrapHandler(addressHandler.UpdateAddressHandler)).Methods("PUT")
	router.HandleFunc("/addresses/{id}", utils.WrapHandler(addressHandler.DeleteAddressHandler)).Methods("DELETE")

	router.HandleFunc("/guests", utils.WrapHandler(guestHandler.GetGuestsHandler)).Methods("GET")
	router.HandleFunc("/guests", utils.WrapHandler(guestHandler.CreateGuestHandler)).Methods("POST")
	router.HandleFunc("/guests/{id}", utils.WrapHandler(guestHandler.GetGuestHandler)).Methods("GET")
	router.HandleFunc("/guests/{id}", utils.WrapHandler(guestHandler.UpdateGuestHandler)).Methods("PUT")
	router.HandleFunc("/guests/{id}", utils.WrapHandler(guestHandler.DeleteGuestHandler)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
