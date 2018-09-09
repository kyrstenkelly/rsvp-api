package main

import (
	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kyrstenkelly/rsvp-api/api"
	"github.com/kyrstenkelly/rsvp-api/db"
)

func main() {
	db.Initialize()
	api.Serve(8000)
}
