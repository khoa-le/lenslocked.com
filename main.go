package main

import (
	"fmt"
	"net/http"
	"lenslocked.com/controllers"
	"lenslocked.com/models"
	"github.com/gorilla/mux"
)
const (
	host = "localhost"
	port = "5432"
	user = "khoa"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
	host, port, user, dbname)
	services,err := models.NewServices(psqlInfo)
	if err !=nil{
		panic(err)
	}

	// TODO: Fix it
	//defer userService.Close()
	//userService.DestructiveReset()

	staticController := controllers.NewStatic()
	userController := controllers.NewUser(services.User)

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("GET")
	r.HandleFunc("/login", userController.DoLogin).Methods("POST")
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
