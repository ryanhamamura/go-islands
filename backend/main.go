package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/torenware/vite-go"
)

//go:embed templates
var templateFS embed.FS

// Templates holds our HTML templates
var Templates *template.Template

// PageData contains data passed to templates
type PageData struct {
	Title       string
	CurrentTime string
	InitialData map[string]interface{}
	ENV         string
	Vite        *vite.ViteAsset
}

// Helper for getting environment variables with defaults
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper for getting boolean environment variables
func getEnvBool(key string, fallback bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		return strings.ToLower(value) == "true" || value == "1"
	}
	return fallback
}

// Middleware for setting cache control headers on static assets
func cacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		
		// Set different cache times based on file type
		if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") {
			// Long cache for JS and CSS (they have content hashes in production)
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") || 
			strings.HasSuffix(path, ".svg") || strings.HasSuffix(path, ".webp") {
			// Images can be cached but not as long
			w.Header().Set("Cache-Control", "public, max-age=86400")
		} else {
			// Short cache for everything else
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		
		next.ServeHTTP(w, r)
	})
}

// Structured application configuration
type Config struct {
	Env         string
	Port        string
	ServerURL   string
	DevServerURL string
	AssetsPath  string
	LogRequests bool
}

func main() {
	// Setup application configuration
	config := Config{
		Env:         getEnv("GO_ENV", "development"),
		Port:        getEnv("PORT", "8080"),
		ServerURL:   getEnv("SERVER_URL", "http://localhost:8080"),
		DevServerURL: getEnv("DEV_SERVER_URL", "http://localhost:5173"),
		AssetsPath:  getEnv("ASSETS_PATH", "./static"),
		LogRequests: getEnvBool("LOG_REQUESTS", true),
	}
	
	// Initialize a structured logger
	logger := log.New(os.Stdout, "[go-react-islands] ", log.LstdFlags)
	
	// Setup templates with custom functions
	templateFuncs := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("Jan 02, 2006 15:04:05")
		},
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	
	Templates = template.Must(template.New("").Funcs(templateFuncs).ParseFS(templateFS, "templates/*.html"))

	r := chi.NewRouter()

	// Vite asset configuration
	viteConfig := &vite.ViteConfig{
		Environment: config.Env,
		FS:          nil, // Use the OS filesystem
	}

	if config.Env == "development" {
		// In development mode, use the dev server
		viteConfig.DevServerURL = config.DevServerURL
	} else {
		// In production mode, use the manifest file
		viteConfig.EntryPoint = "src/islands-client.js"
		viteConfig.Platform = "react"
		viteConfig.AssetsPath = config.AssetsPath
		viteConfig.ManifestPath = filepath.Join(config.AssetsPath, "manifest.json")
	}

	// Create the Vite asset provider
	viteAsset, err := vite.NewViteAsset(viteConfig)
	if err != nil {
		log.Fatalf("Failed to create Vite asset: %v", err)
	}

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(middleware.Compress(5))

	// Setup CORS based on environment
	corsOptions := cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}
	
	if config.Env == "development" {
		corsOptions.AllowedOrigins = []string{config.DevServerURL}
	} else {
		// In production, specify your actual domain or use AllowOriginFunc for more control
		serverURL := config.ServerURL
		if serverURL != "" {
			corsOptions.AllowedOrigins = []string{serverURL}
		} else {
			corsOptions.AllowedOrigins = []string{"https://yourdomain.com"}
		}
		// Alternatively use a function to validate origins
		// corsOptions.AllowOriginFunc = func(r *http.Request, origin string) bool {
		//    return strings.HasSuffix(origin, ".yourdomain.com")
		// }
	}
	
	r.Use(cors.Handler(corsOptions))

	// Add security headers
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			// In production you would want to set a proper CSP
			if config.Env != "development" {
				w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline';")
			}
			next.ServeHTTP(w, r)
		})
	})

	// Static assets - in production, these will be served from the built Vite output
	staticPath := config.AssetsPath
	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		logger.Printf("Static directory %s not found, creating it", staticPath)
		os.MkdirAll(staticPath, 0755)
	}
	
	// Improved static file server with caching headers
	fileServer := http.FileServer(http.Dir(staticPath))
	r.With(cacheControl).Handle("/assets/*", http.StripPrefix("/assets/", fileServer))

	// API response struct for consistent response format
	type APIResponse struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}

	// Helper to send JSON responses
	sendJSONResponse := func(w http.ResponseWriter, statusCode int, success bool, data interface{}, errMsg string) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		
		resp := APIResponse{
			Success: success,
			Data:    data,
			Error:   errMsg,
		}
		
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}
	}

	// API routes with proper error handling and validation
	r.Route("/api", func(r chi.Router) {
		// Add basic rate limiting for API routes (simple example)
		r.Use(middleware.Throttle(100))
		
		r.Get("/time", func(w http.ResponseWriter, r *http.Request) {
			data := map[string]string{"time": time.Now().Format(time.RFC3339)}
			sendJSONResponse(w, http.StatusOK, true, data, "")
		})

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			
			// Basic validation
			if id == "" {
				sendJSONResponse(w, http.StatusBadRequest, false, nil, "User ID is required")
				return
			}
			
			// In a real app, you would fetch from a database
			// This is just a demo with better error handling
			user := map[string]string{
				"id":    id,
				"name":  "John Doe",
				"email": "john@example.com",
				"role":  "Developer",
			}
			
			sendJSONResponse(w, http.StatusOK, true, user, "")
		})
		
		// Example handling errors
		r.Get("/error-demo", func(w http.ResponseWriter, r *http.Request) {
			sendJSONResponse(w, http.StatusInternalServerError, false, nil, "This is a demonstration of error handling")
		})
	})

	// Page routes
	r.Get("/", servePage("home.html", "React Islands Demo", config.Env, viteAsset, logger))
	r.Get("/about", servePage("about.html", "About - React Islands", config.Env, viteAsset, logger))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","version":"1.0.0","time":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Start the server
	addr := fmt.Sprintf(":%s", config.Port)
	logger.Printf("Server starting on port %s in %s mode", config.Port, config.Env)
	
	// Create server with reasonable timeouts
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	logger.Fatal(server.ListenAndServe())
}

// servePage returns a handler function that renders the specified template
func servePage(templateName, title, env string, viteAsset *vite.ViteAsset, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set content type and necessary headers
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		
		// Add reasonable cache policy for HTML pages
		if env != "development" {
			w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		}
		
		// Prepare template data
		data := PageData{
			Title:       title,
			CurrentTime: time.Now().Format(time.RFC3339),
			ENV:         env,
			Vite:        viteAsset,
			InitialData: map[string]interface{}{
				"user": map[string]interface{}{
					"id":   123,
					"name": "John Doe",
				},
			},
		}

		// Execute template with proper error handling
		if err := Templates.ExecuteTemplate(w, templateName, data); err != nil {
			logger.Printf("Template error: %v", err)
			
			// Only show detailed errors in development
			if env == "development" {
				http.Error(w, fmt.Sprintf("Template Error: %v", err), http.StatusInternalServerError)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}
	}
}