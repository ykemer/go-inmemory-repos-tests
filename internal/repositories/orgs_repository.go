package repositories

import (
	"context"
	"errors"
	"tests/internal/models"

	"gorm.io/gorm"
)

type OrgsRepository struct {
	db *gorm.DB
}

func NewOrgsRepository(db *gorm.DB) IOrgsRepository {
	return &OrgsRepository{db: db}
}

func (r *OrgsRepository) List(ctx context.Context) ([]models.Organization, error) {
	var orgs []models.Organization
	if err := r.db.WithContext(ctx).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *OrgsRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	var org models.Organization
	if err := r.db.WithContext(ctx).First(&org, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &org, nil
}

func (r *OrgsRepository) Create(ctx context.Context, data *models.Organization) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *OrgsRepository) Update(ctx context.Context, id string, data *models.Organization) error {
	return r.db.WithContext(ctx).Model(&models.Organization{}).Where("id = ?", id).Updates(data).Error
}

func (r *OrgsRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Organization{}, id).Error
}
