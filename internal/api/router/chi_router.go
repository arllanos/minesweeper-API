package router

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type chiRouter struct{}

var (
	chiDispatcher = chi.NewRouter()
)

func NewChiRouter() Router {
	return &chiRouter{}
}

func (*chiRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Get(uri, WrapHandler(f, chiExtractParams))
}

func (*chiRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Post(uri, WrapHandler(f, chiExtractParams))
}

func (*chiRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	chiDispatcher.Put(uri, WrapHandler(f, chiExtractParams))
}

func (*chiRouter) SERVE(port string) error {
	log.Printf("Chi HTTP server running on port %v", port)
	return http.ListenAndServe(":"+port, chiDispatcher)
}

func chiExtractParams(r *http.Request) map[string]string {
	return map[string]string{
		"gameName": chi.URLParam(r, "gamename"),
		"userName": chi.URLParam(r, "username"),
	}
}
