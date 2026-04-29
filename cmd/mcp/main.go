package main

import (
	"context"
	"encoding/json"
	"log"
	"tests/internal/config"
	"tests/internal/dtos"
	"tests/internal/repositories"
	"tests/internal/services"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	cfg := config.LoadConfig()
	db := cfg.InitDatabase()

	orgRepo := repositories.NewOrgsRepository(db)
	contractRepo := repositories.NewContractsRepository(db)

	orgService := services.NewOrganizationService(orgRepo)
	contractService := services.NewContractService(contractRepo, orgRepo)

	s := server.NewMCPServer("Organizations & Contracts", "1.0.0")

	registerOrgTools(s, orgService)
	registerContractTools(s, contractService)

	if err := server.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}

// ── Organizations ─────────────────────────────────────────────────────────────

func registerOrgTools(s *server.MCPServer, svc *services.OrganizationService) {
	s.AddTool(
		mcp.NewTool("list_organizations",
			mcp.WithDescription("Return all organizations."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.List(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("get_organization",
			mcp.WithDescription("Return a single organization by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Organization ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.Get(ctx, mcp.ParseString(req, "id", ""))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("create_organization",
			mcp.WithDescription("Create a new organization."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Organization name (3–100 chars)")),
			mcp.WithString("email", mcp.Required(), mcp.Description("Contact email address")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.Create(ctx, dtos.CreateOrganizationRequest{
				Name:  mcp.ParseString(req, "name", ""),
				Email: mcp.ParseString(req, "email", ""),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("update_organization",
			mcp.WithDescription("Update an existing organization."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Organization ID")),
			mcp.WithString("name", mcp.Required(), mcp.Description("New name (3–100 chars)")),
			mcp.WithString("email", mcp.Required(), mcp.Description("New email address")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.Update(ctx, mcp.ParseString(req, "id", ""), dtos.UpdateOrganizationRequest{
				Name:  mcp.ParseString(req, "name", ""),
				Email: mcp.ParseString(req, "email", ""),
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("delete_organization",
			mcp.WithDescription("Delete an organization by ID."),
			mcp.WithString("id", mcp.Required(), mcp.Description("Organization ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := svc.Delete(ctx, mcp.ParseString(req, "id", "")); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText("organization deleted"), nil
		},
	)
}

// ── Contracts ─────────────────────────────────────────────────────────────────

func registerContractTools(s *server.MCPServer, svc *services.ContractService) {
	s.AddTool(
		mcp.NewTool("list_contracts",
			mcp.WithDescription("Return all contracts for an organization."),
			mcp.WithString("org_id", mcp.Required(), mcp.Description("Organization ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.ListByOrg(ctx, mcp.ParseString(req, "org_id", ""))
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("get_contract",
			mcp.WithDescription("Return a single contract."),
			mcp.WithString("org_id", mcp.Required(), mcp.Description("Organization ID")),
			mcp.WithString("contract_id", mcp.Required(), mcp.Description("Contract ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			res, err := svc.Get(ctx,
				mcp.ParseString(req, "org_id", ""),
				mcp.ParseString(req, "contract_id", ""),
			)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("create_contract",
			mcp.WithDescription("Create a contract under an organization."),
			mcp.WithString("org_id", mcp.Required(), mcp.Description("Organization ID")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Contract title (5–200 chars)")),
			mcp.WithString("description", mcp.Required(), mcp.Description("Contract description (min 10 chars)")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("Start date (RFC3339, e.g. 2024-01-01T00:00:00Z)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("End date (RFC3339, must be after start_date)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start, err := time.Parse(time.RFC3339, mcp.ParseString(req, "start_date", ""))
			if err != nil {
				return mcp.NewToolResultError("start_date must be RFC3339 (e.g. 2024-01-01T00:00:00Z)"), nil
			}
			end, err := time.Parse(time.RFC3339, mcp.ParseString(req, "end_date", ""))
			if err != nil {
				return mcp.NewToolResultError("end_date must be RFC3339 (e.g. 2024-12-31T23:59:59Z)"), nil
			}
			res, err := svc.Create(ctx, mcp.ParseString(req, "org_id", ""), dtos.CreateContractRequest{
				Title:       mcp.ParseString(req, "title", ""),
				Description: mcp.ParseString(req, "description", ""),
				StartDate:   start,
				EndDate:     end,
			})
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("update_contract",
			mcp.WithDescription("Update an existing contract."),
			mcp.WithString("org_id", mcp.Required(), mcp.Description("Organization ID")),
			mcp.WithString("contract_id", mcp.Required(), mcp.Description("Contract ID")),
			mcp.WithString("title", mcp.Required(), mcp.Description("New title (5–200 chars)")),
			mcp.WithString("description", mcp.Required(), mcp.Description("New description (min 10 chars)")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("New start date (RFC3339)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("New end date (RFC3339, must be after start_date)")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			start, err := time.Parse(time.RFC3339, mcp.ParseString(req, "start_date", ""))
			if err != nil {
				return mcp.NewToolResultError("start_date must be RFC3339"), nil
			}
			end, err := time.Parse(time.RFC3339, mcp.ParseString(req, "end_date", ""))
			if err != nil {
				return mcp.NewToolResultError("end_date must be RFC3339"), nil
			}
			res, err := svc.Update(ctx,
				mcp.ParseString(req, "org_id", ""),
				mcp.ParseString(req, "contract_id", ""),
				dtos.UpdateContractRequest{
					Title:       mcp.ParseString(req, "title", ""),
					Description: mcp.ParseString(req, "description", ""),
					StartDate:   start,
					EndDate:     end,
				},
			)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return jsonResult(res)
		},
	)

	s.AddTool(
		mcp.NewTool("delete_contract",
			mcp.WithDescription("Delete a contract."),
			mcp.WithString("org_id", mcp.Required(), mcp.Description("Organization ID")),
			mcp.WithString("contract_id", mcp.Required(), mcp.Description("Contract ID")),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			if err := svc.Delete(ctx,
				mcp.ParseString(req, "org_id", ""),
				mcp.ParseString(req, "contract_id", ""),
			); err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			return mcp.NewToolResultText("contract deleted"), nil
		},
	)
}

// jsonResult serialises any value as indented JSON text.
// We intentionally avoid mcp.NewToolResultJSON because it sets the
// structuredContent field, which the MCP spec requires to be an object —
// arrays (list responses) would fail protocol validation.
func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return mcp.NewToolResultError("failed to serialise response: " + err.Error()), nil
	}
	return mcp.NewToolResultText(string(b)), nil
}
