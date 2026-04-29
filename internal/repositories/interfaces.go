package repositories

import (
	"context"
	"tests/internal/models"
)

type IOrgsRepository interface {
	List(ctx context.Context) ([]models.Organization, error)
	GetByID(ctx context.Context, id string) (*models.Organization, error)
	Create(ctx context.Context, data *models.Organization) error
	Update(ctx context.Context, id string, data *models.Organization) error
	Delete(ctx context.Context, id string) error
}

type IContractsRepository interface {
	ListByOrgId(ctx context.Context, orgId string) ([]models.Contract, error)
	GetByOrgIdAndContractId(ctx context.Context, orgId string, contractId string) (*models.Contract, error)
	Create(ctx context.Context, data *models.Contract) error
	Update(ctx context.Context, id string, data *models.Contract) error
	Delete(ctx context.Context, id string) error
}
