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
			Login:    "NameName",
			Password: "Password123",
			Email:    "login@gmail.com",
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
			req, err := http.NewRequest("POST", "/signup", bytes.NewReader(data))
			if err != nil {
				log.Printf("requst error: %v", err.Error())
			}

			req.Header.Set("Content-Type", "application/json")
			handler.ServeHTTP(recorder, req)
			assert.Equal(t, tc.Want, recorder.Code)
		})
	}
}
