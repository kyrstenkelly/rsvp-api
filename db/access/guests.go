package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
	"strings"
)

// GuestsPostgresAccess postgres implementation of a CohortsDAO
type GuestsPostgresAccess struct {
	addressAccess AddressesAccess
}

// GuestsAccess interface for a Cohorts data access object
type GuestsAccess interface {
	GetGuests(tx *pg.Tx) ([]models.Guest, error)
	GetGuest(tx *pg.Tx, id int64) (*models.Guest, error)
	CreateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error)
	UpdateGuest(tx *pg.Tx, guest *models.Guest) (*models.Guest, error)
	DeleteGuest(tx *pg.Tx, id int64) (*models.Guest, error)
}

// NewGuestsDAO Create a new guests dao
func NewGuestsDAO() GuestsAccess {
	addressesDAO := NewAddressesDAO()
	return &GuestsPostgresAccess{
		addressAccess: addressesDAO,
	}
}

func _findGuestAddress(guest models.Guest, addresses []models.Address) *models.Address {
	log.WithFields(log.Fields{
		"address_ID": guest.AddressID,
		"guest":      guest,
		"addresses":  addresses,
	}).Debug("Searching for guest address")
	for _, address := range addresses {
		if address.ID == guest.AddressID {
			return &address
		}
	}
	return nil
}

// GetGuests gets all guests
func (a *GuestsPostgresAccess) GetGuests(tx *pg.Tx) ([]models.Guest, error) {
	var guests []models.Guest
	addresses, err := a.addressAccess.GetAddresses(tx)
	err = tx.Model(&guests).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, guest := range guests {
		guest.Address = _findGuestAddress(guest, addresses)
		guests[i] = guest
	}
	return guests, nil
}

// GetGuest gets an guest by id
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
	// TODO: Get address ID
	query :=
		`INSERT INTO
			guests ("first_name", "last_name", "email")
		VALUES
			($1, $2, $3)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var guestID int64
	_, err = stmt.Query(pg.Scan(&guestID), &guest.FirstName, &guest.LastName, &guest.Email)
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
	if guest.FirstName != "" {
		q = append(q, "first_name = ?first_name")
	}
	if guest.LastName != "" {
		q = append(q, "last_name = ?last_name")
	}
	if guest.Email != "" {
		q = append(q, "email = ?email")
	}
	if guest.AddressID > 0 {
		q = append(q, "address_id = ?address_id")
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
