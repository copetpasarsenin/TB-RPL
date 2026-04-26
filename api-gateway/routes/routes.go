package routes

import (
	"api-gateway/controller"
	"api-gateway/middleware"
	"api-gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mengatur semua route API Gateway
func SetupRoutes(
	router *gin.Engine,
	proxyCtrl *controller.ProxyController,
	feeCtrl *controller.FeeController,
	logCtrl *controller.LogController,
	serviceCtrl *controller.ServiceController,
	dashboardCtrl *controller.DashboardController,
	logService *service.LogService,
	rateLimitService *service.RateLimitService,
) {

	// =============================================
	// GLOBAL MIDDLEWARE
	// =============================================
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggingMiddleware(logService))

	// =============================================
	// PUBLIC ROUTES (Tanpa JWT)
	// =============================================

	// Health check gateway
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "API Gateway is running",
			"service": "api-gateway",
		})
	})

	// Generate JWT token (untuk testing/demo)
	router.POST("/auth/token", proxyCtrl.GenerateToken)

	// =============================================
	// GATEWAY MANAGEMENT ROUTES (JWT Required)
	// =============================================
	gateway := router.Group("/gateway")
	gateway.Use(middleware.JWTAuthMiddleware())
	{
		// Dashboard
		gateway.GET("/dashboard", dashboardCtrl.GetDashboard)

		// Stats
		gateway.GET("/stats", logCtrl.GetStats)

		// Logs
		gateway.GET("/logs", logCtrl.GetLogs)
		gateway.GET("/logs/:request_id", logCtrl.GetLogByRequestID)

		// Fee Management
		gateway.POST("/fee/calculate", feeCtrl.CalculateFee)
		gateway.GET("/fee/stats", feeCtrl.GetFeeStats)
		gateway.GET("/fee/:transaction_id", feeCtrl.GetFeeByTransaction)
		gateway.GET("/fees", feeCtrl.GetAllFees)
		gateway.PUT("/fee/:transaction_id/status", feeCtrl.UpdateFeeStatus)

		// Service Registry
		gateway.GET("/services", serviceCtrl.GetAllServices)
		gateway.GET("/services/:name", serviceCtrl.GetServiceByName)
		gateway.POST("/services", serviceCtrl.RegisterService)
		gateway.PUT("/services/:name/status", serviceCtrl.UpdateServiceStatus)
		gateway.DELETE("/services/:name", serviceCtrl.DeleteService)

		// Health Check Services
		gateway.GET("/health", serviceCtrl.HealthCheckAll)
		gateway.GET("/health/:name", serviceCtrl.HealthCheckService)
	}

	// =============================================
	// PROXY ROUTES - Forward ke service lain (JWT + Rate Limit)
	// Sesuai Aturan #5: Semua komunikasi antar-app wajib melalui gateway
	// =============================================
	api := router.Group("/api")
	api.Use(middleware.JWTAuthMiddleware())
	api.Use(middleware.RateLimitMiddleware(rateLimitService))
	{
		// Proxy ke semua service: /api/{service_name}/{path}
		// Contoh:
		//   /api/smartbank/saldo      → SmartBank /saldo
		//   /api/marketplace/products → Marketplace /products
		//   /api/pos/checkout         → POS /checkout
		//   /api/logistikita/shipment → LogistiKita /shipment
		//   /api/supplierhub/materials→ SupplierHub /materials
		//   /api/umkm-insight/reports → UMKM Insight /reports
		api.Any("/:service/*path", proxyCtrl.ProxyToService)
	}
}
