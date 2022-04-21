package main

import (
	"database/sql"
	"github.com/fantarqse/registrationserver/api"
	"log"
)

func main() {
	var db *sql.DB
	app := api.NewStarter(db)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
