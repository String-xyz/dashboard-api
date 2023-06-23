package service

import (
	"context"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/repository"
	"github.com/String-xyz/go-lib/v2/common"
)

type Contract interface {
	Create(ctx context.Context, create model.RequestContractCreate, callerId string, organizationId string) (model.Contract, error)
	GetAll(ctx context.Context, platformId string, organizationId string) ([]model.Contract, error)
	Get(ctx context.Context, contractId string, organizationId string) (model.Contract, error)
	Deactivate(ctx context.Context, contractId string, callerId string, organizationId string) (model.Contract, error)
	Reactivate(ctx context.Context, contractId string, callerId string, organizationId string) (model.Contract, error)
	Update(ctx context.Context, request model.RequestContractUpdate, contractId string, callerId string, organizationId string) (model.Contract, error)
}

type contract struct {
	repos repository.Repositories
}

func NewContract(repos repository.Repositories) Contract {
	return &contract{repos}
}

func (c contract) Create(ctx context.Context, create model.RequestContractCreate, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Create", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	create.OrganizationId = organizationId

	contract, err := c.repos.Contract.Create(ctx, create)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return contract, nil
}

func (c contract) GetAll(ctx context.Context, platformId string, organizationId string) (contracts []model.Contract, err error) {
	_, finish := Span(ctx, "service.contract.GetAll", SpanTag{"organizationId": organizationId})
	defer finish()

	if platformId != "" {
		contracts, err = c.repos.Contract.ListByPlatform(ctx, platformId, 0, 0)
		if err != nil {
			return contracts, common.StringError(err)
		}
	} else {
		contracts, err = c.repos.Contract.ListByOrganization(ctx, organizationId, 0, 0)
		if err != nil {
			return contracts, common.StringError(err)
		}
	}

	return contracts, nil
}

func (c contract) Get(ctx context.Context, contractId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Get", SpanTag{"organizationId": organizationId})
	defer finish()

	contract, err := c.repos.Contract.GetForOrganization(ctx, contractId, organizationId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return contract, nil
}

func (c contract) Deactivate(ctx context.Context, contractId string, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Deactivate", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	contract, err := c.repos.Contract.Deactivate(ctx, contractId, organizationId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return contract, nil
}

func (c contract) Reactivate(ctx context.Context, contractId string, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Reactivate", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	contract, err := c.repos.Contract.Activate(ctx, contractId, organizationId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return contract, nil
}

func (c contract) Update(ctx context.Context, request model.RequestContractUpdate, contractId string, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Update", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	contract, err := c.repos.Contract.Update(ctx, contractId, organizationId, request)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return contract, nil
}
