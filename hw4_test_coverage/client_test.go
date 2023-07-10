package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type xmlUser struct {
	XmlName xml.Name `xml:"root"`
	Row     []struct {
		Id     int    `xml:"id"`
		Name   string `xml:"first_name"`
		Age    int    `xml:"age"`
		About  string `xml:"about"`
		Gender string `xml:"gender"`
	} `xml:"row"`
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")

	xmlData, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		http.Error(w, "Данные не прочитаны", 500)
		return
	}

	var users []xmlUser
	err = xml.Unmarshal(xmlData, &users)
	if err != nil {
		http.Error(w, "Данные не переведены", 500)
		return
	}

	sortUsers(users, orderField)
}

func sortUsers(users []xmlUser, order string) {
	for _, user := range users {
		if order == "id" {
			sort.Slice(user, func(i, j int) bool {
				return user.Row[i].Id < user.Row[j].Id
			})
		}
		if order == "Age" {
			sort.Slice(user, func(i, j int) bool {
				return user.Row[i].Age < user.Row[j].Age
			})
		}
		sort.Slice(user, func(i, j int) bool {
			return user.Row[i].Name < user.Row[j].Name
		})
	}
}

func searchUsers(users []xmlUser, query string) []xmlUser {
	var result []xmlUser
	for i, user := range users {
		if strings.Contains(user.Row[i].Name, query) || strings.Contains(user.Row[i].About, query) {
			result = append(result, user)
		}
	}
	return result
}

func getUsersResponse(users []xmlUser) string {
	var response string
	for i, user := range users {
		response += fmt.Sprintf("ID: %d, Name: %s, About: %s, Age: %d\n", user.Row[i].Id, user.Row[i].Name, user.Row[i].About, user.Row[i].Age)
	}
	return response
}
