package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	emailReg    string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	passwordReg string = "[A-Z]+[a-z]+[^a-zA-Z]{2,}"
)

var db *sql.DB

type Users struct {
	Login      *string `json:"login,omitempty" db:"login"`
	Password   *string `json:"password,omitempty" db:"password"`
	Email      *string `json:"email,omitempty" db:"email"`
	IsVerified bool    `json:"isVerified,omitempty" db:"is_verified"`
}

func (u Users) validator() error {
	email := regexp.MustCompile(emailReg)
	password := regexp.MustCompile(passwordReg)

	if len(*u.Login) < 3 {
		log.Println("Login is not valid!")
		return errors.New("error: login is too short")
	}

	if len(*u.Password) < 8 || !password.MatchString(*u.Password) {
		log.Println("Password is not valid!")
		return errors.New("error: password is not valid")
	}

	if !email.MatchString(*u.Email) {
		log.Println("Email is not valid!")
		return errors.New("error: email is not valid")
	}
	return nil
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
	users := &Users{}
	err := json.NewDecoder(req.Body).Decode(users)
	if err != nil {
		http.Error(w, "Bad request: JSON is not valid", 400)
		return
	}

	validationErr := users.validator()
	if validationErr != nil {
		http.Error(w, validationErr.Error(), 400)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(*users.Password), 8)

	if _, err = db.Query(
		"insert into users (login, password, email) values ($1, $2, $3)",
		users.Login,
		string(hashedPassword),
		users.Email,
	); err != nil {
		log.Println(err.Error())

		if strings.Contains(err.Error(), "повторювані значення ключа") {
			http.Error(w, "Username or Email is not unique.", 500)
			return

		} else if strings.Contains(err.Error(), "порушує not-null") {
			http.Error(w, "Username, Password and Email are required.", 500)
			return

		} else {
			http.Error(w, err.Error(), 500)
			return
		}
	}

	fmt.Fprintf(w, "Sign Up")
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
