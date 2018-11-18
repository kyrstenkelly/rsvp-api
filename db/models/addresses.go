package models

// Address type
type Address struct {
	ID    int64  `json:"id" db:"id" sql:",notnull"`
	Line1 string `json:"line1" db:"line1" sql:",notnull"`
	Line2 string `json:"line2" db:"line2"`
	City  string `json:"city" db:"city" sql:",notnull"`
	State string `json:"state" db:"state" sql:",notnull"`
	Zip   string `json:"zip" db:"zip" sql:",notnull"`
}
