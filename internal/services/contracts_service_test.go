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

func TestContractService_ListByOrg(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		org1 := models.Organization{Name: "Org 1"}
		orgRepo.Create(ctx, &org1)

		contractRepo.Create(ctx, &models.Contract{OrganizationID: org1.ID, Title: "C1"})
		contractRepo.Create(ctx, &models.Contract{OrganizationID: org1.ID, Title: "C2"})

		res, err := service.ListByOrg(ctx, org1.ID)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		_, err := service.ListByOrg(ctx, 999)

		assert.Error(t, err)
	})

	t.Run("RepoError", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		org := models.Organization{Name: "Org"}
		orgRepo.Create(context.Background(), &org)

		contractRepo := repositories.NewInMemoryContractsRepository()
		contractRepo.OnListByOrgId = func(ctx context.Context, orgId uint) ([]models.Contract, error) {
			return nil, fmt.Errorf("db error")
		}
		service := NewContractService(contractRepo, orgRepo)
		_, err := service.ListByOrg(context.Background(), org.ID)
		assert.Error(t, err)
	})
}

func TestContractService_Get(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		org := models.Organization{Name: "Acme"}
		orgRepo.Create(ctx, &org)

		contract := models.Contract{OrganizationID: org.ID, Title: "Contract 1"}
		contractRepo.Create(ctx, &contract)

		resp, err := service.Get(ctx, org.ID, contract.ID)

		assert.NoError(t, err)
		assert.Equal(t, contract.Title, resp.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		_, err := service.Get(ctx, 1, 999)

		assert.Error(t, err)
	})
}

func TestContractService_Create(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		org := models.Organization{Name: "Acme"}
		orgRepo.Create(ctx, &org)

		req := dtos.CreateContractRequest{Title: "Service Agreement", Description: "Testing in-memory description"}

		resp, err := service.Create(ctx, org.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Title, resp.Title)
		assert.Equal(t, org.ID, resp.OrganizationID)
	})

	t.Run("OrgNotFound", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		req := dtos.CreateContractRequest{Title: "Title"}
		_, err := service.Create(ctx, 999, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "organization not found")
	})

	t.Run("RepoCreateError", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		org := models.Organization{Name: "Acme"}
		orgRepo.Create(context.Background(), &org)

		contractRepo := repositories.NewInMemoryContractsRepository()
		contractRepo.OnCreate = func(ctx context.Context, data *models.Contract) error {
			return fmt.Errorf("db error")
		}
		service := NewContractService(contractRepo, orgRepo)
		_, err := service.Create(context.Background(), org.ID, dtos.CreateContractRequest{Title: "T"})
		assert.Error(t, err)
	})
}

func TestContractService_Update(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		org := models.Organization{Name: "Acme"}
		orgRepo.Create(ctx, &org)

		contract := models.Contract{OrganizationID: org.ID, Title: "Old Title"}
		contractRepo.Create(ctx, &contract)

		req := dtos.UpdateContractRequest{Title: "New Title", Description: "New Description"}

		resp, err := service.Update(ctx, org.ID, contract.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Title, resp.Title)

		updated, _ := contractRepo.GetByOrgIdAndContractId(ctx, org.ID, contract.ID)
		assert.Equal(t, req.Title, updated.Title)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		req := dtos.UpdateContractRequest{Title: "New"}
		_, err := service.Update(ctx, 1, 999, req)

		assert.Error(t, err)
	})

	t.Run("RepoUpdateError", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		org := models.Organization{Name: "Acme"}
		orgRepo.Create(context.Background(), &org)

		contractRepo := repositories.NewInMemoryContractsRepository()
		contract := models.Contract{OrganizationID: org.ID, Title: "Old"}
		contractRepo.Create(context.Background(), &contract)

		contractRepo.OnUpdate = func(ctx context.Context, id uint, data *models.Contract) error {
			return fmt.Errorf("db error")
		}
		service := NewContractService(contractRepo, orgRepo)
		_, err := service.Update(context.Background(), org.ID, contract.ID, dtos.UpdateContractRequest{Title: "New"})
		assert.Error(t, err)
	})
}

func TestContractService_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		org := models.Organization{Name: "Acme"}
		orgRepo.Create(ctx, &org)

		contract := models.Contract{OrganizationID: org.ID, Title: "Delete Me"}
		contractRepo.Create(ctx, &contract)

		err := service.Delete(ctx, org.ID, contract.ID)

		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		orgRepo := repositories.NewInMemoryOrgsRepository()
		contractRepo := repositories.NewInMemoryContractsRepository()
		service := NewContractService(contractRepo, orgRepo)
		ctx := context.Background()

		err := service.Delete(ctx, 1, 999)
		assert.Error(t, err)
	})
}
