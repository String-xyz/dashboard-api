package service

import (
	"context"
	"fmt"

	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create(ctx context.Context, request model.RequestPlatformCreate, organizationId string) (platform model.Platform, err error)
	Get(ctx context.Context, platformId string, organizationId string) (platform model.Platform, err error)
	GetAll(ctx context.Context, callerId string, organizationId string) (platforms []model.Platform, err error)
	Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string, organizationId string) (platform model.Platform, err error)
	Deactivate(ctx context.Context, platformId string, callerId string, organizationId string) (model.Platform, error)
	Reactivate(ctx context.Context, platformId string, callerId string, organizationId string) (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewPlatform(repos repository.Repositories, redis database.RedisStore) Platform {
	return &platform{repos, redis}
}

func (p platform) Create(ctx context.Context, request model.RequestPlatformCreate, organizationId string) (platform model.Platform, err error) {
	_, finish := Span(ctx, "service.platform.Create", SpanTag{"organizationId": organizationId})
	defer finish()

	// Generate new Platform with a Name and Description
	platform = model.Platform{Name: request.PlatformName, Description: request.PlatformDescription, OrganizationId: organizationId}
	platform, err = p.repos.Platform.Create(ctx, platform)
	if err != nil {
		return platform, common.StringError(err)
	}

	return platform, nil
}

func (p platform) Get(ctx context.Context, platformId string, organizationId string) (platform model.Platform, err error) {
	_, finish := Span(ctx, "service.platform.Get", SpanTag{"organizationId": organizationId})
	defer finish()

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
	_, finish := Span(ctx, "service.platform.GetAll", SpanTag{"organizationId": organizationId})
	defer finish()

	platforms, err = p.repos.Platform.List(ctx, organizationId, 0, 0)
	if err != nil {
		return platforms, common.StringError(err)
	}

	return platforms, nil
}

// TODO: Ensure multiple platforms do not share the same *domain*
func (p platform) Update(ctx context.Context, request model.RequestPlatformUpdate, platformId string, callerId string, organizationId string) (platform model.Platform, err error) {
	_, finish := Span(ctx, "service.platform.Update", SpanTag{"organizationId": organizationId})
	defer finish()

	err = RequireAuthority(p.repos, callerId, "Admin", "Owner")
	if err != nil {
		return platform, common.StringError(err)
	}
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

func (p platform) Deactivate(ctx context.Context, platformId string, callerId string, organizationId string) (model.Platform, error) {
	_, finish := Span(ctx, "service.platform.Deactivate")
	defer finish()

	result := model.Platform{}

	err := RequireAuthority(p.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	platform, err := p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	if platform.OrganizationId != organizationId {
		return result, common.StringError(serror.FORBIDDEN)
	}

	p.repos.Platform.Deactivate(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	// return deactivated platform
	result, err = p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (p platform) Reactivate(ctx context.Context, platformId string, callerId string, organizationId string) (model.Platform, error) {
	_, finish := Span(ctx, "service.member.Reactivate")
	defer finish()

	result := model.Platform{}

	err := RequireAuthority(p.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	platform, err := p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	if platform.OrganizationId != organizationId {
		return result, common.StringError(serror.FORBIDDEN)
	}

	p.repos.Platform.Activate(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = p.repos.Platform.GetById(ctx, platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}
