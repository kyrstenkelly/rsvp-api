package handlers

import (
	"encoding/json"
	"github.com/go-pg/pg"
	"github.com/gorilla/mux"
	"github.com/kyrstenkelly/rsvp-api/db"
	"github.com/kyrstenkelly/rsvp-api/db/access"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	"github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
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
func (handler *AddressesHandler) GetAddressesHandler(r *http.Request) ([]byte, int, error) {
	log.Info("Getting all addresses")
	// conn := db.GetDBConn()

	// var addresses []models.Address
	// err := dao.Model(&addresses).Select()
	// if err != nil {
	// 	panic(err)
	// }
	// json.NewEncoder(w).Encode(addresses)
	return nil, 0, nil
}

// GetAddressHandler gets an address by id
func (handler *AddressesHandler) GetAddressHandler(r *http.Request) ([]byte, int, error) {
	conn := db.GetDBConn()
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)

	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting address by ID")

	var address *models.Address
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		address, err = handler.dao.GetAddress(tx, id)
		return err
	})
	if err != nil {
		log.Error("Error getting address")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(address, http.StatusOK)
	// address := new(models.Address)
	// err := conn.Model(address).Where("address.id = ?", id).Select()
	// if err != nil {
	// 	panic(err)
	// }
	// json.NewEncoder(w).Encode(address)
}

// CreateAddressHandler handles creating an address
func (handler *AddressesHandler) CreateAddressHandler(r *http.Request) ([]byte, int, error) {
	var address *models.Address
	json.NewDecoder(r.Body).Decode(&address)

	log.WithFields(log.Fields{
		"address": address,
	}).Info("Creating address")

	conn := db.GetDBConn()

	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		address, err = handler.dao.CreateAddress(tx, address)
		return err
	})
	if err != nil {
		log.Error("Error creating address")
		return nil, http.StatusInternalServerError, err
	}

	return utils.SerializeResponse(address, http.StatusOK)
}

// UpdateAddressHandler updates an existing address
func (handler *AddressesHandler) UpdateAddressHandler(r *http.Request) ([]byte, int, error) {
	// TODO
	return nil, 0, nil
}

// DeleteAddressHandler deletes an address
func (handler *AddressesHandler) DeleteAddressHandler(r *http.Request) ([]byte, int, error) {
	// params := mux.Vars(r)
	// for index, item := range guests {
	// 	if item.ID == params["id"] {
	// 		guests = append(guests[:index], guests[index+1:]...)
	// 		break
	// 	}
	// }
	// json.NewEncoder(w).Encode(guests)
	return nil, 0, nil
}
