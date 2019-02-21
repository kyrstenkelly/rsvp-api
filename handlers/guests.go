package handlers

import (
	"encoding/json"
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db"
	"github.com/kyrstenkelly/rsvp-api/db/access"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	"github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// GuestsHandler type
type GuestsHandler struct {
	dao access.GuestsAccess
}

// NewGuestsHandler creates a new handler with the given dao
func NewGuestsHandler(dao access.GuestsAccess) *GuestsHandler {
	return &GuestsHandler{dao: dao}
}

// GetGuestsHandler gets a list of all guests
func (handler *GuestsHandler) GetGuestsHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	log.Info("Getting all guests")
	conn := db.GetDBConn()

	var guests []models.Guest
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		guests, err = handler.dao.GetGuests(tx, nil)
		return err
	})
	if err != nil {
		log.Error("Error getting guests")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(guests, http.StatusOK)
}

// GetGuestHandler gets an guest by id
func (handler *GuestsHandler) GetGuestHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting guest by ID")

	guest, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.GetGuest(tx, id)
	})
	if err != nil {
		log.Error("Error getting guest")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(guest, http.StatusOK)
}

// FindOrCreateGuestHandler handles creating an guest
func (handler *GuestsHandler) FindOrCreateGuestHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	var guest *models.Guest
	json.NewDecoder(r.Body).Decode(&guest)

	log.WithFields(log.Fields{
		"guest": guest,
	}).Info("Creating guest")

	createdGuest, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.FindOrCreateGuest(tx, guest)
	})
	if err != nil {
		log.Error("Error creating guest")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(createdGuest, http.StatusOK)
}

// UpdateGuestHandler updates an existing guest
func (handler *GuestsHandler) UpdateGuestHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)
	var guest *models.Guest
	json.NewDecoder(r.Body).Decode(&guest)
	guest.ID = id

	log.WithFields(log.Fields{
		"guest": guest,
	}).Info("Updating guest")

	updatedGuest, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.UpdateGuest(tx, guest)
	})
	if err != nil {
		log.Error("Error updating guest")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(updatedGuest, http.StatusOK)
}

// DeleteGuestHandler deletes an guest
func (handler *GuestsHandler) DeleteGuestHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting guest")

	_, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.DeleteGuest(tx, id)
	})
	if err != nil {
		log.Error("Error deleting guest")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(nil, http.StatusOK)
}
