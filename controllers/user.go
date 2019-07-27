package controllers

import (
	"lenslocked.com/models"
	"lenslocked.com/rand"
	"lenslocked.com/views"
	"log"
	"net/http"
)

func NewUser(us models.UserService) *User {
	return &User{
		NewView:   views.NewView("bootstrap", "user/new"),
		LoginView: views.NewView("bootstrap", "user/login"),
		us:        us,
	}
}

type User struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// GET /signup
func (u *User) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, r, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// POST /signup
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var signupForm SignupForm
	if err := parseForm(r, &signupForm); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	user := models.User{
		Name:     signupForm.Name,
		Email:    signupForm.Email,
		Password: signupForm.Password,
	}
	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(w, r, vd)
		return
	}
	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w,r , "/gallery/index", http.StatusFound)
}

// GET /login
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, r, nil)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// POST /login
func (u *User) DoLogin(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	var loginForm LoginForm
	if err := parseForm(r, &loginForm); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	user, err := u.us.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.AlertError("Invalid Email Address")
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(w, r, vd)
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(w, r, vd)
		return
	}
	http.Redirect(w,r , "/gallery/index", http.StatusFound)
}

//signIn is used to sign the given user  in via  cookies
func (u *User) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	return nil
}
