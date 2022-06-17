package api

import (
	"bytes"
	"encoding/json"
	"github.com/fantarqse/registrationserver/internal/db"
	"github.com/fantarqse/registrationserver/internal/token"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
)

func TestVerification(t *testing.T) {
	testCase := []struct {
		Name  string `json:"-"`
		Login string `json:"login"`
		Want  int    `json:"-"`
	}{
		{
			Name:  "Test1",
			Login: "SomeValidLogin01",
			Want:  200,
		},
		{
			Name:  "Test2",
			Login: "SomeValidLogin15",
			Want:  200,
		},
		{
			Name:  "Test3",
			Login: "SomeValidLogin16",
			Want:  200,
		},
		{
			Name:  "Test4",
			Login: "SomeValidLogin17",
			Want:  200,
		},
		{
			Name:  "Test5",
			Login: "SomeValidLogin18",
			Want:  200,
		},
		{
			Name:  "Test6",
			Login: "SomeValidLogin19",
			Want:  200,
		},
		{
			Name:  "Test7",
			Login: "SomeValidLogin20",
			Want:  200,
		},
		{
			Name:  "Test8",
			Login: "SomeValidLogin21",
			Want:  200,
		},
		{
			Name:  "Test9",
			Login: "SomeValidLogin22",
			Want:  200,
		},
		{
			Name:  "Test10",
			Login: "SomeValidLogin23",
			Want:  200,
		},
		{
			Name:  "Test11",
			Login: "SomeValidLogin24",
			Want:  200,
		},
		{
			Name:  "Test12",
			Login: "SomeValidLogin25",
			Want:  200,
		},
		{
			Name:  "Test13",
			Login: "SomeValidLogin26",
			Want:  200,
		},
		{
			Name:  "Test14",
			Login: "SomeValidLogin27",
			Want:  200,
		},
		{
			Name:  "Test15",
			Login: "SomeValidLogin28",
			Want:  200,
		},
		{
			Name:  "Test16",
			Login: "SomeValidLogin29",
			Want:  200,
		},
		{
			Name:  "Test17",
			Login: "NotRegisteredUser",
			Want:  409,
		},
		{
			Name:  "Test18",
			Login: "",
			Want:  409,
		},
		{
			Name: "Test19",
			Want: 409,
		},
		{
			Name:  "Test20",
			Login: "TokenIsNotValid",
			Want:  401,
		},
	}

	myDB, _ := db.New()
	router := mux.NewRouter()
	a := &api{
		DB:     myDB,
		router: router,
	}
	handler := http.HandlerFunc(a.verifyHandler)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			data, err := json.Marshal(tc)
			if err != nil {
				log.Printf("json error: %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/verify", bytes.NewReader(data))
			if err != nil {
				log.Printf("requst error: %v", err.Error())
			}

			tokenString, _, err := token.Generate(tc.Login)
			if err != nil {
				log.Printf("token generation error: %v", err.Error())
			}

			if tc.Name == "Test20" {
				tokenString = tokenString + "@"
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Add(Authorization, tokenString)
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tc.Want, recorder.Code)
		})
	}
}
