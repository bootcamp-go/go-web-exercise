package handlers_test

import (
	"app/internal/auth"
	"app/internal/product"
	"app/internal/product/handlers"
	"app/internal/product/repository"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

// Tests for HandlerProducts.Get handler
func TestHandlerProducts_Get_Handler(t *testing.T) {
	t.Run("success to get a list of products", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		exp := time.Now().Add(time.Hour * 24 * 30) // 30 days
		rp := repository.NewRepositoryProductMock()
		rp.FuncGet = func() (p []product.Product, err error) {
			p = []product.Product{
				*product.NewProduct(1, "product 1", 1, "code 1", true, exp, 1.1),
			}
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.Get()

		// act
		req := &http.Request{}
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := fmt.Sprintf(
			`{"message": "products", "data": [{"id": 1, "name": "product 1", "quantity": 1, "code_value": "code 1", "is_published": true, "expiration": "%v", "price": 1.1}]}`,
			exp.Format(time.DateOnly),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})
}

// Tests for HandlerProducts.GetByID handler
func TestHandlerProducts_GetByID_Handler(t *testing.T) {
	t.Run("success to get a product by id", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		exp := time.Now().Add(time.Hour * 24 * 30) // 30 days
		rp := repository.NewRepositoryProductMock()
		rp.FuncGetByID = func(id int) (p product.Product, err error) {
			p = *product.NewProduct(1, "product 1", 1, "code 1", true, exp, 1.1)
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.GetByID()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := fmt.Sprintf(
			`{"message": "product", "data": {"id": 1, "name": "product 1", "quantity": 1, "code_value": "code 1", "is_published": true, "expiration": "%v", "price": 1.1}}`,
			exp.Format(time.DateOnly),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.GetByID()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Invalid id"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - product not found", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		rp := repository.NewRepositoryProductMock()
		rp.FuncGetByID = func(id int) (p product.Product, err error) {
			err = repository.ErrRepositoryProductNotFound
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.GetByID()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Product not found"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})
}

// Tests for HandlerProducts.Create handler
func TestHandlerProducts_Create_Handler(t *testing.T) {
	t.Run("success to create a product", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		rp := repository.NewRepositoryProductMock()
		rp.FuncCreate = func(p *product.Product) (err error) {
			(*p).SetId(1)
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.Create()

		// act
		exp := time.Now().Add(time.Hour * 24 * 30) // 30 days
		req := &http.Request{
			Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
				`{"name": "product 1", "quantity": 1, "code_value": "code 1", "is_published": true, "expiration": "%s", "price": 1.1}`,
				exp.Format(time.DateOnly),
			))),
		}
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusCreated
		expectedBody := fmt.Sprintf(
			`{"message": "product created", "data": {"id": 1, "name": "product 1", "quantity": 1, "code_value": "code 1", "is_published": true, "expiration": "%v", "price": 1.1}}`,
			exp.Format(time.DateOnly),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - invalid token", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			err = auth.ErrAuthTokenInvalid
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.Create()

		// act
		req := &http.Request{}
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusUnauthorized
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Unauthorized"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})
}

// Tests for HandlerProducts.UpdateOrCreate handler
func TestHandlerProducts_UpdateOrCreate_Handler(t *testing.T) {
	t.Run("fail - invalid id", func(t *testing.T) {
		// arrange
		
		// act
		
		// assert
		
	})

	t.Run("fail - invalid token", func(t *testing.T) {
		// arrange

		// act

		// assert
	})
}

// Tests for HandlerProducts.Update handler
func TestHandlerProducts_Update_Handler(t *testing.T) {
	t.Run("fail - invalid id", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.Update()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Invalid id"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - product not found", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		rp := repository.NewRepositoryProductMock()
		rp.FuncUpdate = func(id int, patch map[string]interface{}) (p product.Product, err error) {
			err = repository.ErrRepositoryProductNotFound
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.Update()

		// act
		exp := time.Now().Add(time.Hour * 24 * 30) // 30 days
		req := &http.Request{
			Body: io.NopCloser(strings.NewReader(fmt.Sprintf(
				`{"name": "product 1", "quantity": 1, "code_value": "code 1", "is_published": true, "expiration": "%s", "price": 1.1}`,
				exp.Format(time.DateOnly),
			))),
		}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Product not found"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - invalid token", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			err = auth.ErrAuthTokenInvalid
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.Update()

		// act
		req := &http.Request{}
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusUnauthorized
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Unauthorized"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})
}

// Tests for HandlerProducts.Delete handler
func TestHandlerProducts_Delete_Handler(t *testing.T) {
	t.Run("success to delete a product", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		rp := repository.NewRepositoryProductMock()
		rp.FuncDelete = func(id int) (err error) {
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.Delete()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNoContent
		expectedBody := ""
		expectedHeaders := http.Header{}
		require.Equal(t, expectedCode, res.Code)
		require.Equal(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.Delete()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "invalid")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusBadRequest
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Invalid id"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - product not found", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			return
		}

		// - repository: mock
		rp := repository.NewRepositoryProductMock()
		rp.FuncDelete = func(id int) (err error) {
			err = repository.ErrRepositoryProductNotFound
			return
		}

		// - handler
		hd := handlers.NewHandlerProducts(rp, at)
		hdFunc := hd.Delete()

		// act
		req := &http.Request{}
		chiCtx := chi.NewRouteContext()	// *chi.Context to handle params
		chiCtx.URLParams.Add("id", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx)) // replace *http.Request with the new request having the updated context

		res := httptest.NewRecorder()

		hdFunc(res, req)

		// assert
		expectedCode := http.StatusNotFound
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Product not found"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})

	t.Run("fail - invalid token", func(t *testing.T) {
		// arrange
		// - auth: mock
		at := auth.NewAuthTokenMock()
		at.FuncAuth = func(token string) (err error) {
			err = auth.ErrAuthTokenInvalid
			return
		}

		// - repository: mock
		// ...

		// - handler
		hd := handlers.NewHandlerProducts(nil, at)
		hdFunc := hd.Delete()

		// act
		req := &http.Request{}
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusUnauthorized
		expectedBody := fmt.Sprintf(
			`{"status": "%s", "message": "Unauthorized"}`,
			http.StatusText(expectedCode),
		)
		expectedHeaders := http.Header{"Content-Type": []string{"application/json; charset=utf-8"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeaders, res.Header())
	})
}