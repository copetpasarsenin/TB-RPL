package controller

import (
	"api-gateway/middleware"
	"api-gateway/service"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ProxyController menangani routing/proxy request ke service lain
// Sesuai Aturan #5: Semua komunikasi antar-app wajib melalui gateway
type ProxyController struct {
	ServiceRegistry *service.ServiceRegistryService
	FeeService      *service.FeeService
}

// NewProxyController membuat instance baru
func NewProxyController(sr *service.ServiceRegistryService, fs *service.FeeService) *ProxyController {
	return &ProxyController{
		ServiceRegistry: sr,
		FeeService:      fs,
	}
}

// ProxyToService meneruskan request ke service yang dituju
// Route: ANY /api/:service/*path
func (pc *ProxyController) ProxyToService(c *gin.Context) {
	serviceName := c.Param("service")
	path := c.Param("path")

	// Set target service untuk logging
	c.Set("target_service", serviceName)

	// Resolve service URL
	baseURL, err := pc.ServiceRegistry.GetServiceURL(serviceName)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Service '%s' tidak ditemukan atau tidak aktif", serviceName),
			"error":   err.Error(),
		})
		return
	}

	// Bangun target URL
	targetURL := baseURL + path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Baca request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
	}

	// Buat proxy request
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat request ke service target",
			"error":   err.Error(),
		})
		return
	}

	// Copy headers dari request asli
	for key, values := range c.Request.Header {
		for _, value := range values {
			proxyReq.Header.Add(key, value)
		}
	}

	// Tambahkan header gateway
	proxyReq.Header.Set("X-Gateway-Forwarded", "true")
	proxyReq.Header.Set("X-Request-ID", c.GetString("request_id"))
	if userID, exists := c.Get("user_id"); exists {
		proxyReq.Header.Set("X-User-ID", fmt.Sprintf("%d", userID))
	}
	if username, exists := c.Get("username"); exists {
		proxyReq.Header.Set("X-Username", username.(string))
	}

	// Kirim request ke service target
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": fmt.Sprintf("Gagal menghubungi service '%s'", serviceName),
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// Baca response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membaca response dari service target",
		})
		return
	}

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Kirim response ke client
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

// GenerateToken membuat JWT token untuk testing/demo
// Route: POST /auth/token
func (pc *ProxyController) GenerateToken(c *gin.Context) {
	var req struct {
		UserID   uint   `json:"user_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	token, err := middleware.GenerateJWT(req.UserID, req.Username, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token berhasil dibuat",
		"data": gin.H{
			"token":      token,
			"type":       "Bearer",
			"expires_in": "24h",
		},
	})
}
