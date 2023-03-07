package service

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Login interface {
	Login(ctx context.Context, request model.RequestLogin) (repository.PlatformMemberWithRole, JWT, error)
}

type login struct {
	repos repository.Repositories
	auth  Auth
}

func NewLogin(repos repository.Repositories, auth Auth) Login {
	return &login{repos, auth}
}

func (l login) Login(ctx context.Context, request model.RequestLogin) (repository.PlatformMemberWithRole, JWT, error) {
	jwt := JWT{}
	member, err := l.repos.PlatformMember.GetByEmail(ctx, request.Email)
	if err != nil {
		return member, jwt, common.StringError(err)
	}
	// If member is denied, do not log in
	if member.DeactivatedAt != nil {
		return member, jwt, common.StringError(serror.DEACTIVATED)
	}

	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(request.Password))
	if err != nil {
		return member, jwt, common.StringError(serror.INVALID_PASSWORD)
	}

	platform, err := l.repos.MemberToPlatform.GetByMember(member.ID)

	if err != nil {
		return member, jwt, common.StringError(err)
	}

	jwt, err = l.auth.GenerateJWT(member.ID, platform.PlatformID)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	return member, jwt, nil
}
