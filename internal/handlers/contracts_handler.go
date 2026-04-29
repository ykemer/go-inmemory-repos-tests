package handlers

import (
	"tests/internal/dtos"
	"tests/internal/services"

	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	service *services.ContractService
}

func NewContractHandler(service *services.ContractService) *ContractHandler {
	return &ContractHandler{service: service}
}

// List godoc
// @Summary      List contracts for an organization
// @Tags         Contracts
// @Produce      json
// @Param        orgId  path      int  true  "Organization ID"
// @Success      200    {array}   dtos.ContractResponse
// @Failure      404    {object}  dtos.ProblemDetail
// @Router       /organizations/{orgId}/contracts [get]
func (h *ContractHandler) List(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	res, err := h.service.ListByOrg(c.Context(), orgId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Organization not found")
	}
	return c.JSON(res)
}

// Get godoc
// @Summary      Get a contract
// @Tags         Contracts
// @Produce      json
// @Param        orgId       path      int  true  "Organization ID"
// @Param        contractId  path      int  true  "Contract ID"
// @Success      200         {object}  dtos.ContractResponse
// @Failure      404         {object}  dtos.ProblemDetail
// @Router       /organizations/{orgId}/contracts/{contractId} [get]
func (h *ContractHandler) Get(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	res, err := h.service.Get(c.Context(), orgId, contractId)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}
	return c.JSON(res)
}

// Create godoc
// @Summary      Create a contract
// @Tags         Contracts
// @Accept       json
// @Produce      json
// @Param        orgId            path    int                        true   "Organization ID"
// @Param        Idempotency-Key  header  string                     false  "Client-generated UUID. Repeated requests with the same key return the cached response for 1 hour."
// @Param        body             body    dtos.CreateContractRequest true   "Request body"
// @Success      201  {object}  dtos.ContractResponse
// @Failure      404  {object}  dtos.ProblemDetail
// @Failure      422  {object}  dtos.ProblemDetail
// @Router       /organizations/{orgId}/contracts [post]
func (h *ContractHandler) Create(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	req := new(dtos.CreateContractRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Create(c.Context(), orgId, *req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Update godoc
// @Summary      Update a contract
// @Tags         Contracts
// @Accept       json
// @Produce      json
// @Param        orgId            path    int                        true   "Organization ID"
// @Param        contractId       path    int                        true   "Contract ID"
// @Param        Idempotency-Key  header  string                     false  "Client-generated UUID. Repeated requests with the same key return the cached response for 1 hour."
// @Param        body             body    dtos.UpdateContractRequest true   "Request body"
// @Success      200  {object}  dtos.ContractResponse
// @Failure      404  {object}  dtos.ProblemDetail
// @Failure      422  {object}  dtos.ProblemDetail
// @Router       /organizations/{orgId}/contracts/{contractId} [put]
func (h *ContractHandler) Update(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	req := new(dtos.UpdateContractRequest)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	if err := validate.Struct(req); err != nil {
		return handleValidationError(c, err)
	}

	res, err := h.service.Update(c.Context(), orgId, contractId, *req)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}

	return c.JSON(res)
}

// Delete godoc
// @Summary      Delete a contract
// @Tags         Contracts
// @Produce      json
// @Param        orgId       path  int  true  "Organization ID"
// @Param        contractId  path  int  true  "Contract ID"
// @Success      204
// @Failure      404  {object}  dtos.ProblemDetail
// @Router       /organizations/{orgId}/contracts/{contractId} [delete]
func (h *ContractHandler) Delete(c *fiber.Ctx) error {
	orgId := c.Params("orgId")
	contractId := c.Params("contractId")
	if err := h.service.Delete(c.Context(), orgId, contractId); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Contract not found")
	}
	return c.SendStatus(fiber.StatusNoContent)
}
