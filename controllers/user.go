package controllers

import (
	"fmt"
	"net/http"
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

// POST /signup
func (u *User) Create(w http.ResponseWriter, r *http.Request){
	fmt.Println(w,"This is tempolary response")
}