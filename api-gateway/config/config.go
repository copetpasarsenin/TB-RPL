package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config menyimpan semua konfigurasi aplikasi
type Config struct {
	AppPort string
	AppEnv  string

	DBDriver   string // "sqlite" atau "postgres"
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string

	GatewayFeePercent float64

	SmartBankURL      string
	MarketplaceURL    string
	PosURL            string
	LogistiKitaURL    string
	SupplierHubURL    string
	UMKMInsightURL    string

	RateLimitPerSecond int
	CooldownSeconds    int

	LogLevel string
	LogFile  string
}

var AppConfig Config

// LoadConfig memuat konfigurasi dari file .env
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	feePercent, _ := strconv.ParseFloat(getEnv("GATEWAY_FEE_PERCENT", "0.5"), 64)
	rateLimit, _ := strconv.Atoi(getEnv("RATE_LIMIT_PER_SECOND", "10"))
	cooldown, _ := strconv.Atoi(getEnv("COOLDOWN_SECONDS", "10"))

	AppConfig = Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBDriver:   getEnv("DB_DRIVER", "sqlite"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "api_gateway"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret: getEnv("JWT_SECRET", "your-super-secret-key"),

		GatewayFeePercent: feePercent,

		SmartBankURL:   getEnv("SMARTBANK_BASE_URL", "http://localhost:8081"),
		MarketplaceURL: getEnv("MARKETPLACE_BASE_URL", "http://localhost:8082"),
		PosURL:         getEnv("POS_BASE_URL", "http://localhost:8083"),
		LogistiKitaURL: getEnv("LOGISTIKITA_BASE_URL", "http://localhost:8084"),
		SupplierHubURL: getEnv("SUPPLIERHUB_BASE_URL", "http://localhost:8085"),
		UMKMInsightURL: getEnv("UMKM_INSIGHT_BASE_URL", "http://localhost:8086"),

		RateLimitPerSecond: rateLimit,
		CooldownSeconds:    cooldown,

		LogLevel: getEnv("LOG_LEVEL", "debug"),
		LogFile:  getEnv("LOG_FILE", "logs/gateway.log"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ConnectDatabase menginisialisasi koneksi database (SQLite atau PostgreSQL)
func ConnectDatabase() *gorm.DB {
	logMode := logger.Silent
	if AppConfig.AppEnv == "development" {
		logMode = logger.Info
	}

	var db *gorm.DB
	var err error

	if AppConfig.DBDriver == "postgres" {
		// Mode PostgreSQL (untuk production / integrasi)
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			AppConfig.DBHost,
			AppConfig.DBUser,
			AppConfig.DBPassword,
			AppConfig.DBName,
			AppConfig.DBPort,
			AppConfig.DBSSLMode,
		)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})
	} else {
		// Mode SQLite (untuk testing lokal, tanpa install database)
		db, err = gorm.Open(sqlite.Open("api_gateway.db"), &gorm.Config{
			Logger: logger.Default.LogMode(logMode),
		})
	}

	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	log.Printf("✅ Database connected successfully (driver: %s)", AppConfig.DBDriver)
	return db
}
