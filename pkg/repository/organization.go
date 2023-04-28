package repository

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	librepository "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type OrganizationUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Name          *string    `json:"name" db:"name"`
	Description   *string    `json:"description" db:"description"`
}

type Organization interface {
	database.Transactable
	Create(ctx context.Context, request model.Organization) (org model.Organization, err error)
	GetById(ctx context.Context, id string) (org model.Organization, err error)
	List(ctx context.Context, limit int, offset int) (orgs []model.Organization, err error)
	Update(ctx context.Context, id string, updates any) error
}

type organization[T any] struct {
	librepository.Base[T]
}

func NewOrganization(db database.Queryable) Organization {
	return &organization[model.Organization]{librepository.Base[model.Organization]{Store: db, Table: "organization"}}
}

func (o organization[T]) Create(ctx context.Context, request model.Organization) (org model.Organization, err error) {
	rows, err := o.Store.NamedQuery(`
		INSERT INTO organization (name) 
		VALUES(:name) RETURNING *`, request)

	if err != nil {
		return org, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&org)
		if err != nil {
			return org, common.StringError(err)
		}
	}

	return org, nil
}
