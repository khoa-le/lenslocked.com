package controllers

import (
	"fmt"
	"github.com/gorilla/mux"
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
	"net/http"
	"strconv"
)

const (
	RouteShowGallery   = "show_gallery"
	RouteEditGallery   = "edit_gallery"
	RouteUpdateGallery = "update_gallery"
)

func NewGallery(gs models.GalleryService, r *mux.Router) *Gallery {
	return &Gallery{
		New:       views.NewView("bootstrap", "gallery/new"),
		IndexView: views.NewView("bootstrap", "gallery/index"),
		ShowView:  views.NewView("bootstrap", "gallery/show"),
		EditView:  views.NewView("bootstrap", "gallery/edit"),
		gs:        gs,
		r:         r,
	}
}

type Gallery struct {
	New       *views.View
	IndexView *views.View
	ShowView  *views.View
	EditView  *views.View
	gs        models.GalleryService
	r         *mux.Router
}

type GalleryForm struct {
	Title string `schema."title"`
}

func (g *Gallery) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, vd)

}

//GET /gallery/:id
func (g *Gallery) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

//GET /gallery/:id/edit
func (g *Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	g.EditView.Render(w, vd)
}

//POST /gallery/:id/update
func (g *Gallery) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	gallery.Title = form.Title
	if err := g.gs.Update(gallery); err != nil {
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}

	vd.Alert = &views.Alert{
		Level:   views.AlertLevelSuccess,
		Message: "Gallery successfully updated!",
	}
	g.EditView.Render(w, vd)
}

//POST /gallery/:id/delete
func (g *Gallery) Delete(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if user.ID != gallery.UserID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	if err := g.gs.Delete(gallery.ID); err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(w, vd)
		return
	}
	http.Redirect(w, r, "/gallery/index", http.StatusFound)
}

// POST /gallery/
func (g *Gallery) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}

	url, err := g.r.Get(RouteEditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

func (g *Gallery) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}

	gallery, err := g.gs.ByID(uint(id))

	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong", http.StatusInternalServerError)
		}
		return nil, err
	}

	return gallery, nil
}
