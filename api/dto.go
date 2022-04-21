package api

import (
	"errors"
	"log"
	"regexp"
)

const (
	emailReg    string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	passwordReg string = "[A-Z]+[a-z]+[^a-zA-Z]{2,}"
)

type Validator interface {
	Validate() error
}

type RegistrationUser struct {
	Login      *string `json:"login,omitempty" db:"login"`
	Password   *string `json:"password,omitempty" db:"password"`
	Email      *string `json:"email,omitempty" db:"password"`
	IsVerified bool    `json:"isVerified,omitempty" db:"is_verified"`
}

func (user *RegistrationUser) Validate() error {
	email := regexp.MustCompile(emailReg)
	password := regexp.MustCompile(passwordReg)

	if len(*user.Login) < 3 {
		log.Println("error: login is not valid")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !password.MatchString(*user.Password) {
		log.Println("error: password is not valid")
		return errors.New("error: password is not valid")
	}

	if !email.MatchString(*user.Email) {
		log.Println("error: email is not valid")
		return errors.New("error: email is not valid")
	}
	return nil
}

type AuthenticationUser struct {
	Login    *string `json:"login,omitempty" db:"login"`
	Password *string `json:"password,omitempty" db:"password"`
}

func (user *AuthenticationUser) Validate() error {
	password := regexp.MustCompile(passwordReg)

	if len(*user.Login) < 3 {
		log.Println("error: login is not valid!")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !password.MatchString(*user.Password) {
		log.Println("error: password is not valid!")
		return errors.New("error: password is not valid")
	}
	return nil
}
