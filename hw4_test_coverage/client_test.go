package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
		sendError(w, "BadOrderField", http.StatusInternalServerError)
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

func sendError(w http.ResponseWriter, error string, code int) {
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

type TestServer struct {
	server *httptest.Server
	Search SearchClient
}

func (ts *TestServer) Close() {
	ts.server.Close()
}

func newTestServer(token string) TestServer {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{token, server.URL}

	return TestServer{server, client}
}

func TestLimitLow(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	_, err := ts.Search.FindUsers(SearchRequest{
		Limit: -1,
	})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "limit must be > 0" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestLimitHigh(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	response, _ := ts.Search.FindUsers(SearchRequest{
		Limit: 100,
	})

	if len(response.Users) != 25 {
		t.Errorf("Invalid number of users: %d", len(response.Users))
	}
}

func TestInvalidToken(t *testing.T) {
	ts := newTestServer(AccessToken + "invalid")
	defer ts.Close()

	_, err := ts.Search.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "Bad AccessToken" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestInvalidOrderField(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	_, err := ts.Search.FindUsers(SearchRequest{
		OrderBy:    OrderByAsc,
		OrderField: "Foo",
	})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "OrderFeld Foo invalid" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestOffsetLow(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	_, err := ts.Search.FindUsers(SearchRequest{
		Offset: -1,
	})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "offset must be > 0" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestFindUserByName(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	response, _ := ts.Search.FindUsers(SearchRequest{
		Query: "Annie",
		Limit: 1,
	})

	if len(response.Users) != 1 {
		t.Errorf("Invalid number of users: %d", len(response.Users))
		return
	}

	if response.Users[0].Name != "Annie Osborn" {
		t.Errorf("Invalid user found: %v", response.Users[0])
		return
	}
}

func TestLimitOffset(t *testing.T) {
	ts := newTestServer(AccessToken)
	defer ts.Close()

	response, _ := ts.Search.FindUsers(SearchRequest{
		Limit:  3,
		Offset: 0,
	})

	if len(response.Users) != 3 {
		t.Errorf("Invalid number of users: %d", len(response.Users))
		return
	}

	if response.Users[2].Name != "Brooks Aguilar" {
		t.Errorf("Invalid user at position 3: %v", response.Users[2])
		return
	}

	response, _ = ts.Search.FindUsers(SearchRequest{
		Limit:  5,
		Offset: 2,
	})

	if len(response.Users) != 5 {
		t.Errorf("Invalid number of users: %d", len(response.Users))
		return
	}

	if response.Users[0].Name != "Brooks Aguilar" {
		t.Errorf("Invalid user at position 3: %v", response.Users[0])
		return
	}
}

func TestFatalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Fatal Error", http.StatusInternalServerError)
	}))
	client := SearchClient{AccessToken, server.URL}
	defer server.Close()

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "SearchServer fatal error" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestCantUnpackError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Some Error", http.StatusBadRequest)
	}))
	client := SearchClient{AccessToken, server.URL}
	defer server.Close()

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if !strings.Contains(err.Error(), "cant unpack error json") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestUnknownBadRequestError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendError(w, "Unknown Error", http.StatusBadRequest)
	}))
	client := SearchClient{AccessToken, server.URL}
	defer server.Close()

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if !strings.Contains(err.Error(), "unknown bad request error") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestCantUnpackResultError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "None")
	}))
	client := SearchClient{AccessToken, server.URL}
	defer server.Close()

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if !strings.Contains(err.Error(), "cant unpack result json") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	client := SearchClient{AccessToken, server.URL}
	defer server.Close()

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if !strings.Contains(err.Error(), "timeout for") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestUnknownError(t *testing.T) {
	client := SearchClient{AccessToken, "http://invalid-server/"}

	_, err := client.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if !strings.Contains(err.Error(), "unknown error") {
		t.Errorf("Invalid error: %v", err.Error())
	}
}
