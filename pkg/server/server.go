package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World\n")
}

func FileListHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	content, err := os.ReadFile("files/files.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(w, "%s\n", string(content))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Log file requests
		logrus.Infof("Request made to '%v'", r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func New(dir string) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler)

	// This will serve files under http://localhost:8000/files/<filename>
	logrus.Infof("Serving files from %s\n", dir)
	// fsys := dotFileHidingFileSystem{http.Dir(dir)}
	r.HandleFunc("/files", FileListHandler)
	r.PathPrefix("/files/").Handler(http.StripPrefix("/files/", http.FileServer(http.Dir(dir))))
	r.Use(loggingMiddleware)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logrus.Infof("Created FileServer")
	return srv
}
