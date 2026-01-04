package adapters

import (
	"net/http"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo/core"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/routing"
	"github.com/gorilla/mux"
)

// MuxRouter implements Router interface
type MuxRouter struct {
	*mux.Router
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{
		Router: mux.NewRouter(),
	}
}

func (r *MuxRouter) GET(path string, handler http.HandlerFunc) core.NamedRoute {
	return &routing.NamedRoute{Route: r.HandleFunc(path, handler).Methods("GET")}
}

func (r *MuxRouter) POST(path string, handler http.HandlerFunc) core.NamedRoute {
	return &routing.NamedRoute{Route: r.HandleFunc(path, handler).Methods("POST")}
}

func (r *MuxRouter) PUT(path string, handler http.HandlerFunc) core.NamedRoute {
	return &routing.NamedRoute{Route: r.HandleFunc(path, handler).Methods("PUT")}
}

func (r *MuxRouter) DELETE(path string, handler http.HandlerFunc) core.NamedRoute {
	return &routing.NamedRoute{Route: r.HandleFunc(path, handler).Methods("DELETE")}
}

func (r *MuxRouter) Resource(path string, controller core.Controller) {
	base := path
	r.GET(base, controller.Index)
	r.GET(base+"/new", controller.New)
	r.POST(base, controller.Create)
	r.GET(base+"/{id}", controller.Show)
	r.GET(base+"/{id}/edit", controller.Edit)
	r.HandleFunc(base+"/{id}", controller.Update).Methods("PUT", "PATCH")
	r.HandleFunc(base+"/{id}", controller.Delete).Methods("DELETE")
}

func (r *MuxRouter) Use(middleware core.Middleware) {
	r.Router.Use(mux.MiddlewareFunc(middleware))
}
