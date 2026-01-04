package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// HotReloadScript is the client-side JavaScript that polls for changes and reloads
const HotReloadScript = `
<script>
(function() {
	console.log('🔥 Rebolo hot reload enabled (polling mode)');
	
	let lastCheck = Date.now();
	
	async function checkForChanges() {
		try {
			const response = await fetch('/__rebolo__/changes');
			const data = await response.json();
			
			if (data.changed) {
				console.log('🔄 File changed detected!');
				console.log('⚡ Reloading page...');
				location.reload();
			}
		} catch (err) {
			console.error('Hot reload check error:', err);
		}
	}
	
	// Check for changes every second
	setInterval(checkForChanges, 1000);
	
	console.log('✅ Hot reload polling started');
})();
</script>
`

// responseWriter wraps http.ResponseWriter to capture and modify the response body
type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *responseWriter) Flush() {
	// Remove Content-Length as we're modifying the body
	rw.Header().Del("Content-Length")

	// Write actual headers
	for k, v := range rw.Header() {
		rw.ResponseWriter.Header()[k] = v
	}
	rw.ResponseWriter.WriteHeader(rw.statusCode)

	// Write body
	io.Copy(rw.ResponseWriter, rw.body)
}

// HotReloadMiddleware injects hot reload script into HTML responses in development mode
func HotReloadMiddleware(enabled bool, skipPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Skip SSE endpoints and other paths that need streaming
			for _, path := range skipPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Wrap response writer to capture output
			rw := newResponseWriter(w)

			// Call next handler
			next.ServeHTTP(rw, r)

			// Get content type
			contentType := rw.Header().Get("Content-Type")

			// Only inject script into HTML responses
			if strings.Contains(contentType, "text/html") {
				body := rw.body.String()

				// Inject script before </body>
				if idx := strings.LastIndex(body, "</body>"); idx != -1 {
					body = body[:idx] + HotReloadScript + body[idx:]
					rw.body.Reset()
					rw.body.WriteString(body)
				}
			}

			// Flush response
			rw.Flush()
		})
	}
}
