package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	libcommon "github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type Contract interface {
	database.Transactable
	Create(model.Contract) (model.Contract, error)
	GetById(ctx context.Context, id string) (model.Contract, error)
	ListByPlatformId(ctx context.Context, platformId string, limit int, offset int) ([]model.Contract, error)
	List(ctx context.Context, limit int, offset int) ([]model.Contract, error)
	Update(ctx context.Context, id string, updates any) error
	GetByAddressAndNetworkAndPlatform(address string, networkId string, platformId string) (model.Contract, error)
	Activate(ctx context.Context, ID string) error
}

type contract[T any] struct {
	repository.Base[T]
}

func NewContract(db database.Queryable) Contract {
	return &contract[model.Contract]{repository.Base[model.Contract]{Store: db, Table: "contract"}}
}

func (u contract[T]) Create(insert model.Contract) (model.Contract, error) {
	m := model.Contract{}
	rows, err := u.Store.NamedQuery(`
		INSERT INTO contract (name, address, functions, network_id, platform_id) 
		VALUES(:name, :address, :functions, :network_id, :platform_id) RETURNING *`, insert)
	if err != nil {
		return m, libcommon.StringError(err)
	}
	for rows.Next() {
		err = rows.StructScan(&m)
		if err != nil {
			return m, libcommon.StringError(err)
		}
	}

	defer rows.Close()
	return m, nil
}

func (u contract[T]) ListByPlatformId(ctx context.Context, platformId string, limit int, offset int) ([]model.Contract, error) {
	list := []model.Contract{}
	if limit == 0 {
		limit = 100
	}
	err := u.Store.SelectContext(ctx, &list, fmt.Sprintf("SELECT * FROM %s WHERE platform_id = $1 LIMIT $2 OFFSET $3", u.Table), platformId, limit, offset)
	if err == sql.ErrNoRows {
		return list, nil
	}
	if err != nil {
		return list, err
	}

	return list, nil
}

func (u contract[T]) GetByAddressAndNetworkAndPlatform(address string, networkId string, platformId string) (model.Contract, error) {
	m := model.Contract{}
	err := u.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE address = $1 AND network_id = $2 AND platform_id = $3 LIMIT 1", u.Table), address, networkId, platformId)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	}
	return m, libcommon.StringError(err)
}

func (u contract[T]) GetById(ctx context.Context, ID string) (model.Contract, error) {
	m := model.Contract{}
	err := u.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", u.Table), ID)
	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
