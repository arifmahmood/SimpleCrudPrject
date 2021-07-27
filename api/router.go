package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"simple-crud-project/repo"
)



type Router struct {
	*chi.Mux
	urlRepo   repo.Url
}


func NewRouter(url repo.Url) *Router {
	router := &Router{
		Mux:        chi.NewRouter(),
		urlRepo:   url,
	}
	register(router)
	return router
}

var logger = logrus.New()
func init() {
	logger.SetLevel(logrus.DebugLevel)
}


func register(router *Router) {

	router.Use(middleware.RequestLogger(&middleware.DefaultLogFormatter{Logger: logger}))
	router.Use(recoverer)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err := newAPIError("Not Found", errURINotFound, nil)
		panic(err)
	})

	router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		err := newAPIError("Method Not Allowed", errInvalidMethod, nil)
		resp := response{
			code:   http.StatusMethodNotAllowed,
			Errors: []apiError{*err},
		}
		resp.serveJSON(w)
	})
	router.Route("/", func(r chi.Router) {
		r.Mount("/url", userHandlers(router))
	})
}


func userHandlers(rt *Router) http.Handler {
	h := chi.NewRouter()
	h.Group(func(r chi.Router) {
		r.Post("/create", rt.CreateNewUrl)
		r.Get("/get/{urlName}", rt.GetUrl)
		r.Get("/delete/{urlName}", rt.deleteUrl)
	})

	return h
}


