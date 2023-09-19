package main

import (
	"app/cmd/handlers"
	"app/cmd/middlewares"
	"app/internal/authenticator"
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
	// -> authenticator
	au := authenticator.NewAuthenticatorTokenBasic("token")
	md := middlewares.NewMiddlewareAuthenticator(au)
	
	// -> product
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
	rt.Use(middleware.Recoverer)	// recover from panics without crashing the server
	rt.Use(middleware.Logger)		// log api requests
	rt.Use(md.Auth)					// authenticate via token
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
		// update or create product
		rt.Put("/{id}", ct.UpdateOrCreate())
		// update product
		rt.Patch("/{id}", ct.Update())
		// delete product
		rt.Delete("/{id}", ct.Delete())
	})
		

	// -> run
	if err := http.ListenAndServe(":8080", rt); err != nil {
		log.Println(err)
		return
	}
}