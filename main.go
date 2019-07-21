package main

import (
	"fmt"
	"lenslocked.com/middleware"
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

	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	staticController := controllers.NewStatic()
	userController := controllers.NewUser(services.User)
	galleryController := controllers.NewGallery(services.Gallery)
	
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	r := mux.NewRouter()
	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("GET")
	r.HandleFunc("/login", userController.DoLogin).Methods("POST")

	//Gallery routes
	r.Handle("/gallery/new",requireUserMw.Apply(galleryController.New)).Methods("GET")
	r.HandleFunc("/gallery",requireUserMw.ApplyFn(galleryController.Create)).Methods("POST")
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
