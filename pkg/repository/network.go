package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	librepository "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type NetworkFull struct {
	Id            string     `json:"id" db:"id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string     `json:"name" db:"name"`
	NetworkId     uint64     `json:"networkId" db:"network_id"`
	ChainId       uint64     `json:"chainId" db:"chain_id"`
	GasTokenId    string     `json:"gasTokenId" db:"gas_token_id"`
	GasOracle     string     `json:"gasOracle" db:"gas_oracle"`
	RPCUrl        string     `json:"rpcUrl" db:"rpc_url"`
	ExplorerUrl   string     `json:"explorerUrl" db:"explorer_url"`
}

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
	result := []NetworkFull{}

	err = n.Store.SelectContext(ctx, &result, `SELECT * FROM network LIMIT $1 OFFSET $2;`, limit, offset)

	if err == sql.ErrNoRows {
		return nil, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return nil, common.StringError(err)
	}

	for i := range networks {
		networks = append(networks, model.NetworkData{Id: result[i].Id, Name: result[i].Name})
	}
	return networks, nil
}
