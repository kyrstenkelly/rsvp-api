package models

// Event type
type Event struct {
	ID          int64    `json:"id" db:"id" sql:",notnull"`
	Name        string   `json:"name" db:"name" sql:",notnull"`
	AddressID   int64    `json:"address_id" db:"address_id" sql:",notnull"`
	Address     *Address `json:"address"`
	FoodOptions []string `json:"food_options" db:"food_options"`
}
