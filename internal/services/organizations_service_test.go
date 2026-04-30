package services

import (
	"context"
	"fmt"
	"testing"
	"tests/internal/dtos"
	"tests/internal/models"
	"tests/internal/repositories"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationService_List(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		repo.Create(ctx, &models.Organization{Name: "Org 1", Email: "1@org.com"})
		repo.Create(ctx, &models.Organization{Name: "Org 2", Email: "2@org.com"})

		resp, err := service.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("EmptyList", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		resp, err := service.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, resp, 0)
	})

	t.Run("RepoError", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		repo.OnList = func(ctx context.Context) ([]models.Organization, error) {
			return nil, fmt.Errorf("db error")
		}
		service := NewOrganizationService(repo)
		_, err := service.List(context.Background())
		assert.Error(t, err)
	})
}

func TestOrganizationService_Get(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		org := models.Organization{Name: "Org 1", Email: "1@org.com"}
		repo.Create(ctx, &org)

		resp, err := service.Get(ctx, org.ID)

		assert.NoError(t, err)
		assert.Equal(t, org.Name, resp.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		_, err := service.Get(ctx, 999)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "organization not found")
	})
}

func TestOrganizationService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		req := dtos.CreateOrganizationRequest{
			Name:  "Test Org",
			Email: "test@org.com",
		}

		resp, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
		assert.NotZero(t, resp.ID)
	})

	t.Run("RepoError", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		repo.OnCreate = func(ctx context.Context, data *models.Organization) error {
			return fmt.Errorf("db error")
		}
		service := NewOrganizationService(repo)
		_, err := service.Create(context.Background(), dtos.CreateOrganizationRequest{Name: "Org", Email: "org@test.com"})
		assert.Error(t, err)
	})
}

func TestOrganizationService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		org := models.Organization{Name: "Old Name", Email: "old@org.com"}
		repo.Create(ctx, &org)

		req := dtos.UpdateOrganizationRequest{
			Name:  "New Name",
			Email: "new@org.com",
		}

		resp, err := service.Update(ctx, org.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)

		updated, _ := repo.GetByID(ctx, org.ID)
		assert.Equal(t, req.Name, updated.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		req := dtos.UpdateOrganizationRequest{Name: "New"}
		_, err := service.Update(ctx, 999, req)

		assert.Error(t, err)
	})

	t.Run("RepoUpdateError", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		org := models.Organization{Name: "Old"}
		repo.Create(context.Background(), &org)

		repo.OnUpdate = func(ctx context.Context, id uint, data *models.Organization) error {
			return fmt.Errorf("db update error")
		}
		service := NewOrganizationService(repo)
		_, err := service.Update(context.Background(), org.ID, dtos.UpdateOrganizationRequest{Name: "New"})
		assert.Error(t, err)
	})
}

func TestOrganizationService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		org := models.Organization{Name: "Delete Me"}
		repo.Create(ctx, &org)

		err := service.Delete(ctx, org.ID)

		assert.NoError(t, err)
		res, _ := repo.GetByID(ctx, org.ID)
		assert.Nil(t, res)
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := repositories.NewInMemoryOrgsRepository()
		service := NewOrganizationService(repo)
		ctx := context.Background()

		err := service.Delete(ctx, 999)
		assert.Error(t, err)
	})
}
