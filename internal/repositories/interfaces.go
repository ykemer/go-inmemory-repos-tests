package repositories

import (
	"context"
	"tests/internal/models"
)

type IOrgsRepository interface {
	List(ctx context.Context) ([]models.Organization, error)
	GetByID(ctx context.Context, id uint) (*models.Organization, error)
	Create(ctx context.Context, data *models.Organization) error
	Update(ctx context.Context, id uint, data *models.Organization) error
	Delete(ctx context.Context, id uint) error
}

type IContractsRepository interface {
	ListByOrgId(ctx context.Context, orgId uint) ([]models.Contract, error)
	GetByOrgIdAndContractId(ctx context.Context, orgId uint, contractId uint) (*models.Contract, error)
	Create(ctx context.Context, data *models.Contract) error
	Update(ctx context.Context, id uint, data *models.Contract) error
	Delete(ctx context.Context, id uint) error
}
