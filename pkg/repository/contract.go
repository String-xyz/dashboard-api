package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	"github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

type Contract interface {
	database.Transactable
	Create(ctx context.Context, request model.RequestContractCreate) (contract model.Contract, err error)
	GetForOrganization(ctx context.Context, id string, organizationId string) (contract model.Contract, err error)
	GetById(ctx context.Context, id string) (contract model.Contract, err error)
	GetByAddressAndNetwork(ctx context.Context, address string, networkId string) (contract model.Contract, err error)
	ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (contracts []model.Contract, err error)
	ListByOrganization(ctx context.Context, organizationId string, limit int, offset int) (contracts []model.Contract, err error)
	List(ctx context.Context, limit int, offset int) (contracts []model.Contract, err error)
	Update(ctx context.Context, id string, organizationId string, updates model.RequestContractUpdate) (contract model.Contract, err error)
	Deactivate(ctx context.Context, id string, organizationId string) (model.Contract, error)
	Activate(ctx context.Context, id string, organizationId string) (model.Contract, error)
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

// Theres a lot going on in this function query, so lets break it down.
// with c as ( ... ) is a common table expression (CTE). It allows us to create a temporary table that we can use in the rest of the query.
// The first CTE creates the contract record. It uses the postgres ON CONFLICT clause to prevent duplicate records from being created.
// The second CTE creates the contract_to_platform records. It uses the postgres UNNEST function to create a table from the array of platform ids.
// The third CTE returns the contract record with the array of platform ids.
// The final select statement joins the contract record with the array of platform ids.
// The result is a single contract record with an array of platform ids.
func (c contract[T]) Create(ctx context.Context, request model.RequestContractCreate) (contract model.Contract, err error) {
	rows, err := c.Store.NamedQuery(`
		WITH c AS (
			INSERT INTO contract (name, address, functions, network_id, organization_id)
			VALUES (:name, :address, :functions, :network_id, :organization_id)
			ON CONFLICT (address, organization_id, network_id) DO UPDATE SET name = name WHERE FALSE
			RETURNING *
		),
		platforms AS (
			SELECT UNNEST(:platform_ids::uuid[]) AS platform_id
			FROM platform
			WHERE organization_id = :organization_id
		),
		ctp AS (
			INSERT INTO contract_to_platform (platform_id, contract_id)
			SELECT platforms.platform_id, c.id
			FROM platforms, c
			ON CONFLICT(platform_id, contract_id) DO NOTHING
			RETURNING platform_id, contract_id
		)
		SELECT c.*, array_agg(jctp.platform_id) AS platform_ids
		FROM c
		JOIN contract_to_platform jctp ON c.id = jctp.contract_id
		GROUP BY c.id
	`, request)
	if err != nil {
		return contract, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&contract)
		if err != nil {
			return contract, libcommon.StringError(err)
		}
	}

	// If the contract id is empty, it means the contract already exists.
	if contract.Id == "" {
		return contract, common.StringError(serror.ALREADY_IN_USE)
	}

	defer rows.Close()
	return contract, nil
}

// get the contract by id and ensure it belongs to a platform in the organization
func (c contract[T]) GetForOrganization(ctx context.Context, id string, organizationId string) (contract model.Contract, err error) {
	err = c.Store.GetContext(ctx, &contract, `
		SELECT c.*, array_agg(ctp.platform_id) AS platform_ids
		FROM contract c
		JOIN contract_to_platform ctp ON contract.id = ctp.contract_id
		JOIN platform p ON ctp.platform_id = p.id
		WHERE contract.id = $1
		AND p.organization_id = $2
		GROUP BY contract.id
		`, id, organizationId)

	if err == sql.ErrNoRows {
		return model.Contract{}, nil
	}

	if err != nil {
		return contract, err
	}

	return contract, nil
}

func (c contract[T]) ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (contracts []model.Contract, err error) {
	if limit == 0 {
		limit = 100
	}

	err = c.Store.SelectContext(ctx, &contracts, `
		SELECT c.*, array_agg(ctp.platform_id) AS platform_ids FROM contract c
		LEFT JOIN contract_to_platform ctp 
		ON ctp.contract_id = c.id
		WHERE ctp.platform_id = $1 AND c.deleted_at IS NULL
		GROUP BY c.id
		LIMIT $2 OFFSET $3
		`, platformId, limit, offset)

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
	err = c.Store.SelectContext(ctx, &contracts, `
			SELECT c.*, array_agg(ctp.platform_id) AS platform_ids FROM contract c 
			LEFT JOIN contract_to_platform ctp 
			ON ctp.contract_id = c.id 
			WHERE c.organization_id = $1 AND c.deleted_at IS NULL
		  GROUP BY c.id
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
	err = c.Store.GetContext(ctx, &contract, `
		SELECT c.* array_agg(ctp.platform_id) AS platform_ids FROM contract c
		LEFT JOIN contract_to_platform ctp ON ctp.contract_id = c.id
		WHERE address = $1 AND network_id = $2 AND deleted_at IS NULL LIMIT 1
		GROUP BY c.id 
		`, address, networkId)

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

func (c contract[T]) Deactivate(ctx context.Context, id string, organizationId string) (model model.Contract, err error) {
	err = c.Store.GetContext(ctx, &model, `
		WITH updated_contract AS (
			UPDATE contract c
			SET deactivated_at = NOW()
			WHERE EXISTS (
					SELECT 1
					FROM contract_to_platform ctp
					JOIN platform p ON ctp.platform_id = p.id
					WHERE ctp.contract_id = c.id
					AND p.organization_id = $2
			)
			AND c.id = $1
			RETURNING *
		)
		SELECT uc.*, array_agg(jctp.platform_id) AS platform_ids
		FROM updated_contract uc
		JOIN contract_to_platform jctp ON uc.id = jctp.contract_id
		GROUP BY uc.id
		`, id, organizationId)

	return model, libcommon.StringError(err)
}

func (c contract[T]) Activate(ctx context.Context, id string, organizationId string) (model model.Contract, err error) {
	err = c.Store.GetContext(ctx, &model, `
		WITH updated_contract AS (
			UPDATE contract c
			SET deactivated_at = NULL
			WHERE EXISTS (
				SELECT 1
				FROM contract_to_platform ctp
				JOIN platform p ON ctp.platform_id = p.id
				WHERE ctp.contract_id = c.id
				AND p.organization_id = $2
			)
    	AND c.id = $1
    	RETURNING *
		)
		SELECT uc.*, array_agg(jctp.platform_id) AS platform_ids
		FROM updated_contract uc
		JOIN contract_to_platform jctp ON uc.id = jctp.contract_id
		GROUP BY uc.id
		`, id, organizationId)

	return model, libcommon.StringError(err)
}

func (c contract[T]) Update(ctx context.Context, id string, organizationId string, updates model.RequestContractUpdate) (model model.Contract, err error) {
	names, keyToUpdate := libcommon.KeysAndValues(updates)
	if len(names) == 0 {
		return model, libcommon.StringError(errors.New("no fields to update"))
	}
	query := fmt.Sprintf(`
		WITH updated_contract AS (
		UPDATE contract SET %s  
		WHERE EXISTS (
			SELECT 1
			FROM contract_to_platform
			JOIN platform ON contract_to_platform.platform_id = platform.id
			WHERE contract_to_platform.contract_id = contract.id
			AND platform.organization_id = %s
		)
		AND contract.id = %s AND deleted_at IS NULL
		RETURNING *
		)
		SELECT uc.*, array_agg(jctp.platform_id) AS platform_ids
		FROM updated_contract uc
		JOIN contract_to_platform jctp ON uc.id = jctp.contract_id
		GROUP BY uc.id
		`, strings.Join(names, ", "), organizationId, id)
	err = c.Store.GetContext(ctx, &model, query, keyToUpdate)
	if err != nil {
		return model, err
	}
	return model, err
}
