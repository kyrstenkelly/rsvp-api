package models

// Address type
type Address struct {
	ID    int64  `json:"id" db:"id" sql:",notnull"`
	Line1 string `json:"line_1" db:"line_1" sql:",notnull"`
	Line2 string `json:"line_2" db:"line_2"`
	City  string `json:"city" db:"city" sql:",notnull"`
	State string `json:"state" db:"state" sql:",notnull"`
	Zip   string `json:"zip" db:"zip" sql:",notnull"`
}
