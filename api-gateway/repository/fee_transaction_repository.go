package repository

import (
	"api-gateway/model"

	"gorm.io/gorm"
)

// FeeTransactionRepository menangani operasi database untuk fee transaksi
type FeeTransactionRepository struct {
	DB *gorm.DB
}

// NewFeeTransactionRepository membuat instance baru
func NewFeeTransactionRepository(db *gorm.DB) *FeeTransactionRepository {
	return &FeeTransactionRepository{DB: db}
}

// Create menyimpan transaksi fee baru
func (r *FeeTransactionRepository) Create(tx *model.GatewayFeeTransaction) error {
	return r.DB.Create(tx).Error
}

// FindByID mencari transaksi fee berdasarkan ID
func (r *FeeTransactionRepository) FindByID(id uint) (*model.GatewayFeeTransaction, error) {
	var tx model.GatewayFeeTransaction
	err := r.DB.First(&tx, id).Error
	return &tx, err
}

// FindByTransactionID mencari berdasarkan transaction ID
func (r *FeeTransactionRepository) FindByTransactionID(txID string) (*model.GatewayFeeTransaction, error) {
	var tx model.GatewayFeeTransaction
	err := r.DB.Where("transaction_id = ?", txID).First(&tx).Error
	return &tx, err
}

// FindAll mengambil semua transaksi fee dengan pagination
func (r *FeeTransactionRepository) FindAll(page, perPage int) ([]model.GatewayFeeTransaction, int64, error) {
	var txs []model.GatewayFeeTransaction
	var total int64

	r.DB.Model(&model.GatewayFeeTransaction{}).Count(&total)

	offset := (page - 1) * perPage
	err := r.DB.Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(&txs).Error

	return txs, total, err
}

// UpdateStatus mengubah status transaksi fee
func (r *FeeTransactionRepository) UpdateStatus(txID string, status string, smartBankRef string) error {
	return r.DB.Model(&model.GatewayFeeTransaction{}).
		Where("transaction_id = ?", txID).
		Updates(map[string]interface{}{
			"status":        status,
			"smartbank_ref": smartBankRef,
		}).Error
}

// GetTotalFeeCollected menghitung total fee yang terkumpul
func (r *FeeTransactionRepository) GetTotalFeeCollected() (float64, error) {
	var total float64
	err := r.DB.Model(&model.GatewayFeeTransaction{}).
		Where("status = ?", "success").
		Select("COALESCE(SUM(fee_amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalTransactions menghitung jumlah transaksi
func (r *FeeTransactionRepository) GetTotalTransactions() (int64, error) {
	var count int64
	err := r.DB.Model(&model.GatewayFeeTransaction{}).Count(&count).Error
	return count, err
}
