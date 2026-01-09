package server

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"goscouter/backend/internal/subdomain"
)

// setupRouter wires routes for the HTTP server.
func setupRouter() *gin.Engine {
	r := gin.Default()
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"join": strings.Join,
	}).ParseFS(templatesFS, "templates/*.html"))
	r.SetHTMLTemplate(tmpl)

	finder := subdomain.NewFinder(
		subdomain.WithUserAgent("goscouter-backend/1.0"),
		subdomain.WithHTTPClient(&http.Client{Timeout: 20 * time.Second}),
	)

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api := r.Group("/api")
	api.GET("/subdomains", subdomainScanHandler(finder, 30*time.Second))

	ui := r.Group("/ui")
	ui.GET("/subdomains", subdomainScanUIHandler(finder, 30*time.Second))

	return r
}
