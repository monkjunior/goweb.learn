package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/monkjunior/goweb.learn/context"
	"github.com/monkjunior/goweb.learn/models"
	"github.com/monkjunior/goweb.learn/views"
)

const (
	ShowGallery   = "show_gallery"
	UpdateGallery = "update_gallery"

	// maxMultipartMem = 1MB
	maxMultipartMem = 1 << 20
)

func NewGalleries(gs models.GalleryService, r mux.Router) *Galleries {
	return &Galleries{
		NewView:    views.NewView("bootstrap", "galleries/new"),
		ShowView:   views.NewView("bootstrap", "galleries/show"),
		UpdateView: views.NewView("bootstrap", "galleries/update"),
		IndexView:  views.NewView("bootstrap", "galleries/index"),
		gs:         gs,
		r:          r,
	}
}

type Galleries struct {
	NewView    *views.View
	ShowView   *views.View
	UpdateView *views.View
	IndexView  *views.View
	gs         models.GalleryService
	r          mux.Router
}

// This is used to render the form where a user can create
// a new gallery
//
// GET /galleries/new
func (g *Galleries) New(w http.ResponseWriter, r *http.Request) {
	g.NewView.Render(w, r, nil)
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// Index list all the gallery that user has access to.
//
// GET /galleries
func (g *Galleries) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(w, r, vd)
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
	g.ShowView.Render(w, r, vd)
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
	vd.User = user
	g.UpdateView.Render(w, r, vd)
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
		g.UpdateView.Render(w, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.UpdateView.Render(w, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:   views.AlertLvSuccess,
		Message: "Gallery successfully updated",
	}
	g.UpdateView.Render(w, r, vd)
}

// ImageUpload will upload our images the gallery
//
// POST /galleries/:id/images
func (g *Galleries) ImageUpload(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}

	// TODO: Parse a multi part form
	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartMem)
	if err != nil {
		vd.SetAlert(err)
		g.UpdateView.Render(w, r, vd)
		return
	}

	galleryPath := fmt.Sprintf("images/galleries/%v/", gallery.ID)
	err = os.MkdirAll(galleryPath, 0755)
	if err != nil {
		vd.SetAlert(err)
		g.UpdateView.Render(w, r, vd)
		return
	}

	for _, f := range r.MultipartForm.File["images"] {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.UpdateView.Render(w, r, vd)
			return
		}
		defer file.Close()

		dst, err := os.Create(galleryPath + f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.UpdateView.Render(w, r, vd)
			return
		}

		_, err = io.Copy(dst, file)
		if err != nil {
			vd.SetAlert(err)
			g.UpdateView.Render(w, r, vd)
			return
		}
	}

	fmt.Fprintln(w, "Files successfully uploaded")
}

// Delete will update the gallery edit page
//
// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
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
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.UpdateView.Render(w, r, vd)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
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
		g.NewView.Render(w, r, vd)
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
		g.NewView.Render(w, r, vd)
		return
	}
	url, err := g.r.Get(UpdateGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		http.Redirect(w, r, "/galleries", http.StatusFound)
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
