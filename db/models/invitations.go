package models

// Invitation type
type Invitation struct {
	ID        int64    `json:"id" db:"id" sql:",notnull"`
	Name      string   `json:"name" db:"name" sql:",notnull"`
	Email     string   `json:"email" db:"email"`
	PlusOne   bool     `json:"plus_one" db:"plus_one"`
	AddressID int64    `json:"address_id" db:"address_id"`
	Address   *Address `json:"address"`
}
