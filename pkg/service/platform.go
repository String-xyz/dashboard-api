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

	// TODO: replace this with an invite
	// Generate a new Platform Member with an Email
	member := model.PlatformMember{Email: request.Email, Name: request.Name}
	member, err = a.repos.PlatformMember.Create(ctx, member)
	if err != nil {
		return platform, common.StringError(err)
	}

	// TODO: move this to invite acceptance logic
	// Create Member-To-Platform relationship
	memberToPlatform := model.MemberToPlatform{MemberID: member.ID, PlatformID: platform.ID}
	memberToPlatform, err = a.repos.MemberToPlatform.Create(ctx, memberToPlatform)
	if err != nil {
		return platform, common.StringError(err)
	}

	// TODO: move this to invite acceptance logic
	// Give new Platform Member ownership of the Platform they just created
	ownerId := os.Getenv("MEMBER_ROLE_OWNER_ID")
	if ownerId == "" {
		return platform, common.StringError(errors.New("member role id is not defined in the env"))
	}

	_, err = a.repos.MemberToRole.Create(ctx, model.MemberToRole{MemberID: member.ID, RoleID: ownerId})
	if err != nil {
		return platform, common.StringError(err)
	}

	// Generate Owner invitation
	invite, err := a.repos.MemberInvite.Create(ctx, model.MemberInvite{Email: request.Email, InvitedBy: member.ID, PlatformID: platform.ID, Name: request.Name})
	if err != nil {
		return platform, common.StringError(err)
	}

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have been invited to use the String API</header>" +
		"<br>Dear " + request.PlatformName + " owner," +
		"<br>Thank you for signing up to use the String API.  Please click the link below to set your password and complete your registration process:" +
		"<br><a href='https://platform-api.dev.string-api.xyz/invites/" + invite.ID + "'>Accept Invitation</a>" // TODO: double check :id
	err = SendEmail("String API", "New String API User", "auth@string.xyz", request.Email, "String API Invitation", body)
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
