package model

// =============================================
// API Response Standard
// =============================================

// APIResponse adalah format standar response API Gateway
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse untuk response error
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// =============================================
// Proxy Request / Response
// =============================================

// ProxyRequest menyimpan informasi request yang akan di-forward
type ProxyRequest struct {
	TargetService string            `json:"target_service"`
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Headers       map[string]string `json:"headers,omitempty"`
	Body          interface{}       `json:"body,omitempty"`
}

// =============================================
// Transaction Request (untuk fee calculation)
// =============================================

// TransactionRequest untuk menghitung fee gateway
type TransactionRequest struct {
	Amount         float64 `json:"amount" binding:"required"`
	SourceService  string  `json:"source_service" binding:"required"`
	DestService    string  `json:"destination_service" binding:"required"`
	UserID         uint    `json:"user_id" binding:"required"`
	Description    string  `json:"description"`
}

// TransactionFeeResponse response setelah fee dihitung
type TransactionFeeResponse struct {
	TransactionID  string  `json:"transaction_id"`
	OriginalAmount float64 `json:"original_amount"`
	FeePercent     float64 `json:"fee_percent"`
	FeeAmount      float64 `json:"fee_amount"`
	TotalAmount    float64 `json:"total_amount"`
	Status         string  `json:"status"`
}

// =============================================
// Service Health
// =============================================

// ServiceHealth menyimpan status kesehatan service
type ServiceHealth struct {
	ServiceName string `json:"service_name"`
	Status      string `json:"status"`
	URL         string `json:"url"`
	ResponseMs  int64  `json:"response_time_ms"`
}

// =============================================
// Dashboard / Analytics DTOs
// =============================================

// GatewayStats statistik gateway untuk dashboard
type GatewayStats struct {
	TotalRequests       int64   `json:"total_requests"`
	TotalTransactions   int64   `json:"total_transactions"`
	TotalFeeCollected   float64 `json:"total_fee_collected"`
	AvgResponseTimeMs   float64 `json:"avg_response_time_ms"`
	SuccessRate         float64 `json:"success_rate"`
	ActiveServices      int     `json:"active_services"`
}

// LogFilter untuk filter log
type LogFilter struct {
	Service   string `form:"service"`
	Method    string `form:"method"`
	Status    string `form:"status"`
	DateFrom  string `form:"date_from"`
	DateTo    string `form:"date_to"`
	Page      int    `form:"page,default=1"`
	PerPage   int    `form:"per_page,default=20"`
}

// PaginatedResponse untuk response dengan pagination
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalData  int64       `json:"total_data"`
	TotalPages int         `json:"total_pages"`
}
