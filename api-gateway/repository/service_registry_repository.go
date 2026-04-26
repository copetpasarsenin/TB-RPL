package repository

import (
	"api-gateway/model"

	"gorm.io/gorm"
)

// ServiceRegistryRepository menangani operasi database untuk service registry
type ServiceRegistryRepository struct {
	DB *gorm.DB
}

// NewServiceRegistryRepository membuat instance baru
func NewServiceRegistryRepository(db *gorm.DB) *ServiceRegistryRepository {
	return &ServiceRegistryRepository{DB: db}
}

// Create mendaftarkan service baru
func (r *ServiceRegistryRepository) Create(service *model.ServiceRegistry) error {
	return r.DB.Create(service).Error
}

// FindByName mencari service berdasarkan nama
func (r *ServiceRegistryRepository) FindByName(name string) (*model.ServiceRegistry, error) {
	var service model.ServiceRegistry
	err := r.DB.Where("name = ?", name).First(&service).Error
	return &service, err
}

// FindAll mengambil semua service yang terdaftar
func (r *ServiceRegistryRepository) FindAll() ([]model.ServiceRegistry, error) {
	var services []model.ServiceRegistry
	err := r.DB.Find(&services).Error
	return services, err
}

// FindActive mengambil semua service yang aktif
func (r *ServiceRegistryRepository) FindActive() ([]model.ServiceRegistry, error) {
	var services []model.ServiceRegistry
	err := r.DB.Where("status = ?", "active").Find(&services).Error
	return services, err
}

// Update mengubah data service
func (r *ServiceRegistryRepository) Update(service *model.ServiceRegistry) error {
	return r.DB.Save(service).Error
}

// UpdateStatus mengubah status service
func (r *ServiceRegistryRepository) UpdateStatus(name string, status string) error {
	return r.DB.Model(&model.ServiceRegistry{}).
		Where("name = ?", name).
		Update("status", status).Error
}

// Delete menghapus service (soft delete)
func (r *ServiceRegistryRepository) Delete(name string) error {
	return r.DB.Where("name = ?", name).Delete(&model.ServiceRegistry{}).Error
}

// CountActive menghitung jumlah service aktif
func (r *ServiceRegistryRepository) CountActive() (int64, error) {
	var count int64
	err := r.DB.Model(&model.ServiceRegistry{}).
		Where("status = ?", "active").
		Count(&count).Error
	return count, err
}
