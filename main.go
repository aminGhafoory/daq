package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aminGhafoory/daq/views/layouts"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	fmt.Println("hello")

	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		component := layouts.Show("hello", []string{})
		component.Render(context.Background(), w)
	})

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.FileServer(http.Dir("."))
		fs.ServeHTTP(w, r)
	})

	fmt.Println("server started on http://localhost:8000")
	http.ListenAndServe(":8000", r)

}
