package main

import (
	"fmt"
	"net/http"
	"post/handler"
)

func main() {
	Routes()

	fmt.Println("Сервер на порту 8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Ошибка при запуске сервера")
	}
}

func Routes() {
	var h handler.Handler
	http.HandleFunc("/SignIn", h.SignIn)
	http.HandleFunc("/AddUser", h.AddUser)
	http.HandleFunc("/Users", h.GetUsers)
	http.HandleFunc("/EditUser", h.EditUsers)
	http.HandleFunc("/DeleteUser", h.DeleteUser)
}
