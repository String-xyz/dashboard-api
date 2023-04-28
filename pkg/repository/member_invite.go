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
)

func GetInviteStatus(invite MemberInviteInfo) string {
	if invite.AcceptedAt != nil {
		return "accepted"
	} else if invite.ExpiredAt != nil {
		return "expired"
	} else if invite.DeletedAt != nil {
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
	Create(ctx context.Context, request model.MemberInvite) (invite MemberInviteInfo, err error)
	GetById(ctx context.Context, id string) (invite MemberInviteInfo, err error)
	List(ctx context.Context, limit int, offset int) (invites []model.MemberInvite, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByOrganization(ctx context.Context, organizationId string) (invites []MemberInviteInfo, err error)
	GetByEmail(ctx context.Context, email string) (invite MemberInviteInfo, err error)
	SoftDelete(ctx context.Context, id string) error
}

type memberInvite[T any] struct {
	librepository.Base[T]
}

func NewMemberInvite(db database.Queryable) MemberInvite {
	return &memberInvite[model.MemberInvite]{librepository.Base[model.MemberInvite]{Store: db, Table: "member_invite"}}
}

func (i memberInvite[T]) Create(ctx context.Context, request model.MemberInvite) (invite MemberInviteInfo, err error) {
	rows, err := i.Store.NamedQuery(`
		INSERT INTO member_invite (name, email, invited_by, organization_id, role_id)
		VALUES(:name, :email, :invited_by, :organization_id, :role_id)
		RETURNING *, (SELECT name FROM member_role WHERE id = member_invite.role_id) as role, (SELECT name FROM organization WHERE id = member_invite.organization_id) as organization_name
		`, request)

	if err != nil {
		return invite, common.StringError(err)
	}

	// calculate status
	status := "pending"
	invite.Status = &status

	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&invite)
		if err != nil {
			return invite, common.StringError(err)
		}
	}

	return invite, nil
}

func (i memberInvite[T]) GetByOrganization(ctx context.Context, organizationId string) (invites []MemberInviteInfo, err error) {
	err = i.Store.Select(&invites, getBaseQuery()+`WHERE member_invite.organization_id = $1 AND member_invite.deleted_at IS NULL`, organizationId)

	if err != nil && err == sql.ErrNoRows {
		return invites, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return invites, common.StringError(err)
	}

	// calculate status for each invite
	for i := range invites {
		status := GetInviteStatus(invites[i])
		invites[i].Status = &status
	}

	return invites, nil
}

func (i memberInvite[T]) GetById(ctx context.Context, id string) (invite MemberInviteInfo, err error) {
	err = i.Store.GetContext(ctx, &invite, getBaseQuery()+`WHERE member_invite.id = $1 AND member_invite.deleted_at IS NULL`, id)

	if err == sql.ErrNoRows {
		return invite, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return invite, common.StringError(err)
	}

	// calculate status
	status := GetInviteStatus(invite)
	invite.Status = &status

	return invite, nil
}

func (i memberInvite[T]) GetByEmail(ctx context.Context, email string) (invite MemberInviteInfo, err error) {
	err = i.Store.GetContext(ctx, &invite, getBaseQuery()+`WHERE member_invite.email = $1 AND member_invite.deleted_at IS NULL`, email)

	if err != nil && err == sql.ErrNoRows {
		return invite, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return invite, common.StringError(err)
	}

	// calculate status
	status := GetInviteStatus(invite)
	invite.Status = &status

	return invite, nil
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
