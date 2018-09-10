package migrations

import (
	"github.com/go-pg/migrations"
	"log"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Print("creating table addresses...")
		_, err := db.Exec(`CREATE TABLE CREATE TABLE addresses (
			id          integer PRIMARY KEY,
			line1       varchar(40) NOT NULL,
			line2       varchar(40),
			city        varchar(40) NOT NULL,
			state       varchar(40) NOT NULL,
			zip         varchar(40) NOT NULL
		)`)
		return err
	}, func(db migrations.DB) error {
		log.Print("dropping table addresses...")
		_, err := db.Exec(`DROP TABLE IF EXISTS addresses`)
		return err
	})

	migrations.Register(func(db migrations.DB) error {
		log.Print("creating table guests...")
		_, err := db.Exec(`CREATE TABLE guests (
			id          integer PRIMARY KEY,
			name        varchar(40) NOT NULL,
			email       varchar(40),
			address_id  integer REFERENCES addresses (id) NOT NULL
		)`)
		return err
	}, func(db migrations.DB) error {
		log.Print("dropping table guests...")
		_, err := db.Exec(`DROP TABLE IF EXISTS guests`)
		return err
	})
}
