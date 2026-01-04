package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type ResourceData struct {
	Name       string
	NameLower  string
	NamePlural string
	Fields     []Field
	Timestamp  string
}

type Field struct {
	Name string
	Type string
	GoType string
}

func generateResource(name string, fieldArgs []string) {
	fields := parseFields(fieldArgs)
	
	data := ResourceData{
		Name:       strings.Title(name),
		NameLower:  strings.ToLower(name),
		NamePlural: name + "s", // Simple pluralization
		Fields:     fields,
		Timestamp:  time.Now().Format("20060102150405"),
	}
	
	// Generate model
	generateModel(data)
	
	// Generate controller
	generateController(data)
	
	// Generate migration
	generateMigration(data)
	
	// Generate views
	generateViews(data)
	
	fmt.Printf("✅ Generated resource: %s\n", name)
	fmt.Printf("   - Model: models/%s.go\n", data.NameLower)
	fmt.Printf("   - Controller: controllers/%s_controller.go\n", data.NameLower)
	fmt.Printf("   - Migration: db/migrations/%s_create_%s.sql\n", data.Timestamp, data.NamePlural)
	fmt.Printf("   - Views: views/%s/\n", data.NamePlural)
}

func parseFields(fieldArgs []string) []Field {
	var fields []Field
	
	for _, arg := range fieldArgs {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			continue
		}
		
		fieldName := parts[0]
		fieldType := parts[1]
		
		field := Field{
			Name: fieldName,
			Type: fieldType,
			GoType: mapToGoType(fieldType),
		}
		
		fields = append(fields, field)
	}
	
	return fields
}

func mapToGoType(dbType string) string {
	switch dbType {
	case "string", "text":
		return "string"
	case "int", "integer":
		return "int64"
	case "bool", "boolean":
		return "bool"
	case "float":
		return "float64"
	case "time", "datetime":
		return "time.Time"
	default:
		return "string"
	}
}

const modelTemplate = `package models

import (
	"time"
	"github.com/uptrace/bun"
)

type {{.Name}} struct {
	bun.BaseModel ` + "`" + `bun:"table:{{.NamePlural}}"` + "`" + `
	
	ID        int64     ` + "`" + `bun:",pk,autoincrement"` + "`" + `
{{range .Fields}}	{{.Name | title}}    {{.GoType}}   ` + "`" + `bun:"{{.Name}}"` + "`" + `
{{end}}	CreatedAt time.Time ` + "`" + `bun:",nullzero,notnull,default:current_timestamp"` + "`" + `
	UpdatedAt time.Time ` + "`" + `bun:",nullzero,notnull,default:current_timestamp"` + "`" + `
}
`

const controllerTemplate = `package controllers

import (
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo"
	"../models"
)

type {{.Name}}Controller struct{}

func (c *{{.Name}}Controller) Index(w http.ResponseWriter, r *http.Request) {
	// Fetch all {{.NamePlural}} from database
	var {{.NameLower}}s []models.{{.Name}}
	
	// TODO: Implement database query
	// err := app.DB.NewSelect().Model(&{{.NameLower}}s).Scan(r.Context())
	
	rebolo.Render(w, "{{.NamePlural}}/index.html", map[string]interface{}{
		"{{.Name}}s": {{.NameLower}}s,
	})
}

func (c *{{.Name}}Controller) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var {{.NameLower}} models.{{.Name}}
	
	// TODO: Implement database query
	// err = app.DB.NewSelect().Model(&{{.NameLower}}).Where("id = ?", id).Scan(r.Context())
	
	rebolo.Render(w, "{{.NamePlural}}/show.html", map[string]interface{}{
		"{{.Name}}": {{.NameLower}},
	})
}

func (c *{{.Name}}Controller) New(w http.ResponseWriter, r *http.Request) {
	rebolo.Render(w, "{{.NamePlural}}/new.html", nil)
}

func (c *{{.Name}}Controller) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	{{.NameLower}} := models.{{.Name}}{
{{range .Fields}}		{{.Name | title}}: r.FormValue("{{.Name}}"),
{{end}}	}
	
	// TODO: Implement database insert
	// _, err := app.DB.NewInsert().Model(&{{.NameLower}}).Exec(r.Context())
	
	http.Redirect(w, r, "/{{.NamePlural}}", http.StatusSeeOther)
}

func (c *{{.Name}}Controller) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var {{.NameLower}} models.{{.Name}}
	
	// TODO: Implement database query
	// err = app.DB.NewSelect().Model(&{{.NameLower}}).Where("id = ?", id).Scan(r.Context())
	
	rebolo.Render(w, "{{.NamePlural}}/edit.html", map[string]interface{}{
		"{{.Name}}": {{.NameLower}},
	})
}

func (c *{{.Name}}Controller) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}
	
	{{.NameLower}} := models.{{.Name}}{
		ID: id,
{{range .Fields}}		{{.Name | title}}: r.FormValue("{{.Name}}"),
{{end}}	}
	
	// TODO: Implement database update
	// _, err = app.DB.NewUpdate().Model(&{{.NameLower}}).Where("id = ?", id).Exec(r.Context())
	
	http.Redirect(w, r, "/{{.NamePlural}}", http.StatusSeeOther)
}

func (c *{{.Name}}Controller) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	// TODO: Implement database delete
	// _, err = app.DB.NewDelete().Model((*models.{{.Name}})(nil)).Where("id = ?", id).Exec(r.Context())
	
	http.Redirect(w, r, "/{{.NamePlural}}", http.StatusSeeOther)
}
`

func generateModel(data ResourceData) {
	tmpl := template.Must(template.New("model").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(modelTemplate))
	
	file, _ := os.Create(filepath.Join("models", data.NameLower+".go"))
	defer file.Close()
	
	tmpl.Execute(file, data)
}

func generateController(data ResourceData) {
	os.MkdirAll("controllers", 0755)
	
	tmpl := template.Must(template.New("controller").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(controllerTemplate))
	
	file, _ := os.Create(filepath.Join("controllers", data.NameLower+"_controller.go"))
	defer file.Close()
	
	tmpl.Execute(file, data)
}

func generateMigration(data ResourceData) {
	os.MkdirAll("db/migrations", 0755)
	
	migrationSQL := fmt.Sprintf(`CREATE TABLE %s (
    id BIGSERIAL PRIMARY KEY,
`, data.NamePlural)
	
	for _, field := range data.Fields {
		sqlType := mapToSQLType(field.Type)
		migrationSQL += fmt.Sprintf("    %s %s,\n", field.Name, sqlType)
	}
	
	migrationSQL += `    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);`
	
	filename := fmt.Sprintf("%s_create_%s.sql", data.Timestamp, data.NamePlural)
	file, _ := os.Create(filepath.Join("db", "migrations", filename))
	defer file.Close()
	
	file.WriteString(migrationSQL)
}

func generateViews(data ResourceData) {
	viewsDir := filepath.Join("views", data.NamePlural)
	os.MkdirAll(viewsDir, 0755)
	
	// Generate comprehensive CRUD views
	views := map[string]string{
		"index.html": fmt.Sprintf(`<h1>%s</h1>
<a href="/%s/new" style="background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">New %s</a>

<div id="%s-list" style="margin-top: 2rem;">
    {{range .%s}}
    <div style="border: 1px solid #ddd; padding: 1rem; margin: 1rem 0; border-radius: 5px;">
        <h3><a href="/%s/{{.ID}}">{{.%s}}</a></h3>
        <div>
            <a href="/%s/{{.ID}}/edit">Edit</a> |
            <form method="POST" action="/%s/{{.ID}}" style="display: inline;">
                <input type="hidden" name="_method" value="DELETE">
                <button type="submit" onclick="return confirm('Are you sure?')" style="background: #f44336; color: white; border: none; padding: 5px 10px; border-radius: 3px;">Delete</button>
            </form>
        </div>
    </div>
    {{end}}
</div>`, 
			strings.Title(data.NamePlural), data.NamePlural, data.Name, data.NamePlural, 
			strings.Title(data.NamePlural), data.NamePlural, getFirstStringField(data.Fields), 
			data.NamePlural, data.NamePlural),
		
		"show.html": fmt.Sprintf(`<h1>%s Details</h1>
<div style="background: #f9f9f9; padding: 2rem; border-radius: 5px; margin: 1rem 0;">
%s
</div>
<div>
    <a href="/%s/{{.ID}}/edit" style="background: #2196F3; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px;">Edit</a>
    <a href="/%s" style="background: #666; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-left: 10px;">Back to List</a>
</div>`, 
			data.Name, generateShowFields(data.Fields), data.NamePlural, data.NamePlural),
		
		"new.html": fmt.Sprintf(`<h1>New %s</h1>
<form method="POST" action="/%s" style="max-width: 500px;">
%s
    <div style="margin-top: 1rem;">
        <button type="submit" style="background: #4CAF50; color: white; padding: 10px 20px; border: none; border-radius: 5px;">Create %s</button>
        <a href="/%s" style="background: #666; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-left: 10px;">Cancel</a>
    </div>
</form>`, 
			data.Name, data.NamePlural, generateFormFields(data.Fields), data.Name, data.NamePlural),
		
		"edit.html": fmt.Sprintf(`<h1>Edit %s</h1>
<form method="POST" action="/%s/{{.ID}}" style="max-width: 500px;">
    <input type="hidden" name="_method" value="PUT">
%s
    <div style="margin-top: 1rem;">
        <button type="submit" style="background: #2196F3; color: white; padding: 10px 20px; border: none; border-radius: 5px;">Update %s</button>
        <a href="/%s/{{.ID}}" style="background: #666; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; margin-left: 10px;">Cancel</a>
    </div>
</form>`, 
			data.Name, data.NamePlural, generateFormFields(data.Fields), data.Name, data.NamePlural),
	}
	
	for filename, content := range views {
		file, _ := os.Create(filepath.Join(viewsDir, filename))
		file.WriteString(content)
		file.Close()
	}
}

func getFirstStringField(fields []Field) string {
	for _, field := range fields {
		if field.GoType == "string" {
			return strings.Title(field.Name)
		}
	}
	return "ID"
}

func generateShowFields(fields []Field) string {
	var html string
	for _, field := range fields {
		html += fmt.Sprintf(`    <p><strong>%s:</strong> {{.%s}}</p>
`, strings.Title(field.Name), strings.Title(field.Name))
	}
	return html
}

func generateFormFields(fields []Field) string {
	var html string
	for _, field := range fields {
		inputType := "text"
		if field.Type == "bool" || field.Type == "boolean" {
			html += fmt.Sprintf(`    <div style="margin-bottom: 1rem;">
        <label style="display: block; margin-bottom: 5px;"><strong>%s:</strong></label>
        <input type="checkbox" name="%s" value="true" style="transform: scale(1.2);">
    </div>
`, strings.Title(field.Name), field.Name)
		} else {
			if field.Type == "text" {
				inputType = "textarea"
				html += fmt.Sprintf(`    <div style="margin-bottom: 1rem;">
        <label style="display: block; margin-bottom: 5px;"><strong>%s:</strong></label>
        <textarea name="%s" style="width: 100%%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;" rows="4"></textarea>
    </div>
`, strings.Title(field.Name), field.Name)
			} else {
				if field.Type == "int" || field.Type == "integer" {
					inputType = "number"
				}
				html += fmt.Sprintf(`    <div style="margin-bottom: 1rem;">
        <label style="display: block; margin-bottom: 5px;"><strong>%s:</strong></label>
        <input type="%s" name="%s" style="width: 100%%; padding: 8px; border: 1px solid #ddd; border-radius: 4px;">
    </div>
`, strings.Title(field.Name), inputType, field.Name)
			}
		}
	}
	return html
}

func mapToSQLType(goType string) string {
	switch goType {
	case "string":
		return "VARCHAR(255)"
	case "text":
		return "TEXT"
	case "int", "integer":
		return "BIGINT"
	case "bool", "boolean":
		return "BOOLEAN"
	case "float":
		return "DECIMAL"
	case "time", "datetime":
		return "TIMESTAMP"
	default:
		return "VARCHAR(255)"
	}
}
