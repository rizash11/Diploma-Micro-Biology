package backend

import (
	"net/http"
)

// Keeps all of the handlers of the app.
func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/cultivation", app.cultivation)
	router.HandleFunc("/tests", app.tests)
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

	r.URL.Path = "/cultivation" // redirecting to the cultivation page
	app.Routes().ServeHTTP(w, r)
}

func (app *Application) cultivation(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/cultivation":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.Render(w, r, "cultivation.page.html", nil)
}

func (app *Application) tests(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/tests":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.Render(w, r, "tests.page.html", nil)
}
