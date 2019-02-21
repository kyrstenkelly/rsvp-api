package models

// RSVP Type
type RSVP struct {
	ID           int64       `json:"id" db:"id" sql:",notnull"`
	InvitationID int64       `json:"invitation_id" db:"invitation_id" sql:",notnull"`
	RSVPGuestIds []int64     `json:"-" db:"rsvp_guest_ids"`
	RSVPGuests   []RSVPGuest `json:"rsvp_guests" db:"rsvp_guests"`
}
