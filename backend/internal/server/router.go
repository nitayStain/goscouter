package server

import "github.com/gin-gonic/gin"

// setupRouter wires routes for the HTTP server.
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Go + Gin!"})
	})

	return r
}
