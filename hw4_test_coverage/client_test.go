package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
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
	var users User

	xmlData, err := ioutil.ReadFile("dataset.xml")
	if err != nil {
		http.Error(w, "Данные не прочитаны", 500)
		return
	}

	err = xml.Unmarshal(xmlData, users)
	if err != nil {
		http.Error(w, "Данные не переведены", 500)
		return
	}
}
