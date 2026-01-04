package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed templates
var templates embed.FS

type Generator struct {
	templates   *template.Template
	typeMapping *FieldTypeMapping
}

type AppData struct {
	Name      string
	Module    string
	Framework string
	Title     string
}

type ResourceData struct {
	Name       string
	VarName    string
	TableName  string
	ViewPath   string
	RoutePath  string
	Fields     []Field
	FirstField string
	Timestamp  string
}

type Field struct {
	Name     string
	DBName   string
	FormName string
	GoType   string
	SQLType  string
	HTMLType string
}

func NewGenerator() *Generator {
	// Parse all template files recursively
	tmpl := template.New("").Funcs(template.FuncMap{
		"title": func(s string) string { return cases.Title(language.English).String(s) },
		"lower": strings.ToLower,
	})

	// Parse templates manually to handle nested directories
	tmpl = template.Must(tmpl.ParseFS(templates,
		"templates/app/main.go.tmpl",
		"templates/app/package.json.tmpl",
		"templates/app/src/index.js.tmpl",
		"templates/app/views/layouts/application.html.tmpl",
		"templates/app/views/home/index.html.tmpl",
		"templates/config/config.yml.tmpl",
		"templates/resource/model.go.tmpl",
		"templates/resource/controller.go.tmpl",
		"templates/resource/migration.sql.tmpl",
		"templates/resource/index.html.tmpl",
		"templates/resource/show.html.tmpl",
		"templates/resource/new.html.tmpl",
		"templates/resource/edit.html.tmpl",
	))

	return &Generator{
		templates:   tmpl,
		typeMapping: DefaultFieldTypeMapping(),
	}
}

func (g *Generator) GenerateApp(name string) error {
	data := AppData{
		Name:      name,
		Module:    fmt.Sprintf("github.com/Palaciodiego008/%s", name),
		Framework: "ReboloLang",
		Title:     fmt.Sprintf("Welcome to %s", name),
	}

	// Create directory structure
	dirs := []string{
		name,
		filepath.Join(name, "controllers"),
		filepath.Join(name, "models"),
		filepath.Join(name, "views", "home"),
		filepath.Join(name, "views", "layouts"),
		filepath.Join(name, "public"),
		filepath.Join(name, "src"),
		filepath.Join(name, "db", "migrations"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Generate files from templates
	files := map[string]string{
		filepath.Join(name, "main.go"):                              "app/main.go.tmpl",
		filepath.Join(name, "package.json"):                         "app/package.json.tmpl",
		filepath.Join(name, "config.yml"):                           "config/config.yml.tmpl",
		filepath.Join(name, "src", "index.js"):                      "app/src/index.js.tmpl",
		filepath.Join(name, "views", "layouts", "application.html"): "app/views/layouts/application.html.tmpl",
		filepath.Join(name, "views", "home", "index.html"):          "app/views/home/index.html.tmpl",
	}

	for filePath, tmplName := range files {
		if err := g.renderTemplate(tmplName, filePath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filePath, err)
		}
	}

	fmt.Printf("✅ Generated app: %s\n", name)
	return nil
}

func (g *Generator) GenerateResource(name string, fieldArgs []string) error {
	fields := g.parseFields(fieldArgs)

	data := ResourceData{
		Name:       cases.Title(language.English).String(name),
		VarName:    strings.ToLower(name),
		TableName:  g.pluralize(strings.ToLower(name)),
		ViewPath:   g.pluralize(strings.ToLower(name)),
		RoutePath:  g.pluralize(strings.ToLower(name)),
		Fields:     fields,
		FirstField: g.getFirstStringField(fields),
		Timestamp:  time.Now().Format("20060102150405"),
	}

	// Create directories
	os.MkdirAll("models", 0755)
	os.MkdirAll("controllers", 0755)
	os.MkdirAll("db/migrations", 0755)
	os.MkdirAll(filepath.Join("views", data.ViewPath), 0755)

	// Generate files (models, controllers, migrations, views)
	files := map[string]string{
		filepath.Join("models", data.VarName+".go"):                                        "resource/model.go.tmpl",
		filepath.Join("controllers", data.VarName+"_controller.go"):                        "resource/controller.go.tmpl",
		filepath.Join("db", "migrations", data.Timestamp+"_create_"+data.TableName+".sql"): "resource/migration.sql.tmpl",
		filepath.Join("views", data.ViewPath, "index.html"):                                "resource/index.html.tmpl",
		filepath.Join("views", data.ViewPath, "show.html"):                                 "resource/show.html.tmpl",
		filepath.Join("views", data.ViewPath, "new.html"):                                  "resource/new.html.tmpl",
		filepath.Join("views", data.ViewPath, "edit.html"):                                 "resource/edit.html.tmpl",
	}

	for filePath, tmplName := range files {
		if err := g.renderTemplate(tmplName, filePath, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filePath, err)
		}
	}

	fmt.Printf("✅ Generated resource: %s\n", name)
	fmt.Printf("   - Model: models/%s.go\n", data.VarName)
	fmt.Printf("   - Controller: controllers/%s_controller.go\n", data.VarName)
	fmt.Printf("   - Migration: db/migrations/%s_create_%s.sql\n", data.Timestamp, data.TableName)
	fmt.Printf("   - Views: views/%s/\n", data.ViewPath)

	return nil
}

func (g *Generator) renderTemplate(tmplName, filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Extract just the filename from the template path
	parts := strings.Split(tmplName, "/")
	templateName := parts[len(parts)-1]

	return g.templates.ExecuteTemplate(file, templateName, data)
}

func (g *Generator) parseFields(fieldArgs []string) []Field {
	var fields []Field

	for _, arg := range fieldArgs {
		parts := strings.Split(arg, ":")
		if len(parts) != 2 {
			continue
		}

		name := parts[0]
		fieldType := parts[1]

		field := Field{
			Name:     cases.Title(language.English).String(name),
			DBName:   strings.ToLower(name),
			FormName: strings.ToLower(name),
			GoType:   g.mapToGoType(fieldType),
			SQLType:  g.mapToSQLType(fieldType),
			HTMLType: g.mapToHTMLType(fieldType),
		}

		fields = append(fields, field)
	}

	return fields
}

func (g *Generator) mapToGoType(dbType string) string {
	if goType, ok := g.typeMapping.GoTypes[dbType]; ok {
		return goType
	}
	return "string" // default fallback
}

func (g *Generator) mapToSQLType(goType string) string {
	if sqlType, ok := g.typeMapping.SQLTypes[goType]; ok {
		return sqlType
	}
	return "VARCHAR(255)" // default fallback
}

func (g *Generator) mapToHTMLType(goType string) string {
	if htmlType, ok := g.typeMapping.HTMLTypes[goType]; ok {
		return htmlType
	}
	return "text" // default fallback
}

func (g *Generator) pluralize(word string) string {
	// Enhanced pluralization rules
	switch {
	case strings.HasSuffix(word, "s"), strings.HasSuffix(word, "x"), strings.HasSuffix(word, "z"):
		return word + "es"
	case strings.HasSuffix(word, "ch"), strings.HasSuffix(word, "sh"):
		return word + "es"
	case strings.HasSuffix(word, "y"):
		// Check if preceded by consonant
		if len(word) > 1 && !isVowel(rune(word[len(word)-2])) {
			return word[:len(word)-1] + "ies"
		}
		return word + "s"
	case strings.HasSuffix(word, "f"):
		return word[:len(word)-1] + "ves"
	case strings.HasSuffix(word, "fe"):
		return word[:len(word)-2] + "ves"
	default:
		return word + "s"
	}
}

func isVowel(r rune) bool {
	return strings.ContainsRune("aeiouAEIOU", r)
}

func (g *Generator) getFirstStringField(fields []Field) string {
	for _, field := range fields {
		if field.GoType == "string" {
			return field.Name
		}
	}
	return "ID"
}
