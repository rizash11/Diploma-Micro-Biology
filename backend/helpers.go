package backend

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime/debug"
)

func (app *Application) ParseTemplates(dir string) {
	templateCache := map[string]*template.Template{}
	var err error

	templateCache["base.layout"], err = template.ParseFiles(filepath.Join(dir, "base.layout.html"))
	if err != nil {
		app.Error_log.Println("Error parsing base.layout: " + err.Error())
	}

	app.templateCache = templateCache
}

// This function finds an html template in cache and executes it.
// Enter filename without the html extension as in "base.layout.html" to just "base.layout"
func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string, td *TemplateData) {
	tmpl, ok := app.templateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("the %s page doesn't exist", name))
		return
	}

	err := tmpl.Execute(w, td)
	if err != nil {
		app.ServerError(w, err)
		return
	}
}

func (app *Application) NotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *Application) ClientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.Error_log.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
