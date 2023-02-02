package service

import (
	"context"
	"errors"
	"os"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type MemberCreateResponse struct {
	JWT  JWT                  `json:"authToken"`
	User model.PlatformMember `json:"member"`
}

type Member interface {
	GetAll(ctx context.Context, callerId string, platformId string) ([]model.PlatformMember, error)
	Get(ctx context.Context, callerId string, platformId string, memberId string) (model.PlatformMember, error)
	Update(ctx context.Context, request model.RequestMemberUpdate, callerId string, memberId string) (model.PlatformMember, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	PasswordReset(ctx context.Context, request model.RequestPasswordReset) error
}

type member struct {
	repos repository.Repositories
}

func NewMember(repos repository.Repositories) Member {
	return &member{repos}
}

func (a member) GetAll(ctx context.Context, callerId string, platformId string) ([]model.PlatformMember, error) {
	result := []model.PlatformMember{}
	err := RequireAuthority(a.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	membersToPlatform, err := a.repos.MemberToPlatform.GetByPlatform(platformId)
	if err != nil {
		return result, common.StringError(err)
	}

	for _, mtp := range membersToPlatform {
		member, err := a.repos.PlatformMember.GetById(ctx, mtp.MemberID)
		if err != nil {
			return result, common.StringError(err)
		}
		result = append(result, member)
	}

	return result, nil
}

func (a member) Get(ctx context.Context, callerId string, platformId string, memberId string) (model.PlatformMember, error) {
	result := model.PlatformMember{}
	err := RequireAuthority(a.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	memberToPlatform, err := a.repos.MemberToPlatform.GetByMember(memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if memberToPlatform.PlatformID != platformId {
		return result, common.StringError(errors.New("member not associated with callers platform"))
	}

	result, err = a.repos.PlatformMember.GetById(ctx, memberToPlatform.MemberID)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (a member) Update(ctx context.Context, request model.RequestMemberUpdate, callerId string, memberId string) (model.PlatformMember, error) {
	result := model.PlatformMember{}
	err := RequireAuthority(a.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	// Admin's can't edit Owners or Admins
	callerRole, err := GetRole(a.repos, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	memberRole, err := GetRole(a.repos, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if callerRole == "Admin" && memberRole != "Member" {
		return result, common.StringError(errors.New("admins can only update members"))
	}

	err = a.repos.PlatformMember.Update(ctx, memberId, request)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = a.repos.PlatformMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (a member) SendPasswordResetEmail(ctx context.Context, email string) error {
	// Anyone can request this, so ensure member is associated with 'email'
	member, err := a.repos.PlatformMember.GetByEmail(email)
	if err != nil {
		return common.StringError(err)
	}

	// generate reset token by hashing member id
	secret := os.Getenv("STRING_ENCRYPTION_KEY")
	resetToken, err := common.EncryptString(member.ID, secret)
	if err != nil {
		return common.StringError(err)
	}

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have requested a password reset for the String API</header>" +
		"<br>Dear " + member.Name + "," +
		"<br>If you have forgotten your password, click the link below to reset it:" +
		"<br><a href='https://platform-api.dev.string-api.xyz/members/password-reset/" + resetToken + "'>Reset Password</a>" // TODO: Change URL to Password Reset page and include resetToken

	err = SendEmail("String API", member.Name, "auth@string.xyz", email, "String API Password Reset", body)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (a member) PasswordReset(ctx context.Context, request model.RequestPasswordReset) error {
	// Get the member ID from the password reset token
	secret := os.Getenv("STRING_ENCRYPTION_KEY")
	memberId, err := common.DecryptString(request.ResetToken, secret)
	if err != nil {
		return common.StringError(err)
	}

	// Update their password
	type UpdatePW struct {
		Password string `json:"password" db:"password"`
	}
	updatePW := UpdatePW{Password: request.Password}

	err = a.repos.PlatformMember.Update(ctx, memberId, updatePW)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}
