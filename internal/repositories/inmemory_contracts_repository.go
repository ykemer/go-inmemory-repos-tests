package repositories

import (
	"context"
	"tests/internal/models"
)

type InMemoryContractsRepository struct {
	data   map[uint]models.Contract
	nextID uint

	// Programmable hooks for testing edge cases/errors
	OnListByOrgId func(ctx context.Context, orgId uint) ([]models.Contract, error)
	OnCreate      func(ctx context.Context, data *models.Contract) error
	OnUpdate      func(ctx context.Context, id uint, data *models.Contract) error
}

func NewInMemoryContractsRepository() *InMemoryContractsRepository {
	return &InMemoryContractsRepository{
		data:   make(map[uint]models.Contract),
		nextID: 1,
	}
}

func (r *InMemoryContractsRepository) ListByOrgId(ctx context.Context, orgId uint) ([]models.Contract, error) {
	if r.OnListByOrgId != nil {
		return r.OnListByOrgId(ctx, orgId)
	}

	contracts := make([]models.Contract, 0)
	for _, c := range r.data {
		if c.OrganizationID == orgId {
			contracts = append(contracts, c)
		}
	}
	return contracts, nil
}

func (r *InMemoryContractsRepository) GetByOrgIdAndContractId(ctx context.Context, orgId uint, contractId uint) (*models.Contract, error) {
	c, ok := r.data[contractId]
	if !ok || c.OrganizationID != orgId {
		return nil, nil
	}
	return &c, nil
}

func (r *InMemoryContractsRepository) Create(ctx context.Context, data *models.Contract) error {
	if r.OnCreate != nil {
		return r.OnCreate(ctx, data)
	}

	data.ID = r.nextID
	r.data[r.nextID] = *data
	r.nextID++
	return nil
}

func (r *InMemoryContractsRepository) Update(ctx context.Context, id uint, data *models.Contract) error {
	if r.OnUpdate != nil {
		return r.OnUpdate(ctx, id, data)
	}

	data.ID = id
	r.data[id] = *data
	return nil
}

func (r *InMemoryContractsRepository) Delete(ctx context.Context, id uint) error {
	delete(r.data, id)
	return nil
}
