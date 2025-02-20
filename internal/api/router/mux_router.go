package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct {
	dispatcher *mux.Router
}

func NewMuxRouter() Router {
	return &muxRouter{dispatcher: mux.NewRouter()}
}

func (r *muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodGet)
}

func (r *muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodPost)
}

func (r *muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	r.dispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodPut)
}

func (r *muxRouter) SERVE(port string) error {
	log.Printf("Mux HTTP server running on port %v", port)
	return http.ListenAndServe(":"+port, r.dispatcher)
}

func muxExtractParams(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	return map[string]string{
		"gameName": vars["gamename"],
		"userName": vars["username"],
	}
}
