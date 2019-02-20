package models

// Invitation type
type Invitation struct {
	ID        int64    `json:"id" db:"id" sql:",notnull"`
	Name      string   `json:"name" db:"name" sql:",notnull"`
	Email     string   `json:"email" db:"email" sql:",notnull,unique"`
	PlusOne   bool     `json:"plus_one" db:"plus_one"`
	GuestIds  []int64  `json:"guest_ids" db:"guest_ids" sql:",notnull,array"`
	Guests    *[]Guest `json:"guests"`
	AddressID int64    `json:"address_id" db:"address_id"`
	Address   *Address `json:"address"`
}
