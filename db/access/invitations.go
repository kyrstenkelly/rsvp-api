package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
	"strings"
)

// InvitationsPostgresAccess postgres implementation of a CohortsDAO
type InvitationsPostgresAccess struct {
	addressAccess AddressesAccess
	guestAccess   GuestsAccess
}

// InvitationsAccess interface for a Cohorts data access object
type InvitationsAccess interface {
	GetInvitations(tx *pg.Tx) ([]models.Invitation, error)
	GetInvitation(tx *pg.Tx, id int64) (*models.Invitation, error)
	CreateInvitation(tx *pg.Tx, invitation *models.Invitation) (*models.Invitation, error)
	UpdateInvitation(tx *pg.Tx, invitation *models.Invitation) (*models.Invitation, error)
	DeleteInvitation(tx *pg.Tx, id int64) (*models.Invitation, error)
}

// NewInvitationsDAO Create a new invitations dao
func NewInvitationsDAO() InvitationsAccess {
	addressesDAO := NewAddressesDAO()
	guestsDAO := NewGuestsDAO()
	return &InvitationsPostgresAccess{
		addressAccess: addressesDAO,
		guestAccess:   guestsDAO,
	}
}

// GetInvitations gets all invitations
func (a *InvitationsPostgresAccess) GetInvitations(tx *pg.Tx) ([]models.Invitation, error) {
	var invitations []models.Invitation
	err := tx.Model(&invitations).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return invitations, nil
}

// GetInvitation gets a invitation by id
func (a *InvitationsPostgresAccess) GetInvitation(tx *pg.Tx, id int64) (*models.Invitation, error) {
	invitation := new(models.Invitation)
	err := tx.Model(invitation).Where("invitation.id = ?", id).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return invitation, nil
}

// CreateInvitation creates an invitation
func (a *InvitationsPostgresAccess) CreateInvitation(tx *pg.Tx, invitation *models.Invitation) (*models.Invitation, error) {
	// Create and append address to invitation
	address, err := a.addressAccess.FindOrCreateAddress(tx, invitation.Address)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	invitation.AddressID = address.ID

	invitation.GuestIds, err = a.BuildGuestIDs(tx, invitation.Guests)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	query :=
		`INSERT INTO
			invitations ("name", "email", "plus_one", "guest_ids", "address_id")
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var invitationID int64
	_, err = stmt.Query(pg.Scan(&invitationID), &invitation.Name, &invitation.Email, &invitation.PlusOne, pg.Array(&invitation.GuestIds), &invitation.AddressID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	invitation.ID = invitationID

	return invitation, nil
}

// UpdateInvitation updates an invitation
func (a *InvitationsPostgresAccess) UpdateInvitation(tx *pg.Tx, invitation *models.Invitation) (*models.Invitation, error) {
	var q []string
	if invitation.Name != "" {
		q = append(q, "name = ?name")
	}
	if invitation.Email != "" {
		q = append(q, "email = ?email")
	}
	if invitation.PlusOne {
		q = append(q, "plus_one = ?plus_one")
	}
	if invitation.Address != nil {
		_, err := a.addressAccess.FindOrCreateAddress(tx, invitation.Address)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		q = append(q, "address = ?address")
	}
	if invitation.Guests != nil {
		guestIds, err := a.BuildGuestIDs(tx, invitation.Guests)
		invitation.GuestIds = guestIds
		if err != nil {
			log.Error(err)
			return nil, err
		}
		q = append(q, "guest_ids = ?guest_ids")
	}

	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(invitation).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedInvitation, _ := a.GetInvitation(tx, invitation.ID)
	return updatedInvitation, nil
}

// DeleteInvitation deletes an invitation and the associated guests
func (a *InvitationsPostgresAccess) DeleteInvitation(tx *pg.Tx, id int64) (*models.Invitation, error) {
	invitation, err := a.GetInvitation(tx, id)
	if invitation == nil {
		return nil, nil
	}
	for _, guestID := range invitation.GuestIds {
		a.guestAccess.DeleteGuest(tx, guestID)
	}
	err = tx.Delete(invitation)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}

// BuildGuestIDs takes a list of guests and returns a list of their IDs
// After either finding or creating them
func (a *InvitationsPostgresAccess) BuildGuestIDs(tx *pg.Tx, guests *[]models.Guest) ([]int64, error) {
	var guestIDs []int64
	for _, guest := range *guests {
		newGuest, err := a.guestAccess.FindOrCreateGuest(tx, &guest)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		guestIDs = append(guestIDs, newGuest.ID)
	}
	return guestIDs, nil
}
