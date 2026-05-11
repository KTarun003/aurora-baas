package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime)

		// Log the request details
		method := c.Request.Method
		path := c.Request.URL.Path
		status := c.Writer.Status()
		log.Printf("%s %s %d %v", method, path, status, duration)
	}
}
