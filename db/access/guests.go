package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
	"strings"
)

// GuestsPostgresAccess postgres implementation of a CohortsDAO
type GuestsPostgresAccess struct {
}

// GuestsAccess interface for a Cohorts data access object
type GuestsAccess interface {
	// TODO: Get guests by invitation ID
	GetGuests(tx *pg.Tx) ([]models.Guest, error)
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
func (a *GuestsPostgresAccess) GetGuests(tx *pg.Tx) ([]models.Guest, error) {
	var guests []models.Guest
	err := tx.Model(&guests).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return guests, nil
}

// GetGuest gets a guest by id
func (a *GuestsPostgresAccess) GetGuest(tx *pg.Tx, id int64) (*models.Guest, error) {
	guest := new(models.Guest)
	err := tx.Model(guest).Where("guest.id = ?", id).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return guest, nil
}

// CreateGuest creates an guest
func (a *GuestsPostgresAccess) CreateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error) {
	// TODO: Generate dynamic RSVP code
	query :=
		`INSERT INTO
			guests ("name")
		VALUES
			($1)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var guestID int64
	_, err = stmt.Query(pg.Scan(&guestID), &guest.Name)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	guest.ID = guestID

	return guest, nil
}

// UpdateGuest updates an guest
func (a *GuestsPostgresAccess) UpdateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error) {
	var q []string
	if guest.Name != "" {
		q = append(q, "name = ?name")
	}

	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(guest).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedGuest, _ := a.GetGuest(tx, guest.ID)
	return updatedGuest, nil
}

// DeleteGuest deletes an guest
func (a *GuestsPostgresAccess) DeleteGuest(tx *pg.Tx, id int64) (*models.Guest, error) {
	guest := &models.Guest{
		ID: id,
	}
	err := tx.Delete(guest)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
