package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	libcommon "github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	librepository "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

type ApikeyUpdates struct {
	DeactivatedAt  *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type           *string    `json:"type" db:"type"`
	Data           *string    `json:"data" db:"data"`
	Description    *string    `json:"description" db:"description"`
	CreatedBy      *string    `json:"createdBy" db:"created_by"`
	PlatformId     *string    `json:"platformId" db:"platform_id"`
	OrganizationId *string    `json:"organizationId" db:"organization_id"`
}

type Apikey interface {
	database.Transactable
	Create(ctx context.Context, request model.Apikey) (key model.Apikey, err error)
	GetById(ctx context.Context, id string) (key model.Apikey, err error)
	GetByData(ctx context.Context, data string, keyType string) (model.Apikey, error)
	ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (keys []model.Apikey, err error)
	ListByOrganization(ctx context.Context, organizationId string, limit int, offset int) (keys []model.Apikey, err error)
	Update(ctx context.Context, id string, updates any) error
	SoftDelete(ctx context.Context, ID string) error
}

type apikey[T any] struct {
	librepository.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{librepository.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (a apikey[T]) Create(ctx context.Context, request model.Apikey) (key model.Apikey, err error) {
	var query string
	var args []interface{}

	if request.Type == "secret" {
		query, args, err = a.Named(`
			INSERT INTO apikey (type, data, hint, description, created_by, organization_id) 
			VALUES(:type, :data, :hint, :description, :created_by, :organization_id) RETURNING *`, request)
	} else {
		query, args, err = a.Named(`
			INSERT INTO apikey (type, data, hint, description, created_by, platform_id, organization_id) 
			VALUES(:type, :data, :hint, :description, :created_by, :platform_id, :organization_id) RETURNING *`, request)
	}
	if err != nil {
		return key, libcommon.StringError(err)
	}

	// Use QueryRowxContext to execute the query with the provided context
	err = a.Store.QueryRowxContext(ctx, query, args...).StructScan(&key)
	if err != nil {
		return key, libcommon.StringError(err)
	}

	return key, nil
}

func (a apikey[T]) GetById(ctx context.Context, id string) (key model.Apikey, err error) {
	err = a.Store.GetContext(ctx, &key, fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deleted_at IS NULL", a.Table), id)
	if err == sql.ErrNoRows {
		return key, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return key, common.StringError(err)
	}
	return key, nil
}

func (p apikey[T]) GetByData(ctx context.Context, data string, keyType string) (model.Apikey, error) {
	m := model.Apikey{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE data = $1 AND type = $2", p.Table), data, keyType)
	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (a apikey[T]) ListByPlatform(ctx context.Context, platformId string, limit int, offset int) (keys []model.Apikey, err error) {
	if limit == 0 {
		limit = 20
	}

	err = a.Store.SelectContext(ctx, &keys, `SELECT * FROM apikey WHERE apikey.platform_id = $1 AND apikey.deleted_at IS NULL  LIMIT $2 OFFSET $3;`, platformId, limit, offset)

	if err == sql.ErrNoRows {
		return keys, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return keys, common.StringError(err)
	}

	return keys, nil
}

func (a apikey[T]) ListByOrganization(ctx context.Context, organizationId string, limit int, offset int) (keys []model.Apikey, err error) {
	if limit == 0 {
		limit = 20
	}

	err = a.Store.SelectContext(ctx, &keys, `SELECT * FROM apikey WHERE apikey.organization_id = $1 AND apikey.deleted_at IS NULL LIMIT $2 OFFSET $3;`, organizationId, limit, offset)

	if err == sql.ErrNoRows {
		return keys, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return keys, common.StringError(err)
	}

	return keys, nil
}
