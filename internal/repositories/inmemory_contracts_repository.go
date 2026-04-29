package repositories

import (
	"context"
	"fmt"
	"strconv"
	"tests/internal/models"
)

type InMemoryContractsRepository struct {
	data   map[uint]models.Contract
	nextID uint

	// Programmable hooks for testing edge cases/errors
	OnListByOrgId             func(ctx context.Context, orgId string) ([]models.Contract, error)
	OnGetByOrgIdAndContractId func(ctx context.Context, orgId string, contractId string) (*models.Contract, error)
	OnCreate                  func(ctx context.Context, data *models.Contract) error
	OnUpdate                  func(ctx context.Context, id string, data *models.Contract) error
	OnDelete                  func(ctx context.Context, id string) error
}

func NewInMemoryContractsRepository() *InMemoryContractsRepository {
	return &InMemoryContractsRepository{
		data:   make(map[uint]models.Contract),
		nextID: 1,
	}
}

func (r *InMemoryContractsRepository) ListByOrgId(ctx context.Context, orgId string) ([]models.Contract, error) {
	if r.OnListByOrgId != nil {
		return r.OnListByOrgId(ctx, orgId)
	}

	uintOrgID, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid organization id: %w", err)
	}

	contracts := make([]models.Contract, 0)
	for _, c := range r.data {
		if c.OrganizationID == uint(uintOrgID) {
			contracts = append(contracts, c)
		}
	}
	return contracts, nil
}

func (r *InMemoryContractsRepository) GetByOrgIdAndContractId(ctx context.Context, orgId string, contractId string) (*models.Contract, error) {
	if r.OnGetByOrgIdAndContractId != nil {
		return r.OnGetByOrgIdAndContractId(ctx, orgId, contractId)
	}

	uintOrgID, err := strconv.ParseUint(orgId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid organization id: %w", err)
	}

	uintContractID, err := strconv.ParseUint(contractId, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid contract id: %w", err)
	}

	c, ok := r.data[uint(uintContractID)]
	if !ok || c.OrganizationID != uint(uintOrgID) {
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

func (r *InMemoryContractsRepository) Update(ctx context.Context, id string, data *models.Contract) error {
	if r.OnUpdate != nil {
		return r.OnUpdate(ctx, id, data)
	}

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	data.ID = uint(uintID)
	r.data[uint(uintID)] = *data
	return nil
}

func (r *InMemoryContractsRepository) Delete(ctx context.Context, id string) error {
	if r.OnDelete != nil {
		return r.OnDelete(ctx, id)
	}

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	delete(r.data, uint(uintID))
	return nil
}
