package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type PlaformMemberUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	ActivatedAt   *time.Time `json:"activatedAt" db:"activated_at"`
	Email         *string    `json:"email" db:"email"`
	Password      *string    `json:"password" db:"password"`
}

type OrganizationMemberWithRole struct {
	model.OrganizationMember
	Role string `json:"role" db:"member_role"`
}

type OrganizationMember interface {
	database.Transactable
	Create(ctx context.Context, model model.OrganizationMember) (model.OrganizationMember, error)
	GetById(ctx context.Context, ID string) (OrganizationMemberWithRole, error)
	GetByIdIncludingDeactivated(ctx context.Context, ID string) (OrganizationMemberWithRole, error)
	List(ctx context.Context, organizationId string, limit int, offset int) ([]OrganizationMemberWithRole, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByEmail(ctx context.Context, email string) (OrganizationMemberWithRole, error)
	Deactivate(ctx context.Context, ID string) error
	Activate(ctx context.Context, ID string) error
}

type organizationMember[T any] struct {
	strrepo.Base[T]
}

func NewOrganizationMember(db database.Queryable) OrganizationMember {
	return &organizationMember[model.OrganizationMember]{strrepo.Base[model.OrganizationMember]{Store: db, Table: "organization_member"}}
}

func (o organizationMember[T]) Create(ctx context.Context, m model.OrganizationMember) (model.OrganizationMember, error) {
	newModel := model.OrganizationMember{}
	rows, err := o.Store.NamedQuery(`
		INSERT INTO organization_member (email, name, password) 
		VALUES(:email, :name, :password) RETURNING *`, m)

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

func (o organizationMember[T]) GetByEmail(ctx context.Context, email string) (OrganizationMemberWithRole, error) {
	m := OrganizationMemberWithRole{}
	err := o.Store.GetContext(ctx, &m, `
		SELECT organization_member.*, member_role.name AS member_role
		FROM organization_member
		LEFT JOIN member_to_role
		ON organization_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE organization_member.email = $1`, email)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (o organizationMember[T]) GetById(ctx context.Context, ID string) (OrganizationMemberWithRole, error) {
	m := OrganizationMemberWithRole{}
	err := o.Store.GetContext(ctx, &m, `
		SELECT organization_member.*, member_role.name AS member_role
		FROM organization_member
		LEFT JOIN member_to_role
		ON organization_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE organization_member.id = $1`, ID)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (o organizationMember[T]) GetByIdIncludingDeactivated(ctx context.Context, ID string) (OrganizationMemberWithRole, error) {
	m := OrganizationMemberWithRole{}
	err := o.Store.GetContext(ctx, &m, `
	SELECT organization_member.*, member_role.name AS member_role
	FROM organization_member
	LEFT JOIN member_to_role
	ON organization_member.id = member_to_role.member_id
	LEFT JOIN member_role 
	ON member_role.id = member_to_role.role_id 
	WHERE organization_member.id = $1`, ID)

	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (o organizationMember[T]) List(ctx context.Context, organizationId string, limit int, offset int) ([]OrganizationMemberWithRole, error) {
	m := []OrganizationMemberWithRole{}
	if limit == 0 {
		limit = 20
	}

	err := o.Store.SelectContext(ctx, &m, `
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
		WHERE organization.id = $1
		LIMIT $2
		OFFSET $3;`, organizationId, limit, offset)

	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	return m, nil
}
