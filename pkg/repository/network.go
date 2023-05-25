package repository

import (
	"context"
	"database/sql"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	librepository "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

type Network interface {
	database.Transactable
	List(ctx context.Context, limit int, offset int) (networks []model.NetworkData, err error)
}

type network[T any] struct {
	librepository.Base[T]
}

func NewNetwork(db database.Queryable) Network {
	return &network[model.NetworkData]{librepository.Base[model.NetworkData]{Store: db, Table: "network"}}
}

func (n network[T]) List(ctx context.Context, limit int, offset int) (networks []model.NetworkData, err error) {
	if limit == 0 {
		limit = 20
	}
	results := []model.NetworkFull{}

	err = n.Store.SelectContext(ctx, &results, `SELECT * FROM network WHERE deleted_at IS NULL LIMIT $1 OFFSET $2;`, limit, offset)

	if err == sql.ErrNoRows {
		return nil, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return nil, common.StringError(err)
	}

	for i := range results {
		networks = append(networks, model.NetworkData{Id: results[i].Id, Name: results[i].Name})
	}

	return networks, nil
}
