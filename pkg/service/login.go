package service

import (
	"context"
	"errors"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Login interface {
	Login(ctx context.Context, request model.RequestLogin) (interface{}, error) // TODO: Return JWT
}

type login struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewLogin(repos repository.Repositories, redis database.RedisStore) Login {
	return &login{repos, redis}
}

func (a login) Login(ctx context.Context, request model.RequestLogin) (interface{}, error) {
	jwt := JWT{}
	member, err := a.repos.PlatformMember.GetByEmail(request.Email)
	if err != nil {
		return jwt, common.StringError(err)
	}
	if request.Password != member.Password {
		return jwt, common.StringError(errors.New("wrong password"))
	}

	platform, err := a.repos.MemberToPlatform.GetByMember(member.ID)
	if err != nil {
		return jwt, common.StringError(err)
	}

	auth := NewAuth(a.repos, a.redis)
	auth.GenerateJWT(member.ID, platform.PlatformID)

	return "{}", nil
}
