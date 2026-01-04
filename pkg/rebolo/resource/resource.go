package resource

import (
	"errors"
	"net/http"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo/context"
)

// Resource interface allows for the easy mapping
// of common RESTful actions to a set of paths.
// Similar to Buffalo's Resource interface but using Context.
type Resource interface {
	List(*context.Context) error
	Show(*context.Context) error
	Create(*context.Context) error
	Update(*context.Context) error
	Destroy(*context.Context) error
}

// BaseResource fills in the gaps for any Resource interface
// functions you don't want/need to implement.
// This allows you to only implement the methods you need.
type BaseResource struct{}

// List default implementation. Returns a 404
func (v BaseResource) List(ctx *context.Context) error {
	return ctx.Error(errors.New("resource not implemented"), http.StatusNotFound)
}

// Show default implementation. Returns a 404
func (v BaseResource) Show(ctx *context.Context) error {
	return ctx.Error(errors.New("resource not implemented"), http.StatusNotFound)
}

// Create default implementation. Returns a 404
func (v BaseResource) Create(ctx *context.Context) error {
	return ctx.Error(errors.New("resource not implemented"), http.StatusNotFound)
}

// Update default implementation. Returns a 404
func (v BaseResource) Update(ctx *context.Context) error {
	return ctx.Error(errors.New("resource not implemented"), http.StatusNotFound)
}

// Destroy default implementation. Returns a 404
func (v BaseResource) Destroy(ctx *context.Context) error {
	return ctx.Error(errors.New("resource not implemented"), http.StatusNotFound)
}

// Middler can be implemented to specify additional
// middleware specific to the resource
type Middler interface {
	Use() []interface{} // Middleware functions
}
