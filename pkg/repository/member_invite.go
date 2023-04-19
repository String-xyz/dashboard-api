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

func GetInviteStatus(invite MemberInviteInfo) string {
	if invite.AcceptedAt != nil {
		return "accepted"
	} else if invite.ExpiredAt != nil {
		return "expired"
	} else if invite.DeactivatedAt != nil {
		return "revoked"
	} else if invite.Id != "" {
		return "pending"
	}

	return "invalid"
}

type MemberInviteInfo struct {
	model.MemberInvite

	Role             string  `json:"role" db:"role"`                          // invite role
	OrganizationName *string `json:"organizationName" db:"organization_name"` // organization name
	Status           *string `json:"status"`                                  // invite status
}

type MemberInviteUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	ExpiredAt     *time.Time `json:"expiredAt" db:"expired_at"`
	AcceptedAt    *time.Time `json:"acceptedAt" db:"accepted_at"`
	Email         *string    `json:"email" db:"email"`
	Password      *string    `json:"password" db:"password"`
}

type MemberInvite interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberInvite) (MemberInviteInfo, error)
	GetById(ctx context.Context, ID string) (MemberInviteInfo, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberInvite, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByOrganization(ctx context.Context, organizationId string) ([]MemberInviteInfo, error)
	GetByEmail(ctx context.Context, email string) (MemberInviteInfo, error)
}

type memberInvite[T any] struct {
	strrepo.Base[T]
}

func NewMemberInvite(db database.Queryable) MemberInvite {
	return &memberInvite[model.MemberInvite]{strrepo.Base[model.MemberInvite]{Store: db, Table: "member_invite"}}
}

func (p memberInvite[T]) Create(ctx context.Context, m model.MemberInvite) (MemberInviteInfo, error) {
	newModel := MemberInviteInfo{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_invite (name, email, invited_by, organization_id, role_id)
		VALUES(:name, :email, :invited_by, :organization_id, :role_id)
		RETURNING *, (SELECT name FROM member_role WHERE id = member_invite.role_id) as role, (SELECT name FROM organization WHERE id = member_invite.organization_id) as organization_name
		`, m)

	if err != nil {
		return newModel, common.StringError(err)
	}

	// calculate status
	status := "pending"
	newModel.Status = &status

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&newModel)
		if err != nil {
			return newModel, common.StringError(err)
		}
	}

	return newModel, nil
}

func (p memberInvite[T]) GetByOrganization(ctx context.Context, organizationId string) ([]MemberInviteInfo, error) {
	m := []MemberInviteInfo{}

	err := p.Store.Select(&m, getBaseQuery()+`WHERE member_invite.organization_id = $1`, organizationId)

	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	// calculate status for each invite
	for i := range m {
		status := GetInviteStatus(m[i])
		m[i].Status = &status
	}

	return m, nil
}

func (p memberInvite[T]) GetById(ctx context.Context, ID string) (MemberInviteInfo, error) {
	m := MemberInviteInfo{}

	err := p.Store.GetContext(ctx, &m, getBaseQuery()+`WHERE member_invite.id = $1`, ID)

	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	// calculate status
	status := GetInviteStatus(m)
	m.Status = &status

	return m, nil
}

func (p memberInvite[T]) GetByEmail(ctx context.Context, email string) (MemberInviteInfo, error) {
	m := MemberInviteInfo{}

	err := p.Store.GetContext(ctx, &m, getBaseQuery()+`WHERE member_invite.email = $1`, email)

	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	// calculate status
	status := GetInviteStatus(m)
	m.Status = &status

	return m, nil
}

func getBaseQuery() string {
	return `
		SELECT member_invite.*, member_role.name as role, organization.name as organization_name
		FROM member_invite
		LEFT JOIN organization
		ON member_invite.organization_id = organization.id
		LEFT JOIN member_role
		ON member_invite.role_id = member_role.id
		`
}
