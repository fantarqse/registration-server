package db

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/fantarqse/registrationserver/internal/config"
	_ "github.com/lib/pq"
)

func New(c *config.Config) (*sql.DB, error) {
	log.Println("info: DB is connecting...")
	return sql.Open(
		c.DatabaseDriverName,
		c.DatabaseDataSourceName,
	)
}

func Register(db *sql.DB, login, password, email string) error {
	if _, err := db.Query(
		"insert into users (login, password, email) values ($1, $2, $3)",
		login,
		password,
		email,
	); err != nil {
		if strings.Contains(err.Error(), "повторювані значення ключа") {
			return errors.New("error: username or email is not unique")

		} else if strings.Contains(err.Error(), "порушує not-null") {
			return errors.New("error: username, password and email are required")

		} else {
			return err
		}
	}
	return nil
}

func Authenticate(db *sql.DB, login string, password **string) error {
	result := db.QueryRow("select password from users where login=$1", login)
	if err := result.Scan(password); err != nil {
		log.Printf("error: %v", err.Error())
		return err
	}
	return nil
}
