package main

import (
	"log"
	"net/http"
	"post/database"
	"post/handler"
)

func main() {
	//Нужно добавить нрвую сущность Repository or DBStore которая, работает с бд +
	db, err := database.NewDb()

	h, err := handler.NewHandler(db)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	r := h.Routes()

	//Поменять наименование роутов по спецификации REST
	//Перенести инициализацию роутов в пакет Handler
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err.Error())
	}

}
