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

// InvitationsHandler type
type InvitationsHandler struct {
	dao access.InvitationsAccess
}

// NewInvitationsHandler creates a new handler with the given dao
func NewInvitationsHandler(dao access.InvitationsAccess) *InvitationsHandler {
	return &InvitationsHandler{dao: dao}
}

// GetInvitationsHandler gets a list of all invitations
func (handler *InvitationsHandler) GetInvitationsHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	log.Info("Getting all invitations")
	conn := db.GetDBConn()

	var invitations []models.Invitation
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		invitations, err = handler.dao.GetInvitations(tx)
		return err
	})
	if err != nil {
		log.Error("Error getting invitations")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(invitations, http.StatusOK)
}

// GetInvitationHandler gets an invitation by id
func (handler *InvitationsHandler) GetInvitationHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting invitation by ID")

	invitation, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.GetInvitation(tx, id)
	})
	if err != nil {
		log.Error("Error getting invitation")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(invitation, http.StatusOK)
}

// CreateInvitationHandler handles creating an invitation
func (handler *InvitationsHandler) CreateInvitationHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	var invitation *models.Invitation
	json.NewDecoder(r.Body).Decode(&invitation)

	log.WithFields(log.Fields{
		"invitation": invitation,
	}).Info("Creating invitation")

	createdInvitation, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.CreateInvitation(tx, invitation)
	})
	if err != nil {
		log.Error("Error creating invitation")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(createdInvitation, http.StatusOK)
}

// UpdateInvitationHandler updates an existing invitation
func (handler *InvitationsHandler) UpdateInvitationHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)
	var invitation *models.Invitation
	json.NewDecoder(r.Body).Decode(&invitation)
	invitation.ID = id

	log.WithFields(log.Fields{
		"invitation": invitation,
	}).Info("Updating invitation")

	updatedInvitation, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.UpdateInvitation(tx, invitation)
	})
	if err != nil {
		log.Error("Error updating invitation")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(updatedInvitation, http.StatusOK)
}

// DeleteInvitationHandler deletes an invitation
func (handler *InvitationsHandler) DeleteInvitationHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting invitation")

	_, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.DeleteInvitation(tx, id)
	})
	if err != nil {
		log.Error("Error deleting invitation")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(nil, http.StatusOK)
}
