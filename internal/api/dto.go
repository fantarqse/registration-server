package api

import (
	"errors"
	"log"
	"regexp"
	"unicode"
)

const emailReg string = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

type Validator interface {
	Validate() error
}

type RegistrationData struct {
	Login      *string `json:"login,omitempty" db:"login"`
	Password   *string `json:"password,omitempty" db:"password"`
	Email      *string `json:"email,omitempty" db:"password"`
	IsVerified bool    `json:"isVerified,omitempty" db:"is_verified"`
}

func (user *RegistrationData) Validate() error {
	email := regexp.MustCompile(emailReg)

	if len(*user.Login) < 3 {
		log.Println("error: login is not valid")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !validatePassword(*user.Password) {
		log.Println("error: password is not valid")
		return errors.New("error: password is not valid")
	}

	if !email.MatchString(*user.Email) {
		log.Println("error: email is not valid")
		return errors.New("error: email is not valid")
	}
	return nil
}

type AuthenticationData struct {
	Login    *string `json:"login,omitempty" db:"login"`
	Password *string `json:"password,omitempty" db:"password"`
}

func (user *AuthenticationData) Validate() error {
	if len(*user.Login) < 3 {
		log.Println("error: login is not valid!")
		return errors.New("error: login is too short")
	}

	if len(*user.Password) < 8 || !validatePassword(*user.Password) {
		log.Println("error: password is not valid!")
		return errors.New("error: password is not valid")
	}
	return nil
}

func validatePassword(password string) bool {
	var lower, upper bool
	symbol := 0
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			symbol++
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			symbol++
		case unicode.IsLower(c):
			lower = true
		}
	}
	if symbol > 1 && upper && lower {
		return true
	}
	return false
}
