package repository

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
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
	List(ctx context.Context, limit int, offset int) ([]model.Apikey, error)
	Update(ctx context.Context, ID string, updates any) error
}

type apikey[T any] struct {
	strrepo.Base[T]
}

func NewApikey(db database.Queryable) Apikey {
	return &apikey[model.Apikey]{strrepo.Base[model.Apikey]{Store: db, Table: "apikey"}}
}

func (p apikey[T]) Create(ctx context.Context, m model.Apikey) (model.Apikey, error) {
	newModel := model.Apikey{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_invite (type, data, description, created_by, platform_id) 
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
