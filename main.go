package main

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"lenslocked.com/controllers"
	"lenslocked.com/middleware"
	"lenslocked.com/models"
	"lenslocked.com/rand"
	"net/http"
)

const (
	host   = "localhost"
	port   = "5432"
	user   = "khoa"
	dbname = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}

	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	r := mux.NewRouter()
	staticController := controllers.NewStatic()
	userController := controllers.NewUser(services.User)
	galleryController := controllers.NewGallery(services.Gallery, services.Image, r)

	randomBytes, err := rand.Bytes(32)
	if err != nil{
		panic(err)
	}
	isProd := false
	csrfMw := csrf.Protect(randomBytes, csrf.Secure(isProd))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	r.Handle("/", staticController.Home).Methods("GET")
	r.Handle("/contact", staticController.Contact).Methods("GET")
	r.HandleFunc("/signup", userController.New).Methods("GET")
	r.HandleFunc("/signup", userController.Create).Methods("POST")
	r.HandleFunc("/login", userController.Login).Methods("GET")
	r.HandleFunc("/login", userController.DoLogin).Methods("POST")

	assetsHandler := http.FileServer(http.Dir("./assets/"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assetsHandler))

	//Images route
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	//Gallery routes
	r.HandleFunc("/gallery/index", requireUserMw.ApplyFn(galleryController.Index)).Methods("GET")
	r.Handle("/gallery/new", requireUserMw.Apply(galleryController.New)).Methods("GET")
	r.HandleFunc("/gallery", requireUserMw.ApplyFn(galleryController.Create)).Methods("POST")
	r.HandleFunc("/gallery/{id:[0-9]+}", galleryController.Show).Methods("GET").Name(controllers.RouteShowGallery)
	r.HandleFunc("/gallery/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleryController.Edit)).Methods("GET").Name(controllers.RouteEditGallery)
	r.HandleFunc("/gallery/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleryController.Update)).Methods("POST").Name(controllers.RouteUpdateGallery)
	r.HandleFunc("/gallery/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleryController.ImageUpload)).Methods("POST")
	r.HandleFunc("/gallery/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleryController.ImageDelete)).Methods("POST")

	r.HandleFunc("/gallery/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleryController.Delete)).Methods("POST")

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000",csrfMw(userMw.Apply(r)))
}
