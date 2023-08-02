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
	Id        string `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Birthday  string `json:"birthday"`
}

type Handler struct {
	db *sql.DB
}

func ConnectToDB(db1 *sql.DB) (*Handler, error) {
	handler := &Handler{
		db: db1,
	}

	err := handler.db.Ping()
	if err != nil {
		return nil, err
	}

	return handler, nil
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

	var user User

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signed, err := h.DBSignIn(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	jsonM, err := json.Marshal(signed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
		return
	}

	usere, err := h.DBAddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal(usere)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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

	users, err := h.DBGetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

func (h Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "GetUserById" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	userID := parts[2]

	if r.Method != http.MethodGet {
		http.Error(w, "This request not GET!", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.DbGetUserById(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

func (h Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPut {
		http.Error(w, "This request not PUT!", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "EditUser" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	userID := parts[2]

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Id = userID

	editUser, err := h.DBEditUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal(editUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "DeleteUser" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	userID := parts[2]

	if r.Method != http.MethodDelete {
		http.Error(w, "This request not Delete!", http.StatusMethodNotAllowed)
		return
	}

	var err error

	userDeleted, err := h.DbDeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal(userDeleted)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

// DB
func (h Handler) DBSignIn(login, password string) (tf bool, err error) {
	ctx := context.Background()

	tsql := "SELECT Password FROM GoUser WHERE Login = @Login"

	rows, err := h.db.QueryContext(ctx, tsql, sql.Named("Login", login))
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var Password string

	if rows.Next() {
		err := rows.Scan(&Password)
		if err != nil {
			return false, err
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	if Password != password {
		return false, nil
	}

	return true, nil
}

func (h Handler) DBGetUsers() ([]User, error) {
	ctx := context.Background()

	tsql := "SELECT Id, Login, Password, First_Name, Name, Last_Name, Birthday FROM GoUser"

	rows, err := h.db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.Name, &user.LastName, &user.Birthday)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (h Handler) DbGetUserById(id string) (user User, err error) {
	ctx := context.Background()

	if h.db == nil {
		return user, err
	}

	err = h.db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "SELECT Id, Login, Password, First_Name, Name, Last_Name, Birthday FROM GoUser WHERE Id = @Id"

	rows, err := h.db.QueryContext(ctx, tsql, sql.Named("Id", id))
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.Name, &user.LastName, &user.Birthday)
		if err != nil {
			return user, err
		}
	}

	return user, nil

}

func (h Handler) DBAddUser(user User) (users []User, err error) {
	ctx := context.Background()

	if h.db == nil {
		return nil, err
	}

	err = h.db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := "INSERT INTO GoUser(Login, Password, First_Name, Name, Last_Name, Birthday) VALUES(@Login, @Password, @First_Name, @Name, @Last_Name, @Birthday);"

	stmt, err := h.db.Prepare(tsql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.FirstName),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.LastName),
		sql.Named("Birthday", user.Birthday),
	)
	if err != nil {
		return nil, err
	}

	users = append(users, user)

	return users, nil
}

func (h Handler) DBEditUser(user User) (users []User, err error) {
	ctx := context.Background()

	err = h.db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := fmt.Sprintf("UPDATE GoUser SET Login = @Login, Password = @Password, First_Name = @First_Name, Name = @Name, Last_name = @Last_name, Birthday = @Birthday WHERE Id = @Id")

	stmt, err := h.db.Prepare(tsql)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(
		ctx,
		sql.Named("Id", user.Id),
		sql.Named("Login", user.Login),
		sql.Named("Password", user.Password),
		sql.Named("First_Name", user.FirstName),
		sql.Named("Name", user.Name),
		sql.Named("Last_Name", user.LastName),
		sql.Named("Birthday", user.Birthday),
	)
	if err != nil {
		return nil, err
	}

	users = append(users, user)

	return users, nil
}

func (h Handler) DbDeleteUser(id string) (string, error) {
	ctx := context.Background()

	err := h.db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	tsql := fmt.Sprintf("DELETE FROM GoUser WHERE Id = @Id;")

	_, err = h.db.ExecContext(ctx, tsql, sql.Named("Id", id))
	if err != nil {
		return "", err
	}

	return "Deleted", nil
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}
