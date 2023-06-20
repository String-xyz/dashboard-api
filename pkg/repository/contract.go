package repository

import (
	"context"
	"database/sql"

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
	GetByAddressAndNetwork(ctx context.Context, address string, networkId string) (contract model.Contract, err error)
	Activate(ctx context.Context, id string) error
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

func (c contract[T]) Create(request model.Contract) (contract model.Contract, err error) {
	c.MustBegin()
	defer c.Reset()

	// Insert contract
	rows, err := c.Store.NamedQuery(`
		INSERT INTO contract (name, address, functions, network_id, organization_id) 
		VALUES(:name, :address, :functions, :network_id, :organization_id) RETURNING *`, request)
	if err != nil {
		c.Rollback()
		return contract, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&contract)
		if err != nil {
			c.Rollback()
			return contract, libcommon.StringError(err)
		}
	}
	rows.Close()

	// Insert platforms for contract
	for _, platformId := range request.PlatformIds {
		rows, err = c.Store.NamedQuery(`
			INSERT INTO contract_to_platform (contract_id, platform_id)
			VALUES(:contract_id, :platform_id) RETURNING *
			USING (SELECT organization_id FROM platform WHERE id = :platform_id) org_id
			WHERE contract.organization_id = org_id`,
			map[string]interface{}{"contract_id": contract.Id, "platform_id": platformId})
		if err != nil {
			c.Rollback()
			return contract, libcommon.StringError(err)
		}
		for rows.Next() {
			err = rows.StructScan(&contract)
			if err != nil {
				c.Rollback()
				return contract, libcommon.StringError(err)
			}
		}
		rows.Close()
	}

	if err := c.Commit(); err != nil {
		return contract, libcommon.StringError(err)
	}

	return contract, nil
}

func (c contract[T]) ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (contracts []model.Contract, err error) {
	if limit == 0 {
		limit = 100
	}
	// Query contracts along with their associated platform IDs
	err = c.Store.SelectContext(ctx, &contracts,
		`SELECT contract.*, array_agg(contract_to_platform.platform_id) AS platform_ids
			FROM contract
			JOIN contract_to_platform ON contract.id = contract_to_platform.contract_id
			WHERE contract_to_platform.platform_id = $1 AND contract.deleted_at IS NULL
			GROUP BY contract.id
			LIMIT $2 OFFSET $3`,
		platformId, limit, offset)

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

	// Query contracts along with their associated platform IDs
	err = c.Store.SelectContext(ctx, &contracts,
		`SELECT contract.*, array_agg(distinct contract_to_platform.platform_id) AS platform_ids
			FROM contract
			JOIN contract_to_platform ON contract.id = contract_to_platform.contract_id
			JOIN platform ON contract_to_platform.platform_id = platform.id
			WHERE platform.organization_id = $1 AND contract.deleted_at IS NULL
			GROUP BY contract.id
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

func (c contract[T]) GetByAddressAndNetwork(ctx context.Context, address string, networkId string) (contract model.Contract, err error) {
	err = c.Store.GetContext(ctx, &contract,
		`SELECT contract.*, array_agg(contract_to_platform.platform_id) AS platform_ids
			FROM contract
			JOIN contract_to_platform ON contract.id = contract_to_platform.contract_id
			WHERE address = $1 AND network_id = $2 AND contract.deleted_at IS NULL
			GROUP BY contract.id`,
		address, networkId)

	if err == sql.ErrNoRows {
		return model.Contract{}, serror.NOT_FOUND
	}
	if err != nil {
		return contract, libcommon.StringError(err)
	}

	return contract, nil
}

func (c contract[T]) GetById(ctx context.Context, id string) (contract model.Contract, err error) {
	err = c.Store.GetContext(ctx, &contract,
		`SELECT contract.*, array_agg(contract_to_platform.platform_id) AS platform_ids
			FROM contract
			JOIN contract_to_platform ON contract.id = contract_to_platform.contract_id
			WHERE contract.id = $1 AND contract.deleted_at IS NULL
			GROUP BY contract.id`,
		id)

	if err == sql.ErrNoRows {
		return model.Contract{}, common.StringError(serror.NOT_FOUND)
	}
	if err != nil {
		return contract, common.StringError(err)
	}

	return contract, nil
}
