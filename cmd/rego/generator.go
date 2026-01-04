package main

import (
	"os"
	"path/filepath"
	"text/template"
)

type AppTemplate struct {
	Name   string
	Module string
}

const appMainTemplate = `package main

import (
	"log"
	"net/http"
	
	"github.com/Palaciodiego008/rebololang/pkg/rebolo"
)

func main() {
	app := rebolo.New()
	
	// Routes
	app.GET("/", HomeHandler)
	
	// Static files
	app.router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	rebolo.Render(w, "home/index.html", map[string]string{
		"Title": "Welcome to {{.Name}}",
		"Framework": "ReboloLang",
	})
}
`

const configTemplate = `app:
  name: {{.Name}}
  env: development

server:
  port: 3000
  host: localhost

database:
  url: postgres://localhost/{{.Name}}_development

assets:
  hot_reload: true
`

const packageJsonTemplate = `{
  "name": "{{.Name}}",
  "version": "1.0.0",
  "scripts": {
    "dev": "bun --watch src/index.js",
    "build": "bun build src/index.js --outdir=public"
  },
  "devDependencies": {
    "bun": "latest"
  }
}
`

const bunIndexTemplate = `// {{.Name}} - Frontend Assets powered by ReboloLang
console.log('🚀 {{.Name}} loaded with ReboloLang!');

// Hot reload in development
if (process.env.NODE_ENV === 'development') {
  const eventSource = new EventSource('/dev/reload');
  eventSource.onmessage = () => {
    console.log('🔄 Hot reloading...');
    location.reload();
  };
}

// Add some basic styling
document.addEventListener('DOMContentLoaded', function() {
  const style = document.createElement('style');
  style.textContent = ` + "`" + `
    body { 
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
      margin: 0;
      padding: 2rem;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: white;
      min-height: 100vh;
    }
    h1 { font-size: 3rem; margin-bottom: 1rem; }
    p { font-size: 1.2rem; opacity: 0.9; }
    .container { max-width: 800px; margin: 0 auto; text-align: center; }
  ` + "`" + `;
  document.head.appendChild(style);
});
`

func generateApp(name string) {
	appData := AppTemplate{
		Name:   name,
		Module: "github.com/Palaciodiego008/" + name,
	}
	
	// Create directories
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
		os.MkdirAll(dir, 0755)
	}
	
	// Generate files
	files := map[string]string{
		filepath.Join(name, "main.go"):       appMainTemplate,
		filepath.Join(name, "config.yml"):   configTemplate,
		filepath.Join(name, "package.json"): packageJsonTemplate,
		filepath.Join(name, "src", "index.js"): bunIndexTemplate,
	}
	
	for path, tmplContent := range files {
		tmpl := template.Must(template.New("").Parse(tmplContent))
		file, _ := os.Create(path)
		tmpl.Execute(file, appData)
		file.Close()
	}
	
	// Create basic HTML templates
	createHTMLTemplates(name, appData)
}

func createHTMLTemplates(appName string, data AppTemplate) {
	layoutHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - {{.Framework}}</title>
    <script src="/public/index.js"></script>
</head>
<body>
    <div class="container">
        {{template "content" .}}
    </div>
</body>
</html>`

	homeHTML := `{{define "content"}}
<h1>🚀 {{.Title}}</h1>
<p>Your ReboloLang application is running successfully!</p>
<p>Framework: <strong>{{.Framework}}</strong></p>
<p>Inspired by Rebolo, Barranquilla, Colombia 🇨🇴</p>

<div style="margin-top: 2rem;">
    <h3>Next Steps:</h3>
    <ul style="text-align: left; display: inline-block;">
        <li>Generate resources: <code>rebololang generate resource posts title:string</code></li>
        <li>Add routes in <code>main.go</code></li>
        <li>Configure database in <code>config.yml</code></li>
        <li>Run migrations: <code>rebololang db migrate</code></li>
    </ul>
</div>
{{end}}`

	// Write layout
	layoutFile, _ := os.Create(filepath.Join(appName, "views", "layouts", "application.html"))
	tmpl := template.Must(template.New("").Parse(layoutHTML))
	tmpl.Execute(layoutFile, data)
	layoutFile.Close()
	
	// Write home view
	homeFile, _ := os.Create(filepath.Join(appName, "views", "home", "index.html"))
	homeFile.WriteString(homeHTML)
	homeFile.Close()
}
