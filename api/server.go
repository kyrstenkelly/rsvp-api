package api

import (
	muxHandlers "github.com/gorilla/handlers"
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
	router.Handle("/addresses", buildHandler(addressHandler.GetAddressesHandler, true)).Methods("GET")
	router.Handle("/addresses", buildHandler(addressHandler.CreateAddressHandler, true)).Methods("POST")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.GetAddressHandler, true)).Methods("GET")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.UpdateAddressHandler, true)).Methods("PUT")
	router.Handle("/addresses/{id}", buildHandler(addressHandler.DeleteAddressHandler, true)).Methods("DELETE")

	eventsDAO := access.NewEventsDAO()
	eventsHandler := handlers.NewEventsHandler(eventsDAO)
	router.Handle("/events", buildHandler(eventsHandler.GetEventsHandler, true)).Methods("GET")
	router.Handle("/events", buildHandler(eventsHandler.CreateEventHandler, false)).Methods("POST")
	router.Handle("/events/{id}", buildHandler(eventsHandler.GetEventHandler, false)).Methods("GET")
	router.Handle("/events/{id}", buildHandler(eventsHandler.UpdateEventHandler, false)).Methods("PUT")
	router.Handle("/events/{id}", buildHandler(eventsHandler.DeleteEventHandler, true)).Methods("DELETE")

	invitationsDAO := access.NewInvitationsDAO()
	invitationsHandler := handlers.NewInvitationsHandler(invitationsDAO)
	router.Handle("/invitations", buildHandler(invitationsHandler.GetInvitationsHandler, true)).Methods("GET")
	router.Handle("/invitations", buildHandler(invitationsHandler.CreateInvitationHandler, false)).Methods("POST")
	router.Handle("/invitations/{id}", buildHandler(invitationsHandler.GetInvitationHandler, false)).Methods("GET")
	router.Handle("/invitations/{id}", buildHandler(invitationsHandler.UpdateInvitationHandler, false)).Methods("PUT")
	router.Handle("/invitations/{id}", buildHandler(invitationsHandler.DeleteInvitationHandler, true)).Methods("DELETE")

	rsvpsDAO := access.NewRSVPsDAO()
	rsvpsHandler := handlers.NewRSVPsHandler(rsvpsDAO)
	router.Handle("/rsvps", buildHandler(rsvpsHandler.GetRSVPsHandler, true)).Methods("GET")
	router.Handle("/rsvps", buildHandler(rsvpsHandler.CreateRSVPHandler, false)).Methods("POST")
	router.Handle("/rsvps/{id}", buildHandler(rsvpsHandler.GetRSVPHandler, false)).Methods("GET")
	router.Handle("/rsvps/{id}", buildHandler(rsvpsHandler.UpdateRSVPHandler, false)).Methods("PUT")
	router.Handle("/rsvps/{id}", buildHandler(rsvpsHandler.DeleteRSVPHandler, true)).Methods("DELETE")

	headersOk := muxHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := muxHandlers.AllowedOrigins([]string{"*"})
	methodsOk := muxHandlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Fatal(http.ListenAndServe(":8000", muxHandlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
