package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// HTMLRenderer implements Renderer interface
type HTMLRenderer struct {
	templates *template.Template
}

func NewHTMLRenderer() *HTMLRenderer {
	tmpl := template.New("root")

	// Walk through views and parse each template with its relative path as name
	err := filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Read the template file
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Register template with path relative to "views/" directory
			// e.g., "views/home/index.html" -> "home/index.html"
			relativePath := path[len("views/"):]

			// Create named template
			t := tmpl.New(relativePath)
			_, err = t.Parse(string(content))
			if err != nil {
				log.Printf("⚠️ Failed to parse %s: %v", path, err)
				return err
			}

			log.Printf("   ✓ Loaded: %s (name: %s)", path, relativePath)
		}
		return nil
	})

	if err != nil {
		log.Printf("❌ Error loading templates: %v", err)
		tmpl = template.New("empty")
	}

	log.Printf("📝 Total templates loaded: %d", len(tmpl.Templates())-1) // -1 for root

	return &HTMLRenderer{templates: tmpl}
}

func (r *HTMLRenderer) RenderHTML(w http.ResponseWriter, templateName string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Try multiple template name formats
	names := []string{
		templateName,                // home/index.html
		"views/" + templateName,     // views/home/index.html
		filepath.Base(templateName), // index.html
		filepath.Base(filepath.Dir(templateName)) + "/" + filepath.Base(templateName), // home/index.html
	}

	var err error
	var renderedName string

	// Capture output to a buffer first (for hot reload injection)
	var buf bytes.Buffer

	for _, name := range names {
		buf.Reset()
		err = r.templates.ExecuteTemplate(&buf, name, data)
		if err == nil {
			renderedName = name
			break
		}
	}

	if err != nil {
		log.Printf("❌ Failed to render template: %s (tried: %v)", templateName, names)
		return err
	}

	log.Printf("✅ Rendered template: %s (requested: %s)", renderedName, templateName)

	// Write to actual response
	_, err = w.Write(buf.Bytes())
	return err
}

func (r *HTMLRenderer) RenderJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

func (r *HTMLRenderer) RenderError(w http.ResponseWriter, message string, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(map[string]string{
		"error":  message,
		"status": fmt.Sprintf("%d", status),
	})
}
