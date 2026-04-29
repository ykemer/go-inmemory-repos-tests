package dtos

import "time"

type CreateContractRequest struct {
	Title       string    `json:"title" validate:"required,min=5,max=200"`
	Description string    `json:"description" validate:"required,min=10"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateContractRequest struct {
	Title       string    `json:"title" validate:"required,min=5,max=200"`
	Description string    `json:"description" validate:"required,min=10"`
	StartDate   time.Time `json:"start_date" validate:"required"`
	EndDate     time.Time `json:"end_date" validate:"required,gtfield=StartDate"`
}

type ContractResponse struct {
	ID             uint      `json:"id"`
	OrganizationID uint      `json:"organization_id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
}
