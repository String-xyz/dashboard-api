package repository

import (
	"context"
	"database/sql"
	"fmt"
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

type PlatformMemberWithRole struct {
	model.PlatformMember
	Role string `json:"role" db:"member_role"`
}

type PlatformMember interface {
	database.Transactable
	Create(ctx context.Context, model model.PlatformMember) (model.PlatformMember, error)
	GetById(ctx context.Context, ID string) (PlatformMemberWithRole, error)
	GetByIdIncludingDeactivated(ctx context.Context, ID string) (PlatformMemberWithRole, error)
	List(ctx context.Context, platformId string, limit int, offset int) ([]PlatformMemberWithRole, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByEmail(ctx context.Context, email string) (PlatformMemberWithRole, error)
	Deactivate(ctx context.Context, ID string) error
	Activate(ctx context.Context, ID string) error
}

type platformMember[T any] struct {
	strrepo.Base[T]
}

func NewPlatformMember(db database.Queryable) PlatformMember {
	return &platformMember[model.PlatformMember]{strrepo.Base[model.PlatformMember]{Store: db, Table: "platform_member"}}
}

func (p platformMember[T]) Create(ctx context.Context, m model.PlatformMember) (model.PlatformMember, error) {
	newModel := model.PlatformMember{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO platform_member (email, name, password) 
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

func (p platformMember[T]) GetByEmail(ctx context.Context, email string) (PlatformMemberWithRole, error) {
	m := PlatformMemberWithRole{}
	err := p.Store.GetContext(ctx, &m, `
		SELECT platform_member.*, member_role.name AS member_role
		FROM platform_member
		LEFT JOIN member_to_role
		ON platform_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE platform_member.email = $1`, email)
	if err != nil && err == sql.ErrNoRows {
		return m, serror.NOT_FOUND
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p platformMember[T]) GetById(ctx context.Context, ID string) (PlatformMemberWithRole, error) {
	m := PlatformMemberWithRole{}
	err := p.Store.GetContext(ctx, &m, `
		SELECT platform_member.*, member_role.name AS member_role
		FROM platform_member
		LEFT JOIN member_to_role
		ON platform_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		WHERE platform_member.id = $1`, ID)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p platformMember[T]) GetByIdIncludingDeactivated(ctx context.Context, ID string) (PlatformMemberWithRole, error) {
	m := PlatformMemberWithRole{}
	err := p.Store.GetContext(ctx, &m, fmt.Sprintf("SELECT * FROM %s WHERE id = $1" /* AND deactivated_at IS NULL"*/, p.Table), ID)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p platformMember[T]) List(ctx context.Context, platformId string, limit int, offset int) ([]PlatformMemberWithRole, error) {
	m := []PlatformMemberWithRole{}
	if limit == 0 {
		limit = 20
	}

	err := p.Store.SelectContext(ctx, &m, `
		SELECT platform_member.*, member_role.name AS member_role
		FROM platform_member
		LEFT JOIN member_to_role
		ON platform_member.id = member_to_role.member_id
		LEFT JOIN member_role 
		ON member_role.id = member_to_role.role_id 
		LEFT JOIN member_to_platform
		ON platform_member.id = member_to_platform.member_id 
		LEFT JOIN platform
		ON platform.id = member_to_platform.platform_id
		WHERE platform.id = $1
		LIMIT $2
		OFFSET $3;`, platformId, limit, offset)

	if err == sql.ErrNoRows {
		return m, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return m, common.StringError(err)
	}

	return m, nil
}
