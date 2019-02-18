package access

import (
	"strings"

	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
)

// RSVPsPostgresAccess postgres implementation of a CohortsDAO
type RSVPsPostgresAccess struct {
	guestsAccess GuestsAccess
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
	guestsDAO := NewGuestsDAO()
	return &RSVPsPostgresAccess{
		guestsAccess: guestsDAO,
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
	query :=
		`INSERT INTO
			rsvps ("invitation_id", "guest_id", "attending", "food_choice")
		VALUES
			($1, $2, $3, $4)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var rsvpID int64
	_, err = stmt.Query(pg.Scan(&rsvpID), &rsvp.Attending, &rsvp.FoodChoice, &rsvp.GuestID, &rsvp.InvitationID)
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
	if rsvp.Attending {
		q = append(q, "attending = ?attending")
	}
	if rsvp.FoodChoice != "" {
		q = append(q, "food_choice = ?food_choice")
	}

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
	rsvp := &models.RSVP{
		ID: id,
	}
	err := tx.Delete(rsvp)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
