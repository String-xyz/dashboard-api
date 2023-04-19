package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	strrepo "github.com/String-xyz/go-lib/repository"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
)

// May not be needed, unless platform id changes
type MemberToRoleUpdates struct {
	MemberId *string `json:"memberId" db:"member_id"`
	RoleId   *string `json:"roleId" db:"role_id"`
}

type MemberToRole interface {
	database.Transactable
	Create(ctx context.Context, request model.MemberToRole) (relation model.MemberToRole, err error)
	GetById(ctx context.Context, id string) (relation model.MemberToRole, err error)
	List(ctx context.Context, limit int, offset int) (relations []model.MemberToRole, err error)
	Update(ctx context.Context, id string, updates any) error
	GetByMember(memberId string) (relation model.MemberToRole, err error)
	GetByPlatform(platformId string) (relation model.MemberToRole, err error)
	UpdateRole(memberId string, updates any) error
}

type memberToRole[T any] struct {
	strrepo.Base[T]
}

func NewMemberToRole(db database.Queryable) MemberToRole {
	return &memberToRole[model.MemberToRole]{strrepo.Base[model.MemberToRole]{Store: db, Table: "member_to_role"}}
}

func (m memberToRole[T]) Create(ctx context.Context, request model.MemberToRole) (relation model.MemberToRole, err error) {
	rows, err := m.Store.NamedQuery(`
		INSERT INTO member_to_role (member_id, role_id) 
		VALUES(:member_id, :role_id) RETURNING *`, request)

	if err != nil {
		return relation, common.StringError(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.StructScan(&relation)
		if err != nil {
			return relation, common.StringError(err)
		}
	}

	return relation, nil
}

func (m memberToRole[T]) GetByMember(memberId string) (relation model.MemberToRole, err error) {
	err = m.Store.Get(&relation, fmt.Sprintf("SELECT * FROM %s WHERE member_id = $1", m.Table), memberId)
	if err != nil && err == sql.ErrNoRows {
		return relation, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return relation, common.StringError(err)
	}
	return relation, nil
}

func (m memberToRole[T]) GetByPlatform(platformId string) (relation model.MemberToRole, err error) {
	err = m.Store.Get(&relation, fmt.Sprintf("SELECT * FROM %s WHERE platform_id = $1", m.Table), platformId)
	if err != nil && err == sql.ErrNoRows {
		return relation, common.StringError(serror.NOT_FOUND)
	} else if err != nil {
		return relation, common.StringError(err)
	}
	return relation, nil
}

func (m memberToRole[T]) UpdateRole(memberId string, updates any) error {
	names, keyToUpdate := common.KeysAndValues(updates)
	if len(names) == 0 {
		return common.StringError(errors.New("no fields to update"))
	}
	query := fmt.Sprintf("UPDATE %s SET %s WHERE member_id = '%s'", m.Table, strings.Join(names, ", "), memberId)
	_, err := m.Store.NamedExec(query, keyToUpdate)
	if err != nil {
		return common.StringError(err)
	}
	return nil
}
