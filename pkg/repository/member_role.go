package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	librepository "github.com/String-xyz/go-lib/v2/repository"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
)

type MemberRoleUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          *string    `json:"name" db:"name"`
}

type MemberRole interface {
	database.Transactable
	Create(ctx context.Context, request model.MemberRole) (role model.MemberRole, err error)
	GetById(ctx context.Context, id string) (role model.MemberRole, err error)
	List(ctx context.Context, limit int, offset int) (roles []model.MemberRole, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByName(ctx context.Context, name string) (role model.MemberRole, err error)
}

type memberRole[T any] struct {
	librepository.Base[T]
}

func NewMemberRole(db database.Queryable) MemberRole {
	return &memberRole[model.MemberRole]{librepository.Base[model.MemberRole]{Store: db, Table: "member_role"}}
}

func (r memberRole[T]) Create(ctx context.Context, request model.MemberRole) (role model.MemberRole, err error) {
	rows, err := r.Store.NamedQuery(`
		INSERT INTO member_role (id, name) 
		VALUES(:id, :name) RETURNING *`, request)

	if err != nil {
		return role, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&role)
		if err != nil {
			return role, common.StringError(err)
		}
	}

	return role, nil
}

func (r memberRole[T]) GetByName(ctx context.Context, name string) (role model.MemberRole, err error) {
	err = r.Store.Get(&role, fmt.Sprintf("SELECT * FROM %s WHERE name = $1 LIMIT 1 AND member_invite.deleted_at IS NULL", r.Table), name)
	if err != nil && err == sql.ErrNoRows {
		return role, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return role, common.StringError(err)
	}
	return role, nil
}
