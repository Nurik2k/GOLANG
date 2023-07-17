package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type xmlUser struct {
	ID     int    `xml:"id"`
	Name   string `xml:"first_name"`
	About  string `xml:"about"`
	Age    int    `xml:"age"`
	Gender string `xml:"gender"`
}

type Users struct {
	Users []xmlUser `xml:"row"`
}

type User struct {
	Id     int
	Name   string
	Age    int
	About  string
	Gender string
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	query := q.Get("query")
	orderField := q.Get("order_field")
	limit := q.Get("limit")
	offset := q.Get("offset")

	users, err := loadUsersFromXML("Web/dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var searchResults []User
	if query != "" {
		searchResults = searchUsers(users, query)
	} else {
		searchResults = users
	}

	sortUsers(searchResults, orderField)

	limitedResults := applyLimitAndOffset(searchResults, limit, offset)

	js, err := json.Marshal(limitedResults)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
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

func main() {
	http.HandleFunc("/users", SearchServer)

	port := ":8080"
	fmt.Println("Server in port:", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
