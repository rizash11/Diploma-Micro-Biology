// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"diploma/backend"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	errorLog := log.New(os.Stderr, "ERROR: \t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO \t", log.Ldate|log.Ltime)

	app := backend.Application{
		Info_log:  infoLog,
		Error_log: errorLog,
	}

	app.ParseTemplates("frontend")

	log.Print("starting server...")
	http.HandleFunc("/", handler)

	srv := &http.Server{
		Addr:    "localhost:" + "8080",
		Handler: app.Routes(),
	}

	app.Info_log.Println("Starting a server at http://localhost:" + "8080" + "/")
	app.Error_log.Fatalln(srv.ListenAndServe())
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

func handler(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}
