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
	"github.com/lib/pq"
)

type PlaformUpdates struct {
	DeactivatedAt *time.Time      `json:"deactivatedAt" db:"deactivated_at"`
	Name          *string         `json:"name" db:"name"`
	Description   *string         `json:"description" db:"description"`
	Domains       *pq.StringArray `json:"domains" db:"domains"`
	IPAddresses   *pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type Platform interface {
	database.Transactable
	Create(ctx context.Context, request model.Platform) (platform model.Platform, err error)
	GetById(ctx context.Context, id string) (platform model.Platform, err error)
	List(ctx context.Context, organizationId string, limit int, offset int) (platforms []model.Platform, err error)
	Update(ctx context.Context, id string, updates any) error
}

type platform[T any] struct {
	librepository.Base[T]
}

func NewPlatform(db database.Queryable) Platform {
	return &platform[model.Platform]{librepository.Base[model.Platform]{Store: db, Table: "platform"}}
}

func (p platform[T]) Create(ctx context.Context, request model.Platform) (platform model.Platform, err error) {
	rows, err := p.Store.NamedQuery(`
		INSERT INTO platform (name, description, organization_id) 
		VALUES(:name, :description, :organization_id) RETURNING *`, request)

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

func (p platform[T]) List(ctx context.Context, organizationId string, limit int, offset int) (platforms []model.Platform, err error) {
	if limit == 0 {
		limit = 20
	}

	err = p.Store.SelectContext(ctx, &platforms, `SELECT * FROM platform WHERE platform.organization_id = $1 LIMIT $2 OFFSET $3;`, organizationId, limit, offset)

	if err == sql.ErrNoRows {
		return platforms, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return platforms, common.StringError(err)
	}

	return platforms, nil
}
