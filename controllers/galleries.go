package controllers

import (
	"fmt"
	"net/http"

	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/views"
)

func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		NewView: views.NewView("bootstrap", "galleries/new"),
		gs:      gs,
	}
}

type Galleries struct {
	NewView *views.View
	gs      models.GalleryService
}

// This is used to render the form where a user can create
// a new gallery
//
// GET /galleries/new
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, nil)
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// This is used to process gallery form when a user tries to
// create a new gallery
//
// POST /galleries/new
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}

	gallery := models.Gallery{
		Title: form.Title,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}
