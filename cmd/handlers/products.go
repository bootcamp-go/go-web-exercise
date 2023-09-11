package handlers

import (
	"encoding/json"
	"net/http"
)

// NewControllerProducts returns a new ControllerProducts
func NewControllerProducts(storage map[int]*Product) *ControllerProducts {
	return &ControllerProducts{storage}
}

// Product is a struct that contains the information of a product
type Product struct {
	Id    	 int
	Name  	 string
	Type 	 string
	Quantity int
	Price 	 float64
}

// ControllerProducts is a struct that contains the storage of products
type ControllerProducts struct {
	storage map[int]*Product
}

// Create is a method that creates a new product
type RequestBodyProduct struct {
	Name  	 string	 `json:"name"`
	Type 	 string	 `json:"type"`
	Quantity int	 `json:"quantity"`
	Price 	 float64 `json:"price"`
}
type ResponseBodyProduct struct {
	Message string `json:"message"`
	Data	*struct {
		Id   	 int		`json:"id"`
		Name  	 string		`json:"name"`
		Type 	 string		`json:"type"`
		Quantity int		`json:"quantity"`
		Price 	 float64	`json:"price"`
	} `json:"data"`
	Error 	bool	`json:"error"`
}
func (c *ControllerProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		token := w.Header().Get("Token")
		if token != "123456" {
			code := http.StatusUnauthorized	// 401
			body := &ResponseBodyProduct{Message: "Unauthorized", Data: nil, Error: true}

			w.WriteHeader(code); w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		}

		// ...
		// request
		var reqBody RequestBodyProduct
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			code := http.StatusBadRequest
			body := &ResponseBodyProduct{
				Message: "Bad Request",
				Data: nil,
				Error: true,
			}

			w.WriteHeader(code)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		}

		// process
		// -> deserialization
		pr := &Product{
			Id: len(c.storage) + 1,
			Name: reqBody.Name,
			Type: reqBody.Type,
			Quantity: reqBody.Quantity,
			Price: reqBody.Price,
		}
		// -> save product
		c.storage[pr.Id] = pr

		// response
		code := http.StatusCreated
		body := &ResponseBodyProduct{
			Message: "Product created",
			Data: &struct {
				Id   	 int		`json:"id"`
				Name  	 string		`json:"name"`
				Type 	 string		`json:"type"`
				Quantity int		`json:"quantity"`
				Price 	 float64	`json:"price"`
			}{Id: pr.Id, Name: pr.Name, Type: pr.Type, Quantity: pr.Quantity, Price: pr.Price},
			Error: false,
		}

		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}