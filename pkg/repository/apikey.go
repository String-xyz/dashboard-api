package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	librepository "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type ApikeyUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Data          *string    `json:"data" db:"data"`
	Description   *string    `json:"description" db:"description"`
	CreatedBy     *string    `json:"createdBy" db:"created_by"`
	PlatformId    *string    `json:"platformId" db:"platform_id"`
}

type Apikey interface {
	database.Transactable
	Create(ctx context.Context, request model.Apikey) (key model.Apikey, err error)
	GetById(ctx context.Context, id string) (key model.Apikey, err error)
	List(ctx context.Context, platformId string, limit int, offset int) (keys []model.Apikey, err error)
	Update(ctx context.Context, id string, updates any) error
}

type apikey[T any] struct {
	librepository.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{librepository.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (k apikey[T]) Create(ctx context.Context, request model.Apikey) (key model.Apikey, err error) {
	rows, err := k.Store.NamedQuery(`
		INSERT INTO apikey (type, data, hint, description, created_by, platform_id) 
		VALUES(:type, :data, :hint, :description, :created_by, :platform_id) RETURNING *`, request)

	if err != nil {
		return key, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&key)
		if err != nil {
			return key, common.StringError(err)
		}
	}

	return key, nil
}

func (k apikey[T]) GetById(ctx context.Context, id string) (key model.Apikey, err error) {
	err = k.Store.GetContext(ctx, &key, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", k.Table), id)
	if err == sql.ErrNoRows {
		return key, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return key, common.StringError(err)
	}
	return key, nil
}

func (k apikey[T]) List(ctx context.Context, platformId string, limit int, offset int) (keys []model.Apikey, err error) {
	if limit == 0 {
		limit = 20
	}

	err = k.Store.SelectContext(ctx, &keys, `SELECT * FROM apikey WHERE apikey.platform_id = $1 LIMIT $2 OFFSET $3;`, platformId, limit, offset)

	if err == sql.ErrNoRows {
		return keys, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return keys, common.StringError(err)
	}

	return keys, nil
}
