package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chiRouter struct {
	dispatcher *chi.Mux
}

func NewChiRouter() Router {
	return &chiRouter{dispatcher: chi.NewRouter()}
}

func (r *chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.Get(uri, WrapHandler(f, chiExtractParams))
}

func (r *chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.Post(uri, WrapHandler(f, chiExtractParams))
}

func (r *chiRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.Put(uri, WrapHandler(f, chiExtractParams))
}

func (r *chiRouter) SERVE(port string) error {
	log.Printf("Chi HTTP server running on port %v", port)
	return http.ListenAndServe(":"+port, r.dispatcher)
}

func chiExtractParams(r *http.Request) map[string]string {
	return map[string]string{
		"gameName": chi.URLParam(r, "gamename"),
		"userName": chi.URLParam(r, "username"),
	}
}
