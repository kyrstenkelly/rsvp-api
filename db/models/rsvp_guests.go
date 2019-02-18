package models

// RSVPGuest Type
type RSVPGuest struct {
	ID         int64  `json:"id" db:"id" sql:",notnull"`
	Guest      *Guest `json:"guest"`
	Attending  bool   `json:"attending" db:"attending" sql:",notnull"`
	PlusOne    bool   `json:"plus_one" db:"plus_one" sql:",notnull"`
	FoodChoice string `json:"food_choice" db:"food_choice"`
}
