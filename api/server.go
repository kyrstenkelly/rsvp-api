package api

import (
	"github.com/gorilla/mux"
	"github.com/kyrstenkelly/rsvp-api/db/access"
	"github.com/kyrstenkelly/rsvp-api/handlers"
	"github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func buildHandler(handlerMethod func(request *http.Request, vars map[string]string) ([]byte, int, error), authRequired bool) http.Handler {
	handlerFunc := utils.WrapHandler(handlerMethod)
	if authRequired {
		return authMiddleware(handlerFunc)
	}
	return handlerFunc
}

// Serve sets up handlers and serves at the given port
func Serve(port int64) {
	router := mux.NewRouter()

	router.Handle("/", http.FileServer(http.Dir("./views/")))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	addressDAO := access.NewAddressesDAO()
	addressHandler := handlers.NewAddressesHandler(addressDAO)

	guestsDAO := access.NewGuestsDAO()
	guestHandler := handlers.NewGuestsHandler(guestsDAO)

	router.Handle("/addresses", buildHandler(addressHandler.GetAddressesHandler, true)).Methods("GET")
	router.Handle("/addresses", buildHandler(addressHandler.CreateAddressHandler, true)).Methods("POST")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.GetAddressHandler, true)).Methods("GET")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.UpdateAddressHandler, true)).Methods("PUT")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.DeleteAddressHandler, true)).Methods("DELETE")

	router.Handle("/guests", buildHandler(guestHandler.GetGuestsHandler, true)).Methods("GET")
	router.Handle("/guests", buildHandler(guestHandler.CreateGuestHandler, true)).Methods("POST")
	router.Handle("/guests/{id}", buildHandler(guestHandler.GetGuestHandler, true)).Methods("GET")
	router.Handle("/guests/{id}", buildHandler(guestHandler.UpdateGuestHandler, true)).Methods("PUT")
	router.Handle("/guests/{id}", buildHandler(guestHandler.DeleteGuestHandler, true)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
