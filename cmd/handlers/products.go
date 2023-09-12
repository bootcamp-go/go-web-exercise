package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
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
func (c *ControllerProducts) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		// ...

		// response
		code := http.StatusOK
		body := ResponseBodyProductGet{Message: "products", Data: make([]*ProductJSON, 0, len(c.storage))}
		for k, v := range c.storage {
			body.Data = append(body.Data, &ProductJSON{Id: k, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price})
		}

		w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

// GetByID is a method that returns a product by id
type ResponseBodyGetByID struct {
	Message string		 `json:"message"`
	Data    *ProductJSON `json:"data"`
}
func (c *ControllerProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// -> id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyGetByID{Message: "invalid id", Data: nil}

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
			return
		}

		// process
		// -> get product
		pr, ok := c.storage[id]
		if !ok {
			code := http.StatusNotFound
			body := ResponseBodyGetByID{Message: "product not found", Data: nil}

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
			return
		}

		// response
		code := http.StatusOK
		body := ResponseBodyGetByID{Message: "product", Data: &ProductJSON{Id: id, Name: pr.Name, Quantity: pr.Quantity, CodeValue: pr.CodeValue, IsPublished: pr.IsPublished, Expiration: pr.Expiration.Format("2006-01-02"), Price: pr.Price}}

		w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

// Search is a method that returns a product by id (via query params)
type ResponseBodySearch struct {
	Message string		   `json:"message"`
	Data    []*ProductJSON `json:"data"`
}
func (c *ControllerProducts) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))

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
		body := ResponseBodySearch{Message: "products", Data: make([]*ProductJSON, 0, len(filtered))}
		for k, v := range filtered {
			body.Data = append(body.Data, &ProductJSON{Id: k, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price})
		}

		w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
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
func (c *ControllerProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var reqBody RequestBodyProductCreate
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid request body", Data: nil}

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
			return
		}

		// process
		// -> deserialize
		exp, err := time.Parse("2006-01-02", reqBody.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid date format. Must be yyyy-mm-dd", Data: nil}

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
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

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
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

		w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}