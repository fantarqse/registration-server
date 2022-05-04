package main

import (
	"log"

	"github.com/fantarqse/registrationserver/internal/api"
	"github.com/fantarqse/registrationserver/internal/config"
	"github.com/fantarqse/registrationserver/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	c, err := config.NewConfig(".", ".env")
	if err != nil {
		log.Fatal(err)
	}

	myDB, err := db.New(c)
	if err != nil {
		log.Fatal(err)
	}

	app := api.New(myDB)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
