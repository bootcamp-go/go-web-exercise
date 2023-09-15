package handlers

import (
	"app/internal/product/storage"
	"app/pkg/web/request"
	"app/pkg/web/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewHandlerProducts returns a new HandlerProducts
func NewHandlerProducts(st storage.StorageProduct) *HandlerProducts {
	return &HandlerProducts{st: st}
}

// HandlerProducts is a struct that contains the storage of products
type HandlerProducts struct {
	st storage.StorageProduct
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
func (h *HandlerProducts) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		pr, err := h.st.Get()
		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBodyProductGet{Message: "internal error", Data: nil}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := ResponseBodyProductGet{Message: "products", Data: make([]*ProductJSON, 0, len(pr))}
		for _, v := range pr {
			body.Data = append(body.Data, &ProductJSON{Id: v.Id, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price})
		}

		response.JSON(w, code, body)
	}
}

// GetByID is a method that returns a product by id
type ResponseBodyGetByID struct {
	Message string		 `json:"message"`
	Data    *ProductJSON `json:"data"`
}
func (h *HandlerProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// -> id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyGetByID{Message: "invalid id", Data: nil}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> get product
		pr, err := h.st.GetByID(id)
		if err != nil {
			var code int; var body ResponseBodyGetByID
			switch {
			case errors.Is(err, storage.ErrStorageProductNotFound):
				code = http.StatusNotFound
				body = ResponseBodyGetByID{Message: "product not found", Data: nil}
			default:
				code = http.StatusInternalServerError
				body = ResponseBodyGetByID{Message: "internal error", Data: nil}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusOK
		body := ResponseBodyGetByID{Message: "product", Data: &ProductJSON{Id: id, Name: pr.Name, Quantity: pr.Quantity, CodeValue: pr.CodeValue, IsPublished: pr.IsPublished, Expiration: pr.Expiration.Format("2006-01-02"), Price: pr.Price}}

		response.JSON(w, code, body)
	}
}

// Search is a method that returns a product by id (via query params)
type ResponseBodySearch struct {
	Message string		   `json:"message"`
	Data    []*ProductJSON `json:"data"`
}
func (h *HandlerProducts) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodySearch{Message: "invalid id", Data: nil}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> get product with query
		query := storage.Query{Id: id}
		pr, err := h.st.Search(query)
		if err != nil {
			code := http.StatusInternalServerError
			body := ResponseBodySearch{Message: "internal error", Data: nil}

			response.JSON(w, code, body)
			return
		}
			
		// response
		code := http.StatusOK
		body := ResponseBodySearch{Message: "products", Data: make([]*ProductJSON, 0, len(pr))}
		for _, v := range pr {
			body.Data = append(body.Data, &ProductJSON{Id: v.Id, Name: v.Name, Quantity: v.Quantity, CodeValue: v.CodeValue, IsPublished: v.IsPublished, Expiration: v.Expiration.Format("2006-01-02"), Price: v.Price})
		}

		response.JSON(w, code, body)
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
func (h *HandlerProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		var reqBody RequestBodyProductCreate
		err := request.JSON(r, &reqBody)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid request body", Data: nil}

			response.JSON(w, code, body)
			return
		}

		// process
		// -> deserialize
		exp, err := time.Parse("2006-01-02", reqBody.Expiration)
		if err != nil {
			code := http.StatusBadRequest
			body := ResponseBodyProductCreate{Message: "invalid date format. Must be yyyy-mm-dd", Data: nil}

			response.JSON(w, code, body)
			return
		}
		// -> save
		pr := &storage.Product{
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,	
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  exp,
			Price:       reqBody.Price,
		}
		err = h.st.Create(pr)
		if err != nil {
			var code int; var body ResponseBodyProductCreate
			switch {
			case errors.Is(err, storage.ErrStorageProductInvalid):
				code = http.StatusUnprocessableEntity
				body = ResponseBodyProductCreate{Message: "invalid product", Data: nil}
			default:
				code = http.StatusInternalServerError
				body = ResponseBodyProductCreate{Message: "internal error", Data: nil}
			}

			response.JSON(w, code, body)
			return
		}

		// response
		code := http.StatusCreated
		body := ResponseBodyProductCreate{
			Message: "product created",
			Data: &ProductJSON{										// serialization
				Id:          pr.Id,
				Name:        pr.Name,
				Quantity:    pr.Quantity,
				CodeValue:   pr.CodeValue,
				IsPublished: pr.IsPublished,
				Expiration:  pr.Expiration.Format("2006-01-02"),
				Price:       pr.Price,
			},
		}

		response.JSON(w, code, body)
	}
}