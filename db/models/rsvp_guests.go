package models

// RSVPGuest Type
type RSVPGuest struct {
	ID         int64  `json:"id" db:"id" sql:",notnull"`
	RsvpID     int64  `json:"rsvp_id" db:"rsvp_id" sql:",notnull"`
	RSVP       *RSVP  `json:"rsvp"`
	GuestID    int64  `json:"guest_id" db:"guest_id" sql:",notnull"`
	Guest      *Guest `json:"guest"`
	Attending  bool   `json:"attending" db:"attending" sql:",notnull"`
	IsPlusOne  bool   `json:"is_plus_one" db:"is_plus_one" sql:"default:false"`
	FoodChoice string `json:"food_choice" db:"food_choice"`
}
