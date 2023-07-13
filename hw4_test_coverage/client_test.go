package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")

	xmlData, err := os.Open("Web/dataset.xml")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	reader, err := ioutil.ReadAll(xmlData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var users Users
	err = xml.Unmarshal(reader, &users)
	if err != nil {
		http.Error(w, err.Error(), 500)
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

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	w.Write(js)
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
	if err == nil {
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
	if err == nil {
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
	if err == nil {
		t.Error("expected an error, but got nil")
	}

	// Проверяем, что получили ошибку с соответствующим сообщением
	expectedErrMsg := "offset must be > 0"
	if err.Error() != expectedErrMsg {
		t.Errorf("expected error message '%s', but got '%s'", expectedErrMsg, err.Error())
	}
}
