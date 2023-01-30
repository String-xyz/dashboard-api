package repository

import (
	"context"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

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
		INSERT INTO member_invite (email, invited_by, platform_id) 
		VALUES(:email, :invited_by, :platform_id) RETURNING *`, m)

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
