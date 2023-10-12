package application

import (
	"app/internal/auth"
	"app/internal/product/handlers"
	"app/internal/product/repository"
	"app/internal/product/storage"
	"app/internal/product/validator"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigAppJSON is the configuration of the application
type ConfigAppJSON struct {
	Addr       string
	Token      string
	FilePath   string
	FilePathLastId int
	LayoutDate string
}

// NewApplicationJSON creates a new application
func NewApplicationJSON(cfg *ConfigAppJSON) *ApplicationJSON {
	// default config
	defaultCfg := &ConfigAppJSON{
		Addr:     ":8080",
		Token:    "",
		FilePath: "",
		FilePathLastId: 0,
		LayoutDate: time.DateOnly,
	}
	if cfg != nil {
		if cfg.Addr != "" {
			defaultCfg.Addr = cfg.Addr
		}
		if cfg.Token != "" {
			defaultCfg.Token = cfg.Token
		}
		if cfg.FilePath != "" {
			defaultCfg.FilePath = cfg.FilePath
		}
		if cfg.FilePathLastId != 0 {
			defaultCfg.FilePathLastId = cfg.FilePathLastId
		}
	}

	return &ApplicationJSON{
		rt:       chi.NewRouter(),
		addr:     defaultCfg.Addr,
		token:    defaultCfg.Token,
		filePath: defaultCfg.FilePath,
	}
}

// ApplicationJSON is the main application of the server
type ApplicationJSON struct {
	// rt is the router of the server
	rt *chi.Mux
	// addr is the address of the server
	addr string
	// token is the token of the server
	token string
	// filePath is the file path of the server
	filePath string
	// filePathLastId is the file path of the last id of the server
	filePathLastId int
	// layoutDate is the layout date of the server
	layoutDate string
}

func (a *ApplicationJSON) SetUp() (err error) {
	// dependencies
	// - authenticator
	au := auth.NewAuthTokenBasic("token")

	// - product
	// -- ping
	f, err := os.Open(a.filePath)
	if err != nil {
		return
	}
	f.Close()
	st := storage.NewStorageProductJSON(a.filePath, a.layoutDate)
	vl := validator.NewValidatorProductDefault("")
	rp := repository.NewRepositoryProductStore(st, a.filePathLastId, a.layoutDate)
	rpVl := repository.NewRepositoryProductValidate(rp, vl)
	hd := handlers.NewHandlerProducts(rpVl, au)

	// server
	// - middlewares
	a.rt.Use(middleware.Logger)
	a.rt.Use(middleware.Recoverer)

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

func (a *ApplicationJSON) Run() (err error) {
	err = http.ListenAndServe(a.addr, a.rt)
	return
}
