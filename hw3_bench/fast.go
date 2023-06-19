package main

import (
	"io"
)

type User struct {
	browsers []string `json:"browsers"`
	company  string   `json:"company"`
	country  string   `json:"country"`
	email    string   `json:"email"`
	job      string   `json:"job"`
	name     string   `json:"name"`
	phone    string   `json:"phone"`
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	SlowSearch(out)
}
