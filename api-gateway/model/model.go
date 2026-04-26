package model

import (
	"time"

	"gorm.io/gorm"
)

// =============================================
// REQUEST LOG - Mencatat semua request yang masuk
// =============================================

// RequestLog menyimpan log setiap request yang melalui API Gateway
type RequestLog struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	RequestID     string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"request_id"`
	Method        string         `gorm:"type:varchar(10);not null" json:"method"`
	Path          string         `gorm:"type:varchar(500);not null" json:"path"`
	TargetService string         `gorm:"type:varchar(100);not null" json:"target_service"`
	SourceIP      string         `gorm:"type:varchar(50)" json:"source_ip"`
	UserID        *uint          `gorm:"index" json:"user_id,omitempty"`
	StatusCode    int            `json:"status_code"`
	ResponseTime  int64          `json:"response_time_ms"`
	RequestBody   string         `gorm:"type:text" json:"request_body,omitempty"`
	ResponseBody  string         `gorm:"type:text" json:"response_body,omitempty"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// =============================================
// GATEWAY FEE TRANSACTION - Catat fee 0.5%
// =============================================

// GatewayFeeTransaction menyimpan catatan fee gateway 0.5% dari setiap transaksi
type GatewayFeeTransaction struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	TransactionID     string         `gorm:"type:varchar(36);uniqueIndex;not null" json:"transaction_id"`
	RequestLogID      uint           `gorm:"index" json:"request_log_id"`
	OriginalAmount    float64        `gorm:"not null" json:"original_amount"`
	FeePercent        float64        `gorm:"not null;default:0.5" json:"fee_percent"`
	FeeAmount         float64        `gorm:"not null" json:"fee_amount"`
	SourceService     string         `gorm:"type:varchar(100);not null" json:"source_service"`
	DestinationService string       `gorm:"type:varchar(100);not null" json:"destination_service"`
	UserID            uint           `gorm:"index" json:"user_id"`
	Status            string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, success, failed
	SmartBankRef      string         `gorm:"type:varchar(100)" json:"smartbank_ref,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// =============================================
// SERVICE REGISTRY - Daftar service yang terdaftar
// =============================================

// ServiceRegistry menyimpan daftar service/aplikasi yang terdaftar di gateway
type ServiceRegistry struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	BaseURL   string         `gorm:"type:varchar(500);not null" json:"base_url"`
	Status    string         `gorm:"type:varchar(20);not null;default:'active'" json:"status"` // active, inactive, maintenance
	APIKey    string         `gorm:"type:varchar(255)" json:"api_key,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// =============================================
// RATE LIMIT RECORD - Tracking rate limit per user
// =============================================

// RateLimitRecord mencatat jumlah transaksi user untuk rate limiting
type RateLimitRecord struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"index;not null" json:"user_id"`
	TransactionCount int      `gorm:"not null;default:0" json:"transaction_count"`
	LastTransactionAt time.Time `json:"last_transaction_at"`
	DateKey         string    `gorm:"type:varchar(10);index;not null" json:"date_key"` // format: 2026-04-26
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
