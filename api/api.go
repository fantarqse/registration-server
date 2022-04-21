package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/fantarqse/registrationserver/token"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const authorization string = "Authorization"

type API interface {
	Start() error
}

type api struct {
	DB     *sql.DB
	router *mux.Router
}

func NewStarter(db *sql.DB) API {
	router := mux.NewRouter()
	r := &api{
		DB:     db,
		router: router,
	}

	router.HandleFunc("registration", r.registrationHandler)
	router.HandleFunc("authentication", r.authenticationHandler)
	router.HandleFunc("verify", r.verifyHandler)
	return r
}

func (a *api) Start() error {
	log.Println("info: server was starting")
	return http.ListenAndServe(":8000", a.router)
}

func (a *api) registrationHandler(w http.ResponseWriter, r *http.Request) {
	user := &RegistrationUser{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Println("error: JSON is not valid")
		http.Error(w, "bad request: JSON is not valid", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 8)
	if err != nil {
		log.Println("error: hashed password")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if _, err = a.DB.Query(
		"insert into users (login, password, email) values ($1, $2, $3)",
		user.Login,
		string(hashedPassword),
		user.Email,
	); err != nil {
		log.Printf("error: %v", err.Error())

		if strings.Contains(err.Error(), "повторювані значення ключа") {
			http.Error(w, "Username or Email is not unique.", http.StatusInternalServerError)
			return

		} else if strings.Contains(err.Error(), "порушує not-null") {
			http.Error(w, "Username, Password and Email are required.", http.StatusInternalServerError)
			return

		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tokenString, _, err := token.JWTGeneration(*user.Login)
	if err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add(authorization, tokenString)

	fmt.Fprintf(w, "Sign Up")
}

func (a *api) authenticationHandler(w http.ResponseWriter, r *http.Request) {
	user := &AuthenticationUser{}
	storedUser := &AuthenticationUser{}

	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		log.Println("error: JSON is not valid")
		http.Error(w, "bad request: JSON is not valid", http.StatusBadRequest)
		return
	}

	if err := user.Validate(); err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := a.DB.QueryRow("select password from users where login=$1", user.Login)
	if err := result.Scan(&storedUser.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		log.Printf("error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*storedUser.Password), []byte(*user.Password)); err != nil {
		log.Printf("error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenString, _, err := token.JWTGeneration(*user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add(authorization, tokenString)

	fmt.Fprintf(w, "Sign In")
}

func (a *api) verifyHandler(w http.ResponseWriter, r *http.Request) {
	jwtToken := r.Header.Get(authorization)
	claims, ok, valid := token.JWTVerification(jwtToken)
	if ok && valid {
		res, err := a.DB.Exec("update users set is_verified = true where login = $1", claims["login"])
		if err != nil {
			log.Printf("error: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		count, _ := res.RowsAffected()
		if count <= 0 {
			log.Println("error: field is not updated")
			http.Error(w, "field is not updated", http.StatusConflict)
			return

		} else {
			log.Println("info: field is updated")
		}
	} else {
		log.Println("error: token is not valid")
		http.Error(w, "token is not valid", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Verified")
}
