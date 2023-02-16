package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

func GetInviteStatus(invite model.MemberInvite) string {
	if invite.AcceptedAt != nil {
		return "accepted"
	} else if invite.ExpiredAt != nil {
		return "expired"
	} else if invite.DeactivatedAt != nil {
		return "revoked"
	} else if invite.ID != "" {
		return "pending"
	}
	return "invalid"
}

type MemberInviteInfo struct {
	ID            string     `json:"id" db:"id"`
	Name          string     `json:"name" db:"name"`                  // user name
	Email         string     `json:"email" db:"email"`                // user name
	Role          string     `json:"role" db:"role"`                  // invite role
	PlatformName  string     `json:"platformName" db:"platform_name"` // platform name
	Status        *string    `json:"status"`                          // invite status
	AcceptedAt    *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ExpiredAt     *time.Time `json:"expiredAt,omitempty" db:"expired_at"`
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
	Create(ctx context.Context, model model.MemberInvite) (model.MemberInvite, error)
	GetById(ctx context.Context, ID string) (model.MemberInvite, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberInvite, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByPlatform(platformId string) ([]MemberInviteInfo, error)
	GetMemberAndPlatformName(ctx context.Context, ID string) (MemberInviteInfo, error)
	GetByEmail(email string) (model.MemberInvite, error)
}

type memberInvite[T any] struct {
	strrepo.Base[T]
}

func NewMemberInvite(db database.Queryable) MemberInvite {
	return &memberInvite[model.MemberInvite]{strrepo.Base[model.MemberInvite]{Store: db, Table: "member_invite"}}
}

func (p memberInvite[T]) Create(ctx context.Context, m model.MemberInvite) (model.MemberInvite, error) {
	newModel := model.MemberInvite{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_invite (name, email, invited_by, platform_id, role_id) 
		VALUES(:name, :email, :invited_by, :platform_id, :role_id) RETURNING *`, m)

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

func (p memberInvite[T]) GetByPlatform(platformId string) ([]MemberInviteInfo, error) {
	m := []MemberInviteInfo{}

	err := p.Store.Select(&m, `
	SELECT member_invite.id, member_invite.name, member_invite.email, member_invite.accepted_at, member_invite.deactivated_at, member_invite.expired_at, member_role.name as role, platform.name as platform_name
	FROM member_invite
	LEFT JOIN platform
	ON member_invite.platform_id = platform.id
	LEFT JOIN member_role
	ON member_invite.role_id = member_role.id
	WHERE member_invite.platform_id = $1`, platformId)

	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p memberInvite[T]) GetMemberAndPlatformName(ctx context.Context, ID string) (MemberInviteInfo, error) {
	m := MemberInviteInfo{}
	err := p.Store.GetContext(ctx, &m, `
	SELECT member_invite.id, member_invite.name, member_invite.email, member_role.name as role, platform.name as platform_name
	FROM member_invite
	LEFT JOIN platform
	ON member_invite.platform_id = platform.id
	LEFT JOIN member_role
	ON member_invite.role_id = member_role.id
	WHERE member_invite.id = $1`, ID)
	if err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}

	return m, nil
}

func (p memberInvite[T]) GetById(ctx context.Context, ID string) (model.MemberInvite, error) {
	m := model.MemberInvite{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1" /*AND deactivated_at IS NULL"*/, p.Table), ID)
	if err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p memberInvite[T]) GetByEmail(email string) (model.MemberInvite, error) {
	m := model.MemberInvite{}
	err := p.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE email = $1", p.Table), email)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
