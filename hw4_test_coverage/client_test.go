package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type RowUsers struct {
	Rows []Xmluser `xml:"row"`
}

type Xmluser struct {
	XMLName xml.Name `xml:"row"`
	Id      int      `xml:"id"`
	Name    string   `xml:"first_name" +xml:"last_name"`
	About   string   `xml:"about"`
	Age     int      `xml:"age"`
	Gender  string   `xml:"gender"`
}

var AccessToken = "abc123"

func SearchServer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	query := q.Get("query")                       // Получение значения параметра "query"
	limit, _ := strconv.Atoi(q.Get("limit"))      // Получение значения параметра "limit"
	offset, _ := strconv.Atoi(q.Get("offset"))    // Получение значения параметра "offset"
	orderBy, _ := strconv.Atoi(q.Get("order_by")) // Получение значения параметра "order_by"
	orderField := q.Get("order_field")            // Получение значения параметра "order_field"

	// Проверка валидности access token
	if r.Header.Get("AccessToken") != AccessToken {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	// Чтение и декодирование XML данных из файла
	data, err := readRowUsers("dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filteredRows := filterRowsByQuery(data.Rows, query) // Фильтрация строк данных на основе запроса

	users := convertRowsToUsers(filteredRows) // Преобразование строк данных в структуры пользователей

	// Сортировка пользователей на основе заданного поля и направления
	orderedUsers, err := sortUsers(users, orderBy, orderField, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usersWithLimitOffset := applyLimitOffset(orderedUsers, limit, offset) // Применение ограничения по количеству и смещению к списку пользователей

	// Преобразование списка пользователей в JSON
	js, err := json.Marshal(usersWithLimitOffset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// Чтение и декодирование XML данных из файла
func readRowUsers(filename string) (RowUsers, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return RowUsers{}, fmt.Errorf("failed to open XML file: %w", err)
	}
	defer xmlFile.Close()

	var data RowUsers
	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return RowUsers{}, fmt.Errorf("failed to read XML file: %w", err)
	}
	err = xml.Unmarshal(b, &data)
	if err != nil {
		return RowUsers{}, fmt.Errorf("failed to unmarshal XML data: %w", err)
	}

	return data, nil
}

// Фильтрация строк данных на основе запроса
func filterRowsByQuery(rows []Xmluser, query string) []Xmluser {
	if query == "" {
		return rows
	}

	var filteredRows []Xmluser
	for _, row := range rows {
		// Проверка, содержится ли запрос в поле FirstName, LastName или About
		queryMatch := strings.Contains(row.Name, query) ||
			strings.Contains(row.About, query)

		if queryMatch {
			filteredRows = append(filteredRows, row)
		}
	}

	return filteredRows
}

// Преобразование строк данных в структуры пользователей
func convertRowsToUsers(rows []Xmluser) []User {
	var users []User
	for _, row := range rows {
		user := User{
			Id:     row.Id,
			Name:   row.Name,
			Age:    row.Age,
			About:  row.About,
			Gender: row.Gender,
		}
		users = append(users, user)
	}

	return users
}

// Сортировка пользователей на основе заданного поля и направления
func sortUsers(users []User, orderBy int, orderField string, w http.ResponseWriter) ([]User, error) {
	if orderBy == OrderByAsIs {
		return users, nil
	}

	var isLess func(i, j User) bool

	switch orderField {
	case "Id":
		isLess = func(i, j User) bool {
			return i.Id < j.Id
		}
	case "Age":
		isLess = func(i, j User) bool {
			return i.Age < j.Age
		}
	case "Name", "":
		isLess = func(i, j User) bool {
			return i.Name < j.Name
		}
	default:
		sendSearchErrorResponse(w, "BadOrderField", http.StatusInternalServerError)
	}

	if orderBy == orderDesc {
		sort.Slice(users, func(i, j int) bool {
			return isLess(users[j], users[i])
		})
	} else {
		sort.Slice(users, func(i, j int) bool {
			return isLess(users[i], users[j])
		})
	}

	return users, nil
}

// Применение ограничения по количеству и смещению к списку пользователей
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

func sendSearchErrorResponse(w http.ResponseWriter, error string, code int) {
	js, err := json.Marshal(SearchErrorResponse{error})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

// Tests
func Test_NegativeLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	searchClient := SearchClient{AccessToken, ts.URL}

	req := SearchRequest{
		Limit:  -10,
		Offset: 0,
	}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	expectedErrMsg := "limit must be > 0"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestLimitValidation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	searchClient := SearchClient{AccessToken, ts.URL}

	req := SearchRequest{
		Limit: 30,
	}

	response, _ := searchClient.FindUsers(req)

	if len(response.Users) != 25 {
		t.Errorf("expected limit to be %d, but got %d", 25, len(response.Users))
	}
}

func TestOffsetValidation(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	searchClient := SearchClient{AccessToken, ts.URL}

	req := SearchRequest{
		Offset: -10,
	}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	// Проверяем, что получили ошибку с соответствующим сообщением
	expectedErrMsg := "offset must be > 0"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestTimeOut(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer ts.Close()
	searchClient := SearchClient{AccessToken, ts.URL}

	req := SearchRequest{}

	_, err := searchClient.FindUsers(req)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			t.Errorf("timeout for %s", err.Error())
		}
	}
}

func TestBadAccessToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()
	searchClient := SearchClient{AccessToken + "Invalid", ts.URL}

	req := SearchRequest{}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("Empty error")
	}
	if err.Error() != "Bad AccessToken" {
		t.Error("Invalid error!")
	}
}

func TestSearchServer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Fatal error", http.StatusInternalServerError)
	}))
	defer ts.Close()
	searchClient := SearchClient{AccessToken, ts.URL}

	req := SearchRequest{}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("Empty error")
	}
	if err.Error() != "SearchServer fatal error" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestCantUnpackErrorJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Error", http.StatusBadRequest)
	}))
	searchClient := SearchClient{AccessToken, ts.URL}
	defer ts.Close()

	req := SearchRequest{}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("Empty error")
	}
	if !strings.Contains(err.Error(), "cant unpack error json") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestFindUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	searchClient := SearchClient{AccessToken, ts.URL}
	defer ts.Close()

	req := SearchRequest{
		Query: "Hilda",
		Limit: 1,
	}

	user, err := searchClient.FindUsers(req)
	if err != nil {
		t.Error(err.Error())
	}
	if len(user.Users) != 1 {
		t.Error("Limit is exceeded")
	}
	if user.Users[0].Name != "Hilda" {
		t.Error("No such user")
	}
}

func TestErrorBadOrderField(t *testing.T) {
	ts := httptest.NewServer()
	defer ts.Close()

	req := SearchRequest{
		OrderField: ErrorBadOrderField,
	}

	if req.OrderField == "ErrorBadOrderField" {
		t.Error("OrderField invalid")
	}
}
