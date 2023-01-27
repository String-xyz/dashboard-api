package repository

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

// May not be needed, unless platform ID changes
type MemberToPlatformUpdates struct {
	PlatformID *string `json:"platformId" db:"platform_id"`
	MemberID   *string `json:"memberId" db:"member_id"`
}

type MemberToPlatform interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberToPlatform) (model.MemberToPlatform, error)
	GetById(ctx context.Context, ID string) (model.MemberToPlatform, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberToPlatform, error)
	Update(ctx context.Context, ID string, updates any) error
}

type memberToPlatform[T any] struct {
	strrepo.Base[T]
}

func NewMemberToPlatform(db database.Queryable) MemberToPlatform {
	return &memberToPlatform[model.MemberToPlatform]{strrepo.Base[model.MemberToPlatform]{Store: db, Table: "member_to_platform"}}
}

func (p memberToPlatform[T]) Create(ctx context.Context, m model.MemberToPlatform) (model.MemberToPlatform, error) {
	newModel := model.MemberToPlatform{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_to_platform (platform_id, member_id) 
		VALUES(:platform_id, :member_id) RETURNING *`, m)

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
