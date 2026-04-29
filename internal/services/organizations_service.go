package services

import (
	"context"
	"fmt"
	"tests/internal/dtos"
	"tests/internal/models"
	"tests/internal/repositories"
)

type OrganizationService struct {
	repo repositories.IOrgsRepository
}

func NewOrganizationService(repo repositories.IOrgsRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func (s *OrganizationService) List(ctx context.Context) ([]dtos.OrganizationResponse, error) {
	orgs, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]dtos.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		res[i] = dtos.OrganizationResponse{
			ID:        org.ID,
			Name:      org.Name,
			Email:     org.Email,
			CreatedAt: org.CreatedAt.UTC(),
		}
	}
	return res, nil
}

func (s *OrganizationService) Get(ctx context.Context, id string) (dtos.OrganizationResponse, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dtos.OrganizationResponse{}, err
	}
	if org == nil {
		return dtos.OrganizationResponse{}, fmt.Errorf("organization not found")
	}

	return dtos.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Email:     org.Email,
		CreatedAt: org.CreatedAt.UTC(),
	}, nil
}

func (s *OrganizationService) Create(ctx context.Context, req dtos.CreateOrganizationRequest) (dtos.OrganizationResponse, error) {
	org := models.Organization{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.repo.Create(ctx, &org); err != nil {
		return dtos.OrganizationResponse{}, err
	}

	return dtos.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Email:     org.Email,
		CreatedAt: org.CreatedAt.UTC(),
	}, nil
}

func (s *OrganizationService) Update(ctx context.Context, id string, req dtos.UpdateOrganizationRequest) (dtos.OrganizationResponse, error) {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dtos.OrganizationResponse{}, err
	}
	if org == nil {
		return dtos.OrganizationResponse{}, fmt.Errorf("organization not found")
	}

	org.Name = req.Name
	org.Email = req.Email

	if err := s.repo.Update(ctx, id, org); err != nil {
		return dtos.OrganizationResponse{}, err
	}

	return dtos.OrganizationResponse{
		ID:        org.ID,
		Name:      org.Name,
		Email:     org.Email,
		CreatedAt: org.CreatedAt.UTC(),
	}, nil
}

func (s *OrganizationService) Delete(ctx context.Context, id string) error {
	org, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if org == nil {
		return fmt.Errorf("organization not found")
	}

	return s.repo.Delete(ctx, id)
}
