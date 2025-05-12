package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	lastSeen time.Time
	limiter  int
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimit(maxPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		now := time.Now()
		if !exists || now.Sub(v.lastSeen) > time.Minute {
			visitors[ip] = &visitor{lastSeen: now, limiter: 1}
			mu.Unlock()
			c.Next()
			return
		}
		if v.limiter >= maxPerMinute {
			mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"code": 1, "message": "请求过于频繁，请稍后再试", "data": nil})
			return
		}
		v.limiter++
		v.lastSeen = now
		mu.Unlock()
		c.Next()
	}
}
