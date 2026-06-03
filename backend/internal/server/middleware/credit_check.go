package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// CreditCheckMiddleware blocks API gateway requests when the authenticated user
// has no credits remaining. It reads the userID set by APIKeyAuthMiddleware and
// checks the Redis-cached balance via CreditService. On cache/DB failure it
// fails open (allows the request) to avoid blocking legitimate traffic.
func CreditCheckMiddleware(creditSvc *service.CreditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.Next()
			return
		}
		uid, ok := userID.(int64)
		if !ok || uid == 0 {
			c.Next()
			return
		}

		sufficient, err := creditSvc.HasSufficientCredits(c.Request.Context(), uid)
		if err != nil {
			slog.Warn("credit_check: failed to check credits, failing open", "userID", uid, "err", err)
			c.Next()
			return
		}
		if !sufficient {
			c.AbortWithStatusJSON(http.StatusPaymentRequired, gin.H{
				"type": "error",
				"error": gin.H{
					"type":    "quota_exceeded",
					"message": "积分余额不足，请充值后继续使用",
				},
			})
			return
		}
		c.Next()
	}
}
