package repository

import (
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/jmoiron/sqlx/types"
)

type PlaformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Type          *string         `json:"type" db:"type"`
	Status        *string         `json:"status" db:"status"`
	Tags          *types.JSONText `json:"tags" db:"tags"`
}

type Platform interface {
	database.Transactable
	Create(model.Platform) (model.Platform, error)
	GetById(ID string) (model.Platform, error)
	List(limit int, offset int) ([]model.Platform, error)
	Update(ID string, updates any) error
}

type platform[T any] struct {
	strrepo.Base[T]
}

func NewPlatform(db database.Queryable) Platform {
	return &platform[model.Platform]{strrepo.Base[model.Platform]{Store: db, Table: "platform"}}
}

func (p platform[T]) Create(m model.Platform) (model.Platform, error) {
	plat := model.Platform{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO platform (type, authentication, api_key, status) 
		VALUES(:type, :authentication, :api_key, :status) RETURNING *`, m)

	if err != nil {
		return plat, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&plat)
		if err != nil {
			return plat, common.StringError(err)
		}
	}

	return plat, nil
}
