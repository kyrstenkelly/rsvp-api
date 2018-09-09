package db

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // do something
	"github.com/kelseyhightower/envconfig"
)

// Config describes the database config object
type Config struct {
	Host   string `envconfig:"HOST" required:"true"`
	Port   string `envconfig:"PORT" required:"true"`
	User   string `envconfig:"USER" required:"true"`
	Pass   string `envconfig:"PASS" required:"true"`
	DBName string `envconfig:"NAME" required:"true"`
}

// Initialize initializes the database and runs migrations
func Initialize() {
	var config Config
	err := envconfig.Process("db", &config)
	if err != nil {
		panic(err)
	}
	connectionString := "postgres://" + config.Host + ":" + config.Port + "/" + config.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	m.Steps(2)
}
