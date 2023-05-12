package service

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
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

	// Check if contract already exists
	exists, err := c.repos.Contract.GetByAddressAndNetworkAndPlatform(ctx, create.Address, create.NetworkId, create.PlatformId)
	if err != nil && err != serror.NOT_FOUND {
		return model.Contract{}, common.StringError(err)
	} else if exists.Id != "" {
		return model.Contract{}, common.StringError(serror.ALREADY_IN_USE)
	}

	row := model.Contract{Name: create.Name, PlatformId: create.PlatformId, Address: create.Address, Functions: create.Functions, NetworkId: create.NetworkId}
	row, err = c.repos.Contract.Create(row)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	return row, nil
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

	contract, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	platform, err := c.repos.Platform.GetById(ctx, contract.PlatformId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return model.Contract{}, common.StringError(serror.FORBIDDEN)
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

	deactivate, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	platform, err := c.repos.Platform.GetById(ctx, deactivate.PlatformId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return model.Contract{}, common.StringError(serror.FORBIDDEN)
	}

	type DeactivateUpdate struct {
		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	}

	now := time.Now()
	update := DeactivateUpdate{DeactivatedAt: &now}

	err = c.repos.Contract.Update(ctx, contractId, update)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	deactivate, err = c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return deactivate, nil
}

func (c contract) Reactivate(ctx context.Context, contractId string, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Reactivate", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	reactivate, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	platform, err := c.repos.Platform.GetById(ctx, reactivate.PlatformId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return model.Contract{}, common.StringError(serror.FORBIDDEN)
	}

	err = c.repos.Contract.Activate(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	reactivate, err = c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return reactivate, nil
}

func (c contract) Update(ctx context.Context, request model.RequestContractUpdate, contractId string, callerId string, organizationId string) (model.Contract, error) {
	_, finish := Span(ctx, "service.contract.Update", SpanTag{"organizationId": organizationId})
	defer finish()

	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	update, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	platform, err := c.repos.Platform.GetById(ctx, update.PlatformId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return model.Contract{}, common.StringError(serror.FORBIDDEN)
	}

	err = c.repos.Contract.Update(ctx, contractId, request)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	update, err = c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	return update, nil
}
