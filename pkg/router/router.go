package router

import (
	"log/slog"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	log *slog.Logger
}

func New(log *slog.Logger) *Router {
	return &Router{
		mux: http.NewServeMux(),
		log: log,
	}
}

type Registrable interface {
	register(r *Router)
}

func (r *Router) Add(handlers ...Registrable) {
	for _, v := range handlers {
		v.register(r)
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
