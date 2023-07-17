package main

import (
	json2 "encoding/json"
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
	ID    int    `xml:"id"`
	Name  string `xml:"first_name"`
	About string `xml:"about"`
	Age   int    `xml:"age"`
}

type Users struct {
	Users []xmlUser `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")            // Получение значения параметра "query" из URL запроса
	orderField := r.URL.Query().Get("order_field") // Получение значения параметра "order_field" из URL запроса
	limit := r.URL.Query().Get("limit")            // Получение значения параметра "limit" из URL запроса
	offset := r.URL.Query().Get("offset")          // Получение значения параметра "offset" из URL запроса

	users, err := loadUsersFromXML("Web/dataset.xml") // Загрузка пользователей из XML файла
	if err != nil {
		http.Error(w, err.Error(), 500) // В случае ошибки, возвращаем HTTP ошибку 500 и сообщение об ошибке
		return
	}

	if query == "" { // Если параметр "query" пустой, то выполняем сортировку пользователей и отправляем результат
		sortUsers(users, orderField)

		response := getUsersResponse(users)
		fmt.Fprintf(w, response)
		return
	}

	searchResults := searchUsers(users, query) // Выполнение поиска пользователей по запросу

	sortUsers(searchResults, orderField) // Сортировка результатов поиска

	limitedResults := applyLimitAndOffset(searchResults, limit, offset) // Применение ограничения количества и смещения к результатам

	js, err := json2.Marshal(limitedResults) // Преобразование результатов в формат JSON
	if err != nil {
		http.Error(w, err.Error(), 500) // В случае ошибки, возвращаем HTTP ошибку 500 и сообщение об ошибке
		return
	}

	w.Header().Set("Content-Type", "application/json") // Установка заголовка ответа для указания типа содержимого
	w.Write(js)                                        // Отправка ответа в формате JSON
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
		Id:    xmlUser.ID,
		Name:  xmlUser.Name,
		Age:   xmlUser.Age,
		About: xmlUser.About,
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
	default:
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	}
}

func searchUsers(users []User, query string) []User {
	var results []User
	for _, user := range users { // Поиск пользователей, содержащих запрос в имени или описании
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			results = append(results, user)
		}
	}
	return results // Возвращаем найденные результаты
}

func getUsersResponse(users []User) string {
	var response string

	for _, user := range users { // Создание строки ответа с информацией о пользователях
		response += fmt.Sprintf("ID: %d, Name: %s, About: %s, Age: %d\n", user.Id, user.Name, user.About, user.Age)
	}

	return response // Возвращаем строку ответа
}

// Tests
func Test_NegativeLimit(t *testing.T) {
	searchClient := &SearchClient{
		AccessToken: "testAccessToken",
		URL:         "http://example.com",
	}

	req := SearchRequest{
		Limit:  -10,
		Offset: 0,
	}

	_, err := searchClient.FindUsers(req)
	if err != nil {
		t.Error("expected an error, but got nil")
	}

	expectedErrMsg := "limit must be > 0"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}

func TestLimitValidation(t *testing.T) {
	searchClient := &SearchClient{
		AccessToken: "testAccessToken",
		URL:         "http://example.com",
	}

	req := SearchRequest{
		Limit: 30,
	}

	_, err := searchClient.FindUsers(req)
	if err != nil {
		t.Error("expected an error, but got nil")
	}

	// Проверяем, что значение Limit было ограничено до 25
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
	if err != nil {
		t.Error("expected an error, but got nil")
	}

	// Проверяем, что получили ошибку с соответствующим сообщением
	expectedErrMsg := "offset must be > 0"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}
