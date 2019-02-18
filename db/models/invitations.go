package models

// Invitation type
type Invitation struct {
	ID      int64    `json:"id" db:"id" sql:",notnull"`
	Name    string   `json:"name" db:"name" sql:",notnull"`
	Email   string   `json:"email" db:"email" sql:",notnull,unique"`
	PlusOne bool     `json:"plus_one" db:"plus_one"`
	Guests  *[]Guest `json:"guests"`
	Address *Address `json:"address" db:"address"`
}
