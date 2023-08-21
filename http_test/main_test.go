package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	cases := []TestCase{
		{
			"42",
			`{"status": 200, "resp": {"user": 42}}`,
			http.StatusOK,
		},
		{
			"500",
			`{"status": 500, "err": "db_error"}`,
			http.StatusInternalServerError,
		},
	}

	for caseNum, item := range cases {
		url := "http://example.com/api/user?id=" + item.ID
		req := httptest.NewRequest("GET", url, nil)
		w := httptest.NewRecorder()

		GetUser(w, req)

		if w.Code != item.StatusCode {
			t.Errorf("[%d] wrong status code: got %d, expected %d", caseNum, w.Code, item.StatusCode)
		}

		resp := w.Result()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err.Error())
		}

		bodyStr := string(body)
		if bodyStr != item.Response {
			t.Errorf("[%d] Wrong response: got %+v, expected %+v", caseNum, bodyStr, item.Response)
		}
	}
}
