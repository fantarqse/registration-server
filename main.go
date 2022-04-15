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
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	emailReg      string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	passwordReg   string = "[A-Z]+[a-z]+[^a-zA-Z]{2,}"
	authorization string = "Authorization"
)

var (
	db     *sql.DB
	jwtKey = []byte("secret_key_example")
)

type Claims struct {
	Login string `json:"login" db:"login"`
	jwt.StandardClaims
}

type SignUpUser struct {
	Login      *string `json:"login,omitempty" db:"login"`
	Password   *string `json:"password,omitempty" db:"password"`
	Email      *string `json:"email,omitempty" db:"email"`
	IsVerified bool    `json:"isVerified,omitempty" db:"is_verified"`
}

func (user SignUpUser) validator() error {
	email := regexp.MustCompile(emailReg)
	password := regexp.MustCompile(passwordReg)

	if len(*user.Login) < 3 {
		log.Println("Login is not valid!")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !password.MatchString(*user.Password) {
		log.Println("Password is not valid!")
		return errors.New("error: password is not valid")
	}

	if !email.MatchString(*user.Email) {
		log.Println("Email is not valid!")
		return errors.New("error: email is not valid")
	}
	return nil
}

type SignInUser struct {
	Login    *string `json:"login,omitempty" db:"login"`
	Password *string `json:"password,omitempty" db:"password"`
}

func (user SignInUser) validator() error {
	password := regexp.MustCompile(passwordReg)

	if len(*user.Login) < 3 {
		log.Println("Login is not valid!")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !password.MatchString(*user.Password) {
		log.Println("Password is not valid!")
		return errors.New("error: password is not valid")
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

func JWTVerifier(tokenString string) (jwt.MapClaims, bool, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	valid := token.Valid
	claims, ok := token.Claims.(jwt.MapClaims)

	return claims, ok, valid
}

func JWTGenerator(data string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		Login: data,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, expirationTime, err
}

func SignUpHandler(w http.ResponseWriter, req *http.Request) {
	users := &SignUpUser{}
	err := json.NewDecoder(req.Body).Decode(users)
	if err != nil {
		http.Error(w, "Bad request: JSON is not valid", 400)
		return
	}

	err = users.validator()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*users.Password), 8)
	if err != nil {
		log.Println("hashed password error")
	}

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

	tokenString, _, err := JWTGenerator(*users.Login)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Add(authorization, tokenString)

	fmt.Fprintf(w, "Sign Up")
}

func signInHandler(w http.ResponseWriter, req *http.Request) {
	user := &SignInUser{}
	storedUser := &SignInUser{}

	err := json.NewDecoder(req.Body).Decode(user)
	if err != nil {
		http.Error(w, "Bad request: JSON is not valid", 500)
		return
	}

	err = user.validator()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	result := db.QueryRow("select password from users where login=$1", user.Login)

	err = result.Scan(&storedUser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), 401)
			return
		}
		http.Error(w, err.Error(), 500)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*storedUser.Password), []byte(*user.Password))
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}

	tokenString, _, err := JWTGenerator(*user.Login)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Add(authorization, tokenString)

	fmt.Fprintf(w, "Sign In")
}

func verifyHandler(w http.ResponseWriter, req *http.Request) {
	token := req.Header.Get(authorization)
	claims, ok, valid := JWTVerifier(token)
	if ok && valid {
		res, err := db.Exec("update users set is_verified = true where login = $1", claims["login"])
		if err != nil {
			log.Println(err.Error())
		}

		count, _ := res.RowsAffected()
		if count <= 0 {
			log.Println("not updated")
		} else {
			log.Println("verified")
		}
	} else {
		log.Println("Token is not valid")
	}

	fmt.Fprintf(w, "Verified")
}

func main() {
	http.HandleFunc("/signup", SignUpHandler)
	http.HandleFunc("/signin", signInHandler)
	http.HandleFunc("/verify", verifyHandler)

	initDB()

	log.Println("Server is running...")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
