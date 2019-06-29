package views

import (
	"html/template"
	"path/filepath"
	"net/http"
)

var (
	LayoutDir string = "views/layout/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	layoutFiles := layoutFiles()
	files = append(files, layoutFiles...)
	
	t,err := template.ParseFiles(files...)
	if err != nil{
		panic(err)
	}

	return &View{
		Template: t,
		Layout: layout,
	}
}


type View struct {
	Template *template.Template
	Layout string
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request){
	if err:=v.Render(w,nil); err!=nil{
		panic(err)
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

//layoutFiles return a slice of strings representing
// the layout files used in our application
func layoutFiles() []string{
	files, err := filepath.Glob(LayoutDir+ "*" + TemplateExt)
	if err != nil{
		panic(err)
	}
	return files
}