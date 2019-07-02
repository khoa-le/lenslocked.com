package controllers

import (
	"fmt"
	"net/http"
	"lenslocked.com/views"
	"lenslocked.com/models"
)
func NewUser(us *models.UserService) *User{
	return &User{
		NewView: views.NewView("bootstrap", "user/new"),
		us: us,
	}
}
type User struct{
	NewView *views.View
	us *models.UserService
}

// GET /signup
func (u *User) New(w http.ResponseWriter, r *http.Request){
	if err := u.NewView.Render(w,nil); err !=nil{
		panic(err)
	}
}

type SignupForm struct{
	Name string `schema:"name"`
	Email string `schema:"email"`
	Password string `schema:"password"`
}
// POST /signup
func (u *User) Create(w http.ResponseWriter, r *http.Request){
	
	var signupForm SignupForm
	if err:= parseForm(r,&signupForm); err !=nil{
		panic(err)
	}
	user := models.User{
		Name:signupForm.Name,
		Email: signupForm.Email,
	}
	if err := u.us.Create(&user); err!=nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
	}
	fmt.Fprintln(w, signupForm)
}