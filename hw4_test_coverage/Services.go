package main

import "sort"

func sortUsers(users []User, order string) {
	if order == "id" {
		sort.Slice(users, func(i, j int) bool {
			return users[i].Id < users[j].Id
		})
	}
	if order == "Age" {
		sort.Slice(users, func(i, j int) bool {
			return users[i].Age < users[j].Age
		})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Name < users[j].Name
	})
}
