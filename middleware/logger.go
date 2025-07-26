package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		userID, _ := c.Get("userID")
		log.Printf("[API] %s %s | %d | %v | userID=%v | IP=%s | UA=%s",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
			userID,
			c.ClientIP(),
			c.Request.UserAgent(),
		)
	}
}
