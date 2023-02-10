package service

import (
	"context"
	"errors"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Login interface {
	Login(ctx context.Context, request model.RequestLogin) (JWT, error)
}

type login struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewLogin(repos repository.Repositories, redis database.RedisStore) Login {
	return &login{repos, redis}
}

func (a login) Login(ctx context.Context, request model.RequestLogin) (JWT, error) {
	jwt := JWT{}
	member, err := a.repos.PlatformMember.GetByEmail(ctx, request.Email)
	if err != nil {
		return jwt, common.StringError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(request.Password))
	if err != nil {
		return jwt, common.StringError(errors.New("wrong password"))
	}

	platform, err := a.repos.MemberToPlatform.GetByMember(member.ID)

	if err != nil {
		return jwt, common.StringError(err)
	}

	auth := NewAuth(a.repos, a.redis)
	jwt, err = auth.GenerateJWT(member.ID, platform.PlatformID)
	if err != nil {
		return jwt, common.StringError(err)
	}

	return jwt, nil
}
