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
	err := tx.Model(&invitations).
		Column("invitation.*", "Address", "Event").
		Select()

	var invitationsWithGuests []models.Invitation
	for _, invitation := range invitations {
		guests, err := a.guestAccess.GetGuests(tx, invitation.GuestIds)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		invitation.Guests = &guests
		invitationsWithGuests = append(invitationsWithGuests, invitation)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}
	return invitationsWithGuests, nil
}

// GetInvitation gets a invitation by id
func (a *InvitationsPostgresAccess) GetInvitation(tx *pg.Tx, id int64) (*models.Invitation, error) {
	invitation := new(models.Invitation)
	err := tx.Model(invitation).
		Column("invitation.*", "Address").
		Where("invitation.id = ?", id).
		Select()

	var guests []models.Guest
	for _, guestID := range invitation.GuestIds {
		guest, _ := a.guestAccess.GetGuest(tx, guestID)
		guests = append(guests, *guest)
	}
	invitation.Guests = &guests

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

	// Find and append guest IDs to the invitation
	invitation.GuestIds, err = a.BuildGuestIDs(tx, invitation.Guests)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	query :=
		`INSERT INTO
			invitations ("event_id", "name", "email", "plus_one", "guest_ids", "address_id")
		VALUES
			($1, $2, $3, $4, $5, $6)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var invitationID int64
	_, err = stmt.Query(pg.Scan(&invitationID), &invitation.EventID, &invitation.Name, &invitation.Email, &invitation.PlusOne, pg.Array(&invitation.GuestIds), &invitation.AddressID)
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
