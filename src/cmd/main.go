package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"

	"github.com/thurasw/ProxAuth/src/internal/api"
	"github.com/thurasw/ProxAuth/src/internal/config"
	"github.com/thurasw/ProxAuth/src/internal/db"
)

func main() {
	// Load config
	err := config.Load()
	if err != nil {
		log.Fatalf("error loading config\n%v\n", err)
	}
	config := config.Config

	// Initialize database
	err = db.Init(config.DbPath)
	if err != nil {
		log.Fatalf("error init database\n%v\n", err)
	}

	// Create base chi router
	r := chi.NewRouter()

	// Serve API for /api routes
	r.Mount("/api", api.Router())
	// Serve website for all other routes
	r.HandleFunc("/*", spaHandler(config.WebRootPath, "/index.html"))

	// Run server
	log.Printf("Listening on port: %d...", config.Port)

	addr := fmt.Sprintf(":%d", config.Port)
	log.Fatalln(http.ListenAndServe(addr, r))
}

func spaHandler(webRoot string, indexPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// join internally calls path.Clean to prevent directory traversal
		path := filepath.Join(webRoot, r.URL.Path)

		fi, err := os.Stat(path)
		// file does not exist or path is a directory, serve index.html
		if os.IsNotExist(err) || fi.IsDir() {
			http.ServeFile(w, r, filepath.Join(webRoot, indexPath))
			return
		}

		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		if err != nil {
			log.Printf("error serving web: %v\n", err.Error())
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// otherwise, use http.FileServer to serve the static file
		http.FileServer(http.Dir(webRoot)).ServeHTTP(w, r)
	}
}
