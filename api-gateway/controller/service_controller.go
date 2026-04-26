package controller

import (
	"api-gateway/model"
	"api-gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceController menangani endpoint untuk manajemen service registry
type ServiceController struct {
	ServiceRegistryService *service.ServiceRegistryService
}

// NewServiceController membuat instance baru
func NewServiceController(srs *service.ServiceRegistryService) *ServiceController {
	return &ServiceController{ServiceRegistryService: srs}
}

// GetAllServices mengambil semua service yang terdaftar
// Route: GET /gateway/services
func (sc *ServiceController) GetAllServices(c *gin.Context) {
	services, err := sc.ServiceRegistryService.GetAllServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data service",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Daftar service",
		"data":    services,
	})
}

// GetServiceByName mengambil detail service berdasarkan nama
// Route: GET /gateway/services/:name
func (sc *ServiceController) GetServiceByName(c *gin.Context) {
	name := c.Param("name")

	svc, err := sc.ServiceRegistryService.GetServiceByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Service tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail service",
		"data":    svc,
	})
}

// RegisterService mendaftarkan service baru
// Route: POST /gateway/services
func (sc *ServiceController) RegisterService(c *gin.Context) {
	var svc model.ServiceRegistry
	if err := c.ShouldBindJSON(&svc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	if svc.Status == "" {
		svc.Status = "active"
	}

	if err := sc.ServiceRegistryService.RegisterService(&svc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mendaftarkan service",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Service berhasil didaftarkan",
		"data":    svc,
	})
}

// UpdateServiceStatus mengupdate status service
// Route: PUT /gateway/services/:name/status
func (sc *ServiceController) UpdateServiceStatus(c *gin.Context) {
	name := c.Param("name")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
		})
		return
	}

	validStatuses := map[string]bool{"active": true, "inactive": true, "maintenance": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Status harus: active, inactive, atau maintenance",
		})
		return
	}

	if err := sc.ServiceRegistryService.UpdateServiceStatus(name, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengupdate status service",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status service berhasil diupdate",
	})
}

// DeleteService menghapus service
// Route: DELETE /gateway/services/:name
func (sc *ServiceController) DeleteService(c *gin.Context) {
	name := c.Param("name")

	if err := sc.ServiceRegistryService.DeleteService(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus service",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Service berhasil dihapus",
	})
}

// HealthCheckAll memeriksa kesehatan semua service
// Route: GET /gateway/health
func (sc *ServiceController) HealthCheckAll(c *gin.Context) {
	results := sc.ServiceRegistryService.HealthCheckAll()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Health check semua service",
		"data":    results,
	})
}

// HealthCheckService memeriksa kesehatan satu service
// Route: GET /gateway/health/:name
func (sc *ServiceController) HealthCheckService(c *gin.Context) {
	name := c.Param("name")

	svc, err := sc.ServiceRegistryService.GetServiceByName(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Service tidak ditemukan",
		})
		return
	}

	health := sc.ServiceRegistryService.HealthCheck(svc.Name, svc.BaseURL)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Health check service",
		"data":    health,
	})
}
