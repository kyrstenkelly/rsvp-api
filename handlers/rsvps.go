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

// RSVPsHandler type
type RSVPsHandler struct {
	dao access.RSVPsAccess
}

// NewRSVPsHandler creates a new handler with the given dao
func NewRSVPsHandler(dao access.RSVPsAccess) *RSVPsHandler {
	return &RSVPsHandler{dao: dao}
}

// GetRSVPsHandler gets a list of all rsvps
func (handler *RSVPsHandler) GetRSVPsHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	log.Info("Getting all rsvps")
	conn := db.GetDBConn()

	var rsvps []models.RSVP
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		rsvps, err = handler.dao.GetRSVPs(tx)
		return err
	})
	if err != nil {
		log.Error("Error getting rsvps")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(rsvps, http.StatusOK)
}

// GetRSVPHandler gets an rsvp by id
func (handler *RSVPsHandler) GetRSVPHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting rsvp by ID")

	rsvp, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.GetRSVP(tx, id)
	})
	if err != nil {
		log.Error("Error getting rsvp")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(rsvp, http.StatusOK)
}

// CreateRSVPHandler handles creating an rsvp
func (handler *RSVPsHandler) CreateRSVPHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	var rsvp *models.RSVP
	json.NewDecoder(r.Body).Decode(&rsvp)

	log.WithFields(log.Fields{
		"invitation_id": rsvp.InvitationID,
	}).Info("Creating rsvp")

	createdRSVP, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.CreateRSVP(tx, rsvp)
	})
	if err != nil {
		log.Error("Error creating rsvp")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(createdRSVP, http.StatusOK)
}

// UpdateRSVPHandler updates an existing rsvp
func (handler *RSVPsHandler) UpdateRSVPHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)
	var rsvp *models.RSVP
	json.NewDecoder(r.Body).Decode(&rsvp)
	rsvp.ID = id

	log.WithFields(log.Fields{
		"invitation_id": rsvp.InvitationID,
		"rsvp_guests":   rsvp.RSVPGuests,
	}).Info("Updating rsvp")

	updatedRSVP, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.UpdateRSVP(tx, rsvp)
	})
	if err != nil {
		log.Error("Error updating rsvp")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(updatedRSVP, http.StatusOK)
}

// DeleteRSVPHandler deletes an rsvp
func (handler *RSVPsHandler) DeleteRSVPHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting rsvp")

	_, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.DeleteRSVP(tx, id)
	})
	if err != nil {
		log.Error("Error deleting rsvp")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(nil, http.StatusOK)
}
