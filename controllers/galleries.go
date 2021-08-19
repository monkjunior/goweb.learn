package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/context"
	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/views"
)

const (
	ShowGallery = "show_gallery"
)

func NewGalleries(gs models.GalleryService, r mux.Router) *Galleries {
	return &Galleries{
		NewView:    views.NewView("bootstrap", "galleries/new"),
		ShowView:   views.NewView("bootstrap", "galleries/show"),
		UpdateView: views.NewView("bootstrap", "galleries/update"),
		gs:         gs,
		r:          r,
	}
}

type Galleries struct {
	NewView    *views.View
	ShowView   *views.View
	UpdateView *views.View
	gs         models.GalleryService
	r          mux.Router
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

// Show will look up and show the gallery with specific ID
//
// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(w, vd)
}

// GetUpdate will load the update gallery page
//
// GET /galleries/:id/update
func (g *Galleries) GetUpdate(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.UpdateView.Render(w, vd)
}

// PostUpdate will update the gallery edit page
//
// POST /galleries/:id/update
func (g *Galleries) PostUpdate(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	var form GalleryForm
	vd.Yield = gallery
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.UpdateView.Render(w, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.UpdateView.Render(w, r)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvSuccess,
		Message: "Gallery successfully updated",
	}
	g.UpdateView.Render(w, vd)
}

// Create is used to process gallery form when a user tries to
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
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.NewView.Render(w, vd)
		return
	}
	url, err := g.r.Get(ShowGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		// TODO: make this go to the index page
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.String(), http.StatusFound)
}

func (g *Galleries) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
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
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
