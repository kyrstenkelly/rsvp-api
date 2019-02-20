package access

import (
	"github.com/go-pg/pg"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
	"strings"
)

// EventsPostgresAccess postgres implementation of a CohortsDAO
type EventsPostgresAccess struct {
	addressAccess AddressesAccess
}

// EventsAccess interface for a Cohorts data access object
type EventsAccess interface {
	GetEvents(tx *pg.Tx) ([]models.Event, error)
	GetEvent(tx *pg.Tx, id int64) (*models.Event, error)
	CreateEvent(tx *pg.Tx, event *models.Event) (*models.Event, error)
	UpdateEvent(tx *pg.Tx, event *models.Event) (*models.Event, error)
	DeleteEvent(tx *pg.Tx, id int64) (*models.Event, error)
}

// NewEventsDAO Create a new events dao
func NewEventsDAO() EventsAccess {
	addressesDAO := NewAddressesDAO()
	return &EventsPostgresAccess{
		addressAccess: addressesDAO,
	}
}

// GetEvents gets all events
func (a *EventsPostgresAccess) GetEvents(tx *pg.Tx) ([]models.Event, error) {
	var events []models.Event
	err := tx.Model(&events).Select()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return events, nil
}

// GetEvent gets an event by id
func (a *EventsPostgresAccess) GetEvent(tx *pg.Tx, id int64) (*models.Event, error) {
	event := new(models.Event)
	err := tx.Model(event).Where("event.id = ?", id).Select()
	if err == pg.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Error(err)
		return nil, err
	}
	return event, nil
}

// CreateEvent creates an event
func (a *EventsPostgresAccess) CreateEvent(tx *pg.Tx, event *models.Event) (*models.Event, error) {
	address, err := a.addressAccess.FindOrCreateAddress(tx, event.Address)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	event.AddressID = address.ID

	query :=
		`INSERT INTO
			events ("name", "location", "address_id", "food_options", "date")
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING id`
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var eventID int64
	_, err = stmt.Query(pg.Scan(&eventID), &event.Name, &event.Location, &event.AddressID, &event.FoodOptions, &event.Date)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	event.ID = eventID

	return event, nil
}

// UpdateEvent updates an event
func (a *EventsPostgresAccess) UpdateEvent(tx *pg.Tx, event *models.Event) (*models.Event, error) {
	var q []string
	if event.Name != "" {
		q = append(q, "name = ?name")
	}
	if event.Location != "" {
		q = append(q, "location = ?location")
	}
	if event.Address != nil {
		_, err := a.addressAccess.FindOrCreateAddress(tx, event.Address)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		q = append(q, "address = ?address")
	}
	if event.FoodOptions != nil {
		q = append(q, "food_options = ?food_options")
	}

	qString := strings.Join(q, ", ")
	_, updateErr := tx.Model(event).Set(qString).Where("id = ?id").Update()
	if updateErr != nil {
		log.Error(updateErr)
		return nil, updateErr
	}

	updatedEvent, _ := a.GetEvent(tx, event.ID)
	return updatedEvent, nil
}

// DeleteEvent deletes an event
func (a *EventsPostgresAccess) DeleteEvent(tx *pg.Tx, id int64) (*models.Event, error) {
	event := &models.Event{
		ID: id,
	}
	err := tx.Delete(event)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return nil, nil
}
