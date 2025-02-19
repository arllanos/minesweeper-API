package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type muxRouter struct{}

var (
	muxDispatcher = mux.NewRouter()
)

func NewMuxRouter() Router {
	return &muxRouter{}
}

func (*muxRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodGet)
}

func (*muxRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodPost)
}

func (*muxRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	muxDispatcher.HandleFunc(uri, WrapHandler(f, muxExtractParams)).Methods(http.MethodPut)
}

func (*muxRouter) SERVE(port string) {
	log.Printf("Mux HTTP server running on port %v", port)
	http.ListenAndServe(":"+port, muxDispatcher)
}

func muxExtractParams(r *http.Request) map[string]string {
	vars := mux.Vars(r)
	return map[string]string{
		"gameName": vars["gamename"],
		"userName": vars["username"],
	}
}
