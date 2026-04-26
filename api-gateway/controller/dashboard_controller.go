package controller

import (
	"api-gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DashboardController menangani endpoint dashboard gateway
type DashboardController struct {
	LogService      *service.LogService
	FeeService      *service.FeeService
	ServiceRegistry *service.ServiceRegistryService
}

// NewDashboardController membuat instance baru
func NewDashboardController(ls *service.LogService, fs *service.FeeService, srs *service.ServiceRegistryService) *DashboardController {
	return &DashboardController{
		LogService:      ls,
		FeeService:      fs,
		ServiceRegistry: srs,
	}
}

// GetDashboard mengambil semua statistik untuk dashboard
// Route: GET /gateway/dashboard
func (dc *DashboardController) GetDashboard(c *gin.Context) {
	// Statistik request
	stats, _ := dc.LogService.GetStats()

	// Total fee
	totalFee, _ := dc.FeeService.GetTotalFeeCollected()
	totalTx, _ := dc.FeeService.GetTotalTransactions()

	// Active services
	activeServices, _ := dc.ServiceRegistry.CountActive()

	// Health check
	healthResults := dc.ServiceRegistry.HealthCheckAll()

	// Gabungkan
	if stats != nil {
		stats.TotalFeeCollected = totalFee
		stats.TotalTransactions = totalTx
		stats.ActiveServices = int(activeServices)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Dashboard API Gateway",
		"data": gin.H{
			"stats":           stats,
			"service_health":  healthResults,
		},
	})
}
