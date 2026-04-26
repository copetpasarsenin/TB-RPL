package service

import (
	"api-gateway/config"
	"api-gateway/model"
	"api-gateway/repository"
	"fmt"
	"math"

	"github.com/google/uuid"
)

// FeeService menangani business logic untuk fee gateway (0.5%)
type FeeService struct {
	Repo *repository.FeeTransactionRepository
}

// NewFeeService membuat instance baru
func NewFeeService(repo *repository.FeeTransactionRepository) *FeeService {
	return &FeeService{Repo: repo}
}

// CalculateAndCreateFee menghitung fee 0.5% dan menyimpan transaksi
// Sesuai Aturan Keuangan #10: Fee Gateway = 0.5% dipotong dari transaksi via API
func (s *FeeService) CalculateAndCreateFee(req model.TransactionRequest) (*model.TransactionFeeResponse, error) {
	feePercent := config.AppConfig.GatewayFeePercent
	feeAmount := req.Amount * (feePercent / 100)

	// Pembulatan 2 desimal
	feeAmount = math.Round(feeAmount*100) / 100
	totalAmount := req.Amount + feeAmount

	txID := uuid.New().String()

	feeTx := &model.GatewayFeeTransaction{
		TransactionID:      txID,
		OriginalAmount:     req.Amount,
		FeePercent:         feePercent,
		FeeAmount:          feeAmount,
		SourceService:      req.SourceService,
		DestinationService: req.DestService,
		UserID:             req.UserID,
		Status:             "pending",
	}

	if err := s.Repo.Create(feeTx); err != nil {
		return nil, fmt.Errorf("failed to create fee transaction: %w", err)
	}

	return &model.TransactionFeeResponse{
		TransactionID:  txID,
		OriginalAmount: req.Amount,
		FeePercent:     feePercent,
		FeeAmount:      feeAmount,
		TotalAmount:    totalAmount,
		Status:         "pending",
	}, nil
}

// UpdateFeeStatus mengupdate status fee setelah konfirmasi dari SmartBank
func (s *FeeService) UpdateFeeStatus(txID string, status string, smartBankRef string) error {
	return s.Repo.UpdateStatus(txID, status, smartBankRef)
}

// GetFeeByTransactionID mengambil fee berdasarkan transaction ID
func (s *FeeService) GetFeeByTransactionID(txID string) (*model.GatewayFeeTransaction, error) {
	return s.Repo.FindByTransactionID(txID)
}

// GetAllFees mengambil semua fee dengan pagination
func (s *FeeService) GetAllFees(page, perPage int) (*model.PaginatedResponse, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	fees, total, err := s.Repo.FindAll(page, perPage)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))

	return &model.PaginatedResponse{
		Data:       fees,
		Page:       page,
		PerPage:    perPage,
		TotalData:  total,
		TotalPages: totalPages,
	}, nil
}

// GetTotalFeeCollected menghitung total fee yang terkumpul
func (s *FeeService) GetTotalFeeCollected() (float64, error) {
	return s.Repo.GetTotalFeeCollected()
}

// GetTotalTransactions menghitung total transaksi
func (s *FeeService) GetTotalTransactions() (int64, error) {
	return s.Repo.GetTotalTransactions()
}
