package backend

import (
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Keeps all of the handlers of an app.
func (app *Application) Routes() {
	router := http.NewServeMux()
	router.HandleFunc("/", app.home)
	router.HandleFunc("/results/", app.results)
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

	// For every new request new handler is created, so that a single instance of this application could serve many requests at the same time.
	// This way, every instance of the requests is served separately from each other and their data do not mingle.
	newInstance := &TemplateDataStruct{
		Links: make(map[string]string),
	}

	// To identify different instances of handlers a key is assigned to each of them
	newInstance.Key = app.GenerateRandomString(5, w)

	// There are links that are used by html templates, they are initialized here
	newInstance.Links["CultivationKz"] = "/" + newInstance.Key + "/cultivation/kz"
	newInstance.Links["CultivationRu"] = "/" + newInstance.Key + "/cultivation/ru"
	newInstance.Links["CultivationEn"] = "/" + newInstance.Key + "/cultivation/en"
	newInstance.Links["TestsKz"] = "/" + newInstance.Key + "/tests/kz"
	newInstance.Links["TestsRu"] = "/" + newInstance.Key + "/tests/ru"
	newInstance.Links["TestsEn"] = "/" + newInstance.Key + "/tests/en"

	// every instance and its data is stored in a map within the application structure
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

	if len(myInstance.Results[5]) != 0 {
		f, err := os.OpenFile("./backend/results.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			app.ServerError(w, err)
			return
		}
		var ResStr string
		defer f.Close()

		dt := time.Now()
		ResStr = ResStr + "\n" + dt.Format("01-02-2006 15:04:05") + "\n"
		ResStr = ResStr + "First name: " + myInstance.Fname + "\n" +
			"Last name: " + myInstance.Lname + "\n"

		for i := 0; i < 6; i++ {
			ResStr = ResStr + "Question #" + strconv.Itoa(i+1) + ": " + myInstance.Results[i] + "\n"
		}
		ResStr = ResStr + "\n"

		_, err = fmt.Fprintln(f, ResStr)
		if err != nil {
			app.ServerError(w, err)
			return
		}
	}

	app.Render(w, r, "tests.page.html", myInstance)
}

func (app *Application) results(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, string(os.PathSeparator))

	myInstance := &TemplateDataStruct{
		Key: "results",
	}

	switch {
	case split[1] != "results":
		app.NotFound(w)
		return

	case r.Method != http.MethodGet:
		app.ClientError(w, http.StatusMethodNotAllowed)
		return
	}

	b, err := os.ReadFile("./backend/results.txt")
	myInstance.ResultsTxt = string(b)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	savable := r.FormValue("savable")
	if savable == "true" {
		newRes := r.FormValue("results-window")
		err := os.WriteFile("./backend/results.txt", []byte(newRes), fs.FileMode(os.O_TRUNC))
		if err != nil {
			app.ServerError(w, err)
			return
		}
		myInstance.ResultsTxt = newRes
	}

	app.Render(w, r, "results.page.html", myInstance)
}
