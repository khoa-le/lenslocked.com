package views

import (
	"bytes"
	"errors"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"lenslocked.com/context"
)

var (
	LayoutDir   string = "views/layout/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	t, err := template.New("").Funcs(template.FuncMap{
		csrf.TemplateTag: func() (template.HTML, error) {
			return "", errors.New("csrfField is not implemented")
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

//Render is used to render the view with  the predefined layout
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buff bytes.Buffer

	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		csrf.TemplateTag: func() template.HTML {
			return csrfField
		},
	})

	if err := tpl.ExecuteTemplate(&buff, v.Layout, vd); err != nil {
		http.Error(w, "Something went w  rong.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buff)
}

//layoutFiles return a slice of strings representing
// the layout files used in our application
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}

//
func addTemplatePath(files [] string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}
func addTemplateExt(files [] string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}
