package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	librepository "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

// May not be needed, unless organization id changes
type MemberToOrganizationUpdates struct {
	OrganizationId *string `json:"organizationId" db:"organization_id"`
	MemberId       *string `json:"memberId" db:"member_id"`
}

type MemberToOrganization interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberToOrganization) (relation model.MemberToOrganization, err error)
	GetById(ctx context.Context, id string) (relation model.MemberToOrganization, err error)
	List(ctx context.Context, limit int, offset int) (relations []model.MemberToOrganization, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByMember(memberId string) (relation model.MemberToOrganization, err error)
	GetByOrganization(organizationId string) (relations []model.MemberToOrganization, err error)
}

type memberToOrganization[T any] struct {
	librepository.Base[T]
}

func NewMemberToOrganization(db database.Queryable) MemberToOrganization {
	return &memberToOrganization[model.MemberToOrganization]{librepository.Base[model.MemberToOrganization]{Store: db, Table: "member_to_organization"}}
}

func (m memberToOrganization[T]) Create(ctx context.Context, request model.MemberToOrganization) (relation model.MemberToOrganization, err error) {
	rows, err := m.Store.NamedQuery(`
		INSERT INTO member_to_organization (organization_id, member_id) 
		VALUES(:organization_id, :member_id) RETURNING *`, request)

	if err != nil {
		return relation, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&relation)
		if err != nil {
			return relation, common.StringError(err)
		}
	}

	return relation, nil
}

func (m memberToOrganization[T]) GetByMember(memberId string) (relation model.MemberToOrganization, err error) {
	err = m.Store.Get(&relation, fmt.Sprintf("SELECT * FROM %s WHERE member_id = $1", m.Table), memberId)
	if err != nil && err == sql.ErrNoRows {
		return relation, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return relation, common.StringError(err)
	}
	return relation, nil
}

func (m memberToOrganization[T]) GetByOrganization(organizationId string) (relations []model.MemberToOrganization, err error) {
	err = m.Store.Select(&relations, fmt.Sprintf("SELECT * FROM %s WHERE organization_id = $1", m.Table), organizationId)
	if err != nil && err == sql.ErrNoRows {
		return relations, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return relations, common.StringError(err)
	}
	return relations, nil
}
