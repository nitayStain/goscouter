package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"goscouter/internal/subdomain"
)

var DebugMode bool

// setupRouter wires routes for the HTTP server.
func setupRouter() *gin.Engine {
	// Set to release mode and disable console color
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	// Disable Gin's default logging
	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard

	// Create engine without default middleware
	r := gin.New()

	// Add recovery middleware (catches panics)
	r.Use(gin.Recovery())

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Create subdomain finder with debug mode if enabled
	finderOpts := []subdomain.FinderOption{
		subdomain.WithUserAgent("goscouter-backend/1.0"),
		subdomain.WithHTTPClient(&http.Client{Timeout: 20 * time.Second}),
	}

	if DebugMode {
		finderOpts = append(finderOpts, subdomain.WithDebug(true))
	}

	finder := subdomain.NewFinder(finderOpts...)

	// API routes
	api := r.Group("/api")
	api.GET("/subdomains", subdomainScanHandler(finder, 30*time.Second))

	// Get frontend path from environment or use default
	frontendPath := os.Getenv("GOSCOUTER_FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "frontend/out"
	}

	// Serve Next.js static assets (_next directory)
	r.Static("/_next", filepath.Join(frontendPath, "_next"))

	// Serve other static files
	r.Static("/static", filepath.Join(frontendPath, "static"))

	// Serve favicon
	r.StaticFile("/favicon.ico", filepath.Join(frontendPath, "favicon.ico"))

	// Serve index.html for root and all other routes (SPA fallback)
	indexPath := filepath.Join(frontendPath, "index.html")
	r.GET("/", func(c *gin.Context) {
		c.File(indexPath)
	})

	r.NoRoute(func(c *gin.Context) {
		c.File(indexPath)
	})

	return r
}
