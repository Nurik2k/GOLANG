package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
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
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")

	xmlData, err := os.Open("Web/dataset.xml")
	if err != nil {
		http.Error(w, "File not found", 500)
	}

	reader, err := ioutil.ReadAll(xmlData)
	if err != nil {
		http.Error(w, "Error reading data file", 500)
		return
	}

	var users Users
	err = xml.Unmarshal(reader, &users)
	if err != nil {
		http.Error(w, "Error parsing XML data", 500)
		return
	}

	if query == "" {
		sortUsers(users.Users, orderField)

		response := getUsersResponse(users.Users)
		fmt.Fprintf(w, response)
		return
	}

	searchResults := searchUsers(users.Users, query)

	sortUsers(searchResults, orderField)

	response := getUsersResponse(searchResults)
	fmt.Fprintf(w, response)

}

func sortUsers(users []xmlUser, orderField string) {
	switch orderField {
	case "Id":
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID < users[j].ID
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

func searchUsers(users []xmlUser, query string) []xmlUser {
	var results []xmlUser
	for _, user := range users {
		if strings.Contains(user.Name, query) || strings.Contains(user.About, query) {
			results = append(results, user)
		}
	}

	return results
}

func getUsersResponse(users []xmlUser) string {
	var response string

	for _, user := range users {
		response += fmt.Sprintf("ID: %d, Name: %s, About: %s, Age: %d\n", user.ID, user.Name, user.About, user.Age)
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
