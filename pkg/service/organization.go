package service

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Organization interface {
	Create(ctx context.Context, request model.RequestOrganizationCreate) (model.Organization, error)
	Get(ctx context.Context, organizationId string) (model.Organization, error)
	Update(ctx context.Context, request model.RequestOrganizationUpdate, organizationId string, callerId string) (model.Organization, error)
}

type organization struct {
	repos repository.Repositories
	redis database.RedisStore
}

func NewOrganization(repos repository.Repositories, redis database.RedisStore) Organization {
	return &organization{repos, redis}
}

// TODO: Ensure valid email is provided
func (a organization) Create(ctx context.Context, request model.RequestOrganizationCreate) (model.Organization, error) {
	// Ensure there are no duplicate emails
	preexisting, err := a.repos.OrganizationMember.GetByEmail(ctx, request.Email)

	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Organization{}, err
	} else if preexisting.Email == request.Email {
		return model.Organization{}, common.StringError(serror.ALREADY_IN_USE)
	}

	pendingInvite, err := a.repos.MemberInvite.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Organization{}, common.StringError(err)
	} else if pendingInvite.Email == request.Email {
		return model.Organization{}, common.StringError(serror.ALREADY_IN_USE)
	}

	// Generate new Organization with a Name
	organization := model.Organization{Name: request.OrganizationName}
	organization, err = a.repos.Organization.Create(ctx, organization)
	if err != nil {
		return organization, common.StringError(err)
	}

	inviteReq := model.RequestInviteSend{Name: request.Name, Email: request.Email, Role: "Owner"}

	// Get String Organization Id
	// Generate Owner invitation
	Invite := NewInvite(a.repos, a.redis)
	_, err = Invite.Send(ctx, inviteReq, nil, organization.Id)
	if err != nil {
		return organization, common.StringError(err)
	}

	return organization, nil
}

func (a organization) Get(ctx context.Context, organizationId string) (model.Organization, error) {
	result, err := a.repos.Organization.GetById(ctx, organizationId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

// TODO: Ensure multiple organizations do not share the same *domain*
func (a organization) Update(ctx context.Context, request model.RequestOrganizationUpdate, organizationId string, callerId string) (model.Organization, error) {
	result := model.Organization{}
	err := RequireAuthority(a.repos, callerId, "Owner")
	if err != nil {
		return result, common.StringError(err)
	}
	err = a.repos.Organization.Update(ctx, organizationId, request)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = a.repos.Organization.GetById(ctx, organizationId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}
