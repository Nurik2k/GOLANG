package main

import (
	"database/sql"
	"log"
	"net/http"
	"post/handler"
)

func main() {
	db, err := sql.Open("sqlserver", "Server=localhost;Database=Users;User Id=sa;Password=yourStrong(!)Password;port=1433;MultipleActiveResultSets=true;TrustServerCertificate=true;")
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	h, err := handler.ConnectToDB(db)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	http.HandleFunc("/SignIn", h.SignIn)
	http.HandleFunc("/AddUser", h.AddUser)
	http.HandleFunc("/Users", h.GetUsers)
	http.HandleFunc("/GetUserById/", h.GetUserById)
	http.HandleFunc("/EditUser/", h.EditUser)
	http.HandleFunc("/DeleteUser/", h.DeleteUser)

	log.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
