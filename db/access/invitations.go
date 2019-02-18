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
	return &InvitationsPostgresAccess{
		addressAccess: addressesDAO,
	}
}

func _findInvitationAddress(invitation models.Invitation, addresses []models.Address) *models.Address {
	log.WithFields(log.Fields{
		"address_ID": invitation.AddressID,
		"invitation": invitation,
		"addresses":  addresses,
	}).Debug("Searching for invitation address")
	for _, address := range addresses {
		if address.ID == invitation.AddressID {
			return &address
		}
	}
	return nil
}

// GetInvitations gets all invitations
func (a *InvitationsPostgresAccess) GetInvitations(tx *pg.Tx) ([]models.Invitation, error) {
	var invitations []models.Invitation
	addresses, err := a.addressAccess.GetAddresses(tx)
	err = tx.Model(&invitations).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, invitation := range invitations {
		invitation.Address = _findInvitationAddress(invitation, addresses)
		invitations[i] = invitation
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
	address, err := a.addressAccess.CreateAddress(tx, invitation.Address)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// TODO: create guests upon invitation creation

	query :=
		`INSERT INTO
			invitations ("name", "email", "plus_one", "address_id")
		VALUES
			($1, $2, $3, $4)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var invitationID int64
	_, err = stmt.Query(pg.Scan(&invitationID), &invitation.Name, &invitation.Email, &invitation.PlusOne, &address.ID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	invitation.ID = invitationID
	invitation.AddressID = address.ID

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
	if invitation.AddressID > 0 {
		q = append(q, "address_id = ?address_id")
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

// DeleteInvitation deletes an invitation
func (a *InvitationsPostgresAccess) DeleteInvitation(tx *pg.Tx, id int64) (*models.Invitation, error) {
	invitation := &models.Invitation{
		ID: id,
	}
	err := tx.Delete(invitation)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
