package models

// RSVP Types
type RSVP struct {
	ID        int64 `json:"id" db:"id" sql:",notnull"`
	HeadCount int   `json:"head_count" db:"head_count" sql:",notnull"`
	GuestID   int64 `json:"guest_id" db:"guest_id" sql:",notnull"`
	Guest     *Guest
}
