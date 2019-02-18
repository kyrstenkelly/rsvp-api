package access

import (
	"strings"

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
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return rsvps, nil
}

// GetRSVP gets an rsvp by id
func (a *RSVPsPostgresAccess) GetRSVP(tx *pg.Tx, id int64) (*models.RSVP, error) {
	rsvp := new(models.RSVP)
	err := tx.Model(rsvp).Where("rsvp.id = ?", id).Select()
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
	// Create and append RSVPGuests to the RSVP
	var rsvpGuests []*models.RSVPGuest
	for _, rsvpGuest := range rsvp.RSVPGuests {
		newRSVPGuest, err := a.rsvpGuestAccess.CreateRSVPGuest(tx, rsvpGuest)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		rsvpGuests = append(rsvpGuests, newRSVPGuest)
	}
	rsvp.RSVPGuests = rsvpGuests

	query :=
		`INSERT INTO
			rsvps ("invitation_id", "rsvp_guests")
		VALUES
			($1, $2)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var rsvpID int64
	_, err = stmt.Query(pg.Scan(&rsvpID), &rsvp.InvitationID, &rsvp.RSVPGuests)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rsvp.ID = rsvpID

	return rsvp, nil
}

// UpdateRSVP updates an rsvp
func (a *RSVPsPostgresAccess) UpdateRSVP(tx *pg.Tx, rsvp *models.RSVP) (*models.RSVP, error) {
	var q []string
	// TODO: update each rsvp guest

	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(rsvp).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedRSVP, _ := a.GetRSVP(tx, rsvp.ID)
	return updatedRSVP, nil
}

// DeleteRSVP deletes an rsvp
func (a *RSVPsPostgresAccess) DeleteRSVP(tx *pg.Tx, id int64) (*models.RSVP, error) {
	// First delete the RSVP guests, then the RSVP
	rsvp, err := a.GetRSVP(tx, id)
	for _, rsvpGuest := range rsvp.RSVPGuests {
		a.rsvpGuestAccess.DeleteRSVPGuest(tx, rsvpGuest.ID)
	}
	err = tx.Delete(rsvp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
