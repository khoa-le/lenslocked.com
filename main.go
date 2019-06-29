package main

import (
	"fmt"
	"net/http"
	"lenslocked.com/views"
	"lenslocked.com/controllers"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
	signupView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	homeView.Render(w,nil)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	contactView.Render(w,nil)
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>My Faq page</h1>")
}
func pageNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadGateway)
	fmt.Fprint(w, "<h1>We could not find the page you are looking for.</h1>")
}

func main() {
	homeView = views.NewView("bootstrap","views/home.gohtml")
	contactView = views.NewView("bootstrap","views/contact.gohtml")
	userController := controllers.NewUser()

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(pageNotFound)
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/contact", contact).Methods("GET")
	r.HandleFunc("/faq", faq).Methods("GET")
	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
