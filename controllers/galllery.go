package controllers

import (
	"lenslocked.com/models"
	"lenslocked.com/views"
	"net/http"
)

func NewGallery(gs models.GalleryService) *Gallery {
	return &Gallery{
		NewView:   views.NewView("bootstrap", "gallery/new"),
		gs:        gs,
	}
}

type Gallery struct {
	NewView   *views.View
	gs        models.GalleryService
}

func (g *Gallery) New(w http.ResponseWriter,r *http.Request){
	g.NewView.Render(w, nil)
}