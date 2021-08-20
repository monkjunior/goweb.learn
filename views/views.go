package views

import (
	"bytes"
	"errors"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/monkjunior/goweb.learn/context"
)

var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	files = addTemplatePath(files)
	files = addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", errors.New("csrfField is not implemented yet")
			},
		}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view with predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch data.(type) {
	case Data:
		vd = data.(Data)
	default:
		vd = Data{
			Yield: data,
		}
	}
	if alert := getAlert(r); alert != nil {
		vd.Alert = alert
		clearAlert(w)
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	if err := tpl.ExecuteTemplate(&buf, v.Layout, vd); err != nil {
		log.Println(err)
		http.Error(w, AlertMsgGeneric, http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// layoutFiles return a slice of strings representing
// the layout files used in our application
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

// addTemplatePath takes in a slice of strings
// representing file paths for templates, and it prepends
// the TemplateDir directory to each string in the slice
func addTemplatePath(files []string) []string {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
	return files
}

// addTemplateExt takes in a slice of strings
// representing file paths for templates, and it appends
// the TemplateExt extension to each string in the slice
func addTemplateExt(files []string) []string {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
	return files
}
