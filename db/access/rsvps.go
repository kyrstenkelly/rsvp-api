package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
)

// RSVPsPostgresAccess postgres implementation of a CohortsDAO
type RSVPsPostgresAccess struct {
	rsvpGuestAccess RSVPGuestsAccess
}

// RSVPsAccess interface for a Cohorts data access object
type RSVPsAccess interface {
	GetRSVPs(tx *pg.Tx) ([]models.RSVP, error)
	GetRSVP(tx *pg.Tx, id int64) (*models.RSVP, error)
	CreateRSVP(tx *pg.Tx, rsvp *models.RSVP) (*models.RSVP, error)
	UpdateRSVP(tx *pg.Tx, rsvp *models.RSVP) (*models.RSVP, error)
	DeleteRSVP(tx *pg.Tx, id int64) (*models.RSVP, error)
}

// NewRSVPsDAO Create a new rsvps dao
func NewRSVPsDAO() RSVPsAccess {
	rsvpGuestsDAO := NewRSVPGuestsDAO()
	return &RSVPsPostgresAccess{
		rsvpGuestAccess: rsvpGuestsDAO,
	}
}

// GetRSVPs gets all rsvps
func (a *RSVPsPostgresAccess) GetRSVPs(tx *pg.Tx) ([]models.RSVP, error) {
	var rsvps []models.RSVP
	err := tx.Model(&rsvps).Select()

	var rsvpsWithGuests []models.RSVP
	for _, rsvp := range rsvps {
		rsvpGuests, err := a.rsvpGuestAccess.GetRSVPGuests(tx, rsvp.ID)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		rsvp.RSVPGuests = rsvpGuests
		rsvpsWithGuests = append(rsvpsWithGuests, rsvp)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}
	return rsvpsWithGuests, nil
}

// GetRSVP gets an rsvp by id
func (a *RSVPsPostgresAccess) GetRSVP(tx *pg.Tx, id int64) (*models.RSVP, error) {
	rsvp := new(models.RSVP)

	err := tx.Model(rsvp).
		Where("rsvp.id = ?", id).
		Select()

	rsvpGuests, err := a.rsvpGuestAccess.GetRSVPGuests(tx, rsvp.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rsvp.RSVPGuests = rsvpGuests

	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return rsvp, nil
}

// CreateRSVP creates an rsvp
func (a *RSVPsPostgresAccess) CreateRSVP(tx *pg.Tx, rsvp *models.RSVP) (*models.RSVP, error) {
	query :=
		`INSERT INTO rsvps ("invitation_id") VALUES ($1)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var rsvpID int64
	_, err = stmt.Query(pg.Scan(&rsvpID), &rsvp.InvitationID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rsvp.ID = rsvpID

	// Create and append RSVPGuests to the RSVP
	var rsvpGuestIDs []int64
	for _, rsvpGuest := range rsvp.RSVPGuests {
		newRSVPGuest, err := a.rsvpGuestAccess.CreateRSVPGuest(tx, rsvpID, &rsvpGuest)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		rsvpGuestIDs = append(rsvpGuestIDs, newRSVPGuest.ID)
	}
	rsvp.RSVPGuestIds = rsvpGuestIDs
	_, updateErr := tx.Model(rsvp).Set("rsvp_guest_ids = ?rsvp_guest_ids").Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	return rsvp, nil
}

// UpdateRSVP updates an rsvp
func (a *RSVPsPostgresAccess) UpdateRSVP(tx *pg.Tx, rsvp *models.RSVP) (*models.RSVP, error) {
	var updatedRSVPGuests []models.RSVPGuest
	log.Debug("about to update rsvp guests: ")
	log.Debug(rsvp.RSVPGuests)
	for _, rsvpGuest := range rsvp.RSVPGuests {
		log.WithFields(log.Fields{
			"foodChoice": rsvpGuest.FoodChoice,
			"guest":      rsvpGuest.Guest,
		}).Info("about to update rsvp guest")
		updated, err := a.rsvpGuestAccess.UpdateRSVPGuest(tx, &rsvpGuest)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		updatedRSVPGuests = append(updatedRSVPGuests, *updated)
	}
	rsvp.RSVPGuests = updatedRSVPGuests
	return rsvp, nil
}

// DeleteRSVP deletes an rsvp
func (a *RSVPsPostgresAccess) DeleteRSVP(tx *pg.Tx, id int64) (*models.RSVP, error) {
	// First delete the RSVP guests, then the RSVP
	rsvp, err := a.GetRSVP(tx, id)
	for _, rsvpGuestID := range rsvp.RSVPGuestIds {
		a.rsvpGuestAccess.DeleteRSVPGuest(tx, rsvpGuestID)
	}
	err = tx.Delete(rsvp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
