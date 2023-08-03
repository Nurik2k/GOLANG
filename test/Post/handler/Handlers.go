package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
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

// Есть пакет handler который отвечает только за Handlers, нужен пакет DBStore
type Handler struct {
	dbs *DbStore //DBStore
}

type DbStore struct {
	db *sql.DB
}

func NewHandler(dbs *DbStore) (*Handler, error) {
	handler := &Handler{
		dbs: dbs,
	}

	err := handler.dbs.db.Ping()
	if err != nil {
		return nil, err
	}

	return handler, nil
}

func NewDb() (*DbStore, error) {
	db1, err := sql.Open("sqlserver", "Server=localhost;Database=Users;User Id=sa;Password=yourStrong(!)Password;port=1433;MultipleActiveResultSets=true;TrustServerCertificate=true;")
	if err != nil {
		return &DbStore{}, nil
	}

	dbs := &DbStore{
		db: db1,
	}

	return dbs, nil
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.URL.Path != "/login" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "This request not POST!", http.StatusMethodNotAllowed)
		return
	}

	var err error
	var user User

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	signed, err := h.dbs.SignIn(user.Login, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if signed == false {
		http.Error(w, "Unsigned", http.StatusUnauthorized)
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

func (h *Handler) addUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//Path == /user
	if r.URL.Path != "/user" {
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

	//Возвращать как User
	err = h.dbs.Create(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal([]byte("Added"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	enableCors(w)

	if r.URL.Path != "/users" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "This request not GET!", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.dbs.Get()
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

func (h *Handler) getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "user" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	userID := parts[2]

	if r.Method != http.MethodGet {
		http.Error(w, "This request not GET!", http.StatusMethodNotAllowed)
		return
	}

	user, err := h.dbs.GetById(userID)
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

func (h *Handler) editUser(w http.ResponseWriter, r *http.Request) {
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
	if len(parts) != 3 || parts[1] != "user" {
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

	err := h.dbs.Edit(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal([]byte("Изменено"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

func (h *Handler) deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	enableCors(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "user" {
		http.Error(w, "404 Page not found!", http.StatusNotFound)
		return
	}

	userID := parts[2]

	if r.Method != http.MethodDelete {
		http.Error(w, "This request not Delete!", http.StatusMethodNotAllowed)
		return
	}

	var err error

	err = h.dbs.Delete(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonM, err := json.Marshal([]byte("Deleted"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonM)
}

// DB

type DbInterface interface {
	SignIn(login, password string) (tf bool, err error)
	Get() ([]User, error)
	GetById(id string) (user User, err error)
	Create(user *User) (User, error)
	Edit(user *User) (err error)
	Delete(id string) error
}

func (db *DbStore) SignIn(login, password string) (tf bool, err error) {
	ctx := context.Background()

	err = db.db.PingContext(ctx)
	if err != nil {
		return false, err
	}

	tsql := "SELECT Password FROM GoUser WHERE Login = @Login"

	row := db.db.QueryRowContext(ctx, tsql, sql.Named("Login", login))
	if err != nil {
		return false, err
	}

	var Password string

	err = row.Scan(&Password)
	if err != nil {
		return false, err
	}

	if err := row.Err(); err != nil {
		return false, err
	}

	if Password != password {
		return false, nil
	}

	return true, nil
}

func (db *DbStore) Get() ([]User, error) {
	ctx := context.Background()

	err := db.db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tsql := "SELECT Id, Login, Password, First_Name, Name, Last_Name, Birthday FROM GoUser"

	rows, err := db.db.QueryContext(ctx, tsql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//make
	//limit and offset, Погинация
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

	limitedUsers := applyLimitOffset(users, 25, 0)

	return limitedUsers, nil
}

func (db *DbStore) GetById(id string) (user *User, err error) {
	ctx := context.Background()

	err = db.db.PingContext(ctx)
	if err != nil {
		return user, err
	}

	if db.db == nil {
		return user, err
	}

	tsql := "SELECT Id, Login, Password, First_Name, Name, Last_Name, Birthday FROM GoUser WHERE Id = @Id"

	//Посмотреть возвращение только одной операции row
	row := db.db.QueryRowContext(ctx, tsql, sql.Named("Id", id))

	err = row.Scan(&user.Id, &user.Login, &user.Password, &user.FirstName, &user.Name, &user.LastName, &user.Birthday)
	if err != nil {
		return user, err
	}

	return user, nil
}

// return err
func (db *DbStore) Create(user *User) (err error) {
	ctx := context.Background()

	err = db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	if db.db == nil {
		return err
	}

	//Запрос return @Id. user.Id = @Id
	tsql := "INSERT INTO GoUser(Login, Password, First_Name, Name, Last_Name, Birthday) OUTPUT INSERTED.@Id VALUES(@Login, @Password, @First_Name, @Name, @Last_Name, @Birthday);"

	stmt, err := db.db.Prepare(tsql)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (db *DbStore) Edit(user *User) (err error) {
	ctx := context.Background()

	err = db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("UPDATE GoUser SET Login = @Login, Password = @Password, First_Name = @First_Name, Name = @Name, Last_name = @Last_name, Birthday = @Birthday WHERE Id = @Id")

	stmt, err := db.db.Prepare(tsql)
	if err != nil {
		return err
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
		return err
	}

	return nil
}

func (db *DbStore) Delete(id string) error {
	ctx := context.Background()

	err := db.db.PingContext(ctx)
	if err != nil {
		return err
	}

	tsql := fmt.Sprintf("DELETE FROM GoUser WhERE Id = @Id;")

	_, err = db.db.ExecContext(ctx, tsql, sql.Named("Id", id))
	if err != nil {
		return err
	}

	return nil
}

// other Functions
func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (h Handler) Routes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/login", h.signIn).Methods(http.MethodPost)
	r.HandleFunc("/user", h.addUser).Methods(http.MethodPost)
	r.HandleFunc("/users", h.getUsers).Methods(http.MethodGet)
	r.HandleFunc("/user", h.getUserById).Methods(http.MethodGet)
	r.HandleFunc("/user", h.editUser).Methods(http.MethodPut)
	r.HandleFunc("/user", h.deleteUser).Methods(http.MethodDelete)

	log.Println("Server listening on port 8080")
	//log
	return r
}

func applyLimitOffset(users []User, limit, offset int) []User {
	if limit <= 0 {
		return users
	}

	from := offset
	if from > len(users)-1 {
		return []User{}
	}

	to := offset + limit
	if to > len(users) {
		to = len(users)
	}

	return users[from:to]
}
