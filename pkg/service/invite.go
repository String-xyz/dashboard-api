package service

import (
	"context"
	"net/url"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/env"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Invite interface {
	Send(ctx context.Context, request model.RequestInviteSend, callerId *string, organizationId string) (repository.MemberInviteInfo, error)
	Accept(ctx context.Context, inviteId string, requestBody model.RequestInviteAcceptance) (model.OrganizationMember, JWT, error)
	List(ctx context.Context, status string, organizationId string) ([]repository.MemberInviteInfo, error)
	Resend(ctx context.Context, inviteId string, callerId string) (repository.MemberInviteInfo, error)
	Update(ctx context.Context, request model.RequestInviteUpdate, inviteId string, callerId string) (repository.MemberInviteInfo, error)
	Deactivate(ctx context.Context, inviteId string, callerId string) (repository.MemberInviteInfo, error)
	Get(ctx context.Context, id string) (repository.MemberInviteInfo, error)
}

type invite struct {
	repos repository.Repositories
	redis database.RedisStore
}

type TokenPayload struct {
	ExpiresAt int64  `json:"exp"`
	Email     string `json:"email"`
}

func NewInvite(repos repository.Repositories, redis database.RedisStore) Invite {
	return &invite{repos, redis}
}

func (i invite) Send(ctx context.Context, request model.RequestInviteSend, callerId *string, organizationId string) (repository.MemberInviteInfo, error) {
	// Ensure there are no duplicate emails
	preexisting, err := i.repos.OrganizationMember.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return repository.MemberInviteInfo{}, common.StringError(err)
	} else if preexisting.Email == request.Email {
		return repository.MemberInviteInfo{}, common.StringError(serror.ALREADY_IN_USE)
	}

	// If there is a pending invite, update the role
	pendingInvite, err := i.repos.MemberInvite.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return repository.MemberInviteInfo{}, common.StringError(err)
	} else if pendingInvite.Email == request.Email && callerId != nil && pendingInvite.DeactivatedAt == nil {
		// Update invite and resend it
		newRequest := model.RequestInviteUpdate{Role: request.Role, Name: request.Name}
		newInvite, err := i.Update(ctx, newRequest, pendingInvite.Id, *callerId)
		if err != nil {
			return repository.MemberInviteInfo{}, common.StringError(err)
		}
		return i.Resend(ctx, newInvite.Id, *callerId)
	}

	roleId := GetRoleId(request.Role)
	// TODO: VULNERABILITY! Ensure Owner can only be set as role if no other users exist!
	invite, err := i.repos.MemberInvite.Create(ctx, model.MemberInvite{Email: request.Email, InvitedBy: callerId, OrganizationId: organizationId, Name: request.Name, RoleId: roleId})
	if err != nil {
		return invite, common.StringError(err)
	}

	// Create encrypted token so that only the email owner can accept the invite
	key := env.Var.STRING_ENCRYPTION_KEY
	token, err := common.Encrypt(TokenPayload{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Email: request.Email}, key)
	if err != nil {
		return invite, common.StringError(err)
	}

	body := i.createEmailBody(invite.Id, request.Name, token)

	err = SendEmail("String API", "New String API User", request.Email, "String API Invitation", body)
	if err != nil {
		return invite, common.StringError(err)
	}

	return invite, nil
}

func (i invite) Accept(ctx context.Context, inviteId string, requestBody model.RequestInviteAcceptance) (model.OrganizationMember, JWT, error) {
	member := model.OrganizationMember{}
	jwt := JWT{}

	// Ensure that an invite exists
	invite, err := i.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Check invite status
	if repository.GetInviteStatus(invite) != "pending" {
		return member, jwt, common.StringError(serror.ALREADY_IN_USE)
	}

	// Ensure that the token is valid
	key := env.Var.STRING_ENCRYPTION_KEY
	token := requestBody.Token
	payload, err := common.Decrypt[TokenPayload](token, key)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Ensure that the token is not expired
	if payload.ExpiresAt < time.Now().Unix() {
		return member, jwt, common.StringError(serror.FORBIDDEN)
	}

	// Ensure that the token is for the correct email
	if payload.Email != invite.Email {
		return member, jwt, common.StringError(serror.FORBIDDEN)
	}

	// Generate a new Organization Member with an Email
	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 8)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create new member
	member = model.OrganizationMember{Email: invite.Email, Name: invite.Name, Password: string(hash)}
	member, err = i.repos.OrganizationMember.Create(ctx, member)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create Member-To-Organization relationship
	memberToOrganization := model.MemberToOrganization{MemberId: member.Id, OrganizationId: invite.OrganizationId}
	memberToOrganization, err = i.repos.MemberToOrganization.Create(ctx, memberToOrganization)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create Member-To-Role relationship
	_, err = i.repos.MemberToRole.Create(ctx, model.MemberToRole{MemberId: member.Id, RoleId: invite.RoleId})
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Update the invitation
	now := time.Now()
	update := repository.MemberInviteUpdates{AcceptedAt: &now}
	err = i.repos.MemberInvite.Update(ctx, invite.Id, update)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create a JWT
	auth := NewAuth(i.repos, i.redis)
	jwt, err = auth.GenerateJWT(member.Id, invite.OrganizationId)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	return member, jwt, nil
}

func (i invite) List(ctx context.Context, status string, organizationId string) ([]repository.MemberInviteInfo, error) {
	result, err := i.repos.MemberInvite.GetByOrganization(ctx, organizationId)
	if err != nil {
		return result, common.StringError(err)
	}

	// TODO: Filter by status
	// // If a status filter is provided, remove Invites which do not have the filter
	// if status != "" {
	// 	for i, j := range result {
	// 		if !strings.EqualFold(status, repository.GetInviteStatus(j)) { // case insensitive
	// 			result = append(result[:i], result[i+1:]...) // remove from slice
	// 		}
	// 	}
	// }

	return result, nil
}

func (i invite) Resend(ctx context.Context, inviteId string, callerId string) (repository.MemberInviteInfo, error) {
	result := repository.MemberInviteInfo{}
	err := RequireAuthority(i.repos, callerId, "Admin", "Owner")
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = i.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}

	// Create encrypted token so that only the email owner can accept the invite
	key := env.Var.STRING_ENCRYPTION_KEY
	token, err := common.Encrypt(TokenPayload{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Email: result.Email}, key)
	if err != nil {
		return result, common.StringError(err)
	}

	body := i.createEmailBody(result.Id, result.Name, token)

	err = SendEmail("String API", "New String API User", result.Email, "String API Invitation", body)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (i invite) Update(ctx context.Context, request model.RequestInviteUpdate, inviteId string, callerId string) (repository.MemberInviteInfo, error) {
	result := repository.MemberInviteInfo{}
	err := RequireAuthority(i.repos, callerId, "Admin", "Owner")
	if err != nil {
		return result, common.StringError(err)
	}

	if request.Role == "Owner" || request.Role == "owner" {
		return result, common.StringError(serror.FORBIDDEN)
	}

	type RoleUpdate struct {
		RoleId string  `json:"roleId" db:"role_id"`
		Name   *string `json:"name" db:"name"`
	}

	result, err = i.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}

	if request.Name == "" {
		request.Name = result.Name
	}

	update := RoleUpdate{RoleId: GetRoleId(request.Role), Name: &request.Name}
	err = i.repos.MemberInvite.Update(ctx, inviteId, update)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = i.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (i invite) Deactivate(ctx context.Context, inviteId string, callerId string) (repository.MemberInviteInfo, error) {
	result := repository.MemberInviteInfo{}
	err := RequireAuthority(i.repos, callerId, "Admin", "Owner")
	if err != nil {
		return result, common.StringError(err)
	}

	type DeactivateUpdate struct {
		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	}

	now := time.Now()
	update := DeactivateUpdate{DeactivatedAt: &now}
	err = i.repos.MemberInvite.Update(ctx, inviteId, update)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = i.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (i invite) Get(ctx context.Context, id string) (repository.MemberInviteInfo, error) {
	result, err := i.repos.MemberInvite.GetById(ctx, id)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (i invite) createEmailBody(inviteId, userName, token string) string {
	token = url.QueryEscape(token) // make

	href := env.Var.BASE_DASHBOARD_URL + "/invite/" + inviteId + "?token=" + token

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have been invited to use the String API</header>" +
		"<br>Dear " + userName + "," +
		"<br>Thank you for signing up to use the String API.  Please click the link below to set your password and complete your registration process:" +
		"<br><a href='" + href + "'>Accept Invitation</a>"

	return body
}
