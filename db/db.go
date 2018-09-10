package db

import (
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/kelseyhightower/envconfig"
	"github.com/kyrstenkelly/rsvp-api/db/models"
	log "github.com/sirupsen/logrus"
)

var conn *pg.DB

// GetConfig loads the config object from env vars and returns it
func GetConfig() (*Config, error) {
	var config Config

	if err := envconfig.Process("db", &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func createSchema(db *pg.DB) error {
	for _, model := range models.Models {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// InitDb initializes the database
func InitDb() error {
	log.Info("InitDb: Started")

	config, err := GetConfig()
	if err != nil {
		log.Error("Unable to configure the database")
		return err
	}

	options, err := pg.ParseURL(config.ConnectionString())
	if err != nil {
		log.Error("Unable to parse connection string")
		return err
	}

	conn = pg.Connect(options)

	var oldVersion, newVersion int64
	err = createSchema(conn)
	err = conn.RunInTransaction(func(tx *pg.Tx) (err error) {
		oldVersion, newVersion, err = migrations.RunMigrations(tx, migrations.RegisteredMigrations())
		return
	})
	if err != nil {
		log.Error("Unable to run database migrations")
		return err
	}
	if newVersion != oldVersion {
		log.WithFields(log.Fields{
			"oldVersion": oldVersion,
			"newVersion": newVersion,
		}).Debug("Migrated to new version")
	} else {
		log.WithFields(log.Fields{
			"version": oldVersion,
		}).Debug("Database version")
	}

	log.Info("InitDb: Complete")

	return nil
}

// GetDBConn get the active connection
func GetDBConn() *pg.DB {
	if conn == nil {
		log.Fatal("Database connection not initialized. You must call `InitDB()` first")
	}
	return conn
}
