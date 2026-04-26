package repository

import (
	"api-gateway/model"
	"time"

	"gorm.io/gorm"
)

// RequestLogRepository menangani operasi database untuk log request
type RequestLogRepository struct {
	DB *gorm.DB
}

// NewRequestLogRepository membuat instance baru RequestLogRepository
func NewRequestLogRepository(db *gorm.DB) *RequestLogRepository {
	return &RequestLogRepository{DB: db}
}

// Create menyimpan log request baru ke database
func (r *RequestLogRepository) Create(log *model.RequestLog) error {
	return r.DB.Create(log).Error
}

// FindByID mencari log berdasarkan ID
func (r *RequestLogRepository) FindByID(id uint) (*model.RequestLog, error) {
	var log model.RequestLog
	err := r.DB.First(&log, id).Error
	return &log, err
}

// FindByRequestID mencari log berdasarkan Request ID (UUID)
func (r *RequestLogRepository) FindByRequestID(requestID string) (*model.RequestLog, error) {
	var log model.RequestLog
	err := r.DB.Where("request_id = ?", requestID).First(&log).Error
	return &log, err
}

// FindAll mengambil semua log dengan filter dan pagination
func (r *RequestLogRepository) FindAll(filter model.LogFilter) ([]model.RequestLog, int64, error) {
	var logs []model.RequestLog
	var total int64

	query := r.DB.Model(&model.RequestLog{})

	// Apply filters
	if filter.Service != "" {
		query = query.Where("target_service = ?", filter.Service)
	}
	if filter.Method != "" {
		query = query.Where("method = ?", filter.Method)
	}
	if filter.Status != "" {
		if filter.Status == "success" {
			query = query.Where("status_code >= 200 AND status_code < 300")
		} else if filter.Status == "error" {
			query = query.Where("status_code >= 400")
		}
	}
	if filter.DateFrom != "" {
		query = query.Where("created_at >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query = query.Where("created_at <= ?", filter.DateTo+" 23:59:59")
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (filter.Page - 1) * filter.PerPage
	err := query.Order("created_at DESC").
		Offset(offset).
		Limit(filter.PerPage).
		Find(&logs).Error

	return logs, total, err
}

// GetStats mengambil statistik request
func (r *RequestLogRepository) GetStats() (*model.GatewayStats, error) {
	stats := &model.GatewayStats{}

	// Total requests
	r.DB.Model(&model.RequestLog{}).Count(&stats.TotalRequests)

	// Average response time
	r.DB.Model(&model.RequestLog{}).
		Select("COALESCE(AVG(response_time), 0)").
		Scan(&stats.AvgResponseTimeMs)

	// Success rate
	var successCount int64
	r.DB.Model(&model.RequestLog{}).
		Where("status_code >= 200 AND status_code < 300").
		Count(&successCount)

	if stats.TotalRequests > 0 {
		stats.SuccessRate = float64(successCount) / float64(stats.TotalRequests) * 100
	}

	return stats, nil
}

// DeleteOlderThan menghapus log yang lebih lama dari durasi tertentu
func (r *RequestLogRepository) DeleteOlderThan(duration time.Duration) error {
	cutoff := time.Now().Add(-duration)
	return r.DB.Where("created_at < ?", cutoff).Delete(&model.RequestLog{}).Error
}
