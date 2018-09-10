package models

// Guest type
type Guest struct {
	ID        int64  `json:"id" db:"id" sql:",notnull"`
	FirstName string `json:"first_name" db:"first_name" sql:",notnull"`
	LastName  string `json:"last_name" db:"last_name" sql:",notnull"`
	Email     string `json:"email" db:"email"`
	AddressID int64  `json:"address_id" db:"address_id"`
	Address   *Address
}
