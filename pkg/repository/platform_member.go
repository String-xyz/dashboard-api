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

type PlaformMemberUpdates struct {
	DeactivatedAt *time.Time `json:"deactivatedAt" db:"deactivated_at"`
	ActivatedAt   *time.Time `json:"activatedAt" db:"activated_at"`
	Email         *string    `json:"email" db:"email"`
	Password      *string    `json:"password" db:"password"`
}

type PlatformMember interface {
	database.Transactable
	Create(ctx context.Context, model model.PlatformMember) (model.PlatformMember, error)
	GetById(ctx context.Context, ID string) (model.PlatformMember, error)
	List(ctx context.Context, limit int, offset int) ([]model.PlatformMember, error)
	Update(ctx context.Context, ID string, updates any) error
	GetByEmail(email string) (model.PlatformMember, error)
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
		INSERT INTO platform_member (email) 
		VALUES(:email) RETURNING *`, m)

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

func (p platformMember[T]) GetByEmail(email string) (model.PlatformMember, error) {
	m := model.PlatformMember{}
	err := p.Store.Get(&m, fmt.Sprintf("SELECT * FROM %s WHERE email = $1", p.Table), email)
	if err != nil && err == sql.ErrNoRows {
		return m, common.StringError(ErrNotFound)
	} else if err != nil {
		return m, common.StringError(err)
	}
	return m, nil
}
