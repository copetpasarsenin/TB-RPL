package controller

import (
	"api-gateway/model"
	"api-gateway/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FeeController menangani endpoint fee gateway
type FeeController struct {
	FeeService *service.FeeService
}

// NewFeeController membuat instance baru
func NewFeeController(fs *service.FeeService) *FeeController {
	return &FeeController{FeeService: fs}
}

// CalculateFee menghitung fee gateway 0.5% untuk sebuah transaksi
// Route: POST /gateway/fee/calculate
func (fc *FeeController) CalculateFee(c *gin.Context) {
	var req model.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validasi amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Amount harus lebih dari 0",
		})
		return
	}

	result, err := fc.FeeService.CalculateAndCreateFee(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghitung fee",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Fee berhasil dihitung",
		"data":    result,
	})
}

// GetFeeByTransaction mengambil detail fee berdasarkan transaction ID
// Route: GET /gateway/fee/:transaction_id
func (fc *FeeController) GetFeeByTransaction(c *gin.Context) {
	txID := c.Param("transaction_id")

	fee, err := fc.FeeService.GetFeeByTransactionID(txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Transaksi fee tidak ditemukan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Detail fee ditemukan",
		"data":    fee,
	})
}

// GetAllFees mengambil semua transaksi fee dengan pagination
// Route: GET /gateway/fees
func (fc *FeeController) GetAllFees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	result, err := fc.FeeService.GetAllFees(page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data fee",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data fee berhasil diambil",
		"data":    result,
	})
}

// UpdateFeeStatus mengupdate status fee (callback dari SmartBank)
// Route: PUT /gateway/fee/:transaction_id/status
func (fc *FeeController) UpdateFeeStatus(c *gin.Context) {
	txID := c.Param("transaction_id")

	var req struct {
		Status       string `json:"status" binding:"required"`
		SmartBankRef string `json:"smartbank_ref"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data tidak valid",
			"error":   err.Error(),
		})
		return
	}

	// Validasi status
	validStatuses := map[string]bool{"pending": true, "success": true, "failed": true}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Status harus: pending, success, atau failed",
		})
		return
	}

	if err := fc.FeeService.UpdateFeeStatus(txID, req.Status, req.SmartBankRef); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengupdate status fee",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Status fee berhasil diupdate",
	})
}

// GetFeeStats mengambil statistik fee
// Route: GET /gateway/fee/stats
func (fc *FeeController) GetFeeStats(c *gin.Context) {
	totalFee, _ := fc.FeeService.GetTotalFeeCollected()
	totalTx, _ := fc.FeeService.GetTotalTransactions()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Statistik fee",
		"data": gin.H{
			"total_fee_collected":  totalFee,
			"total_transactions":   totalTx,
			"fee_percent":          0.5,
		},
	})
}
