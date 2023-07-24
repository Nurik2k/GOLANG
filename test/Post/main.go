package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", SignIn)
	fmt.Println("Сервер на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера")
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("User: ", user)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data allowed!"))
}
