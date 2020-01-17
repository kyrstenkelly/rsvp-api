package access

import (
	"errors"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
)

// RSVPGuestsPostgresAccess postgres implementation of a RSVPGuestDAO
type RSVPGuestsPostgresAccess struct {
	guestAccess GuestsAccess
}

// RSVPGuestsAccess interface for a Cohorts data access object
type RSVPGuestsAccess interface {
	GetRSVPGuests(tx *pg.Tx, invitationID int64) ([]models.RSVPGuest, error)
	GetRSVPGuest(tx *pg.Tx, id int64) (*models.RSVPGuest, error)
	CreateRSVPGuest(tx *pg.Tx, rsvpID int64, rsvpGuest *models.RSVPGuest) (*models.RSVPGuest, error)
	UpdateRSVPGuest(tx *pg.Tx, rsvpGuest *models.RSVPGuest) (*models.RSVPGuest, error)
	DeleteRSVPGuest(tx *pg.Tx, id int64) (*models.RSVPGuest, error)
}

// NewRSVPGuestsDAO Create a new rsvpguests dao
func NewRSVPGuestsDAO() RSVPGuestsAccess {
	guestsDAO := NewGuestsDAO()
	return &RSVPGuestsPostgresAccess{
		guestAccess: guestsDAO,
	}
}

// GetRSVPGuests gets all rsvps
func (a *RSVPGuestsPostgresAccess) GetRSVPGuests(tx *pg.Tx, rsvpID int64) ([]models.RSVPGuest, error) {
	var rsvpGuests []models.RSVPGuest
	// TODO: Use invitation ID to get guests
	err := tx.Model(&rsvpGuests).
		Column("rsvp_guest.*", "Guest").
		Where("rsvp_guest.rsvp_id = ?", rsvpID).
		Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return rsvpGuests, nil
}

// GetRSVPGuest gets an rsvpGuest by id
func (a *RSVPGuestsPostgresAccess) GetRSVPGuest(tx *pg.Tx, id int64) (*models.RSVPGuest, error) {
	rsvpGuest := new(models.RSVPGuest)
	err := tx.Model(rsvpGuest).Where("rsvp_guest.id = ?", id).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return rsvpGuest, nil
}

// CreateRSVPGuest creates an rsvpGuest
func (a *RSVPGuestsPostgresAccess) CreateRSVPGuest(tx *pg.Tx, rsvpID int64, rsvpGuest *models.RSVPGuest) (*models.RSVPGuest, error) {
	var guest *models.Guest
	var err error
	// If it's a plus one, create a new guest. Otherwise verify that the guest is in the DB.
	if rsvpGuest.IsPlusOne {
		guest, err = a.guestAccess.FindOrCreateGuest(tx, rsvpGuest.Guest)
	} else {
		guest, err = a.guestAccess.GetGuestByName(tx, rsvpGuest.Guest.Name)
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if guest == nil {
		err := errors.New("Cannot create RSVP for a guest does not exist")
		return nil, err
	}
	rsvpGuest.GuestID = guest.ID

	query :=
		`INSERT INTO
			rsvp_guests ("rsvp_id", "guest_id", "attending", "food_choice", "is_plus_one")
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var rsvpGuestID int64
	_, err = stmt.Query(pg.Scan(&rsvpGuestID), &rsvpID, &rsvpGuest.GuestID, &rsvpGuest.Attending, &rsvpGuest.FoodChoice, &rsvpGuest.IsPlusOne)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rsvpGuest.ID = rsvpGuestID

	return rsvpGuest, nil
}

// UpdateRSVPGuest updates an rsvpGuest
func (a *RSVPGuestsPostgresAccess) UpdateRSVPGuest(tx *pg.Tx, rsvpGuest *models.RSVPGuest) (*models.RSVPGuest, error) {
	log.Debug("top of UpdateRSVPGuest")
	var q []string
	q = append(q, "attending = ?attending")
	q = append(q, "is_plus_one = ?is_plus_one")
	if rsvpGuest.FoodChoice != "" {
		q = append(q, "food_choice = ?food_choice")
	}

	log.Debug("updating rsvp guest: ")
	log.Debug(rsvpGuest)
	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(rsvpGuest).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedRSVPGuest, err := a.GetRSVPGuest(tx, rsvpGuest.ID)
	updatedRSVPGuest.Guest, err = a.guestAccess.GetGuest(tx, updatedRSVPGuest.GuestID)
	if err != nil {
		log.Error(updateErr)
		return nil, updateErr
	}
	return updatedRSVPGuest, nil
}

// DeleteRSVPGuest deletes an rsvpGuest
func (a *RSVPGuestsPostgresAccess) DeleteRSVPGuest(tx *pg.Tx, id int64) (*models.RSVPGuest, error) {
	rsvpGuest := &models.RSVPGuest{
		ID: id,
	}
	err := tx.Delete(rsvpGuest)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
