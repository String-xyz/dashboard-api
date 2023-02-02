package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

// May not be needed, unless platform ID changes
type MemberToRoleUpdates struct {
	MemberID *string `json:"memberId" db:"member_id"`
	RoleID   *string `json:"roleId" db:"role_id"`
}

type MemberToRole interface {
	database.Transactable
	Create(ctx context.Context, model model.MemberToRole) (model.MemberToRole, error)
	GetById(ctx context.Context, ID string) (model.MemberToRole, error)
	List(ctx context.Context, limit int, offset int) ([]model.MemberToRole, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByMember(memberId string) (model.MemberToRole, error)
	GetByPlatform(platformId string) (model.MemberToRole, error)
}

type memberToRole[T any] struct {
	strrepo.Base[T]
}

func NewMemberToRole(db database.Queryable) MemberToRole {
	return &memberToRole[model.MemberToRole]{strrepo.Base[model.MemberToRole]{Store: db, Table: "member_to_role"}}
}

func (p memberToRole[T]) Create(ctx context.Context, m model.MemberToRole) (model.MemberToRole, error) {
	newModel := model.MemberToRole{}
	rows, err := p.Store.NamedQuery(`
		INSERT INTO member_to_role (member_id, role_id) 
		VALUES(:member_id, :role_id) RETURNING *`, m)

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

func (p memberToRole[T]) GetByMember(memberId string) (model.MemberToRole, error) {
	m := model.MemberToRole{}
	err := p.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE member_id = $1", p.Table), memberId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}

func (p memberToRole[T]) GetByPlatform(platformId string) (model.MemberToRole, error) {
	m := model.MemberToRole{}
	err := p.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE platform_id = $1", p.Table), platformId)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
