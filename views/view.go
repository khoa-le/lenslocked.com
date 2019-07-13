package views

import (
	"html/template"
	"path/filepath"
	"net/http"
)

var (
	LayoutDir string = "views/layout/"
	TemplateDir string = "views/"
	TemplateExt string = ".gohtml"
)

func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
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

//Render is used to render the view with  the predefined layout
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type){
	case Data:
		//do nothing
	default:
		data = Data{
			Yield: data,
		}
	}


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

//
func addTemplatePath(files[] string){
	for i,f := range files{
		files[i] = TemplateDir + f
	}
}
func addTemplateExt(files[] string){
	for i,f := range files{
		files[i] = f + TemplateExt
	}
}