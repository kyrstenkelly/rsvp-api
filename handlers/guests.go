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

// GuestsHandler type
type GuestsHandler struct {
	dao access.GuestsAccess
}

// NewGuestsHandler creates a new handler with the given dao
func NewGuestsHandler(dao access.GuestsAccess) *GuestsHandler {
	return &GuestsHandler{dao: dao}
}

// GetGuestsHandler gets a list of all guests
func (handler *GuestsHandler) GetGuestsHandler(r *http.Request) ([]byte, int, error) {
	log.Info("Getting all guests")
	// conn := db.GetDBConn()

	// var guests []models.Guest
	// err := dao.Model(&guests).Select()
	// if err != nil {
	// 	panic(err)
	// }
	// json.NewEncoder(w).Encode(guests)
	return nil, 0, nil
}

// GetGuestHandler gets an guest by id
func (handler *GuestsHandler) GetGuestHandler(r *http.Request) ([]byte, int, error) {
	conn := db.GetDBConn()

	params := mux.Vars(r)
	id, _ := strconv.ParseInt(params["id"], 10, 64)
	log.WithFields(log.Fields{
		"id": id,
	}).Info("Getting guest by ID")

	var guest *models.Guest
	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		guest, err = handler.dao.GetGuest(tx, id)
		return err
	})
	if err != nil {
		log.Error("Error getting guest")
		return nil, http.StatusInternalServerError, err
	}
	return utils.SerializeResponse(guest, http.StatusOK)
}

// CreateGuestHandler handles creating an guest
func (handler *GuestsHandler) CreateGuestHandler(r *http.Request) ([]byte, int, error) {
	conn := db.GetDBConn()

	var guest *models.Guest
	json.NewDecoder(r.Body).Decode(&guest)

	err := conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		_, err = handler.dao.CreateGuest(tx, guest)
		return err
	})
	if err != nil {
		log.Error("Error creating guest")
		return nil, http.StatusInternalServerError, err
	}

	return utils.SerializeResponse(guest, http.StatusOK)
}

// UpdateGuestHandler updates an existing guest
func (handler *GuestsHandler) UpdateGuestHandler(r *http.Request) ([]byte, int, error) {
	// TODO
	return nil, 0, nil
}

// DeleteGuestHandler deletes an guest
func (handler *GuestsHandler) DeleteGuestHandler(r *http.Request) ([]byte, int, error) {
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
