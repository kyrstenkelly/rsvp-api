package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
	"strings"
)

// AddressesPostgresAccess postgres implementation of a AddressesDAO
type AddressesPostgresAccess struct {
}

// AddressesAccess interface for a Cohorts data access object
type AddressesAccess interface {
	GetAddresses(tx *pg.Tx) ([]models.Address, error)
	GetAddress(tx *pg.Tx, id int64) (*models.Address, error)
	FindOrCreateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error)
	UpdateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error)
	DeleteAddress(tx *pg.Tx, id int64) (*models.Address, error)
}

// NewAddressesDAO Create a new addresses dao
func NewAddressesDAO() AddressesAccess {
	return &AddressesPostgresAccess{}
}

// CheckForDuplicate checks for an existing address
func CheckForDuplicate(tx *pg.Tx, address *models.Address) (int64, error) {
	query :=
		`SELECT id FROM addresses
		WHERE line1 = $1 AND line2 = $2 AND city = $3
			AND state = $4 AND zip = $5`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	var addressID int64
	_, err = stmt.Query(pg.Scan(&addressID), &address.Line1, &address.Line2, &address.City, &address.State, &address.Zip)
	if err != nil {
		return 0, err
	}
	return addressID, nil
}

// GetAddresses gets all addresses
func (a *AddressesPostgresAccess) GetAddresses(tx *pg.Tx) ([]models.Address, error) {
	var addresses []models.Address
	err := tx.Model(&addresses).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return addresses, nil
}

// GetAddress gets an address by id
func (a *AddressesPostgresAccess) GetAddress(tx *pg.Tx, id int64) (*models.Address, error) {
	address := new(models.Address)
	err := tx.Model(address).Where("address.id = ?", id).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return address, nil
}

// FindOrCreateAddress creates an address
func (a *AddressesPostgresAccess) FindOrCreateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error) {
	existingAddressID, err := CheckForDuplicate(tx, address)
	if err != nil {
		return nil, err
	}
	if existingAddressID > 0 {
		return a.GetAddress(tx, existingAddressID)
	}
	query :=
		`INSERT INTO
			addresses ("line1", "line2", "city", "state", "zip")
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var addressID int64
	_, err = stmt.Query(pg.Scan(&addressID), &address.Line1, &address.Line2, &address.City, &address.State, &address.Zip)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	address.ID = addressID

	return address, nil
}

// UpdateAddress updates an address
func (a *AddressesPostgresAccess) UpdateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error) {
	var q []string
	if address.Line1 != "" {
		q = append(q, "line1 = ?line1")
	}
	if address.Line2 != "" {
		q = append(q, "line2 = ?line2")
	}
	if address.City != "" {
		q = append(q, "city = ?city")
	}
	if address.State != "" {
		q = append(q, "state = ?state")
	}
	if address.Zip != "" {
		q = append(q, "zip = ?zip")
	}

	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(address).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedAddress, _ := a.GetAddress(tx, address.ID)
	return updatedAddress, nil
}

// DeleteAddress deletes an address
func (a *AddressesPostgresAccess) DeleteAddress(tx *pg.Tx, id int64) (*models.Address, error) {
	address := &models.Address{
		ID: id,
	}
	err := tx.Delete(address)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
