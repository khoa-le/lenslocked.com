package controllers

import (
	"lenslocked.com/context"
	"lenslocked.com/models"
	"lenslocked.com/views"
	"net/http"
)

func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		New:   views.NewView("bootstrap", "gallery/new"),
		gs:        gs,
	}
}

type Gallery struct {
	New   *views.View
	gs        models.GalleryService
}

type GalleryForm struct{
	Title string `schema."title"`
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
	if user ==nil{
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	gallery := models.Gallery{
		Title:     form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
}