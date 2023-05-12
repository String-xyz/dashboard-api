package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	librepository "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type PlaformMemberUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	Email         *string    `json:"email" db:"email"`
	Password      *string    `json:"password" db:"password"`
}

type OrganizationMemberWithRole struct {
	model.OrganizationMember
	Role string `json:"role" db:"member_role"`
}

type OrganizationMember interface {
	database.Transactable
	Create(ctx context.Context, request model.OrganizationMember) (member model.OrganizationMember, err error)
	GetById(ctx context.Context, id string) (member OrganizationMemberWithRole, err error)
	List(ctx context.Context, organizationId string, limit int, offset int) (members []OrganizationMemberWithRole, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByEmail(ctx context.Context, email string) (member OrganizationMemberWithRole, err error)
	Deactivate(ctx context.Context, id string) error
	Activate(ctx context.Context, id string) error
}

type organizationMember[T any] struct {
	librepository.Base[T]
}

func NewOrganizationMember(db database.Queryable) OrganizationMember {
	return &organizationMember[model.OrganizationMember]{librepository.Base[model.OrganizationMember]{Store: db, Table: "organization_member"}}
}

func (m organizationMember[T]) Create(ctx context.Context, request model.OrganizationMember) (member model.OrganizationMember, err error) {
	rows, err := m.Store.NamedQuery(`
		INSERT INTO organization_member (email, name, password) 
		VALUES(:email, :name, :password) RETURNING *`, request)

	if err != nil {
		return member, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&member)
		if err != nil {
			return member, common.StringError(err)
		}
	}

	return member, nil
}

func (m organizationMember[T]) GetByEmail(ctx context.Context, email string) (member OrganizationMemberWithRole, err error) {
	err = m.Store.GetContext(ctx, &member, `
		SELECT organization_member.*, member_role.name AS member_role
		FROM organization_member
		LEFT JOIN member_to_role
		ON organization_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE organization_member.email = $1 AND organization_member.deleted_at IS NULL`, email)
	if err != nil && err == sql.ErrNoRows {
		return member, serror.NOT_FOUND
	} else if err != nil {
		return member, common.StringError(err)
	}
	return member, nil
}

func (m organizationMember[T]) GetById(ctx context.Context, id string) (member OrganizationMemberWithRole, err error) {
	err = m.Store.GetContext(ctx, &member, `
		SELECT organization_member.*, member_role.name AS member_role
		FROM organization_member
		LEFT JOIN member_to_role
		ON organization_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE organization_member.id = $1 AND organization_member.deleted_at IS NULL`, id)
	if err != nil && err == sql.ErrNoRows {
		return member, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return member, common.StringError(err)
	}
	return member, nil
}

func (m organizationMember[T]) List(ctx context.Context, organizationId string, limit int, offset int) (members []OrganizationMemberWithRole, err error) {
	if limit == 0 {
		limit = 20
	}

	err = m.Store.SelectContext(ctx, &members, `
		SELECT organization_member.*, member_role.name AS member_role
		FROM organization_member
		LEFT JOIN member_to_role
		ON organization_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		LEFT JOIN member_to_organization
		ON organization_member.id = member_to_organization.member_id 
		LEFT JOIN organization
		ON organization.id = member_to_organization.organization_id
		WHERE organization.id = $1 AND organization_member.deleted_at IS NULL
		LIMIT $2
		OFFSET $3;`, organizationId, limit, offset)

	if err == sql.ErrNoRows {
		return members, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return members, common.StringError(err)
	}

	return members, nil
}
