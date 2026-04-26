package service

import (
	"api-gateway/config"
	"api-gateway/repository"
	"fmt"
	"time"
)

// RateLimitService menangani business logic rate limiting
// Sesuai Aturan Keuangan:
// - #15: Cooldown Transaksi = 10-30 detik
// - #16: Max Transaksi Harian = 10 transaksi
type RateLimitService struct {
	Repo *repository.RateLimitRepository
}

// NewRateLimitService membuat instance baru
func NewRateLimitService(repo *repository.RateLimitRepository) *RateLimitService {
	return &RateLimitService{Repo: repo}
}

// CheckRateLimit memeriksa apakah user boleh melakukan transaksi
func (s *RateLimitService) CheckRateLimit(userID uint) error {
	record, err := s.Repo.GetOrCreate(userID)
	if err != nil {
		return fmt.Errorf("failed to get rate limit record: %w", err)
	}

	// Cek max transaksi harian (10 transaksi per hari)
	maxDaily := 10
	if record.TransactionCount >= maxDaily {
		return fmt.Errorf("batas transaksi harian tercapai (%d/%d transaksi). Coba lagi besok", record.TransactionCount, maxDaily)
	}

	// Cek cooldown (10-30 detik antar transaksi)
	cooldown := time.Duration(config.AppConfig.CooldownSeconds) * time.Second
	timeSinceLast := time.Since(record.LastTransactionAt)
	if record.TransactionCount > 0 && timeSinceLast < cooldown {
		remaining := cooldown - timeSinceLast
		return fmt.Errorf("cooldown aktif. Tunggu %.0f detik lagi sebelum transaksi berikutnya", remaining.Seconds())
	}

	return nil
}

// RecordTransaction mencatat transaksi baru untuk rate limiting
func (s *RateLimitService) RecordTransaction(userID uint) error {
	return s.Repo.IncrementCount(userID)
}

// CleanupOldRecords membersihkan record rate limit yang sudah lewat hari
func (s *RateLimitService) CleanupOldRecords() error {
	return s.Repo.ResetDaily()
}
