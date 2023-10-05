package backend

import "net/http"

// Keeps all of the handlers of the app.
func (app *Application) Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", app.Home)

	return router
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	app.Render(w, r, "base.layout", nil)
}
