package main

import (
	"app/cmd/handlers"
	"log"

	"github.com/gin-gonic/gin"
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
	rt := gin.New()
	// -> middleware
	rt.Use(gin.Recovery())
	rt.Use(gin.Logger())
	// -> routes
	// -> -> products group
	pr := rt.Group("/products")
	{
		// get all products
		pr.GET("/", ct.Get())
		// get product by id
		pr.GET("/:id", ct.GetByID())
		// search product by id (query params)
		pr.GET("/search", ct.Search())
		// create product
		pr.POST("/", ct.Create())
	}

	// -> run
	if err := rt.Run(":8080"); err != nil {
		log.Println(err)
		return
	}
}