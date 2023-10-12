package handlers

import (
	"app/internal/auth"
	"app/internal/product"
	"app/internal/product/repository"
	"app/platform/web/request"
	"app/platform/web/response"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

// NewHandlerProducts returns a new HandlerProducts
func NewHandlerProducts(st repository.RepositoryProduct, au auth.AuthToken) *HandlerProducts {
	return &HandlerProducts{st: st, au: au}
}

// HandlerProducts is a struct that contains the repository of products
type HandlerProducts struct {
	// st is the repository of products
	st repository.RepositoryProduct
	// au is the authenticator
	au auth.AuthToken
}

type ProductJSON struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Quantity    int       `json:"quantity"`
	CodeValue   string    `json:"code_value"`
	IsPublished bool      `json:"is_published"`
	Expiration  string	  `json:"expiration"`
	Price       float64   `json:"price"`
}

// Get is a method that returns all products
func (h *HandlerProducts) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// request
		// ...

		// process
		pr, err := h.st.Get()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Internal error")
			return
		}

		// response
		data := make([]ProductJSON, 0, len(pr))
		for _, v := range pr {
			data = append(data, ProductJSON{Id: v.Id(), Name: v.Name(), Quantity: v.Quantity(), CodeValue: v.CodeValue(), IsPublished: v.IsPublished(), Expiration: v.Expiration().Format("2006-01-02"), Price: v.Price()})
		}
		response.JSON(w, http.StatusOK, map[string]any{"message": "products", "data": data})
	}
}

// GetByID is a method that returns a product by id
func (h *HandlerProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		// request
		// - id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		// process
		// - get product
		pr, err := h.st.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRepositoryProductNotFound):
				response.Error(w, http.StatusNotFound, "Product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product",
			"data": ProductJSON{Id: pr.Id(), Name: pr.Name(), Quantity: pr.Quantity(), CodeValue: pr.CodeValue(), IsPublished: pr.IsPublished(), Expiration: pr.Expiration().Format("2006-01-02"), Price: pr.Price()},
		})
	}
}

// Search is a method that returns a product by id (via query params)
func (h *HandlerProducts) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		// request
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		// process
		// - get product with query
		query := repository.Query{Id: id}
		pr, err := h.st.Search(query)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Internal error")
			return
		}
			
		// response
		data := make([]ProductJSON, 0, len(pr))
		for _, v := range pr {
			data = append(data, ProductJSON{Id: v.Id(), Name: v.Name(), Quantity: v.Quantity(), CodeValue: v.CodeValue(), IsPublished: v.IsPublished(), Expiration: v.Expiration().Format("2006-01-02"), Price: v.Price()})
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product",
			"data": data,
		})
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
func (h *HandlerProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		// request
		var reqBody RequestBodyProductCreate
		err = request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// process
		// - deserialize
		exp, err := time.Parse("2006-01-02", reqBody.Expiration)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid date format. Must be yyyy-mm-dd")
			return
		}
		// - save
		pr := product.NewProduct(0, reqBody.Name, reqBody.Quantity, reqBody.CodeValue, reqBody.IsPublished, exp, reqBody.Price)
		err = h.st.Create(pr)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRepositoryProductInvalid):
				response.Error(w, http.StatusUnprocessableEntity, "Invalid product")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "product created",
			"data": ProductJSON{Id: pr.Id(), Name: pr.Name(), Quantity: pr.Quantity(), CodeValue: pr.CodeValue(), IsPublished: pr.IsPublished(), Expiration: pr.Expiration().Format("2006-01-02"), Price: pr.Price()},
		})
	}
}

// UpdateOrCreate is a method that updates or creates a product
type RequestBodyProductUpdateOrCreate struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}
func (h *HandlerProducts) UpdateOrCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		var reqBody RequestBodyProductUpdateOrCreate
		err = request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// process
		// - deserialize
		exp, err := time.Parse("2006-01-02", reqBody.Expiration)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid date format. Must be yyyy-mm-dd")
			return
		}
		pr := product.NewProduct(id, reqBody.Name, reqBody.Quantity, reqBody.CodeValue, reqBody.IsPublished, exp, reqBody.Price)
		// - update or create
		err = h.st.UpdateOrCreate(pr)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRepositoryProductInvalid):
				response.Error(w, http.StatusUnprocessableEntity, "Invalid product")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "product updated or created",
			"data": ProductJSON{Id: pr.Id(), Name: pr.Name(), Quantity: pr.Quantity(), CodeValue: pr.CodeValue(), IsPublished: pr.IsPublished(), Expiration: pr.Expiration().Format("2006-01-02"), Price: pr.Price()},
		})
	}
}

// Update is a method that updates a product
func (h *HandlerProducts) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}
		patch := make(map[string]any)
		err = request.JSON(r, &patch)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		// process
		pr, err := h.st.Update(id, patch)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRepositoryProductNotFound):
				response.Error(w, http.StatusNotFound, "Product not found")
			case errors.Is(err, repository.ErrRepositoryProductInvalid):
				response.Error(w, http.StatusUnprocessableEntity, "Invalid product")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product updated - patched",
			"data": ProductJSON{Id: pr.Id(), Name: pr.Name(), Quantity: pr.Quantity(), CodeValue: pr.CodeValue(), IsPublished: pr.IsPublished(), Expiration: pr.Expiration().Format("2006-01-02"), Price: pr.Price()},
		})
	}
}

// Delete is a method that deletes a product by id
func (h *HandlerProducts) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// auth
		token := r.Header.Get("Token")
		err := h.au.Auth(token)
		if err != nil {
			response.Error(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		
		// request
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid id")
			return
		}

		// process
		// - delete
		err = h.st.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrRepositoryProductNotFound):
				response.Error(w, http.StatusNotFound, "Product not found")
			default:
				response.Error(w, http.StatusInternalServerError, "Internal error")
			}
			return
		}

		// response
		response.JSON(w, http.StatusNoContent, nil)
	}
}