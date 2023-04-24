package service

import (
	"context"
	"fmt"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create(ctx context.Context, request model.RequestPlatformCreate, organizationId string) (platform model.Platform, err error)
	Get(ctx context.Context, platformId string, organizationId string) (platform model.Platform, err error)
	GetAll(ctx context.Context, callerId string, organizationId string) (platforms []model.Platform, err error)
	Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string, organizationId string) (platform model.Platform, err error)
}

type platform struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewPlatform(repos repository.Repositories, redis database.RedisStore) Platform {
	return &platform{repos, redis}
}

func (p platform) Create(ctx context.Context, request model.RequestPlatformCreate, organizationId string) (platform model.Platform, err error) {
	// Generate new Platform with a Name and Description
	platform = model.Platform{Name: request.PlatformName, Description: request.PlatformDescription, OrganizationId: organizationId}
	platform, err = p.repos.Platform.Create(ctx, platform)
	if err != nil {
		return platform, common.StringError(err)
	}

	return platform, nil
}

func (p platform) Get(ctx context.Context, platformId string, organizationId string) (platform model.Platform, err error) {
	platform, err = p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return platform, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return platform, common.StringError(fmt.Errorf("Platform does not belong to organization"))
	}
	return platform, nil
}

func (p platform) GetAll(ctx context.Context, callerId string, organizationId string) (platforms []model.Platform, err error) {
	platforms, err = p.repos.Platform.List(ctx, organizationId, 0, 0)
	if err != nil {
		return platforms, common.StringError(err)
	}

	return platforms, nil
}

// TODO: Ensure multiple platforms do not share the same *domain*
func (p platform) Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string, organizationId string) (platform model.Platform, err error) {
	err = RequireAuthority(p.repos, callerId, "Admin", "Owner")
	if err != nil {
		return platform, common.StringError(err)
	}
	fmt.Printf("\n\n>>>>> platformId: %s, callerId: %+v, organizationId: %+v\n\n", platformId, callerId, organizationId)
	err = p.repos.Platform.Update(ctx, platformId, request)
	if err != nil {
		return platform, common.StringError(err)
	}
	platform, err = p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return platform, common.StringError(err)
	}
	if platform.OrganizationId != organizationId {
		return platform, common.StringError(fmt.Errorf("Platform does not belong to organization"))
	}
	return platform, nil
}
