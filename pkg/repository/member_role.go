package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serrors "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

type MemberRoleUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          *string    `json:"name" db:"name"`
}

type MemberRole interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberRole) (model.MemberRole, error)
	GetById(ctx context.Context, ID string) (model.MemberRole, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberRole, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByName(ctx context.Context, m model.MemberRole) (model.MemberRole, error)
}

type memberRole[T any] struct {
	strrepo.Base[T]
}

func NewMemberRole(db database.Queryable) MemberRole {
	return &memberRole[model.MemberRole]{strrepo.Base[model.MemberRole]{Store: db, Table: "member_role"}}
}

func (p memberRole[T]) Create(ctx context.Context, m model.MemberRole) (model.MemberRole, error) {
	newModel := model.MemberRole{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_role (id, name) 
		VALUES(:id, :name) RETURNING *`, m)

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

func (p memberRole[T]) GetByName(ctx context.Context, m model.MemberRole) (model.MemberRole, error) {
	result := model.MemberRole{}
	err := p.Store.Get(&result, fmt.Sprintf("SELECT * FROM %s WHERE name = $1 LIMIT 1", p.Table), m.Name)
	if err != nil && err == sql.ErrNoRows {
		return result, common.StringError(serrors.ERR_NOT_FOUND)
	} else if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}
