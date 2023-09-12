package main

import (
	"app/cmd/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// env
	// ...

	// dependencies
	db, err := handlers.LoaderProducts("./docs/db/json/products.json")
	if err != nil {
		log.Println(err)
		return
	}
	ct := handlers.NewControllerProducts(db, len(db))

	// server
	rt := chi.NewRouter()
	// -> middleware
	rt.Use(middleware.Recoverer)
	rt.Use(middleware.Logger)
	// -> routes
	// -> -> products group
	rt.Route("/products", func(rt chi.Router) {
		// get all products
		rt.Get("/", ct.Get())
		// get product by id
		rt.Get("/{id}", ct.GetByID())
		// search products
		rt.Get("/search", ct.Search())
		// create product
		rt.Post("/", ct.Create())
	})
		

	// -> run
	if err := http.ListenAndServe(":8080", rt); err != nil {
		log.Println(err)
		return
	}
}