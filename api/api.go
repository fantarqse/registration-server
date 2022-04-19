package api

import (
	"database/sql"
	_ "github.com/gorilla/mux"
)

type Api interface {
	Start()
}

type api struct {
	DB *sql.DB
	//TODO Route manager
}

func (a *api) Start() {}
