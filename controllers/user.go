package controllers

import (
	"fmt"
	"net/http"
	"github.com/gorilla/schema"
	"lenslocked.com/views"
)
func NewUser() *User{
	return &User{
		NewView: views.NewView("bootstrap", "views/user/new.gohtml"),
	}
}
type User struct{
	NewView *views.View
}

// GET /signup
func (u *User) New(w http.ResponseWriter, r *http.Request){
	if err := u.NewView.Render(w,nil); err !=nil{
		panic(err)
	}
}

type SignupForm struct{
	Email string `schema:"email"`
	Password string `schema:"password"`
}
// POST /signup
func (u *User) Create(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err !=nil	{
		panic(err)
	}
	var signupForm SignupForm
	decoder := schema.NewDecoder()
	if err := decoder.Decode(&signupForm,r.PostForm); err !=nil{
		panic(err)
	}
	fmt.Fprintln(w, signupForm)
	fmt.Fprintln(w, r.PostForm["email"])
	fmt.Fprintln(w, r.PostFormValue( "email"))
	// fmt.Println(w,"This is tempolary response")
}