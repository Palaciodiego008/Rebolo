package routing

import (
	"fmt"

	"github.com/gorilla/mux"
)

// NamedRoute wraps a mux.Route to provide a fluent API
type NamedRoute struct {
	*mux.Route
}

// Name sets the name for the route
func (r *NamedRoute) Name(name string) *NamedRoute {
	r.Route.Name(name)
	return r
}

// URLFor generates a URL for a named route with the given parameters
func URLFor(router *mux.Router, name string, params map[string]string) (string, error) {
	route := router.Get(name)
	if route == nil {
		return "", fmt.Errorf("route %s not found", name)
	}

	// Build URL
	url, err := route.URL(pairsFromMap(params)...)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

// pairsFromMap converts a map to key-value pairs for mux.URL()
func pairsFromMap(params map[string]string) []string {
	pairs := make([]string, 0, len(params)*2)
	for k, v := range params {
		pairs = append(pairs, k, v)
	}
	return pairs
}

// URLForString is a convenience function that returns the URL as a string
// or returns an empty string if there's an error
func URLForString(router *mux.Router, name string, params map[string]string) string {
	url, err := URLFor(router, name, params)
	if err != nil {
		return ""
	}
	return url
}
