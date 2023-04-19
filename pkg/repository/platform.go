package repository

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/lib/pq"
)

type PlaformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	ActivatedAt   *time.Time      `json:"activatedAt" db:"activated_at"`
	Name          *string         `json:"name" db:"name"`
	Description   *string         `json:"description" db:"description"`
	Domains       *pq.StringArray `json:"domains" db:"domains"`
	IPAddresses   *pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type Platform interface {
	database.Transactable
	Create(ctx context.Context, request model.Platform) (platform model.Platform, err error)
	GetById(ctx context.Context, id string) (platform model.Platform, err error)
	List(ctx context.Context, limit int, offset int) (platforms []model.Platform, err error)
	Update(ctx context.Context, id string, updates any) error
}

type platform[T any] struct {
	strrepo.Base[T]
}

func NewPlatform(db database.Queryable) Platform {
	return &platform[model.Platform]{strrepo.Base[model.Platform]{Store: db, Table: "platform"}}
}

func (p platform[T]) Create(ctx context.Context, request model.Platform) (platform model.Platform, err error) {
	rows, err := p.Store.NamedQuery(`
		INSERT INTO platform (name) 
		VALUES(:name) RETURNING *`, request)

	if err != nil {
		return platform, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&platform)
		if err != nil {
			return platform, common.StringError(err)
		}
	}

	return platform, nil
}
