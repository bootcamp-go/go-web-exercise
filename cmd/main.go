package main

import (
	"app/cmd/handlers"
	"app/internal/product/storage"
	"app/internal/product/storage/loader"
	"app/internal/product/validator"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// env

	
	// dependencies
	ld := loader.NewLoaderJSON("./docs/db/json/products.json")
	db, err := ld.Load()
	if err != nil {
		log.Println(err)
		return
	}
	st := storage.NewStorageProductMap(db.Db, db.LastId)
	vl := validator.NewValidatorProductDefault("")
	stVl := storage.NewStorageProductValidate(st, vl)
	ct := handlers.NewHandlerProducts(stVl)

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