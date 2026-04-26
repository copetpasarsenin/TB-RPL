package middleware

import (
	"api-gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware membatasi jumlah transaksi per user
// Sesuai Aturan Keuangan:
// - #15: Cooldown Transaksi = 10-30 detik
// - #16: Max Transaksi Harian = 10 transaksi
func RateLimitMiddleware(rateLimitService *service.RateLimitService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Hanya berlaku untuk request yang memiliki user (authenticated)
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userID := userIDVal.(uint)

		// Cek rate limit
		if err := rateLimitService.CheckRateLimit(userID); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Rate limit exceeded",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		c.Next()

		// Catat transaksi setelah berhasil (status 2xx)
		if c.Writer.Status() >= 200 && c.Writer.Status() < 300 {
			go rateLimitService.RecordTransaction(userID)
		}
	}
}
