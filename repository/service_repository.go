package repository

import (
	"api-gateway/model"

	"gorm.io/gorm"
)

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) Create(service *model.Service) error {
	return r.db.Create(service).Error
}

func (r *ServiceRepository) FindByID(id uint) (*model.Service, error) {
	var service model.Service
	err := r.db.First(&service, id).Error
	return &service, err
}

func (r *ServiceRepository) FindAll(page, limit int) ([]model.Service, int64, error) {
	var services []model.Service
	var total int64

	r.db.Model(&model.Service{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.Offset(offset).Limit(limit).Order("id ASC").Find(&services).Error
	return services, total, err
}

func (r *ServiceRepository) FindAllActive() ([]model.Service, error) {
	var services []model.Service
	err := r.db.Where("is_active = ?", true).Find(&services).Error
	return services, err
}

func (r *ServiceRepository) FindByRoutePrefix(prefix string) (*model.Service, error) {
	var service model.Service
	err := r.db.Where("route_prefix = ? AND is_active = ?", prefix, true).First(&service).Error
	return &service, err
}

func (r *ServiceRepository) Update(service *model.Service) error {
	return r.db.Save(service).Error
}

func (r *ServiceRepository) Delete(id uint) error {
	return r.db.Delete(&model.Service{}, id).Error
}
