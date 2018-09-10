package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kyrstenkelly/rsvp-api/api"
	"github.com/kyrstenkelly/rsvp-api/db"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
	err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}
	api.Serve(8000)
}
