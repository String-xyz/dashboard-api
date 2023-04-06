package service

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/go-lib/database"
	serror "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type Invite interface {
	Send(ctx context.Context, request model.RequestInviteSend, callerId *string, platformId string) (repository.MemberInviteInfo, error)
	Accept(ctx context.Context, inviteId string, requestBody model.RequestInviteAcceptance) (model.PlatformMember, JWT, error)
	List(ctx context.Context, status string, platformId string) ([]repository.MemberInviteInfo, error)
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

func (a invite) Send(ctx context.Context, request model.RequestInviteSend, callerId *string, platformId string) (repository.MemberInviteInfo, error) {
	// Ensure there are no duplicate emails
	preexisting, err := a.repos.PlatformMember.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return repository.MemberInviteInfo{}, common.StringError(err)
	} else if preexisting.Email == request.Email {
		return repository.MemberInviteInfo{}, common.StringError(serror.ALREADY_IN_USE)
	}

	// If there is a pending invite, update the role
	pendingInvite, err := a.repos.MemberInvite.GetByEmail(ctx, request.Email)
	if err != nil && !serror.Is(err, serror.NOT_FOUND) {
		return repository.MemberInviteInfo{}, common.StringError(err)
	} else if pendingInvite.Email == request.Email && callerId != nil && pendingInvite.DeactivatedAt == nil {
		// Update invite and resend it
		newRequest := model.RequestInviteUpdate{Role: request.Role, Name: request.Name}
		newInvite, err := a.Update(ctx, newRequest, pendingInvite.ID, *callerId)
		if err != nil {
			return repository.MemberInviteInfo{}, common.StringError(err)
		}
		return a.Resend(ctx, newInvite.ID, *callerId)
	}

	roleId := GetRoleId(request.Role)
	// TODO: VULNERABILITY! Ensure Owner can only be set as role if no other users exist!
	invite, err := a.repos.MemberInvite.Create(ctx, model.MemberInvite{Email: request.Email, InvitedBy: callerId, PlatformId: platformId, Name: request.Name, RoleID: roleId})
	if err != nil {
		return invite, common.StringError(err)
	}

	// Create encrypted token so that only the email owner can accept the invite
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	token, err := common.Encrypt(TokenPayload{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Email: request.Email}, key)
	if err != nil {
		return invite, common.StringError(err)
	}

	body := a.createEmailBody(invite.ID, request.Name, token)

	err = SendEmail("String API", "New String API User", "auth@string.xyz", request.Email, "String API Invitation", body)
	if err != nil {
		return invite, common.StringError(err)
	}

	return invite, nil
}

func (a invite) Accept(ctx context.Context, inviteId string, requestBody model.RequestInviteAcceptance) (model.PlatformMember, JWT, error) {
	member := model.PlatformMember{}
	jwt := JWT{}

	// Ensure that an invite exists
	invite, err := a.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Check invite status
	if repository.GetInviteStatus(invite) != "pending" {
		return member, jwt, common.StringError(serror.ALREADY_IN_USE)
	}

	// Ensure that the token is valid
	key := os.Getenv("STRING_ENCRYPTION_KEY")
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

	// Generate a new Platform Member with an Email
	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 8)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create new member
	member = model.PlatformMember{Email: invite.Email, Name: invite.Name, Password: string(hash)}
	member, err = a.repos.PlatformMember.Create(ctx, member)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create Member-To-Platform relationship
	memberToPlatform := model.MemberToPlatform{MemberID: member.ID, PlatformId: invite.PlatformId}
	memberToPlatform, err = a.repos.MemberToPlatform.Create(ctx, memberToPlatform)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create Member-To-Role relationship
	_, err = a.repos.MemberToRole.Create(ctx, model.MemberToRole{MemberID: member.ID, RoleID: invite.RoleID})
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Update the invitation
	now := time.Now()
	update := repository.MemberInviteUpdates{AcceptedAt: &now}
	err = a.repos.MemberInvite.Update(ctx, invite.ID, update)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	// Create a JWT
	auth := NewAuth(a.repos, a.redis)
	jwt, err = auth.GenerateJWT(member.ID, invite.PlatformId)
	if err != nil {
		return member, jwt, common.StringError(err)
	}

	return member, jwt, nil
}

func (a invite) List(ctx context.Context, status string, platformId string) ([]repository.MemberInviteInfo, error) {
	result, err := a.repos.MemberInvite.GetByPlatform(ctx, platformId)
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
	key := os.Getenv("STRING_ENCRYPTION_KEY")
	token, err := common.Encrypt(TokenPayload{ExpiresAt: time.Now().Add(time.Hour * 24 * 30).Unix(), Email: result.Email}, key)
	if err != nil {
		return result, common.StringError(err)
	}

	body := i.createEmailBody(result.ID, result.Name, token)

	err = SendEmail("String API", "New String API User", "auth@string.xyz", result.Email, "String API Invitation", body)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (a invite) Update(ctx context.Context, request model.RequestInviteUpdate, inviteId string, callerId string) (repository.MemberInviteInfo, error) {
	result := repository.MemberInviteInfo{}
	err := RequireAuthority(a.repos, callerId, "Admin", "Owner")
	if err != nil {
		return result, common.StringError(err)
	}

	if request.Role == "Owner" || request.Role == "owner" {
		return result, common.StringError(serror.FORBIDDEN)
	}

	type RoleUpdate struct {
		RoleID string  `json:"roleId" db:"role_id"`
		Name   *string `json:"name" db:"name"`
	}

	result, err = a.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}

	if request.Name == "" {
		request.Name = result.Name
	}

	update := RoleUpdate{RoleID: GetRoleId(request.Role), Name: &request.Name}
	err = a.repos.MemberInvite.Update(ctx, inviteId, update)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = a.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (a invite) Deactivate(ctx context.Context, inviteId string, callerId string) (repository.MemberInviteInfo, error) {
	result := repository.MemberInviteInfo{}
	err := RequireAuthority(a.repos, callerId, "Admin", "Owner")
	if err != nil {
		return result, common.StringError(err)
	}

	type DeactivateUpdate struct {
		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	}

	now := time.Now()
	update := DeactivateUpdate{DeactivatedAt: &now}
	err = a.repos.MemberInvite.Update(ctx, inviteId, update)
	if err != nil {
		return result, common.StringError(err)
	}
	result, err = a.repos.MemberInvite.GetById(ctx, inviteId)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (a invite) Get(ctx context.Context, id string) (repository.MemberInviteInfo, error) {
	result, err := a.repos.MemberInvite.GetById(ctx, id)
	if err != nil {
		return result, common.StringError(err)
	}
	return result, nil
}

func (i invite) createEmailBody(inviteId, userName, token string) string {
	token = url.QueryEscape(token) // make

	href := os.Getenv("BASE_DASHBOARD_URL") + "/invite/" + inviteId + "?token=" + token

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have been invited to use the String API</header>" +
		"<br>Dear " + userName + "," +
		"<br>Thank you for signing up to use the String API.  Please click the link below to set your password and complete your registration process:" +
		"<br><a href='" + href + "'>Accept Invitation</a>"

	return body
}
