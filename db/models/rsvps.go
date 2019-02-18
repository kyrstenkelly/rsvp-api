package models

// RSVP Type
type RSVP struct {
	ID           int64       `json:"id" db:"id" sql:",notnull"`
	InvitationID int64       `json:"invitation_id" db:"invitation_id" sql:",notnull"`
	Invitation   *Invitation `json:"invitation"`
	GuestID      int64       `json:"guest_id" db:"guest_id" sql:",notnull"`
	Guest        *Guest      `json:"guest"`
	Attending    bool        `json:"attending" db:"attending" sq:",notnull"`
	FoodChoice   string      `json:"food_choice" db:"food_choice"`
}
