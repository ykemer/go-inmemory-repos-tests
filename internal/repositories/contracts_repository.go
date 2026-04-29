package repositories

import (
	"context"
	"errors"
	"tests/internal/models"

	"gorm.io/gorm"
)

type ContractsRepository struct {
	db *gorm.DB
}

func NewContractsRepository(db *gorm.DB) IContractsRepository {
	return &ContractsRepository{db: db}
}

func (r *ContractsRepository) ListByOrgId(ctx context.Context, orgId string) ([]models.Contract, error) {
	var contracts []models.Contract
	if err := r.db.WithContext(ctx).Where("organization_id = ?", orgId).Find(&contracts).Error; err != nil {
		return nil, err
	}
	return contracts, nil
}

func (r *ContractsRepository) GetByOrgIdAndContractId(ctx context.Context, orgId string, contractId string) (*models.Contract, error) {
	var contract models.Contract
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND id = ?", orgId, contractId).First(&contract).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &contract, nil
}

func (r *ContractsRepository) Create(ctx context.Context, data *models.Contract) error {
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *ContractsRepository) Update(ctx context.Context, id string, data *models.Contract) error {
	return r.db.WithContext(ctx).Model(&models.Contract{}).Where("id = ?", id).Updates(data).Error
}

func (r *ContractsRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&models.Contract{}, id).Error
}
