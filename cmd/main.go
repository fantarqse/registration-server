package main

import (
	"log"

	"github.com/fantarqse/registrationserver/api"
	"github.com/fantarqse/registrationserver/db"
	_ "github.com/lib/pq"
)

func main() {
	myDB, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	app := api.NewStarter(myDB)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
