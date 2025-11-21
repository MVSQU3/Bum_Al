package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimited struct {
	limiter *rate.Limiter
}

func NewRateLimited(requestsPerSecond rate.Limit, burst int) *RateLimited {
	return &RateLimited{
		limiter: rate.NewLimiter(requestsPerSecond, burst),
	}
}

// Allow vérifie si la requête est autorisée
func (rl *RateLimited) Allow() bool {
	return rl.limiter.Allow()
}

// Wait attend que le token soit disponible
func (rl *RateLimited) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}

// GinMiddleware crée un middleware Gin compatible
func (rl *RateLimited) RateLimite() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow() {
			c.JSON(429, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
