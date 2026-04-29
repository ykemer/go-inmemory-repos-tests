package repositories

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"tests/internal/models"
	"time"
)

type InMemoryOrgsRepository struct {
	mu     sync.RWMutex
	data   map[uint]models.Organization
	nextID uint

	// Programmable hooks for testing edge cases/errors
	OnList    func(ctx context.Context) ([]models.Organization, error)
	OnGetByID func(ctx context.Context, id string) (*models.Organization, error)
	OnCreate  func(ctx context.Context, data *models.Organization) error
	OnUpdate  func(ctx context.Context, id string, data *models.Organization) error
	OnDelete  func(ctx context.Context, id string) error
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

	r.mu.RLock()
	defer r.mu.RUnlock()

	orgs := make([]models.Organization, 0, len(r.data))
	for _, org := range r.data {
		orgs = append(orgs, org)
	}
	return orgs, nil
}

func (r *InMemoryOrgsRepository) GetByID(ctx context.Context, id string) (*models.Organization, error) {
	if r.OnGetByID != nil {
		return r.OnGetByID(ctx, id)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %w", err)
	}

	org, ok := r.data[uint(uintID)]
	if !ok {
		return nil, nil
	}
	return &org, nil
}

func (r *InMemoryOrgsRepository) Create(ctx context.Context, data *models.Organization) error {
	if r.OnCreate != nil {
		return r.OnCreate(ctx, data)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	data.ID = r.nextID
	data.CreatedAt = time.Now()
	r.data[r.nextID] = *data
	r.nextID++
	return nil
}

func (r *InMemoryOrgsRepository) Update(ctx context.Context, id string, data *models.Organization) error {
	if r.OnUpdate != nil {
		return r.OnUpdate(ctx, id, data)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	data.ID = uint(uintID)
	r.data[uint(uintID)] = *data
	return nil
}

func (r *InMemoryOrgsRepository) Delete(ctx context.Context, id string) error {
	if r.OnDelete != nil {
		return r.OnDelete(ctx, id)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	uintID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid id: %w", err)
	}

	delete(r.data, uint(uintID))
	return nil
}

func (r *InMemoryOrgsRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data = make(map[uint]models.Organization)
	r.nextID = 1
	r.OnList = nil
	r.OnGetByID = nil
	r.OnCreate = nil
	r.OnUpdate = nil
	r.OnDelete = nil
}
