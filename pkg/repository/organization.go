package repository

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type OrganizationUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	ActivatedAt   *time.Time `json:"activatedAt" db:"activated_at"`
	Name          *string    `json:"name" db:"name"`
	Description   *string    `json:"description" db:"description"`
}

type Organization interface {
	database.Transactable
	Create(ctx context.Context, model model.Organization) (model.Organization, error)
	GetById(ctx context.Context, ID string) (model.Organization, error)
	List(ctx context.Context, limit int, offset int) ([]model.Organization, error)
	Update(ctx context.Context, ID string, updates any) error
}

type organization[T any] struct {
	strrepo.Base[T]
}

func NewOrganization(db database.Queryable) Organization {
	return &organization[model.Organization]{strrepo.Base[model.Organization]{Store: db, Table: "organization"}}
}

func (o organization[T]) Create(ctx context.Context, m model.Organization) (model.Organization, error) {
	newModel := model.Organization{}
	rows, err := o.Store.NamedQuery(`
		INSERT INTO organization (name) 
		VALUES(:name) RETURNING *`, m)

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
