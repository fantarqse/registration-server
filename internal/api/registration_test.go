package api

import (
	"bytes"
	"encoding/json"
	"github.com/fantarqse/registrationserver/internal/db"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {
	testCase := []struct {
		Name     string `json:"-"`
		Login    string `json:"login"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Want     int    `json:"-"`
	}{
		{
			Name:     "Test1",
			Login:    "name",
			Password: "pAssword123",
			Email:    "name@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test2",
			Login:    "abc",
			Password: "Abcabc11",
			Email:    "abc@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test3",
			Login:    "Name",
			Password: "NameName1!",
			Email:    "Name@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test4",
			Login:    "NameName",
			Password: "NameName2*",
			Email:    "NameName@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test5",
			Login:    "Alex_Grey",
			Password: "Alex_Grey23",
			Email:    "alexgrey23@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test6",
			Login:    "Roman_Tkachenko",
			Password: "Roman_Tkachenko09",
			Email:    "RT09@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test7",
			Login:    "name_name_name",
			Password: "Name_namename77",
			Email:    "name_name_name@gmail.com",
			Want:     200,
		},

		{
			Name:  "Test8",
			Login: "NameName",
			Email: "login@gmail.com",
			Want:  400,
		},
		{
			Name:     "Test9",
			Login:    "NameName",
			Password: "Password123",
			Want:     400,
		},
		{
			Name:     "Test10",
			Login:    "Ab",
			Password: "Password123",
			Email:    "login@gmail.com",
			Want:     400,
		},
		{
			Name:     "Test11",
			Login:    "NameName",
			Password: "Password",
			Email:    "login@gmail.com",
			Want:     400,
		},
		{
			Name:     "Test12",
			Login:    "NameName",
			Password: "Password123",
			Email:    "logingmail.com",
			Want:     400,
		},
		{
			Name:     "Test13",
			Login:    "a_b_c",
			Password: "SomeValidPassword123",
			Email:    "SomeValidEmail@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test14",
			Login:    "SomeValidLogin01",
			Password: "a________Aa",
			Email:    "SomeValidLogin01@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test15",
			Login:    "SomeValidLogin15",
			Password: "A______a",
			Email:    "SomeValidLogin15@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test16",
			Login:    "SomeValidLogin16",
			Password: "a______A",
			Email:    "SomeValidLogin16@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test17",
			Login:    "SomeValidLogin17",
			Password: "Aa________",
			Email:    "SomeValidLogin17@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test18",
			Login:    "SomeValidLogin18",
			Password: "aA________",
			Email:    "SomeValidLogin18@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test19",
			Login:    "SomeValidLogin19",
			Password: "________Aa",
			Email:    "SomeValidLogin19@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test20",
			Login:    "SomeValidLogin20",
			Password: "________aA",
			Email:    "SomeValidLogin20@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test21",
			Login:    "SomeValidLogin21",
			Password: "____Aa____",
			Email:    "SomeValidLogin21@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test22",
			Login:    "SomeValidLogin22",
			Password: "____aA____",
			Email:    "SomeValidLogin22@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test23",
			Login:    "SomeValidLogin23",
			Password: "A___A___a",
			Email:    "SomeValidLogin23@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test24",
			Login:    "SomeValidLogin24",
			Password: "A___a___A",
			Email:    "SomeValidLogin24@gmail.com",
			Want:     200,
		},
		{
			Name:     "Test25",
			Login:    "SomeValidLogin25",
			Password: "a___A___A",
			Email:    "SomeValidLogin25@gmail.com",
			Want:     200,
		}, {
			Name:     "Test26",
			Login:    "SomeValidLogin26",
			Password: "_Aa______",
			Email:    "SomeValidLogin26@gmail.com",
			Want:     200,
		}, {
			Name:     "Test27",
			Login:    "SomeValidLogin27",
			Password: "_aA______",
			Email:    "SomeValidLogin27@gmail.com",
			Want:     200,
		}, {
			Name:     "Test28",
			Login:    "SomeValidLogin28",
			Password: "__Aa_____",
			Email:    "SomeValidLogin28@gmail.com",
			Want:     200,
		}, {
			Name:     "Test29",
			Login:    "SomeValidLogin29",
			Password: "__aA_____",
			Email:    "SomeValidLogin29@gmail.com",
			Want:     200,
		}, {
			Name:     "Test30",
			Login:    "SomeValidLogin30",
			Password: "_____Aa__",
			Email:    "SomeValidLogin30@gmail.com",
			Want:     200,
		}, {
			Name:     "Test31",
			Login:    "SomeValidLogin31",
			Password: "_____aA__",
			Email:    "SomeValidLogin31@gmail.com",
			Want:     200,
		}, {
			Name:     "Test32",
			Login:    "SomeValidLogin32",
			Password: "______Aa_",
			Email:    "SomeValidLogin32@gmail.com",
			Want:     200,
		}, {
			Name:     "Test33",
			Login:    "SomeValidLogin33",
			Password: "______aA_",
			Email:    "SomeValidLogin33@gmail.com",
			Want:     200,
		},
	}

	myDB, _ := db.New()
	router := mux.NewRouter()
	a := &api{
		DB:     myDB,
		router: router,
	}
	handler := http.HandlerFunc(a.registrationHandler)

	for _, tc := range testCase {
		t.Run(tc.Name, func(t *testing.T) {
			data, err := json.Marshal(tc)
			if err != nil {
				log.Printf("json error: %v", err.Error())
			}

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/registration", bytes.NewReader(data))
			if err != nil {
				log.Printf("requst error: %v", err.Error())
			}

			req.Header.Set("Content-Type", "application/json")
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tc.Want, recorder.Code)
		})
	}
}
