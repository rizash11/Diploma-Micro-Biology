package backend

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Keeps all of the handlers of an app.
func (app *Application) Routes() {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.Handle("/frontend/", http.StripPrefix("/frontend", http.FileServer(http.Dir("./frontend/"))))

	app.SrvMux = router
}

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path != "/":
		fmt.Println("Not found")
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	newInstance := &TemplateDataStruct{
		Links: make(map[string]string),
	}
	newInstance.Key = app.GenerateRandomString(5, w)
	newInstance.Links["CultivationKz"] = "/" + newInstance.Key + "/cultivation/kz"
	newInstance.Links["CultivationRu"] = "/" + newInstance.Key + "/cultivation/ru"
	newInstance.Links["CultivationEn"] = "/" + newInstance.Key + "/cultivation/en"
	newInstance.Links["TestsKz"] = "/" + newInstance.Key + "/tests/kz"
	newInstance.Links["TestsRu"] = "/" + newInstance.Key + "/tests/ru"
	newInstance.Links["TestsEn"] = "/" + newInstance.Key + "/tests/en"
	app.ReqInstances[newInstance.Key] = newInstance

	uniqueCultivation := "/" + newInstance.Key + "/cultivation/"
	uniqueTests := "/" + newInstance.Key + "/tests/"

	r.URL.Path = uniqueCultivation + "kz"
	app.SrvMux.HandleFunc(uniqueCultivation, app.cultivation)
	app.SrvMux.HandleFunc(uniqueTests, app.tests)
	app.SrvMux.ServeHTTP(w, r)
}

func (app *Application) cultivation(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))
	myInstance := app.ReqInstances[split[1]]

	switch {
	case split[2] != "cultivation":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if len(myInstance.Fname) == 0 || len(myInstance.Lname) == 0 {
		myInstance.Fname = r.FormValue("fname")
		myInstance.Lname = r.FormValue("lname")
	}

	app.Render(w, r, "cultivation.page.html", myInstance)
}

func (app *Application) tests(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))
	myInstance := app.ReqInstances[split[1]]

	switch {
	case split[2] != "tests":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	if len(myInstance.Fname) > 0 && len(myInstance.Lname) > 0 {
		for i := 0; i < 6; i++ {
			if len(myInstance.Results[i]) == 0 {
				myInstance.Results[i] = r.FormValue("q" + strconv.Itoa(i+1))
			}
		}
	}

	app.Render(w, r, "tests.page.html", myInstance)
}
