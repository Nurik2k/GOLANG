package main

import (
	json2 "encoding/json"
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
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	users, err := loadUsersFromXML("Web/dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if query == "" {
		sortUsers(users, orderField)

		response := getUsersResponse(users)
		fmt.Fprintf(w, response)
		return
	}

	searchResults := searchUsers(users, query)

	sortUsers(searchResults, orderField)

	limitedResults := applyLimitAndOffset(searchResults, limit, offset)

	js, err := json2.Marshal(limitedResults)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func loadUsersFromXML(filename string) ([]User, error) {
	xmlData, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlData.Close()

	var users Users

	file, err := ioutil.ReadAll(xmlData)
	if err != nil {
		return nil, err
	}
	xml.Unmarshal(file, &users)

	return convertXMLUsersToUsers(users.Users), nil
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
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit > 0 {
		from := offset
		if from > len(result)-1 {
			return []User{}
		} else {
			to := offset + limit
			if to > len(result) {
				to = len(result)
			}

			return result[from:to]
		}
	}

	return result
}

// todo: order by
func sortUsers(users []User, orderField string) {
	switch orderField {
	case "Id":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})
	case "Name", "":
		// todo: check sort by strings desc and asc
	case "Age":
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})
	default:
		// todo: return error
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	}
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

func getUsersResponse(users []User) string {
	var response string

	for _, user := range users {
		response += fmt.Sprintf("ID: %d, Name: %s, About: %s, Age: %d, Gender: %s\n", user.Id, user.Name, user.About, user.Age, user.Gender)
	}

	return response
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
