package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// Parsing HTML templates and storing them in a cache to be executed later
func (app *Application) ParseTemplates(dir string) (map[string]*template.Template, error) {
	// templateCache := map[string]*template.Template{}
	// var err error

	// templateCache["base.layout"], err = template.ParseFiles(filepath.Join(dir, "base.layout.html"))
	// if err != nil {
	// 	app.Error_log.Println("Error parsing base.layout: " + err.Error())
	// }

	// app.templateCache = templateCache
	templateCache := map[string]*template.Template{}
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*partial.html"))
		if err != nil {
			return nil, err
		}

		name := filepath.Base(page)
		templateCache[name] = ts
	}

	return templateCache, nil
}

// Receiving texts from json files
func (app *Application) ParseTemplateData(dir string) error {
	// Unmarshalling Kz texts
	dataKz, err := os.ReadFile(filepath.Join(dir, "kz.json"))
	if err != nil {
		return err
	}

	tdKz := make(map[string]map[string]string)
	err = json.Unmarshal(dataKz, &tdKz)
	if err != nil {
		return err
	}
	app.TemplateTextKz = tdKz

	// Unmarshalling Ru texts
	dataRu, err := os.ReadFile(filepath.Join(dir, "ru.json"))
	if err != nil {
		return err
	}

	tdRu := make(map[string]map[string]string)
	err = json.Unmarshal(dataRu, &tdRu)
	if err != nil {
		return err
	}
	app.TemplateTextRu = tdRu

	// Unmarshalling En texts
	dataEn, err := os.ReadFile(filepath.Join(dir, "en.json"))
	if err != nil {
		return err
	}

	tdEn := make(map[string]map[string]string)
	err = json.Unmarshal(dataEn, &tdEn)
	if err != nil {
		return err
	}
	app.TemplateTextEn = tdEn

	return err
}

// This function finds an html template in cache and executes it.
// Enter filename without the html extension as in "base.layout.html" to just "base.layout"
func (app *Application) Render(w http.ResponseWriter, r *http.Request, name string) {
	tmpl, ok := app.TemplateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("the %s page doesn't exist", name))
		return
	}

	split := strings.Split(r.URL.Path, string(os.PathSeparator))
	var err error

	switch {
	case split[2] == "kz":
		app.TemplateData.TemplateText = app.TemplateTextKz
	case split[2] == "ru":
		app.TemplateData.TemplateText = app.TemplateTextRu
	case split[2] == "en":
		app.TemplateData.TemplateText = app.TemplateTextEn
	default:
		err = errors.New(split[0] + "/" + split[1] + "/" + split[2] + " - requested language is not found")
		app.ServerError(w, err)
		return
	}

	err = tmpl.Execute(w, app.TemplateData)

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
