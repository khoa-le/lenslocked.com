package main

import (
	"net/http"
	"lenslocked.com/controllers"
	"github.com/gorilla/mux"
)

func main() {
	staticController := controllers.NewStatic()
	userController := controllers.NewUser()

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
