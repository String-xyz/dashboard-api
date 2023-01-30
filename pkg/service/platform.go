package service

import (
	"context"
	"errors"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
}

func NewPlatform(repos repository.Repositories) Platform {
	return &platform{repos}
}

// TODO: Ensure valid email is provided
// TODO: Can multiple platforms share the same name?  Maybe we should prevent this.
func (a platform) Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error) {
	// Generate new Platform with a Name
	result := model.Platform{Name: request.PlatformName}
	result, err := a.repos.Platform.Create(ctx, result)
	if err != nil {
		return result, common.StringError(err)
	}

	// Generate a new Platform Member with an Email
	member := model.PlatformMember{Email: request.Email}
	member, err = a.repos.PlatformMember.Create(ctx, member)
	if err != nil {
		return result, common.StringError(err)
	}

	ownerId := os.Getenv("MEMBER_ROLE_OWNER_ID")
	if ownerId == "" {
		return result, common.StringError(errors.New("member role id is not defined in the env"))
	}

	_, err = a.repos.MemberToRole.Create(ctx, model.MemberToRole{MemberID: member.ID, RoleID: ownerId})
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}
