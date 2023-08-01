package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

type User struct {
	Id         string `json:"id"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	First_Name string `json:"first_name"`
	Name       string `json:"name"`
	Last_Name  string `json:"last_name"`
	Birthday   string `json:"birthday"`
}

// Handlers
type Handler struct {
	connectionString string
	ctx              context.Context
	db               *sql.DB
}

func NewHandler(db *sql.DB) (Handler, error) {
	var err error

	nh := Handler{
		connectionString: "Server=localhost;Database=Users;User Id=sa;Password=yourStrong(!)Password;port=1433;MultipleActiveResultSets=true;TrustServerCertificate=true;",
	}

	db, err = sql.Open("sqlserver", nh.connectionString)
	if err != nil {
		return nh, err
	}
	defer db.Close()
	ctx := context.Background()

	err = db.PingContext(ctx)
	if err != nil {
		return nh, err
	}

	return nh, nil
}

func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.URL.Path != "/SignIn" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "This request not GET!", http.StatusMethodNotAllowed)
		return
	}

	var err error

	h, err = NewHandler(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var user User

	signed, err := h.DBSignIn(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	jsonM, err := json.Marshal(signed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}

func (h Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.URL.Path != "/AddUser" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "This request not POST!", http.StatusMethodNotAllowed)
		return
	}

	var user User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	h, err = NewHandler(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	usere, err := h.DBAddUser(user)
	if err != nil {
		http.Error(w, err.Error(), 404)
	}

	jsonM, err := json.Marshal(usere)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}

func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	enableCors(w)

	if r.URL.Path != "/Users" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "This request not GET!", http.StatusMethodNotAllowed)
		return
	}

	var err error

	nh, err := NewHandler(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	users, err := nh.DBGetUsers()
	if err != nil {
		http.Error(w, err.Error(), 404)
	}

	jsonM, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}

func (h Handler) EditUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	enableCors(w)

	if r.URL.Path != "/EditUser" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "This request not PUT!", http.StatusMethodNotAllowed)
		return
	}

	var user User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	h, err = NewHandler(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	editUser, err := h.DBEditUser(user)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}

	jsonM, err := json.Marshal(editUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	enableCors(w)

	if r.URL.Path != "/DeleteUser" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "This request not Delete!", http.StatusMethodNotAllowed)
		return
	}

	var user User
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	h, err = NewHandler(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	userDeleted, err := h.DbDeleteUser(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	jsonM, err := json.Marshal(userDeleted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonM)
}

// DB

func (h Handler) DBSignIn(login, password string) (tf bool, err error) {
	err = h.db.PingContext(h.ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "Select Password from GoUser where Login = @Login"

	rows, err := h.db.QueryContext(h.ctx, tsql)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		err := rows.Scan(&user.Login, &user.Password)
		if err != nil {
			return false, err
		}
	}

	if !strings.Contains(user.Login, login) && !strings.Contains(user.Password, password) {
		return false, err
	}

	return true, nil
}

func (h Handler) DBGetUsers() (users []User, err error) {
	err = h.db.PingContext(h.ctx)
	if err != nil {
		return nil, err
	}
	tsql := "SELECT * FROM GoUser"

	rows, err := h.db.QueryContext(h.ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user User

	for rows.Next() {

		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.First_Name, &user.Name, &user.Last_Name, &user.Birthday)
		if err != nil {
			return nil, err
		}
	}

	users = append(users, user)

	return users, nil
}

func (h Handler) DBAddUser(user User) (users []User, err error) {

	if h.db == nil {
		return nil, err
	}

	err = h.db.PingContext(h.ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "INSERT INTO GoUser(Login, Password, First_Name, Name, Last_Name, Birthday) VALUES(@Login, @Password, @First_Name, @Name, @Last_Name, @Birthday);"

	stmt, err := h.db.Prepare(tsql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(
		h.ctx,
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.First_Name),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.Last_Name),
		sql.Named("Birthday", user.Birthday),
	)

	fmt.Println(row)

	users = append(users, user)

	return users, nil
}

func (h Handler) DBEditUser(user User) (users []User, err error) {

	err = h.db.PingContext(h.ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := fmt.Sprintf("UPDATE GoUser SET Login = @Login, Password = @Password, First_Name = @First_Name, Name = @Name, Last_name = @Last_name, Birthday = @Birthday WHERE Id = @Id")

	// Execute non-query with named parameters
	result, err := h.db.ExecContext(
		h.ctx,
		tsql,
		sql.Named("Id", user.Id),
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.First_Name),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.Last_Name),
		sql.Named("Birthday", user.Birthday),
	)
	if err != nil {
		return nil, err
	}

	result.RowsAffected()

	users = append(users, user)

	return users, nil
}

func (h Handler) DbDeleteUser(id string) (string, error) {
	err := h.db.PingContext(h.ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := fmt.Sprintf("DELETE FROM GoUser WHERE Id = @Id;")

	// Execute non-query with named parameters
	result, err := h.db.ExecContext(h.ctx, tsql, sql.Named("Id", id))
	if err != nil {
		return "", err
	}
	result.RowsAffected()

	return "Deleted", nil
}

// Functions
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
