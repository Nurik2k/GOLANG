package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/login" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "This request not POST!", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Not Json!", http.StatusBadRequest)
		return
	}

	fmt.Println("Login: ", user.Login, "\n",
		"Password: ", user.Password)

	jsonM, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}
