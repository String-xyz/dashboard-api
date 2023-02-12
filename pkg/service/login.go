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
	Login(ctx context.Context, request model.RequestLogin) (repository.PlatformMemberWithRole, JWT, error)
	RefreshToken(refreshToken string) (MemberCreateResponse, error)
	InvalidateRefreshToken(refreshToken string) error
}

type login struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewLogin(repos repository.Repositories, redis database.RedisStore) Login {
	return &login{repos, redis}
}

func (l login) Login(ctx context.Context, request model.RequestLogin) (repository.PlatformMemberWithRole, JWT, error) {
	jwt := JWT{}
	member, err := l.repos.PlatformMember.GetByEmail(ctx, request.Email)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(request.Password))
	if err != nil {
		return member, jwt, common.StringError(errors.New("wrong password"))
	}

	platform, err := l.repos.MemberToPlatform.GetByMember(member.ID)

	if err != nil {
		return member, jwt, common.StringError(err)
	}

	auth := NewAuth(l.repos, l.redis)
	jwt, err = auth.GenerateJWT(member.ID, platform.PlatformID)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	return member, jwt, nil
}

func (l login) RefreshToken(refreshToken string) (MemberCreateResponse, error) {
	auth := NewAuth(l.repos, l.redis)
	return auth.RefreshToken(refreshToken)
}

func (l login) InvalidateRefreshToken(refreshToken string) error {
	auth := NewAuth(l.repos, l.redis)
	return auth.InvalidateRefreshToken(refreshToken)
}
