package main

import (
	"app/cmd/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// env
	// ...

	// dependencies
	db := make(map[int]*handlers.Product)
	ct := handlers.NewControllerProducts(db)

	// server
	rt := chi.NewRouter()
	// -> middleware
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// -> routes
	// -> -> products group
	rt.Route("/products", func(rt chi.Router) {
		// create
		rt.Post("/", ct.Create())
	})

	// -> run
	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}