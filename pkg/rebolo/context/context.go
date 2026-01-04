package context

import (
	"encoding/json"
	"net/http"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo/session"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/validation"
	"github.com/gorilla/mux"
)

// AppContext defines the interface for application dependencies
type AppContext interface {
	GetSession(r *http.Request, w http.ResponseWriter) (*session.Session, error)
	Bind(r *http.Request, v interface{}) error
	RenderHTML(w http.ResponseWriter, template string, data interface{}) error
}

// Context wraps http.Request and http.ResponseWriter with convenient helpers
type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	App      AppContext
	params   map[string]string // URL params from gorilla/mux
}

// NewContext creates a new Context instance
func NewContext(w http.ResponseWriter, r *http.Request, app AppContext) *Context {
	return &Context{
		Request:  r,
		Response: w,
		App:      app,
		params:   mux.Vars(r),
	}
}

// Session retrieves the session for the current request
func (c *Context) Session() (*session.Session, error) {
	return c.App.GetSession(c.Request, c.Response)
}

// Flash retrieves flash messages helper
func (c *Context) Flash() (*session.Flash, error) {
	sess, err := c.Session()
	if err != nil {
		return nil, err
	}
	return session.NewFlash(sess), nil
}

// Param retrieves a URL parameter by name (from gorilla/mux)
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Query retrieves a query parameter by name
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// FormValue retrieves a form value by name
func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

// Bind binds request data to a struct with validation
func (c *Context) Bind(v interface{}) error {
	return c.App.Bind(c.Request, v)
}

// Render renders an HTML template with data
func (c *Context) Render(template string, data interface{}) error {
	return c.App.RenderHTML(c.Response, template, data)
}

// JSON sends a JSON response
func (c *Context) JSON(status int, data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewEncoder(c.Response).Encode(data)
}

// String sends a plain text response
func (c *Context) String(status int, text string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(text))
	return err
}

// Redirect redirects to a URL
func (c *Context) Redirect(url string, code int) {
	http.Redirect(c.Response, c.Request, url, code)
}

// Status sets the HTTP status code
func (c *Context) Status(code int) *Context {
	c.Response.WriteHeader(code)
	return c
}

// Set sets a response header
func (c *Context) Set(key, value string) *Context {
	c.Response.Header().Set(key, value)
	return c
}

// Get gets a request header
func (c *Context) Get(key string) string {
	return c.Request.Header.Get(key)
}

// Method returns the HTTP method
func (c *Context) Method() string {
	return c.Request.Method
}

// Path returns the request path
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// IsAjax returns true if the request is an AJAX request
func (c *Context) IsAjax() bool {
	return c.Get("X-Requested-With") == "XMLHttpRequest"
}

// IsJSON returns true if the request content type is JSON
func (c *Context) IsJSON() bool {
	return c.Get("Content-Type") == "application/json"
}

// Error sends an error response
func (c *Context) Error(err error, code int) error {
	http.Error(c.Response, err.Error(), code)
	return err
}

// SaveSession is a helper to save the session
func (c *Context) SaveSession() error {
	sess, err := c.Session()
	if err != nil {
		return err
	}
	return sess.Save()
}

// BindAndValidate binds request data and validates it
func (c *Context) BindAndValidate(v interface{}) error {
	// Bind data
	if err := c.Bind(v); err != nil {
		return err
	}

	// Validate
	return validation.ValidateStruct(v)
}

// ContextHandler is a handler function that accepts Context
type ContextHandler func(*Context) error
