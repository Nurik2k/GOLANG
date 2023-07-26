package main

import (
	"fmt"
	"net/http"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/login", SignUp)
	fmt.Println("Сервер на порту 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера")
	}

}
