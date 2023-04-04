package service

import (
	"context"
	"errors"
	"time"

	"github.com/String-xyz/go-lib/common"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Contract interface {
	Create(ctx context.Context, create model.RequestContractCreate, callerId string, platformId string) (model.Contract, error)
	GetAll(ctx context.Context, platformId string) ([]model.Contract, error)
	Get(ctx context.Context, platformId string, contractId string) (model.Contract, error)
	Deactivate(ctx context.Context, callerId string, platformId string, contractId string) (model.Contract, error)
	Reactivate(ctx context.Context, callerId string, platformId string, contractId string) (model.Contract, error)
	Update(ctx context.Context, request model.RequestContractUpdate, callerId string, platformId string, contractId string) (model.Contract, error)
}

type contract struct {
	repos repository.Repositories
}

func NewContract(repos repository.Repositories) Contract {
	return &contract{repos}
}

func (c contract) Create(ctx context.Context, create model.RequestContractCreate, callerId string, platformId string) (model.Contract, error) {
	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	// Check if contract already exists
	exists, err := c.repos.Contract.GetByAddressAndNetworkAndPlatform(ctx, create.Address, create.NetworkID, platformId)
	if err != nil && err != serror.NOT_FOUND {
		return model.Contract{}, common.StringError(err)
	} else if exists.ID != "" {
		return model.Contract{}, common.StringError(errors.New("contract already exists"))
	}

	row := model.Contract{Name: create.Name, PlatformID: platformId, Address: create.Address, Functions: create.Functions, NetworkID: create.NetworkID}
	row, err = c.repos.Contract.Create(row)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	return row, nil
}

func (c contract) GetAll(ctx context.Context, platformId string) ([]model.Contract, error) {
	contracts, err := c.repos.Contract.ListByPlatformId(ctx, platformId, 0, 0)
	if err != nil {
		return nil, common.StringError(err)
	}

	return contracts, nil
}

func (c contract) Get(ctx context.Context, platformId string, contractId string) (model.Contract, error) {
	contract, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if contract.PlatformID != platformId {
		return model.Contract{}, common.StringError(errors.New("contract not maintained by accessing platform"))
	}
	return contract, nil
}

func (c contract) Deactivate(ctx context.Context, callerId string, platformId string, contractId string) (model.Contract, error) {
	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	deactivate, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if deactivate.PlatformID != platformId {
		return model.Contract{}, common.StringError(errors.New("contract not maintained by accessing platform"))
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

func (c contract) Reactivate(ctx context.Context, callerId string, platformId string, contractId string) (model.Contract, error) {
	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	reactivate, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if reactivate.PlatformID != platformId {
		return model.Contract{}, common.StringError(errors.New("contract not maintained by accessing platform"))
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

func (c contract) Update(ctx context.Context, request model.RequestContractUpdate, callerId string, platformId string, contractId string) (model.Contract, error) {
	err := RequireAuthority(c.repos, callerId, "Admin", "Owner")
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}

	update, err := c.repos.Contract.GetById(ctx, contractId)
	if err != nil {
		return model.Contract{}, common.StringError(err)
	}
	if update.PlatformID != platformId {
		return model.Contract{}, common.StringError(errors.New("contract not maintained by accessing platform"))
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
