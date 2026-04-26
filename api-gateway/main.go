package main

import (
	"api-gateway/config"
	"api-gateway/controller"
	"api-gateway/model"
	"api-gateway/repository"
	"api-gateway/routes"
	"api-gateway/service"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// =============================================
	// 1. Load Configuration
	// =============================================
	config.LoadConfig()

	// Set Gin mode
	if config.AppConfig.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// =============================================
	// 2. Connect Database
	// =============================================
	db := config.ConnectDatabase()

	// =============================================
	// 3. Auto Migrate Models
	// =============================================
	log.Println("🔄 Running database migrations...")
	err := db.AutoMigrate(
		&model.RequestLog{},
		&model.GatewayFeeTransaction{},
		&model.ServiceRegistry{},
		&model.RateLimitRecord{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}
	log.Println("✅ Database migration completed")

	// =============================================
	// 4. Initialize Repositories
	// =============================================
	requestLogRepo := repository.NewRequestLogRepository(db)
	feeRepo := repository.NewFeeTransactionRepository(db)
	serviceRegistryRepo := repository.NewServiceRegistryRepository(db)
	rateLimitRepo := repository.NewRateLimitRepository(db)

	// =============================================
	// 5. Initialize Services
	// =============================================
	logService := service.NewLogService(requestLogRepo)
	feeService := service.NewFeeService(feeRepo)
	serviceRegistryService := service.NewServiceRegistryService(serviceRegistryRepo)
	rateLimitService := service.NewRateLimitService(rateLimitRepo)

	// =============================================
	// 6. Seed Default Services
	// =============================================
	serviceRegistryService.SeedDefaultServices()
	log.Println("✅ Default services seeded")

	// =============================================
	// 7. Initialize Controllers
	// =============================================
	proxyCtrl := controller.NewProxyController(serviceRegistryService, feeService)
	feeCtrl := controller.NewFeeController(feeService)
	logCtrl := controller.NewLogController(logService)
	serviceCtrl := controller.NewServiceController(serviceRegistryService)
	dashboardCtrl := controller.NewDashboardController(logService, feeService, serviceRegistryService)

	// =============================================
	// 8. Setup Router & Routes
	// =============================================
	router := gin.Default()
	routes.SetupRoutes(
		router,
		proxyCtrl,
		feeCtrl,
		logCtrl,
		serviceCtrl,
		dashboardCtrl,
		logService,
		rateLimitService,
	)

	// =============================================
	// 9. Create logs directory
	// =============================================
	os.MkdirAll("logs", os.ModePerm)

	// =============================================
	// 10. Start Server
	// =============================================
	port := config.AppConfig.AppPort
	fmt.Println("========================================")
	fmt.Println("🚀 API Gateway - Tugas Besar RPL 2")
	fmt.Println("========================================")
	fmt.Printf("📡 Server running on port: %s\n", port)
	fmt.Printf("🌍 Environment: %s\n", config.AppConfig.AppEnv)
	fmt.Printf("💰 Gateway Fee: %.1f%%\n", config.AppConfig.GatewayFeePercent)
	fmt.Println("========================================")
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET  /health                    - Health check")
	fmt.Println("   POST /auth/token                - Generate JWT token")
	fmt.Println("   GET  /gateway/dashboard          - Dashboard")
	fmt.Println("   GET  /gateway/stats              - Statistics")
	fmt.Println("   GET  /gateway/logs               - Request logs")
	fmt.Println("   POST /gateway/fee/calculate      - Calculate fee")
	fmt.Println("   GET  /gateway/services           - Service registry")
	fmt.Println("   GET  /gateway/health             - Health check all services")
	fmt.Println("   ANY  /api/:service/*path          - Proxy to services")
	fmt.Println("========================================")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
