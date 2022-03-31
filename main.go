package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

type Users struct {
	Login    string `json:"login,omitempty", db: "login"`
	Password string `json:"password,omitempty", db: "password"`
	Email    string `json:"email,omitempty", db: "email"`
}

func initDB() {
	var err error

	db, err = sql.Open("postgres", "user=postgres password=qwerty dbname=registration_db sslmode=disable")
	if err != nil {
		panic(err)
	}
	log.Println("DB connected!")
}

func signUpHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Sign Up")

	users := &Users{}
	err := json.NewDecoder(req.Body).Decode(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err = db.Query("insert into users values ($1, $2, $3)", users.Login, users.Password, users.Email); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(http.StatusInternalServerError)
		return
	}
}

func signInHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Sign In")
}

func main() {
	http.HandleFunc("/signup", signUpHandler)
	http.HandleFunc("/signin", signInHandler)

	initDB()

	log.Println("Server is running...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
