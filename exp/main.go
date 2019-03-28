package main

import (
	"html/template"
	"os"
)

type User struct {
	Name   string
	Family []string
}

func main() {
	t, err := template.ParseFiles("index.gohtml")

	if err != nil {
		panic(err)
	}

	data := make(map[string]interface{})
	data["User"] = User{Name: "Khoa Le", Family: []string{
		"Hong An",
		"Khanh Minh",
	}}
	data["Dog"] = "Morty"

	err = t.Execute(os.Stdout, data)

	if err != nil {
		panic(err)
	}
}
