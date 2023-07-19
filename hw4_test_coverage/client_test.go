package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
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
	orderedUsers, err := sortUsers(users, orderBy, orderField)
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
func sortUsers(users []User, orderBy int, orderField string) ([]User, error) {
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
		return nil, fmt.Errorf("invalid order field: %s", orderField)
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

// Tests

func Test_NegativeLimit(t *testing.T) {
	searchClient := SearchClient{
		AccessToken: "testAccessToken",
		URL:         "http://example.com",
	}

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
	searchClient := SearchClient{
		AccessToken: "testAccessToken",
		URL:         "http://example.com",
	}

	req := SearchRequest{
		Limit: 30,
	}

	_, err := searchClient.FindUsers(req)
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	expectedLimit := 25
	if req.Limit != expectedLimit {
		t.Errorf("expected limit to be %d, but got %d", expectedLimit, req.Limit)
	}
}

func TestOffsetValidation(t *testing.T) {
	searchClient := &SearchClient{
		AccessToken: "testAccessToken",
		URL:         "http://example.com",
	}

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
