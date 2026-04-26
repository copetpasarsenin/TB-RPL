package controller

import (
	"api-gateway/model"
	"api-gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LogController menangani endpoint untuk log request
type LogController struct {
	LogService *service.LogService
}

// NewLogController membuat instance baru
func NewLogController(ls *service.LogService) *LogController {
	return &LogController{LogService: ls}
}

// GetLogs mengambil semua log dengan filter dan pagination
// Route: GET /gateway/logs
func (lc *LogController) GetLogs(c *gin.Context) {
	var filter model.LogFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Parameter filter tidak valid",
			"error":   err.Error(),
		})
		return
	}

	result, err := lc.LogService.GetLogs(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil log",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Log berhasil diambil",
		"data":    result,
	})
}

// GetLogByRequestID mengambil log berdasarkan Request ID
// Route: GET /gateway/logs/:request_id
func (lc *LogController) GetLogByRequestID(c *gin.Context) {
	requestID := c.Param("request_id")

	log, err := lc.LogService.GetLogByRequestID(requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Log tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Log ditemukan",
		"data":    log,
	})
}

// GetStats mengambil statistik gateway
// Route: GET /gateway/stats
func (lc *LogController) GetStats(c *gin.Context) {
	stats, err := lc.LogService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil statistik",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Statistik gateway",
		"data":    stats,
	})
}
