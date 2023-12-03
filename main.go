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

	app := backend.Application{
		Info_log:  infoLog,
		Error_log: errorLog,
	}

	// Parsing HTML templates and storing them in a cache to be executed later
	var err error
	app.TemplateCache, err = app.ParseTemplates(filepath.Join("frontend", "html"))
	if err != nil {
		app.Error_log.Println("Error parsing html templates: " + err.Error())
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

	// log.Print("starting server...")
	// http.HandleFunc("/", handler)

	// // Determine port for HTTP service.
	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// 	app.Info_log.Printf("defaulting to port %s", port)
	// }

	// // Start HTTP server.
	// app.Info_log.Printf("listening on port %s", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	app.Info_log.Fatal(err)
	// }
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	name := os.Getenv("NAME")
// 	if name == "" {
// 		name = "World"
// 	}
// 	fmt.Fprintf(w, "Hello %s!\n", name)
// }
