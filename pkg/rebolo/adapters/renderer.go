package adapters

import (
	"html/template"
	"net/http"
	"encoding/json"
	"fmt"
)

// HTMLRenderer implements Renderer interface
type HTMLRenderer struct {
	templates *template.Template
}

func NewHTMLRenderer() *HTMLRenderer {
	tmpl, err := template.ParseGlob("views/**/*.html")
	if err != nil {
		// Fallback to empty template if no views found
		tmpl = template.New("empty")
	}
	return &HTMLRenderer{templates: tmpl}
}

func (r *HTMLRenderer) RenderHTML(w http.ResponseWriter, templateName string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return r.templates.ExecuteTemplate(w, templateName, data)
}

func (r *HTMLRenderer) RenderJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func (r *HTMLRenderer) RenderError(w http.ResponseWriter, message string, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]string{
		"error": message,
		"status": fmt.Sprintf("%d", status),
	})
}
