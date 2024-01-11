package backend

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Keeps all of the handlers of the app.
func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/cultivation/", app.cultivation)
	router.HandleFunc("/tests/", app.tests)
	router.Handle("/frontend/", http.StripPrefix("/frontend", http.FileServer(http.Dir("./frontend/"))))

	return router
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	r.URL.Path = "/cultivation/kz" // redirecting to the cultivation page
	app.Routes().ServeHTTP(w, r)
}

func (app *Application) cultivation(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))

	switch {
	case split[1] != "cultivation":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if len(app.TemplateData.Fname) == 0 || len(app.TemplateData.Lname) == 0 {
		app.TemplateData.Fname = r.FormValue("fname")
		app.TemplateData.Lname = r.FormValue("lname")
	}

	app.Render(w, r, "cultivation.page.html")
}

func (app *Application) tests(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))

	switch {
	case split[1] != "tests":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if len(app.TemplateData.Fname) > 0 && len(app.TemplateData.Lname) > 0 {
		for i := 0; i < 6; i++ {
			if len(app.TemplateData.Results[i]) == 0 {
				app.TemplateData.Results[i] = r.FormValue("q" + strconv.Itoa(i+1))
			}
		}
	}

	fmt.Println(app.TemplateData.Results)

	app.Render(w, r, "tests.page.html")
}
