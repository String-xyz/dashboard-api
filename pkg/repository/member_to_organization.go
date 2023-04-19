package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

// May not be needed, unless organization ID changes
type MemberToOrganizationUpdates struct {
	OrganizationId *string `json:"organizationId" db:"organization_id"`
	MemberID       *string `json:"memberId" db:"member_id"`
}

type MemberToOrganization interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberToOrganization) (newModel model.MemberToOrganization, err error)
	GetById(ctx context.Context, id string) (model model.MemberToOrganization, err error)
	List(ctx context.Context, limit int, offset int) (models []model.MemberToOrganization, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByMember(memberId string) (model model.MemberToOrganization, err error)
	GetByOrganization(organizationId string) (models []model.MemberToOrganization, err error)
}

type memberToOrganization[T any] struct {
	strrepo.Base[T]
}

func NewMemberToOrganization(db database.Queryable) MemberToOrganization {
	return &memberToOrganization[model.MemberToOrganization]{strrepo.Base[model.MemberToOrganization]{Store: db, Table: "member_to_organization"}}
}

func (m memberToOrganization[T]) Create(ctx context.Context, model model.MemberToOrganization) (newModel model.MemberToOrganization, err error) {
	rows, err := m.Store.NamedQuery(`
		INSERT INTO member_to_organization (organization_id, member_id) 
		VALUES(:organization_id, :member_id) RETURNING *`, model)

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

func (m memberToOrganization[T]) GetByMember(memberId string) (model model.MemberToOrganization, err error) {
	err = m.Store.Get(&model, fmt.Sprintf("SELECT * FROM %s WHERE member_id = $1", m.Table), memberId)
	if err != nil && err == sql.ErrNoRows {
		return model, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return model, common.StringError(err)
	}
	return model, nil
}

func (m memberToOrganization[T]) GetByOrganization(organizationId string) (results []model.MemberToOrganization, err error) {
	err = m.Store.Select(&results, fmt.Sprintf("SELECT * FROM %s WHERE organization_id = $1", m.Table), organizationId)
	if err != nil && err == sql.ErrNoRows {
		return results, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return results, common.StringError(err)
	}
	return results, nil
}
