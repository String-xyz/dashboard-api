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

type MemberRoleUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          *string    `json:"name" db:"name"`
}

type MemberRole interface {
	database.Transactable
	Create(ctx context.Context, request model.MemberRole) (result model.MemberRole, err error)
	GetById(ctx context.Context, id string) (result model.MemberRole, err error)
	List(ctx context.Context, limit int, offset int) (results []model.MemberRole, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByName(ctx context.Context, name string) (result model.MemberRole, err error)
}

type memberRole[T any] struct {
	strrepo.Base[T]
}

func NewMemberRole(db database.Queryable) MemberRole {
	return &memberRole[model.MemberRole]{strrepo.Base[model.MemberRole]{Store: db, Table: "member_role"}}
}

func (m memberRole[T]) Create(ctx context.Context, request model.MemberRole) (result model.MemberRole, err error) {
	rows, err := m.Store.NamedQuery(`
		INSERT INTO member_role (id, name) 
		VALUES(:id, :name) RETURNING *`, request)

	if err != nil {
		return result, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return result, common.StringError(err)
		}
	}

	return result, nil
}

func (m memberRole[T]) GetByName(ctx context.Context, name string) (result model.MemberRole, err error) {
	err = m.Store.Get(&result, fmt.Sprintf("SELECT * FROM %s WHERE name = $1 LIMIT 1", m.Table), name)
	if err != nil && err == sql.ErrNoRows {
		return result, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}
