package errors

import (
	"fmt"
	"net/http"
)

// ErrorHandler is a function that handles HTTP errors
type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, code int)

// ErrorHandlers stores custom error handlers by status code
type ErrorHandlers map[int]ErrorHandler

// NewErrorHandlers creates a new ErrorHandlers with defaults
func NewErrorHandlers() ErrorHandlers {
	handlers := make(ErrorHandlers)

	// Default 404 handler
	handlers[404] = func(w http.ResponseWriter, r *http.Request, err error, code int) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(404)
		html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>404 - Página no encontrada</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%);
            color: white;
            margin: 0;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }
        .container {
            text-align: center;
            padding: 40px;
            background: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            backdrop-filter: blur(10px);
        }
        h1 { font-size: 6em; margin: 0; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        h2 { font-size: 2em; margin: 20px 0; }
        p { font-size: 1.2em; opacity: 0.9; }
        a {
            display: inline-block;
            margin-top: 20px;
            padding: 15px 30px;
            background: rgba(255,255,255,0.3);
            color: white;
            text-decoration: none;
            border-radius: 5px;
            transition: transform 0.2s;
        }
        a:hover { transform: translateY(-2px); }
        .path { font-family: monospace; opacity: 0.7; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>404</h1>
        <h2>🔍 Página no encontrada</h2>
        <p>La página que buscas no existe</p>
        <p class="path">%s</p>
        <a href="/">← Volver al inicio</a>
    </div>
</body>
</html>
`, r.URL.Path)
		w.Write([]byte(html))
	}

	// Default 500 handler
	handlers[500] = func(w http.ResponseWriter, r *http.Request, err error, code int) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(500)

		errorMsg := ""
		if err != nil {
			errorMsg = err.Error()
		}

		html := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>500 - Error del servidor</title>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%);
            color: white;
            margin: 0;
            padding: 0;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }
        .container {
            text-align: center;
            padding: 40px;
            background: rgba(255, 255, 255, 0.1);
            border-radius: 10px;
            backdrop-filter: blur(10px);
            max-width: 600px;
        }
        h1 { font-size: 6em; margin: 0; text-shadow: 2px 2px 4px rgba(0,0,0,0.3); }
        h2 { font-size: 2em; margin: 20px 0; }
        p { font-size: 1.2em; opacity: 0.9; }
        .error {
            background: rgba(0,0,0,0.3);
            padding: 15px;
            border-radius: 5px;
            margin: 20px 0;
            font-family: monospace;
            font-size: 0.9em;
            word-break: break-word;
        }
        a {
            display: inline-block;
            margin-top: 20px;
            padding: 15px 30px;
            background: rgba(255,255,255,0.3);
            color: white;
            text-decoration: none;
            border-radius: 5px;
            transition: transform 0.2s;
        }
        a:hover { transform: translateY(-2px); }
    </style>
</head>
<body>
    <div class="container">
        <h1>500</h1>
        <h2>⚠️ Error del servidor</h2>
        <p>Ha ocurrido un error inesperado</p>
        %s
        <a href="/">← Volver al inicio</a>
    </div>
</body>
</html>
`, func() string {
			if errorMsg != "" {
				return fmt.Sprintf(`<div class="error">%s</div>`, errorMsg)
			}
			return ""
		}())
		w.Write([]byte(html))
	}

	return handlers
}
