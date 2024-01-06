// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"diploma/backend"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	errorLog := log.New(os.Stderr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)
	TemplateData := backend.TemplateDataStruct{}

	app := backend.Application{
		Info_log:     infoLog,
		Error_log:    errorLog,
		TemplateData: &TemplateData,
	}

	// Parsing HTML templates and storing them in a cache to be executed later
	var err error
	app.TemplateCache, err = app.ParseTemplates(filepath.Join("frontend", "html"))
	if err != nil {
		app.Error_log.Println("Error parsing html templates: " + err.Error())
	}

	// Getting the texts that will be passed to html templates
	err = app.ParseTemplateData(filepath.Join("backend", "template_data"))
	if err != nil {
		app.Error_log.Println("Error parsing templates' data: " + err.Error())
	}

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		app.Info_log.Printf("defaulting to port %s", port)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.Routes(),
	}

	app.Info_log.Println("Starting a server at http://localhost:" + port)
	app.Error_log.Fatalln(srv.ListenAndServe())

}
