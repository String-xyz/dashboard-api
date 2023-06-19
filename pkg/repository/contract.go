package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	"github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

type Contract interface {
	database.Transactable
	Create(request model.Contract) (contract model.Contract, err error)
	GetById(ctx context.Context, id string) (contract model.Contract, err error)
	ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (contracts []model.Contract, err error)
	ListByOrganization(ctx context.Context, organizationId string, limit int, offset int) (contracts []model.Contract, err error)
	List(ctx context.Context, limit int, offset int) (contracts []model.Contract, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByAddressAndNetworkAndPlatform(ctx context.Context, address string, networkId string, platformId string) (contract model.Contract, err error)
	Activate(ctx context.Context, id string) error
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

func (c contract[T]) Create(request model.Contract) (contract model.Contract, err error) {
	rows, err := c.Store.NamedQuery(`
		INSERT INTO contract (name, address, functions, network_id, platform_id) 
		VALUES(:name, :address, :functions, :network_id, :platform_id) RETURNING *`, request)
	if err != nil {
		return contract, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&contract)
		if err != nil {
			return contract, libcommon.StringError(err)
		}
	}

	defer rows.Close()
	return contract, nil
}

func (c contract[T]) ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (contracts []model.Contract, err error) {
	if limit == 0 {
		limit = 100
	}
	err = c.Store.SelectContext(ctx, &contracts, fmt.Sprintf("SELECT * FROM %s WHERE platform_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3", c.Table), platformId, limit, offset)
	if err == sql.ErrNoRows {
		return []model.Contract{}, nil
	}
	if err != nil {
		return contracts, err
	}

	return contracts, nil
}

func (c contract[T]) ListByOrganization(ctx context.Context, organizationId string, limit int, offset int) (contracts []model.Contract, err error) {
	if limit == 0 {
		limit = 100
	}
	err = c.Store.SelectContext(ctx, &contracts,
		`SELECT contract.* FROM contract 
			LEFT JOIN platform 
			ON contract.platform_id = platform.id 
			WHERE platform.organization_id = $1 AND contract.deleted_at IS NULL
			LIMIT $2 OFFSET $3`,
		organizationId, limit, offset)

	if err == sql.ErrNoRows {
		return []model.Contract{}, nil
	}
	if err != nil {
		return contracts, err
	}

	return contracts, nil
}

func (c contract[T]) GetByAddressAndNetworkAndPlatform(ctx context.Context, address string, networkId string, platformId string) (contract model.Contract, err error) {
	err = c.Store.GetContext(ctx, &contract, fmt.Sprintf("SELECT * FROM %s WHERE address = $1 AND network_id = $2 AND platform_id = $3 AND deleted_at IS NULL LIMIT 1", c.Table), address, networkId, platformId)
	if err != nil && err == sql.ErrNoRows {
		return contract, serror.NOT_FOUND
	}
	return contract, libcommon.StringError(err)
}

func (c contract[T]) GetById(ctx context.Context, id string) (contract model.Contract, err error) {
	err = c.Store.GetContext(ctx, &contract, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deleted_at IS NULL", c.Table), id)
	if err == sql.ErrNoRows {
		return contract, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return contract, common.StringError(err)
	}
	return contract, nil
}
