package service

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error)
	Get(ctx context.Context, request string) (model.Platform, error)
	Update(ctx context.Context, request model.RequestPlatformUpdate, id string) (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
}

func NewPlatform(repos repository.Repositories) Platform {
	return &platform{repos}
}

// TODO: Ensure valid email is provided
func (a platform) Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error) {
	// Generate new Platform with a Name
	platform := model.Platform{Name: request.PlatformName}
	platform, err := a.repos.Platform.Create(ctx, platform)
	if err != nil {
		return platform, common.StringError(err)
	}

	inviteReq := model.RequestInviteSend{Name: request.Name, Email: request.Email, Role: "Owner"}

	// Get String Platform Id
	// Generate Owner invitation
	Invite := NewInvite(a.repos)
	_, err = Invite.Send(ctx, inviteReq, platform)
	if err != nil {
		return platform, common.StringError(err)
	}

	return platform, nil
}

func (a platform) Get(ctx context.Context, request string) (model.Platform, error) {
	result := model.Platform{}
	return result, nil
}

// TODO: Ensure multiple platforms do not share the same *domain*
func (a platform) Update(ctx context.Context, request model.RequestPlatformUpdate, id string) (model.Platform, error) {
	result := model.Platform{}
	return result, nil
}
