package access

import (
	// "encoding/json"
	"github.com/go-pg/pg"
	// "github.com/gorilla/mux"
	// "github.com/kyrstenkelly/rsvp-api/db"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	// "github.com/kyrstenkelly/rsvp-api/utils"
	log "github.com/sirupsen/logrus"
	// "net/http"
	// "strconv"
)

// AddressesPostgresAccess postgres implementation of a CohortsDAO
type AddressesPostgresAccess struct {
}

// AddressesAccess interface for a Cohorts data access object
type AddressesAccess interface {
	GetAddresses(tx *pg.Tx) (*models.Address, error)
	GetAddress(tx *pg.Tx, id int64) (*models.Address, error)
	CreateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error)
	UpdateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error)
	DeleteAddress(tx *pg.Tx, id int64) (*models.Address, error)
}

// NewAddressesDAO Create a new addresses dao
func NewAddressesDAO() AddressesAccess {
	return &AddressesPostgresAccess{}
}

// GetAddresses gets all addresses
func (a *AddressesPostgresAccess) GetAddresses(tx *pg.Tx) (*models.Address, error) {
	return nil, nil
}

// GetAddress gets an address by id
func (a *AddressesPostgresAccess) GetAddress(tx *pg.Tx, id int64) (*models.Address, error) {
	return nil, nil
}

// CreateAddress creates an address
func (a *AddressesPostgresAccess) CreateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error) {
	query :=
		`INSERT INTO
			addresses ("line1", "line2", "city", "state", "zip")
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}

	var addressID int64
	_, err = stmt.Query(pg.Scan(&addressID), &address.Line1, &address.Line2, &address.City, &address.State, &address.Zip)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	address.ID = addressID

	return address, nil
}

// UpdateAddress updates an address
func (a *AddressesPostgresAccess) UpdateAddress(tx *pg.Tx, address *models.Address) (*models.Address, error) {
	return nil, nil
}

// DeleteAddress deletes an address
func (a *AddressesPostgresAccess) DeleteAddress(tx *pg.Tx, id int64) (*models.Address, error) {
	return nil, nil
}
