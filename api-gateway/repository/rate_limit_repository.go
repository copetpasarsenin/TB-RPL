package repository

import (
	"api-gateway/model"
	"time"

	"gorm.io/gorm"
)

// RateLimitRepository menangani operasi database untuk rate limiting
type RateLimitRepository struct {
	DB *gorm.DB
}

// NewRateLimitRepository membuat instance baru
func NewRateLimitRepository(db *gorm.DB) *RateLimitRepository {
	return &RateLimitRepository{DB: db}
}

// GetOrCreate mengambil atau membuat record rate limit untuk user pada hari ini
func (r *RateLimitRepository) GetOrCreate(userID uint) (*model.RateLimitRecord, error) {
	dateKey := time.Now().Format("2006-01-02")
	var record model.RateLimitRecord

	err := r.DB.Where("user_id = ? AND date_key = ?", userID, dateKey).First(&record).Error
	if err == gorm.ErrRecordNotFound {
		record = model.RateLimitRecord{
			UserID:            userID,
			TransactionCount:  0,
			LastTransactionAt: time.Now(),
			DateKey:           dateKey,
		}
		if createErr := r.DB.Create(&record).Error; createErr != nil {
			return nil, createErr
		}
		return &record, nil
	}
	return &record, err
}

// IncrementCount menambah counter transaksi
func (r *RateLimitRepository) IncrementCount(userID uint) error {
	dateKey := time.Now().Format("2006-01-02")
	return r.DB.Model(&model.RateLimitRecord{}).
		Where("user_id = ? AND date_key = ?", userID, dateKey).
		Updates(map[string]interface{}{
			"transaction_count":   gorm.Expr("transaction_count + 1"),
			"last_transaction_at": time.Now(),
		}).Error
}

// ResetDaily menghapus record rate limit yang sudah lewat hari
func (r *RateLimitRepository) ResetDaily() error {
	dateKey := time.Now().Format("2006-01-02")
	return r.DB.Where("date_key < ?", dateKey).Delete(&model.RateLimitRecord{}).Error
}
