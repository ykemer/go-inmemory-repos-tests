package services

import (
	"context"
	"fmt"
	"tests/internal/dtos"
	"tests/internal/models"
	"tests/internal/repositories"
)

type ContractService struct {
	repo     repositories.IContractsRepository
	orgsRepo repositories.IOrgsRepository
}

func NewContractService(repo repositories.IContractsRepository, orgsRepo repositories.IOrgsRepository) *ContractService {
	return &ContractService{repo: repo, orgsRepo: orgsRepo}
}

func (s *ContractService) ListByOrg(ctx context.Context, orgId uint) ([]dtos.ContractResponse, error) {
	org, err := s.orgsRepo.GetByID(ctx, orgId)
	if err != nil {
		return nil, err
	}
	if org == nil {
		return nil, fmt.Errorf("organization not found")
	}

	contracts, err := s.repo.ListByOrgId(ctx, orgId)
	if err != nil {
		return nil, err
	}

	res := make([]dtos.ContractResponse, len(contracts))
	for i, c := range contracts {
		res[i] = dtos.ContractResponse{
			ID:             c.ID,
			OrganizationID: c.OrganizationID,
			Title:          c.Title,
			Description:    c.Description,
			StartDate:      c.StartDate.UTC(),
			EndDate:        c.EndDate.UTC(),
		}
	}
	return res, nil
}

func (s *ContractService) Get(ctx context.Context, orgId uint, contractId uint) (dtos.ContractResponse, error) {
	contract, err := s.repo.GetByOrgIdAndContractId(ctx, orgId, contractId)
	if err != nil {
		return dtos.ContractResponse{}, err
	}
	if contract == nil {
		return dtos.ContractResponse{}, fmt.Errorf("contract not found")
	}

	return dtos.ContractResponse{
		ID:             contract.ID,
		OrganizationID: contract.OrganizationID,
		Title:          contract.Title,
		Description:    contract.Description,
		StartDate:      contract.StartDate.UTC(),
		EndDate:        contract.EndDate.UTC(),
	}, nil
}

func (s *ContractService) Create(ctx context.Context, orgId uint, req dtos.CreateContractRequest) (dtos.ContractResponse, error) {
	org, err := s.orgsRepo.GetByID(ctx, orgId)
	if err != nil {
		return dtos.ContractResponse{}, err
	}
	if org == nil {
		return dtos.ContractResponse{}, fmt.Errorf("organization not found")
	}

	contract := models.Contract{
		OrganizationID: orgId,
		Title:          req.Title,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
	}

	if err := s.repo.Create(ctx, &contract); err != nil {
		return dtos.ContractResponse{}, err
	}

	return dtos.ContractResponse{
		ID:             contract.ID,
		OrganizationID: contract.OrganizationID,
		Title:          contract.Title,
		Description:    contract.Description,
		StartDate:      contract.StartDate.UTC(),
		EndDate:        contract.EndDate.UTC(),
	}, nil
}

func (s *ContractService) Update(ctx context.Context, orgId uint, contractId uint, req dtos.UpdateContractRequest) (dtos.ContractResponse, error) {
	contract, err := s.repo.GetByOrgIdAndContractId(ctx, orgId, contractId)
	if err != nil {
		return dtos.ContractResponse{}, err
	}
	if contract == nil {
		return dtos.ContractResponse{}, fmt.Errorf("contract not found")
	}

	contract.Title = req.Title
	contract.Description = req.Description
	contract.StartDate = req.StartDate
	contract.EndDate = req.EndDate

	if err := s.repo.Update(ctx, contractId, contract); err != nil {
		return dtos.ContractResponse{}, err
	}

	return dtos.ContractResponse{
		ID:             contract.ID,
		OrganizationID: contract.OrganizationID,
		Title:          contract.Title,
		Description:    contract.Description,
		StartDate:      contract.StartDate.UTC(),
		EndDate:        contract.EndDate.UTC(),
	}, nil
}

func (s *ContractService) Delete(ctx context.Context, orgId uint, contractId uint) error {
	contract, err := s.repo.GetByOrgIdAndContractId(ctx, orgId, contractId)
	if err != nil {
		return err
	}
	if contract == nil {
		return fmt.Errorf("contract not found")
	}

	return s.repo.Delete(ctx, contractId)
}
