package service

import (
	"api-gateway/config"
	"api-gateway/model"
	"api-gateway/repository"
	"fmt"
	"net/http"
	"time"
)

// ServiceRegistryService menangani business logic untuk service registry
type ServiceRegistryService struct {
	Repo *repository.ServiceRegistryRepository
}

// NewServiceRegistryService membuat instance baru
func NewServiceRegistryService(repo *repository.ServiceRegistryRepository) *ServiceRegistryService {
	return &ServiceRegistryService{Repo: repo}
}

// RegisterService mendaftarkan service baru ke gateway
func (s *ServiceRegistryService) RegisterService(service *model.ServiceRegistry) error {
	return s.Repo.Create(service)
}

// GetServiceByName mengambil service berdasarkan nama
func (s *ServiceRegistryService) GetServiceByName(name string) (*model.ServiceRegistry, error) {
	return s.Repo.FindByName(name)
}

// GetAllServices mengambil semua service
func (s *ServiceRegistryService) GetAllServices() ([]model.ServiceRegistry, error) {
	return s.Repo.FindAll()
}

// GetActiveServices mengambil semua service aktif
func (s *ServiceRegistryService) GetActiveServices() ([]model.ServiceRegistry, error) {
	return s.Repo.FindActive()
}

// UpdateService mengupdate data service
func (s *ServiceRegistryService) UpdateService(service *model.ServiceRegistry) error {
	return s.Repo.Update(service)
}

// UpdateServiceStatus mengubah status service
func (s *ServiceRegistryService) UpdateServiceStatus(name string, status string) error {
	return s.Repo.UpdateStatus(name, status)
}

// DeleteService menghapus service
func (s *ServiceRegistryService) DeleteService(name string) error {
	return s.Repo.Delete(name)
}

// GetServiceURL mendapatkan URL service berdasarkan nama
// Cek dari registry DB dulu, kalau tidak ada fallback ke config
func (s *ServiceRegistryService) GetServiceURL(serviceName string) (string, error) {
	service, err := s.Repo.FindByName(serviceName)
	if err == nil && service.Status == "active" {
		return service.BaseURL, nil
	}

	// Fallback ke config
	switch serviceName {
	case "smartbank":
		return config.AppConfig.SmartBankURL, nil
	case "marketplace":
		return config.AppConfig.MarketplaceURL, nil
	case "pos":
		return config.AppConfig.PosURL, nil
	case "logistikita":
		return config.AppConfig.LogistiKitaURL, nil
	case "supplierhub":
		return config.AppConfig.SupplierHubURL, nil
	case "umkm-insight":
		return config.AppConfig.UMKMInsightURL, nil
	default:
		return "", fmt.Errorf("service '%s' not found", serviceName)
	}
}

// HealthCheck memeriksa status kesehatan service
func (s *ServiceRegistryService) HealthCheck(serviceName, url string) *model.ServiceHealth {
	health := &model.ServiceHealth{
		ServiceName: serviceName,
		URL:         url,
		Status:      "unknown",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	start := time.Now()

	resp, err := client.Get(url + "/health")
	elapsed := time.Since(start).Milliseconds()
	health.ResponseMs = elapsed

	if err != nil {
		health.Status = "down"
		return health
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		health.Status = "healthy"
	} else {
		health.Status = "unhealthy"
	}

	return health
}

// HealthCheckAll memeriksa kesehatan semua service yang terdaftar
func (s *ServiceRegistryService) HealthCheckAll() []model.ServiceHealth {
	services, err := s.Repo.FindAll()
	if err != nil {
		return nil
	}

	var results []model.ServiceHealth
	for _, svc := range services {
		health := s.HealthCheck(svc.Name, svc.BaseURL)
		results = append(results, *health)
	}
	return results
}

// CountActive menghitung service aktif
func (s *ServiceRegistryService) CountActive() (int64, error) {
	return s.Repo.CountActive()
}

// SeedDefaultServices mendaftarkan service default dari config
func (s *ServiceRegistryService) SeedDefaultServices() {
	defaultServices := []model.ServiceRegistry{
		{Name: "smartbank", BaseURL: config.AppConfig.SmartBankURL, Status: "active"},
		{Name: "marketplace", BaseURL: config.AppConfig.MarketplaceURL, Status: "active"},
		{Name: "pos", BaseURL: config.AppConfig.PosURL, Status: "active"},
		{Name: "logistikita", BaseURL: config.AppConfig.LogistiKitaURL, Status: "active"},
		{Name: "supplierhub", BaseURL: config.AppConfig.SupplierHubURL, Status: "active"},
		{Name: "umkm-insight", BaseURL: config.AppConfig.UMKMInsightURL, Status: "active"},
	}

	for _, svc := range defaultServices {
		existing, err := s.Repo.FindByName(svc.Name)
		if err != nil || existing.ID == 0 {
			s.Repo.Create(&svc)
		}
	}
}
