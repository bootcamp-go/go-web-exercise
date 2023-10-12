package application

import (
	"app/internal/auth"
	middlewareAuth "app/internal/auth/middleware"
	"app/internal/product/handlers"
	"app/internal/product/repository"
	"app/internal/product/validator"
	middlewareLog "app/platform/web/middleware"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigAppMap is the configuration of the application
type ConfigAppMap struct {
	Addr  string
	Token string
	LayoutDate string
}
// NewApplicationMap creates a new application
func NewApplicationMap(cfg *ConfigAppMap) *ApplicationMap {
	// default config
	defaultCfg := &ConfigAppMap{
		Addr:  ":8080",
		Token: "",
		LayoutDate: time.DateOnly,
	}
	if cfg != nil {
		if cfg.Addr != "" {
			defaultCfg.Addr = cfg.Addr
		}
		if cfg.Token != "" {
			defaultCfg.Token = cfg.Token
		}
		if cfg.LayoutDate != "" {
			defaultCfg.LayoutDate = cfg.LayoutDate
		}
	}

	return &ApplicationMap{
		rt:    chi.NewRouter(),
		addr:  defaultCfg.Addr,
		token: defaultCfg.Token,
		layoutDate: defaultCfg.LayoutDate,
	}
}

// ApplicationMap is the main application of the server
type ApplicationMap struct {
	// rt is the router of the server
	rt *chi.Mux
	// addr is the address of the server
	addr string
	// token is the token of the server
	token string
	// layoutDate is the layout date of the server
	layoutDate string
}

func (a *ApplicationMap) SetUp() (err error) {
	// dependencies
	// - logger
	lgMd := middlewareLog.NewLogger()

	// - authenticator
	au := auth.NewAuthTokenBasic(a.token)
	auMd := middlewareAuth.NewAuthenticator(au)


	// - product
	vl := validator.NewValidatorProductDefault("")
	rp := repository.NewRepositoryProductMap(make(map[int]repository.ProductAttributesMap), 0, a.layoutDate)
	rpVl := repository.NewRepositoryProductValidate(rp, vl)
	hd := handlers.NewHandlerProducts(rpVl)

	// server
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)
	a.rt.Use(lgMd.Log)
	a.rt.Use(auMd.Auth)

	// - routes
	a.rt.Route("/products", func(rt chi.Router) {
		// get all products
		rt.Get("/", hd.Get())
		// get product by id
		rt.Get("/{id}", hd.GetByID())
		// search products
		rt.Get("/search", hd.Search())
		// create product
		rt.Post("/", hd.Create())
		// update or create product
		rt.Put("/{id}", hd.UpdateOrCreate())
		// update product
		rt.Patch("/{id}", hd.Update())
		// delete product
		rt.Delete("/{id}", hd.Delete())
	})

	return
}

func (a *ApplicationMap) Run() (err error) {
	err = http.ListenAndServe(a.addr, a.rt)
	return
}