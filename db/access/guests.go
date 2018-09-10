package access

import (
	// "encoding/json"
	"github.com/go-pg/pg"
	// "github.com/gorilla/mux"
	// "github.com/kyrstenkelly/rsvp-api/db"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	// "github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	// "net/http"
	// "strconv"
)

// GuestsPostgresAccess postgres implementation of a CohortsDAO
type GuestsPostgresAccess struct {
}

// GuestsAccess interface for a Cohorts data access object
type GuestsAccess interface {
	GetGuests(tx *pg.Tx) (*models.Guest, error)
	GetGuest(tx *pg.Tx, id int64) (*models.Guest, error)
	CreateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error)
	UpdateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error)
	DeleteGuest(tx *pg.Tx, id int64) (*models.Guest, error)
}

// NewGuestsDAO Create a new guests dao
func NewGuestsDAO() GuestsAccess {
	return &GuestsPostgresAccess{}
}

// GetGuests gets all guests
func (a *GuestsPostgresAccess) GetGuests(tx *pg.Tx) (*models.Guest, error) {
	return nil, nil
}

// GetGuest gets an guest by id
func (a *GuestsPostgresAccess) GetGuest(tx *pg.Tx, id int64) (*models.Guest, error) {
	return nil, nil
}

// CreateGuest creates an guest
func (a *GuestsPostgresAccess) CreateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error) {
	log.WithFields(log.Fields{
		"guest": guest,
	}).Info("Creating guest")
	query :=
		`INSERT INTO
			guests ("first_name", "last_name", "email", "address_id")
		VALUES
			($1, $2, $3, $4)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	var guestID int64
	_, err = stmt.Query(pg.Scan(&guestID), &guest.FirstName, &guest.LastName, &guest.Email, &guest.AddressID)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	guest.ID = guestID

	return guest, nil
}

// UpdateGuest updates an guest
func (a *GuestsPostgresAccess) UpdateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error) {
	return nil, nil
}

// DeleteGuest deletes an guest
func (a *GuestsPostgresAccess) DeleteGuest(tx *pg.Tx, id int64) (*models.Guest, error) {
	return nil, nil
}
