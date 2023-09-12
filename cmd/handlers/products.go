package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// NewControllerProducts returns a new ControllerProducts
func NewControllerProducts(storage map[int]*Product, lastId int) *ControllerProducts {
	return &ControllerProducts{storage: storage, lastId: lastId}
}

// ControllerProducts is a struct that contains the storage of products
type ControllerProducts struct {
	storage map[int]*Product
	lastId  int
}

// Get is a method that returns all products
type ProductJSON struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  string	  `json:"expiration"`
	Price       float64   `json:"price"`
}
type ResponseBodyProductGet struct {
	Message string		   `json:"message"`
	Data    []*ProductJSON `json:"data"`
}
func (c *ControllerProducts) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// ...

		// process
		// ...

		// response
		code := http.StatusOK
		body := ResponseBodyProductGet{Message: "products", Data: make([]*ProductJSON, len(c.storage))}
		for k, v := range c.storage {
			body.Data[k] = &ProductJSON{Id: k, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price}
		}

		ctx.JSON(code, body)
	}
}

// GetByID is a method that returns a product by id
type ResponseBodyGetByID struct {
	Message string		 `json:"message"`
	Data    *ProductJSON `json:"data"`
}
func (c *ControllerProducts) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		// -> id
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyGetByID{Message: "invalid id", Data: nil}

			ctx.JSON(code, body)
			return
		}

		// process
		// -> get product
		pr, ok := c.storage[id]
		if !ok {
			code := http.StatusNotFound
			body := ResponseBodyGetByID{Message: "product not found", Data: nil}

			ctx.JSON(code, body)
			return
		}

		// response
		code := http.StatusOK
		body := ResponseBodyGetByID{Message: "product", Data: &ProductJSON{Id: id, Name: pr.Name, Quantity: pr.Quantity, CodeValue: pr.CodeValue, IsPublished: pr.IsPublished, Expiration: pr.Expiration.Format("2006-01-02"), Price: pr.Price}}

		ctx.JSON(code, body)
	}
}

// Search is a method that returns a product by id (via query params)
type ResponseBodySearch struct {
	Message string		   `json:"message"`
	Data    []*ProductJSON `json:"data"`
}
func (c *ControllerProducts) Search() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		id, _ := strconv.Atoi(ctx.Query("id"))

		// process
		filtered := make(map[int]*Product)
		for k, v := range c.storage {
			// filter: check if query is set
			if id > 0 {
				if k != id {
					continue
				}
			}

			// default: add to filtered
			filtered[k] = v
		}

		// response
		code := http.StatusOK
		body := ResponseBodySearch{Message: "products", Data: make([]*ProductJSON, len(filtered))}
		for k, v := range filtered {
			body.Data[k] = &ProductJSON{Id: k, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price}
		}

		ctx.JSON(code, body)
	}
}

// Create is a method that creates a new product
type RequestBodyProductCreate struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}
type ResponseBodyProductCreate struct {
	Message string		 `json:"message"`
	Data    *ProductJSON `json:"data"`
}
func (c *ControllerProducts) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// request
		var reqBody RequestBodyProductCreate
		if err := ctx.ShouldBindJSON(&reqBody); err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid request body", Data: nil}

			ctx.JSON(code, body)
			return
		}

		// process
		// -> deserialize
		exp, err := time.Parse("2006-01-02", reqBody.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid date format. Must be yyyy-mm-dd", Data: nil}

			ctx.JSON(code, body)
			return
		}
		pr := &Product{
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  exp,
			Price:       reqBody.Price,
		}
		// -> validate
		if err := Validate(pr); err != nil {
			code := http.StatusConflict
			body := ResponseBodyProductCreate{Message: "invalid product", Data: nil}

			ctx.JSON(code, body)
			return
		}
		// -> save
		c.lastId++
		c.storage[c.lastId] = pr

		// response
		code := http.StatusCreated
		body := ResponseBodyProductCreate{
			Message: "product created",
			Data:    &ProductJSON{Id: c.lastId, Name: pr.Name, Quantity: pr.Quantity, CodeValue: pr.CodeValue, IsPublished: pr.IsPublished, Expiration: pr.Expiration.Format("2006-01-02"), Price: pr.Price},
		}

		ctx.JSON(code, body)
	}
}