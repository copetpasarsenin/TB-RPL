package service

import (
	"api-gateway/model"
	"api-gateway/repository"
	"math"
)

// LogService menangani business logic untuk logging
type LogService struct {
	Repo *repository.RequestLogRepository
}

// NewLogService membuat instance baru LogService
func NewLogService(repo *repository.RequestLogRepository) *LogService {
	return &LogService{Repo: repo}
}

// SaveLog menyimpan log request
func (s *LogService) SaveLog(log *model.RequestLog) error {
	return s.Repo.Create(log)
}

// GetLogByID mengambil log berdasarkan ID
func (s *LogService) GetLogByID(id uint) (*model.RequestLog, error) {
	return s.Repo.FindByID(id)
}

// GetLogByRequestID mengambil log berdasarkan Request ID
func (s *LogService) GetLogByRequestID(requestID string) (*model.RequestLog, error) {
	return s.Repo.FindByRequestID(requestID)
}

// GetLogs mengambil semua log dengan filter dan pagination
func (s *LogService) GetLogs(filter model.LogFilter) (*model.PaginatedResponse, error) {
	// Set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PerPage <= 0 {
		filter.PerPage = 20
	}

	logs, total, err := s.Repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.PerPage)))

	return &model.PaginatedResponse{
		Data:       logs,
		Page:       filter.Page,
		PerPage:    filter.PerPage,
		TotalData:  total,
		TotalPages: totalPages,
	}, nil
}

// GetStats mengambil statistik log
func (s *LogService) GetStats() (*model.GatewayStats, error) {
	return s.Repo.GetStats()
}
