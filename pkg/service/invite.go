package service

import (
	"context"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Invite interface {
	Send(ctx context.Context, request model.RequestInviteSend, platform model.Platform) (model.MemberInvite, error)
	Accept(ctx context.Context, requestBody model.RequestInviteAcceptance) (model.PlatformMember, error)
	List(ctx context.Context, request string) ([]model.MemberInvite, error)
	Resend(ctx context.Context, request string) (model.MemberInvite, error)
	Update(ctx context.Context, request model.RequestInviteUpdate, id string) (model.MemberInvite, error)
}

type invite struct {
	repos repository.Repositories
}

func NewInvite(repos repository.Repositories) Invite {
	return &invite{repos}
}

func (a invite) Send(ctx context.Context, request model.RequestInviteSend, platform model.Platform) (model.MemberInvite, error) {
	var roleId string
	if request.Role == "Member" {
		roleId = os.Getenv("MEMBER_ROLE_MEMBER_ID")
	} else if request.Role == "Admin" {
		roleId = os.Getenv("MEMBER_ROLE_ADMIN_ID")
	} else if request.Role == "Owner" {
		roleId = os.Getenv("MEMBER_ROLE_OWNER_ID")
	}
        // TODO: Use Role Helper Function
        // TODO: VULNERABILITY! Ensure Owner can only be set as role if no other users exist!
	invite, err := a.repos.MemberInvite.Create(ctx, model.MemberInvite{Email: request.Email, InvitedBy: request.Invitee, PlatformID: platform.ID, Name: request.Name, RoleID: roleId})
	if err != nil {
		return invite, common.StringError(err)
	}

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have been invited to use the String API</header>" +
		"<br>Dear " + request.Name + "," +
		"<br>Thank you for signing up to use the String API.  Please click the link below to set your password and complete your registration process:" +
		"<br><a href='" + os.Getenv("BASE_APP_URL") + "invites/" + invite.ID + "'>Accept Invitation</a>" // TODO: double check :id

	err = SendEmail("String API", "New String API User", "auth@string.xyz", request.Email, "String API Invitation", body)
	if err != nil {
		return invite, common.StringError(err)
	}

	return invite, nil
}

func (a invite) Accept(ctx context.Context, requestBody model.RequestInviteAcceptance) (model.PlatformMember, error) {

	invite, err := a.repos.MemberInvite.GetById(ctx, *requestBody.Id)
	// Generate a new Platform Member with an Email
	member := model.PlatformMember{Email: invite.Email, Name: invite.Name, Password: requestBody.Password}
	member, err = a.repos.PlatformMember.Create(ctx, member)
	if err != nil {
		return member, common.StringError(err)
	}

	// Create Member-To-Platform relationship
	memberToPlatform := model.MemberToPlatform{MemberID: member.ID, PlatformID: invite.PlatformID}
	memberToPlatform, err = a.repos.MemberToPlatform.Create(ctx, memberToPlatform)
	if err != nil {
		return member, common.StringError(err)
	}

	// Create Member-To-Role relationship
	_, err = a.repos.MemberToRole.Create(ctx, model.MemberToRole{MemberID: member.ID, RoleID: invite.RoleID})
	if err != nil {
		return member, common.StringError(err)
	}

	return member, nil
}

func (a invite) List(ctx context.Context, request string) ([]model.MemberInvite, error) {
	result := []model.MemberInvite{}
	return result, nil
}

func (a invite) Resend(ctx context.Context, request string) (model.MemberInvite, error) {
	result := model.MemberInvite{}
	return result, nil
}

func (a invite) Update(ctx context.Context, request model.RequestInviteUpdate, id string) (model.MemberInvite, error) {
	result := model.MemberInvite{}
	return result, nil
}
