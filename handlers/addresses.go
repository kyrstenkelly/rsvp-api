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

// AddressesHandler type
type AddressesHandler struct {
	dao access.AddressesAccess
}

// NewAddressesHandler creates a new handler with the given dao
func NewAddressesHandler(dao access.AddressesAccess) *AddressesHandler {
	return &AddressesHandler{dao: dao}
}

// GetAddressesHandler gets a list of all addresses
func (handler *AddressesHandler) GetAddressesHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	log.Info("Getting all addresses")
	conn := db.GetDBConn()

	var addresses []models.Address
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		addresses, err = handler.dao.GetAddresses(tx)
		return err
	})
	if err != nil {
		log.Error("Error getting addresses")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(addresses, http.StatusOK)
}

// GetAddressHandler gets an address by id
func (handler *AddressesHandler) GetAddressHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting address by ID")

	address, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.GetAddress(tx, id)
	})
	if err != nil {
		log.Error("Error getting address")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(address, http.StatusOK)
}

// FindOrCreateAddressHandler handles creating an address
func (handler *AddressesHandler) FindOrCreateAddressHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	var address *models.Address
	json.NewDecoder(r.Body).Decode(&address)

	log.WithFields(log.Fields{
		"address": address,
	}).Info("Creating address")

	createdAddress, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.FindOrCreateAddress(tx, address)
	})
	if err != nil {
		log.Error("Error creating address")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(createdAddress, http.StatusOK)
}

// UpdateAddressHandler updates an existing address
func (handler *AddressesHandler) UpdateAddressHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)
	var address *models.Address
	json.NewDecoder(r.Body).Decode(&address)
	address.ID = id

	log.WithFields(log.Fields{
		"address": address,
	}).Info("Updating address")

	updatedAddress, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.UpdateAddress(tx, address)
	})
	if err != nil {
		log.Error("Error updating address")
		return nil, http.StatusBadRequest, err
	}

	return utils.SerializeResponse(updatedAddress, http.StatusOK)
}

// DeleteAddressHandler deletes an address
func (handler *AddressesHandler) DeleteAddressHandler(r *http.Request, vars map[string]string) ([]byte, int, error) {
	id := utils.GetIDFromVars(vars)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Deleting address")

	_, err := utils.RunWithTransaction(func(tx *pg.Tx) (interface{}, error) {
		return handler.dao.DeleteAddress(tx, id)
	})
	if err != nil {
		log.Error("Error deleting address")
		return nil, http.StatusBadRequest, err
	}
	return utils.SerializeResponse(nil, http.StatusOK)
}
