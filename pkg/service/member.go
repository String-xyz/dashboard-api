package service

import (
	"context"
	"errors"
	"net/url"
	"os"
	"time"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type MemberCreateResponse struct {
	JWT  JWT                               `json:"authToken"`
	User repository.PlatformMemberWithRole `json:"member"`
}

type Member interface {
	GetAll(ctx context.Context, callerId string, platformId string) ([]repository.PlatformMemberWithRole, error)
	Get(ctx context.Context, callerId string, platformId string, memberId string) (repository.PlatformMemberWithRole, error)
	UpdateMember(ctx context.Context, request model.RequestMemberUpdateOther, callerId string, memberId string) (model.MemberToRole, error)
	UpdateSelf(ctx context.Context, request model.RequestMemberUpdateSelf, callerId string) (repository.PlatformMemberWithRole, error)
	SendPasswordResetEmail(ctx context.Context, email string) error
	PasswordReset(ctx context.Context, request model.RequestPasswordReset) error
	Deactivate(ctx context.Context, callerId string, memberId string) (repository.PlatformMemberWithRole, error)
}

type member struct {
	repos repository.Repositories
}

func NewMember(repos repository.Repositories) Member {
	return &member{repos}
}

func (a member) GetAll(ctx context.Context, callerId string, platformId string) ([]repository.PlatformMemberWithRole, error) {
	result := []repository.PlatformMemberWithRole{}
	err := RequireAuthority(a.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = a.repos.PlatformMember.List(ctx, platformId, 0, 0)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (a member) Get(ctx context.Context, callerId string, platformId string, memberId string) (repository.PlatformMemberWithRole, error) {
	result := repository.PlatformMemberWithRole{}
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

func (a member) UpdateMember(ctx context.Context, request model.RequestMemberUpdateOther, callerId string, memberId string) (model.MemberToRole, error) {
	result := model.MemberToRole{}
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

	if request.Role == "Owner" || request.Role == "owner" {
		return result, common.StringError(errors.New("cannot elevate member to owner"))
	}

	role, err := a.repos.MemberToRole.GetByMember(memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	role.RoleID = GetRoleId(request.Role)
	a.repos.MemberToRole.UpdateRole(memberId, role)

	result = role

	return result, nil
}

func (a member) UpdateSelf(ctx context.Context, request model.RequestMemberUpdateSelf, callerId string) (repository.PlatformMemberWithRole, error) {
	result := repository.PlatformMemberWithRole{}

	// Actual DB request
	type UpdateRequest struct {
		Name     string `json:"name" db:"name"`
		Password string `json:"password" db:"password"`
	}
	updateRequest := UpdateRequest{}

	// If password change is requested, verify the current password
	if request.NewPassword != nil && request.OldPassword != nil {
		if len(*request.NewPassword) < 8 {
			return result, common.StringError(errors.New("invalid password length")) // TODO: intensify sophistication
		}
		m, err := a.repos.PlatformMember.GetById(ctx, callerId)
		if err != nil {
			return result, common.StringError(err)
		}
		if bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(*request.OldPassword)) != nil {
			return result, common.StringError(errors.New("invalid password"))
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*request.NewPassword), 8)
		if err != nil {
			return result, common.StringError(err)
		}
		updateRequest.Password = string(hash)
	}

	member, err := a.repos.PlatformMember.GetById(ctx, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	// Keep old name if "" or none was passed in
	if request.Name == nil || *request.Name == "" {
		request.Name = &member.Name
	}
	updateRequest.Name = *request.Name

	err = a.repos.PlatformMember.Update(ctx, callerId, updateRequest)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = a.repos.PlatformMember.GetById(ctx, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (a member) SendPasswordResetEmail(ctx context.Context, email string) error {
	// Anyone can request this, so ensure member is associated with 'email'
	member, err := a.repos.PlatformMember.GetByEmail(ctx, email)
	if err != nil {
		return common.StringError(err)
	}

	// generate reset token by hashing member id
	secret := os.Getenv("STRING_ENCRYPTION_KEY")
	resetToken, err := common.EncryptString(member.ID, secret)

	// url encode token
	resetToken = url.QueryEscape(resetToken)

	if err != nil {
		return common.StringError(err)
	}

	body := "" +
		"<a href='https://www.string.xyz'>" +
		"<img src='https://uploads-ssl.webflow.com/63163482142485bcffc0cd47/6318c58524a46f188e0adef6_Logo-dark-lg-p-500.png'></img></a>" +
		"<header>You have requested a password reset for the String API</header>" +
		"<br>Dear " + member.Name + "," +
		"<br>If you have forgotten your password, click the link below to reset it:" +
		"<br><a href='" + os.Getenv("BASE_DASHBOARD_URL") + "/members/password-reset/" + resetToken + "'>Reset Password</a>" // TODO: Change URL to Password Reset page and include resetToken

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

	encrypted, err := bcrypt.GenerateFromPassword([]byte(request.Password), 8)
	if err != nil {
		return common.StringError(err)
	}
	// Update their password
	type UpdatePW struct {
		Password string `json:"password" db:"password"`
	}
	updatePW := UpdatePW{Password: string(encrypted)}

	err = a.repos.PlatformMember.Update(ctx, memberId, updatePW)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (a member) Deactivate(ctx context.Context, callerId string, memberId string) (repository.PlatformMemberWithRole, error) {
	result := repository.PlatformMemberWithRole{}
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

	// Ensure not already deactivated || never existed
	result, err = a.repos.PlatformMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	type DeactivateMember struct {
		DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	}
	now := time.Now()
	deactivateMember := DeactivateMember{DeactivatedAt: &now}

	err = a.repos.PlatformMember.Update(ctx, memberId, deactivateMember)
	if err != nil {
		return result, common.StringError(err)
	}

	// Once a member has been deactivated,
	result, err = a.repos.PlatformMember.GetByIdIncludingDeactivated(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}
