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
		LoginView: views.NewView("bootstrap", "user/login"),
		us: us,
	}
}
type User struct{
	NewView *views.View
	LoginView *views.View
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
		Password: signupForm.Password,
	}
	if err := u.us.Create(&user); err!=nil {
		http.Error(w, err.Error(),http.StatusInternalServerError)
	}
	signIn(w, &user)
}

// GET /login
func (u *User) Login(w http.ResponseWriter, r *http.Request){
	if err := u.LoginView.Render(w,nil); err !=nil{
		panic(err)
	}
}

type LoginForm struct{
	Email string `schema:"email"`
	Password string `schema:"password"`
}

// POST /login
func (u *User) DoLogin(w http.ResponseWriter, r *http.Request){
	var loginForm LoginForm
	if err:= parseForm( r, &loginForm); err != nil{
		panic(err)
	}
	user,err := u.us.Authenticate(loginForm.Email, loginForm.Password)
	if err !=nil{
		switch err{
		case models.ErrNotFound:
			fmt.Fprintln(w,"Invalid Email Addess")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w,"Invalid password provided")
		default:
			http.Error(w, err.Error(),http.StatusInternalServerError)
		}
		return
	}
	signIn(w, user)
}

//signIn is used to sign the given user  in via  cookies
func signIn(w http.ResponseWriter,  user *models.User){
	cookie := http.Cookie{
		Name: "email",
		Value: user.Email,
	}
	http.SetCookie(w, &cookie)
}