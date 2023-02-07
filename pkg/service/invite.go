package service

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Invite interface {
	Send(ctx context.Context, request model.RequestInviteSend, callerId *string, platformId string) (model.MemberInvite, error)
	Accept(ctx context.Context, requestBody model.RequestInviteAcceptance) (model.PlatformMember, error)
	List(ctx context.Context, status string, platformId string) ([]model.MemberInvite, error)
	Resend(ctx context.Context, request string) (model.MemberInvite, error)
	Update(ctx context.Context, request model.RequestInviteUpdate, id string) (model.MemberInvite, error)
}

type invite struct {
	repos repository.Repositories
}

func NewInvite(repos repository.Repositories) Invite {
	return &invite{repos}
}

func (a invite) Send(ctx context.Context, request model.RequestInviteSend, callerId *string, platformId string) (model.MemberInvite, error) {
	roleId := GetRoleId(request.Role)
	// TODO: VULNERABILITY! Ensure Owner can only be set as role if no other users exist!
	invite, err := a.repos.MemberInvite.Create(ctx, model.MemberInvite{Email: request.Email, InvitedBy: callerId, PlatformID: platformId, Name: request.Name, RoleID: roleId})
	if err != nil {
		return invite, common.StringError(err)
	}

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have been invited to use the String API</header>" +
		"<br>Dear " + request.Name + "," +
		"<br>Thank you for signing up to use the String API.  Please click the link below to set your password and complete your registration process:" +
		"<br><a href='" + os.Getenv("BASE_APP_URL") + "/invites/" + invite.ID + "'>Accept Invitation</a>" // TODO: double check :id

	err = SendEmail("String API", "New String API User", "auth@string.xyz", request.Email, "String API Invitation", body)
	if err != nil {
		return invite, common.StringError(err)
	}

	return invite, nil
}

func (a invite) Accept(ctx context.Context, requestBody model.RequestInviteAcceptance) (model.PlatformMember, error) {
	member := model.PlatformMember{}

	// Ensure that an invite exists
	invite, err := a.repos.MemberInvite.GetById(ctx, *requestBody.Id)
	if err != nil {
		return member, common.StringError(err)
	}

	// Check invite status
	if repository.GetInviteStatus(invite) != "Pending" {
		return member, common.StringError(errors.New("invite is not pending"))
	}

	// Generate a new Platform Member with an Email
	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 8)
	if err != nil {
		return member, common.StringError(err)
	}

	// Create new member
	member = model.PlatformMember{Email: invite.Email, Name: invite.Name, Password: string(hash)}
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

	// Update the invitation
	now := time.Now()
	update := repository.MemberInviteUpdates{AcceptedAt: &now}
	err = a.repos.MemberInvite.Update(ctx, invite.ID, update)
	if err != nil {
		return member, common.StringError(err)
	}

	return member, nil
}

func (a invite) List(ctx context.Context, status string, platformId string) ([]model.MemberInvite, error) {
	result, err := a.repos.MemberInvite.GetByPlatform(platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	// If a status filter is provided, remove Invites which do not have the filter
	if status != "" {
		for i, j := range result {
			if !strings.EqualFold(status, repository.GetInviteStatus(j)) { // case insensitive
				result = append(result[:i], result[i+1:]...) // remove from slice
			}
		}
	}

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
