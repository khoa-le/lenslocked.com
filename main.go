package main

import (
	"fmt"
	"net/http"
	"lenslocked.com/views"

	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := homeView.Template.ExecuteTemplate(w,homeView.Layout, nil)
	if err != nil {
		panic(err)
	}

}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := contactView.Template.ExecuteTemplate(w,homeView.Layout, nil)
	if err != nil {
		panic(err)
	}
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
	homeView = views.NewView("boostrap","views/home.gohtml")
	contactView = views.NewView("boostrap","views/contact.gohtml")

	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(pageNotFound)
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	http.ListenAndServe(":3000", r)
}
