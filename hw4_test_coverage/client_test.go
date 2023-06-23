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
)

const AccessToken = "abc123"

type XmlData struct {
	XMLName xml.Name `xml:"root"`
	Rows    []XmlRow `xml:"row"`
}

type XmlRow struct {
	XMLName   xml.Name `xml:"row"`
	Id        int      `xml:"id"`
	FirstName string   `xml:"first_name"`
	LastName  string   `xml:"last_name"`
	About     string   `xml:"about"`
	Age       int      `xml:"age"`
	Gender    string   `xml:"gender"`
}

func SearchServer(writer http.ResponseWriter, req *http.Request) {
	if req.Header.Get("AccessToken") != AccessToken {
		http.Error(writer, "Invalid access token", http.StatusUnauthorized)
		return
	}

	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer xmlFile.Close()

	var (
		data   XmlData
		result []User
	)

	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	xml.Unmarshal(b, &data)

	q := req.URL.Query()

	query := q.Get("query")
	for _, row := range data.Rows {
		if query != "" {
			queryMatch := strings.Contains(row.FirstName, query) ||
				strings.Contains(row.LastName, query) ||
				strings.Contains(row.About, query)

			if !queryMatch {
				continue
			}
		}

		result = append(result, User{
			Id:     row.Id,
			Name:   row.FirstName + " " + row.LastName,
			Age:    row.Age,
			About:  row.About,
			Gender: row.Gender,
		})
	}

	orderBy, _ := strconv.Atoi(q.Get("order_by"))
	if orderBy != OrderByAsIs {
		var isLess func(u1, u2 User) bool

		switch q.Get("order_field") {
		case "Id":
			isLess = func(u1, u2 User) bool {
				return u1.Id < u2.Id
			}

		case "Age":
			isLess = func(u1, u2 User) bool {
				return u1.Age < u2.Age
			}

		case "Name":
			fallthrough

		case "":
			isLess = func(u1, u2 User) bool {
				return u1.Name < u2.Name
			}

		default:
			sendError(writer, "ErrorBadOrderField", http.StatusBadRequest)
			return
		}

		sort.Slice(result, func(i, j int) bool {
			return isLess(result[i], result[j]) && (orderBy == orderDesc)
		})
	}

	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit > 0 {
		from := offset
		if from > len(result)-1 {
			result = []User{}
		} else {
			to := offset + limit
			if to > len(result) {
				to = len(result)
			}

			result = result[from:to]
		}
	}

	js, err := json.Marshal(result)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func sendError(writer http.ResponseWriter, error string, code int) {
	js, err := json.Marshal(SearchErrorResponse{error})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	fmt.Fprintln(writer, string(js))
}

type TestServer struct {
	server *httptest.Server
	Search SearchClient
}

func newTestServer(token string) TestServer {
	server := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{token, server.URL}

	return TestServer{server, client}
}

func TestLimitLow(t *testing.T) {
	testServer := newTestServer(AccessToken)
	defer testServer.server.Close()

	_, err := testServer.Search.FindUsers(SearchRequest{
		Limit: -1,
	})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "limit must be > 0" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestLimitHigh(t *testing.T) {
	testServer := newTestServer(AccessToken)
	defer testServer.server.Close()

	response, _ := testServer.Search.FindUsers(SearchRequest{
		Limit: 100,
	})

	if len(response.Users) != 25 {
		t.Errorf("Invalid number of users: %d", len(response.Users))
	}
}

func TestInvalidToken(t *testing.T) {
	testServer := newTestServer(AccessToken + "invalid")
	defer testServer.server.Close()

	_, err := testServer.Search.FindUsers(SearchRequest{})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "Bad AccessToken" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestInvalidOrderField(t *testing.T) {
	testServer := newTestServer(AccessToken)
	defer testServer.server.Close()

	_, err := testServer.Search.FindUsers(SearchRequest{
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
	testServer := newTestServer(AccessToken)
	defer testServer.server.Close()

	_, err := testServer.Search.FindUsers(SearchRequest{
		Offset: -1,
	})

	if err == nil {
		t.Errorf("Empty error")
	} else if err.Error() != "offset must be > 0" {
		t.Errorf("Invalid error: %v", err.Error())
	}
}

func TestFindUserByName(t *testing.T) {
	testServer := newTestServer(AccessToken)
	defer testServer.server.Close()

	response, _ := testServer.Search.FindUsers(SearchRequest{
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
