package handlers

import (
	"encoding/json"
	"github.com/go-pg/pg/v9"
	"github.com/kyrstenkelly/rsvp-api/db"
	"github.com/kyrstenkelly/rsvp-api/db/access"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	"github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// EventsHandler type
type EventsHandler struct {
	dao access.EventsAccess
}

// NewEventsHandler creates a new handler with the given dao
func NewEventsHandler(dao access.EventsAccess) *EventsHandler {
	return &EventsHandler{dao: dao}
}

// GetEventsHandler gets a list of all events
func (handler *EventsHandler) GetEventsHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	log.Info("Getting all events")
	conn := db.GetDBConn()

	var events []models.Event
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		events, err = handler.dao.GetEvents(tx)
		return err
	})
	if err != nil {
		log.Error("Error getting events")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(events, http.StatusOK)
}

// GetEventHandler gets an event by id
func (handler *EventsHandler) GetEventHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting event by ID")

	event, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.GetEvent(tx, id)
	})
	if err != nil {
		log.Error("Error getting event")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(event, http.StatusOK)
}

// CreateEventHandler handles creating an event
func (handler *EventsHandler) CreateEventHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	var event *models.Event
	json.NewDecoder(r.Body).Decode(&event)

	log.WithFields(log.Fields{
		"event": event,
	}).Info("Creating event")

	createdEvent, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.CreateEvent(tx, event)
	})
	if err != nil {
		log.Error("Error creating event")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(createdEvent, http.StatusOK)
}

// UpdateEventHandler updates an existing event
func (handler *EventsHandler) UpdateEventHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)
	var event *models.Event
	json.NewDecoder(r.Body).Decode(&event)
	event.ID = id

	log.WithFields(log.Fields{
		"event": event,
	}).Info("Updating event")

	updatedEvent, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.UpdateEvent(tx, event)
	})
	if err != nil {
		log.Error("Error updating event")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(updatedEvent, http.StatusOK)
}

// DeleteEventHandler deletes an event
func (handler *EventsHandler) DeleteEventHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting event")

	_, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.DeleteEvent(tx, id)
	})
	if err != nil {
		log.Error("Error deleting event")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(nil, http.StatusOK)
}
