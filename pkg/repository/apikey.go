package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type ApikeyUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string    `json:"type" db:"type"`
	Data          *string    `json:"data" db:"data"`
	Description   *string    `json:"description" db:"description"`
	CreatedBy     *string    `json:"createdBy" db:"created_by"`
	PlatformID    *string    `json:"platformId" db:"platform_id"`
}

type Apikey interface {
	database.Transactable
	Create(ctx context.Context, model model.Apikey) (model.Apikey, error)
	GetById(ctx context.Context, ID string) (model.Apikey, error)
	List(ctx context.Context, platformId string, limit int, offset int) ([]model.Apikey, error)
	Update(ctx context.Context, ID string, updates any) error
}

type apikey[T any] struct {
	strrepo.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{strrepo.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (a apikey[T]) Create(ctx context.Context, m model.Apikey) (model.Apikey, error) {
	newModel := model.Apikey{}
	rows, err := a.Store.NamedQuery(`
		INSERT INTO apikey (type, data, description, created_by, platform_id) 
		VALUES(:type, :data, :description, :created_by, :platform_id) RETURNING *`, m)

	if err != nil {
		return newModel, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&newModel)
		if err != nil {
			return newModel, common.StringError(err)
		}
	}

	return newModel, nil
}

func (p apikey[T]) GetById(ctx context.Context, ID string) (model.Apikey, error) {
	m := model.Apikey{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1", p.Table), ID)
	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p apikey[T]) List(ctx context.Context, platformId string, limit int, offset int) ([]model.Apikey, error) {
	m := []model.Apikey{}
	if limit == 0 {
		limit = 20
	}

	err := p.Store.SelectContext(ctx, &m, `SELECT * FROM apikey WHERE apikey.platform_id = $1 LIMIT $2 OFFSET $3;`, platformId, limit, offset)

	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	return m, nil
}
