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

type xmlUser struct {
	ID     int    `xml:"id"`
	Name   string `xml:"first_name" xml:"last_name"`
	About  string `xml:"about"`
	Age    int    `xml:"age"`
	Gender string `xml:"gender"`
}

type Users struct {
	Users []xmlUser `xml:"row"`
}

var AccessToken = "abc123"

func SearchServer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := q.Get("query")            // Получаем значение параметра "query" из URL-запроса
	orderField := q.Get("order_field") // Получаем значение параметра "order_field" из URL-запроса
	limit := q.Get("limit")            // Получаем значение параметра "limit" из URL-запроса
	offset := q.Get("offset")          // Получаем значение параметра "offset" из URL-запроса

	users, err := loadUsersFromXML("Web/dataset.xml") // Загружаем данные пользователей из XML-файла
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Если произошла ошибка, возвращаем её как HTTP-ответ с кодом 500
		return
	}

	var searchResults []User
	if query != "" {
		searchResults = searchUsers(users, query) // Если значение параметра "query" не пустое, выполняем поиск пользователей
	} else {
		searchResults = users // Иначе возвращаем всех пользователей
	}

	sortUsers(searchResults, orderField) // Сортируем результаты по указанному полю

	limitedResults := applyLimitAndOffset(searchResults, limit, offset) // Применяем ограничение и смещение к результатам поиска

	js, err := json.Marshal(limitedResults) // Преобразуем ограниченные результаты в формат JSON
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Если произошла ошибка при маршалинге в JSON, возвращаем её как HTTP-ответ с кодом 500
		return
	}

	resp, err := http.Get("http://external-api.com") // Выполняем GET-запрос к внешнему API
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Если произошла ошибка при выполнении запроса, возвращаем её как HTTP-ответ с кодом 500
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // Читаем тело ответа
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // Если произошла ошибка при чтении тела ответа, возвращаем её как HTTP-ответ с кодом 500
		return
	}

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		http.Error(w, "Bad AccessToken", http.StatusUnauthorized) // Если получен статус-код 401 Unauthorized, возвращаем ошибку с кодом 401
		return
	case http.StatusInternalServerError:
		http.Error(w, "SearchServer fatal error", http.StatusInternalServerError) // Если получен статус-код 500 Internal Server Error, возвращаем ошибку с кодом 500
		return
	case http.StatusBadRequest:
		errRespBytes, err := json.Marshal(body) // Преобразуем тело ответа в формат JSON
		if err != nil {
			http.Error(w, fmt.Sprintf("cant pack error json: %s", err), http.StatusInternalServerError) // Если произошла ошибка при маршалинге в JSON, возвращаем её как HTTP-ответ с кодом 500
			return
		}
		errResp := string(errRespBytes)
		if errResp == "ErrorBadOrderField" {
			http.Error(w, fmt.Sprintf("OrderField %s invalid", orderField), http.StatusBadRequest) // Если получен статус-код 400 Bad Request с ошибкой "ErrorBadOrderField", возвращаем ошибку с кодом 400
			return
		}
		http.Error(w, fmt.Sprintf("unknown bad request error: %s", errResp), http.StatusBadRequest) // Если получен статус-код 400 Bad Request с неизвестной ошибкой, возвращаем ошибку с кодом 400
		return
	}

	w.Header().Set("Content-Type", "application/json") // Устанавливаем заголовок HTTP-ответа для указания типа контента как JSON
	w.Write(js)                                        // Отправляем ограниченные результаты в формате JSON как HTTP-ответ
}

func searchUsers(users []User, query string) []User {
	var results []User
	for _, user := range users {
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			results = append(results, user)
		}
	}
	return results
}

func loadUsersFromXML(filename string) ([]User, error) {
	xmlData, err := os.Open(filename) // Открытие XML файла
	if err != nil {
		return nil, err // В случае ошибки, возвращаем ошибку
	}
	defer xmlData.Close() // Закрытие файла после окончания работы функции

	var users Users

	file, err := ioutil.ReadAll(xmlData) // Чтение содержимого файла
	if err != nil {
		return nil, err // В случае ошибки, возвращаем ошибку
	}
	xml.Unmarshal(file, &users) // Разбор XML данных и сохранение результатов в структуру users

	return convertXMLUsersToUsers(users.Users), nil // Преобразование пользователей из формата XML в формат User и возвращение результатов
}

func convertXMLUsersToUsers(xmlUsers []xmlUser) []User {
	var users []User
	for _, xmlUser := range xmlUsers {
		user := convertXMLUserToUser(xmlUser)
		users = append(users, user)
	}
	return users
}

func convertXMLUserToUser(xmlUser xmlUser) User {
	return User{
		Id:     xmlUser.ID,
		Name:   xmlUser.Name,
		Age:    xmlUser.Age,
		About:  xmlUser.About,
		Gender: xmlUser.Gender,
	}
}

func applyLimitAndOffset(result []User, limitStr, offsetStr string) []User {
	limit, _ := strconv.Atoi(limitStr)   // Преобразование строкового значения "limitStr" в целое число
	offset, _ := strconv.Atoi(offsetStr) // Преобразование строкового значения "offsetStr" в целое число

	if limit > 0 { // Если задано ограничение
		from := offset // Определение начального индекса
		if from > len(result)-1 {
			return []User{} // Если начальный индекс превышает длину результата, возвращаем пустой массив
		} else {
			to := offset + limit // Определение конечного индекса
			if to > len(result) {
				to = len(result) // Если конечный индекс превышает длину результата, устанавливаем его равным длине результата
			}

			return result[from:to] // Возвращаем срез результата с примененными ограничениями
		}
	}

	return result // Если ограничение не задано, возвращаем весь результат
}

func sortUsers(users []User, orderField string) {
	switch orderField { // В зависимости от значения "orderField" выполняем сортировку пользователей
	case "Id":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})
	case "Age":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})

	case "Name", " ":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	default:
		return
	}
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
