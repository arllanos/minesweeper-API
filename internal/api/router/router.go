package router

import (
	"context"
	"net/http"
)

type Router interface {
	GET(uri string, f func(w http.ResponseWriter, r *http.Request))
	POST(uri string, f func(w http.ResponseWriter, r *http.Request))
	PUT(uri string, f func(w http.ResponseWriter, r *http.Request))
	SERVE(port string) error
}

func WrapHandler(f func(w http.ResponseWriter, r *http.Request), getParams func(r *http.Request) map[string]string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := getParams(r)
		ctx := r.Context()
		for k, v := range params {
			ctx = context.WithValue(ctx, k, v)
		}
		f(w, r.WithContext(ctx))
	}
}
