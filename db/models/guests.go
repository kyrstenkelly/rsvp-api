package models

// Guest type
type Guest struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
