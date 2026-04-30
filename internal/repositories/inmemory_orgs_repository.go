package repositories

import (
	"context"
	"tests/internal/models"
	"time"
)

type InMemoryOrgsRepository struct {
	data   map[uint]models.Organization
	nextID uint

	// Programmable hooks for testing edge cases/errors
	OnList   func(ctx context.Context) ([]models.Organization, error)
	OnCreate func(ctx context.Context, data *models.Organization) error
	OnUpdate func(ctx context.Context, id uint, data *models.Organization) error
}

func NewInMemoryOrgsRepository() *InMemoryOrgsRepository {
	return &InMemoryOrgsRepository{
		data:   make(map[uint]models.Organization),
		nextID: 1,
	}
}

func (r *InMemoryOrgsRepository) List(ctx context.Context) ([]models.Organization, error) {
	if r.OnList != nil {
		return r.OnList(ctx)
	}

	orgs := make([]models.Organization, 0, len(r.data))
	for _, org := range r.data {
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r *InMemoryOrgsRepository) GetByID(ctx context.Context, id uint) (*models.Organization, error) {
	org, ok := r.data[id]
	if !ok {
		return nil, nil
	}
	return &org, nil
}

func (r *InMemoryOrgsRepository) Create(ctx context.Context, data *models.Organization) error {
	if r.OnCreate != nil {
		return r.OnCreate(ctx, data)
	}

	data.ID = r.nextID
	data.CreatedAt = time.Now()
	r.data[r.nextID] = *data
	r.nextID++
	return nil
}

func (r *InMemoryOrgsRepository) Update(ctx context.Context, id uint, data *models.Organization) error {
	if r.OnUpdate != nil {
		return r.OnUpdate(ctx, id, data)
	}

	data.ID = id
	r.data[id] = *data
	return nil
}

func (r *InMemoryOrgsRepository) Delete(ctx context.Context, id uint) error {
	delete(r.data, id)
	return nil
}
