package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fantarqse/registrationserver/db"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
)

func TestAuthentication(t *testing.T) {
	testCase := []struct {
		Name     string `json:"-"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Want     int    `json:"-"`
	}{
		{
			Name:     "Test1",
			Login:    "Aa",
			Password: "SomeValidPassword00",
			Want:     400,
		},
		{
			Name:     "Test2",
			Login:    "SomeValidLogin",
			Password: "Aa",
			Want:     400,
		},
		{
			Name:  "Test3",
			Login: "SomeValidLogin",
			Want:  400,
		},
		{
			Name:     "Test4",
			Password: "SomeValidPassword00",
			Want:     400,
		},
		{
			Name:     "Test5",
			Login:    "SomeValidLogin01",
			Password: "a________Aa",
			Want:     200,
		},
		{
			Name:     "Test6",
			Login:    "SomeValidLogin15",
			Password: "A______a",
			Want:     200,
		},
		{
			Name:     "Test7",
			Login:    "SomeValidLogin16",
			Password: "a______A",
			Want:     200,
		},
		{
			Name:     "Test8",
			Login:    "SomeValidLogin17",
			Password: "Aa________",
			Want:     200,
		},
		{
			Name:     "Test9",
			Login:    "SomeValidLogin18",
			Password: "aA________",
			Want:     200,
		},
		{
			Name:     "Test10",
			Login:    "SomeValidLogin19",
			Password: "________Aa",
			Want:     200,
		},
		{
			Name:     "Test11",
			Login:    "SomeValidLogin20",
			Password: "________aA",
			Want:     200,
		},
		{
			Name:     "Test12",
			Login:    "SomeValidLogin21",
			Password: "____Aa____",
			Want:     200,
		},
		{
			Name:     "Test13",
			Login:    "SomeValidLogin22",
			Password: "____aA____",
			Want:     200,
		},
		{
			Name:     "Test14",
			Login:    "SomeValidLogin23",
			Password: "A___A___a",
			Want:     200,
		},
		{
			Name:     "Test15",
			Login:    "SomeValidLogin24",
			Password: "A___a___A",
			Want:     200,
		},
		{
			Name:     "Test16",
			Login:    "SomeValidLogin25",
			Password: "a___A___A",
			Want:     200,
		},
		{
			Name:     "Test17",
			Login:    "SomeValidLogin26",
			Password: "_Aa______",
			Want:     200,
		},
		{
			Name:     "Test18",
			Login:    "SomeValidLogin27",
			Password: "_aA______",
			Want:     200,
		},
		{
			Name:     "Test19",
			Login:    "SomeValidLogin28",
			Password: "__Aa_____",
			Want:     200,
		},
		{
			Name:     "Test20",
			Login:    "SomeValidLogin29",
			Password: "__aA_____",
			Want:     200,
		},
		{
			Name:     "Test21",
			Login:    "NotRegisteredUser",
			Password: "SomeValidPassword00",
			Want:     401,
		},
		{
			Name:     "Test22",
			Login:    "SomeValidLogin29",
			Password: "IncorrectPassword00",
			Want:     401,
		},
	}

	myDB, _ := db.New()
	router := mux.NewRouter()
	a := &api{
		DB:     myDB,
		router: router,
	}
	handler := http.HandlerFunc(a.authenticationHandler)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			data, err := json.Marshal(tc)
			if err != nil {
				log.Printf("json error: %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/authentication", bytes.NewReader(data))
			if err != nil {
				log.Printf("requst error: %v", err.Error())
			}

			req.Header.Set("Content-Type", "application/json")
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tc.Want, recorder.Code)
		})
	}
}
