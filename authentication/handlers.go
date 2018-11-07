package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// Credentials Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

func responseWithMessage(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//HomeHandler handlers
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	message := map[string]string{"message": "Yey"}
	response, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(response)
}

//SignupHandler handlers
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		responseWithMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if _, err = db.Query("insert into users values (?, ?)", creds.Username, string(hashedPassword)); err != nil {
		responseWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	responseWithMessage(w, http.StatusBadRequest, "Your registration success!")
	return
}

//SigninHandler handlers
func SigninHandler(w http.ResponseWriter, r *http.Request) {
	creds := &Credentials{}
	err := json.NewDecoder(r.Body).Decode(creds)
	if err != nil {
		responseWithMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	result := db.QueryRow("select password from users where username=?", creds.Username)
	if err != nil {
		responseWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			responseWithMessage(w, http.StatusUnauthorized, err.Error())
			return
		}
		responseWithMessage(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		responseWithMessage(w, http.StatusUnauthorized, err.Error())
		return
	}

	responseWithMessage(w, http.StatusOK, "Welcome!")
	return
}
