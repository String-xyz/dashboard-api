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
	Create(ctx context.Context, model model.MemberToOrganization) (model.MemberToOrganization, error)
	GetById(ctx context.Context, ID string) (model.MemberToOrganization, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberToOrganization, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByMember(memberId string) (model.MemberToOrganization, error)
	GetByOrganization(organizationId string) ([]model.MemberToOrganization, error)
}

type memberToOrganization[T any] struct {
	strrepo.Base[T]
}

func NewMemberToOrganization(db database.Queryable) MemberToOrganization {
	return &memberToOrganization[model.MemberToOrganization]{strrepo.Base[model.MemberToOrganization]{Store: db, Table: "member_to_organization"}}
}

func (p memberToOrganization[T]) Create(ctx context.Context, m model.MemberToOrganization) (model.MemberToOrganization, error) {
	newModel := model.MemberToOrganization{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_to_organization (organization_id, member_id) 
		VALUES(:organization_id, :member_id) RETURNING *`, m)

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

func (p memberToOrganization[T]) GetByMember(memberId string) (model.MemberToOrganization, error) {
	m := model.MemberToOrganization{}
	err := p.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE member_id = $1", p.Table), memberId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p memberToOrganization[T]) GetByOrganization(organizationId string) ([]model.MemberToOrganization, error) {
	m := []model.MemberToOrganization{}
	err := p.Store.Select(&m, fmt.Sprintf("SELECT * FROM %s WHERE organization_id = $1", p.Table), organizationId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
