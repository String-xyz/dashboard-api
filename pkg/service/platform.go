package service

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error)
	Get(ctx context.Context, platformId string) (model.Platform, error)
	Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string) (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewPlatform(repos repository.Repositories, redis database.RedisStore) Platform {
	return &platform{repos, redis}
}

// TODO: Ensure valid email is provided
func (p platform) Create(ctx context.Context, request model.RequestPlatformCreate) (model.Platform, error) {
	// Ensure there are no duplicate emails
	preexisting, err := p.repos.OrganizationMember.GetByEmail(ctx, request.Email)

	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Platform{}, err
	} else if preexisting.Email == request.Email {
		return model.Platform{}, common.StringError(serror.ALREADY_IN_USE)
	}

	pendingInvite, err := p.repos.MemberInvite.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Platform{}, common.StringError(err)
	} else if pendingInvite.Email == request.Email {
		return model.Platform{}, common.StringError(serror.ALREADY_IN_USE)
	}

	// Generate new Platform with a Name
	platform := model.Platform{Name: request.PlatformName}
	platform, err = p.repos.Platform.Create(ctx, platform)
	if err != nil {
		return platform, common.StringError(err)
	}

	inviteReq := model.RequestInviteSend{Name: request.Name, Email: request.Email, Role: "Owner"}

	// Get String Platform Id
	// Generate Owner invitation
	Invite := NewInvite(p.repos, p.redis)
	_, err = Invite.Send(ctx, inviteReq, nil, platform.Id)
	if err != nil {
		return platform, common.StringError(err)
	}

	return platform, nil
}

func (p platform) Get(ctx context.Context, platformId string) (model.Platform, error) {
	result, err := p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

// TODO: Ensure multiple platforms do not share the same *domain*
func (p platform) Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string) (model.Platform, error) {
	result := model.Platform{}
	err := RequireAuthority(p.repos, callerId, "Owner")
	if err != nil {
		return result, common.StringError(err)
	}
	err = p.repos.Platform.Update(ctx, platformId, request)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}
