package service

import (
	"context"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/repository"
	"github.com/String-xyz/go-lib/v2/common"
	"github.com/String-xyz/go-lib/v2/database"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
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
func (o organization) Create(ctx context.Context, request model.RequestOrganizationCreate) (model.Organization, error) {
	_, finish := Span(ctx, "service.organization.Create")
	defer finish()

	// Ensure there are no duplicate emails
	preexisting, err := o.repos.OrganizationMember.GetByEmail(ctx, request.Email)

	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Organization{}, err
	} else if preexisting.Email == request.Email {
		return model.Organization{}, common.StringError(serror.ALREADY_IN_USE)
	}

	pendingInvite, err := o.repos.MemberInvite.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return model.Organization{}, common.StringError(err)
	} else if pendingInvite.Email == request.Email {
		return model.Organization{}, common.StringError(serror.ALREADY_IN_USE)
	}

	// Generate new Organization with a Name
	organization := model.Organization{Name: request.OrganizationName}
	organization, err = o.repos.Organization.Create(ctx, organization)
	if err != nil {
		return organization, common.StringError(err)
	}

	inviteReq := model.RequestInviteSend{Name: request.Name, Email: request.Email, Role: "Owner"}

	// Get String Organization Id
	// Generate Owner invitation
	Invite := NewInvite(o.repos, o.redis)
	_, err = Invite.Send(ctx, inviteReq, nil, organization.Id)
	if err != nil {
		return organization, common.StringError(err)
	}

	return organization, nil
}

func (o organization) Get(ctx context.Context, organizationId string) (model.Organization, error) {
	_, finish := Span(ctx, "service.organization.Get", SpanTag{"organizationId": organizationId})
	defer finish()

	result, err := o.repos.Organization.GetById(ctx, organizationId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

// TODO: Ensure multiple organizations do not share the same *domain*
func (o organization) Update(ctx context.Context, request model.RequestOrganizationUpdate, organizationId string, callerId string) (model.Organization, error) {
	_, finish := Span(ctx, "service.organization.Update", SpanTag{"organizationId": organizationId})
	defer finish()

	result := model.Organization{}
	err := RequireAuthority(o.repos, callerId, "Owner")
	if err != nil {
		return result, common.StringError(err)
	}
	err = o.repos.Organization.Update(ctx, organizationId, request)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = o.repos.Organization.GetById(ctx, organizationId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}
