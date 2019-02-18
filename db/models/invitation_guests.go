package models

// InvitationGuest type
type InvitationGuest struct {
	InvitationID int64       `json:"invitation_id" db:"invitation_id"`
	Invitation   *Invitation `json:"invitation"`
	GuestID      int64       `json:"guest_id" db:"guest_id"`
	Guest        *Guest      `json:"guest"`
}
