package middleware

import (
	"api-gateway/model"
	"api-gateway/service"
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// responseWriter adalah custom response writer untuk menangkap response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware mencatat setiap request dan response ke database
// Sesuai Aturan #6: Logging wajib pada setiap request
func LoggingMiddleware(logService *service.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate Request ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Capture request body
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// Reset body agar bisa dibaca lagi oleh handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Wrap response writer untuk menangkap response
		wrappedWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = wrappedWriter

		// Record start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(startTime).Milliseconds()

		// Get target service dari context
		targetService, _ := c.Get("target_service")
		targetServiceStr := ""
		if targetService != nil {
			targetServiceStr = targetService.(string)
		}

		// Get user ID jika ada
		var userID *uint
		if uid, exists := c.Get("user_id"); exists {
			id := uid.(uint)
			userID = &id
		}

		// Limit response body length untuk storage
		responseBody := wrappedWriter.body.String()
		if len(responseBody) > 2000 {
			responseBody = responseBody[:2000] + "...[truncated]"
		}
		if len(requestBody) > 2000 {
			requestBody = requestBody[:2000] + "...[truncated]"
		}

		// Get error message
		errorMsg := ""
		if len(c.Errors) > 0 {
			errorMsg = c.Errors.String()
		}

		// Save log ke database (async)
		go func() {
			log := &model.RequestLog{
				RequestID:     requestID,
				Method:        c.Request.Method,
				Path:          c.Request.URL.Path,
				TargetService: targetServiceStr,
				SourceIP:      c.ClientIP(),
				UserID:        userID,
				StatusCode:    c.Writer.Status(),
				ResponseTime:  responseTime,
				RequestBody:   requestBody,
				ResponseBody:  responseBody,
				ErrorMessage:  errorMsg,
			}
			logService.SaveLog(log)
		}()
	}
}

// CORSMiddleware menangani Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// RecoveryMiddleware menangani panic dan mengembalikan response error
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "Internal server error",
					"error":   "An unexpected error occurred",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
