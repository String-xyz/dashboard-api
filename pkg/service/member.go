package service

import (
	"context"
	"net/url"
	"strings"

	"github.com/String-xyz/go-lib/v2/common"
	serror "github.com/String-xyz/go-lib/v2/stringerror"
	"golang.org/x/crypto/bcrypt"

	"github.com/String-xyz/dashboard-api/config"
	"github.com/String-xyz/dashboard-api/pkg/internal/emailer"
	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/repository"
)

type MemberCreateResponse struct {
	JWT    JWT                                   `json:"authToken"`
	Member repository.OrganizationMemberWithRole `json:"member"`
}

type Member interface {
	GetAll(ctx context.Context, organizationId string) ([]repository.OrganizationMemberWithRole, error)
	Get(ctx context.Context, callerId string, organizationId string, memberId string) (repository.OrganizationMemberWithRole, error)
	UpdateMember(ctx context.Context, request model.RequestMemberUpdateOther, callerId string, memberId string) (repository.OrganizationMemberWithRole, error)
	UpdateSelf(ctx context.Context, request model.RequestMemberUpdateSelf, callerId string) (repository.OrganizationMemberWithRole, error)
	TransferOwnership(ctx context.Context, request model.RequestTransferOwnership, callerId string, memberId string) (repository.OrganizationMemberWithRole, error)
	SendPasswordResetEmail(ctx context.Context, request model.RequestPasswordResetEmail) error
	PasswordReset(ctx context.Context, request model.RequestPasswordReset) error
	Deactivate(ctx context.Context, callerId string, memberId string) (repository.OrganizationMemberWithRole, error)
	Reactivate(ctx context.Context, callerId string, memberId string) (repository.OrganizationMemberWithRole, error)
}

type member struct {
	repos repository.Repositories
}

func NewMember(repos repository.Repositories) Member {
	return &member{repos}
}

func (m member) GetAll(ctx context.Context, organizationId string) ([]repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.GetAll", SpanTag{"organizationId": organizationId})
	defer finish()

	result, err := m.repos.OrganizationMember.List(ctx, organizationId, 0, 0)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (m member) Get(ctx context.Context, callerId string, organizationId string, memberId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.Get", SpanTag{"organizationId": organizationId})
	defer finish()

	result := repository.OrganizationMemberWithRole{}
	err := RequireAuthority(m.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	memberToOrganization, err := m.repos.MemberToOrganization.GetByMember(memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if memberToOrganization.OrganizationId != organizationId {
		return result, common.StringError(serror.FORBIDDEN)
	}

	result, err = m.repos.OrganizationMember.GetById(ctx, memberToOrganization.MemberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (m member) UpdateSelf(ctx context.Context, request model.RequestMemberUpdateSelf, callerId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.UpdateSelf")
	defer finish()

	result := repository.OrganizationMemberWithRole{}

	// Actual DB request
	type UpdateRequest struct {
		Name     string `json:"name" db:"name"`
		Password string `json:"password" db:"password"`
	}
	updateRequest := UpdateRequest{}

	// If password change is requested, verify the current password
	if request.NewPassword != nil && request.OldPassword != nil {
		if len(*request.NewPassword) < 8 {
			return result, common.StringError(serror.INVALID_PASSWORD)
		}
		m, err := m.repos.OrganizationMember.GetById(ctx, callerId)
		if err != nil {
			return result, common.StringError(err)
		}
		if bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(*request.OldPassword)) != nil {
			return result, common.StringError(serror.INVALID_PASSWORD)
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*request.NewPassword), 8)
		if err != nil {
			return result, common.StringError(err)
		}
		updateRequest.Password = string(hash)
	}

	member, err := m.repos.OrganizationMember.GetById(ctx, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	// Keep old name if "" or none was passed in
	if request.Name == nil || *request.Name == "" {
		request.Name = &member.Name
	}
	updateRequest.Name = *request.Name

	err = m.repos.OrganizationMember.Update(ctx, callerId, updateRequest)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = m.repos.OrganizationMember.GetById(ctx, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (m member) SendPasswordResetEmail(ctx context.Context, request model.RequestPasswordResetEmail) error {
	_, finish := Span(ctx, "service.member.SendPasswordResetEmail")
	defer finish()

	email := request.Email

	// Anyone can request this, so ensure member is associated with 'email'
	member, err := m.repos.OrganizationMember.GetByEmail(ctx, email)
	if err != nil {
		return common.StringError(err)
	}

	// generate reset token by hashing member id
	secret := config.Var.STRING_ENCRYPTION_KEY
	resetToken, err := common.EncryptString(member.Id, secret)
	if err != nil {
		return common.StringError(err)
	}

	// url encode token
	resetToken = url.QueryEscape(resetToken)

	emailer := emailer.New()
	return emailer.SendPasswordResetEmail(ctx, email, resetToken, member.Name)
}

func (m member) PasswordReset(ctx context.Context, request model.RequestPasswordReset) error {
	_, finish := Span(ctx, "service.member.PasswordReset")
	defer finish()

	// Get the member Id from the password reset token
	secret := config.Var.STRING_ENCRYPTION_KEY
	memberId, err := common.DecryptString(request.ResetToken, secret)
	if err != nil {
		if strings.Contains(err.Error(), "illegal base64") {
			return common.StringError(serror.INVALID_RESET_TOKEN)
		}
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

	err = m.repos.OrganizationMember.Update(ctx, memberId, updatePW)
	if err != nil {
		return common.StringError(err)
	}

	return nil
}

func (m member) UpdateMember(ctx context.Context, request model.RequestMemberUpdateOther, callerId string, memberId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.UpdateMember")
	defer finish()

	result := repository.OrganizationMemberWithRole{}

	err := RequireAuthority(m.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	// Admin's can't edit Owners or Admins
	callerRole, err := GetRole(m.repos, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	memberRole, err := GetRole(m.repos, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if callerRole == "Admin" && memberRole != "Member" {
		return result, common.StringError(serror.FORBIDDEN)
	}

	if strings.ToLower(request.Role) == "owner" {
		return result, common.StringError(serror.FORBIDDEN)
	}

	role, err := m.repos.MemberToRole.GetByMember(memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	role.RoleId = GetRoleId(request.Role)
	err = m.repos.MemberToRole.UpdateRole(memberId, role)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (m member) TransferOwnership(ctx context.Context, request model.RequestTransferOwnership, callerId string, memberId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.TransferOwnership")
	defer finish()

	result := repository.OrganizationMemberWithRole{}

	caller, err := m.repos.OrganizationMember.GetById(ctx, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	if caller.Role != "Owner" {
		return result, common.StringError(serror.FORBIDDEN)
	}

	// Check for redundancy
	if callerId == memberId {
		return caller, nil
	}

	// Check member exists
	_, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	// Verify password
	if request.Password == "" {
		return result, common.StringError(serror.INVALID_PASSWORD)
	}

	if bcrypt.CompareHashAndPassword([]byte(caller.Password), []byte(request.Password)) != nil {
		return result, common.StringError(serror.INVALID_PASSWORD)
	}

	// Execute role updates in single db transaction
	m.repos.MemberToRole.MustBegin()

	defer m.repos.MemberToRole.Reset()

	// Promote member to Owner
	roleObjMember := model.MemberToRole{MemberId: memberId, RoleId: GetRoleId("Owner")}

	err = m.repos.MemberToRole.UpdateRole(memberId, roleObjMember)
	if err != nil {
		m.repos.MemberToRole.Rollback()
		return result, common.StringError(serror.NOT_FOUND)
	}

	// Demote caller to Admin
	roleObjCaller := model.MemberToRole{MemberId: callerId, RoleId: GetRoleId("Admin")}

	err = m.repos.MemberToRole.UpdateRole(callerId, roleObjCaller)
	if err != nil {
		m.repos.MemberToRole.Rollback()
		return result, common.StringError(err)
	}

	if err := m.repos.MemberToRole.Commit(); err != nil {
		return result, common.StringError(err)
	}

	result, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}

func (m member) Deactivate(ctx context.Context, callerId string, memberId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.Deactivate")
	defer finish()

	result := repository.OrganizationMemberWithRole{}

	err := RequireAuthority(m.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	// Admin's can't edit Owners or Admins

	/***** TODO: This code is repeated multiple times, refactor into a function *****/
	callerRole, err := GetRole(m.repos, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	memberRole, err := GetRole(m.repos, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if callerRole == "Admin" && memberRole != "Member" {
		return result, common.StringError(serror.FORBIDDEN)
	}
	/**********************************************************************************/

	// Ensure not already deactivated || never existed
	result, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	m.repos.OrganizationMember.Deactivate(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	// Once a member has been deactivated,
	result, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	// add member to denylist
	go m.repos.DenyList.AddMember(memberId)

	return result, nil
}

func (m member) Reactivate(ctx context.Context, callerId string, memberId string) (repository.OrganizationMemberWithRole, error) {
	_, finish := Span(ctx, "service.member.Reactivate")
	defer finish()

	result := repository.OrganizationMemberWithRole{}

	err := RequireAuthority(m.repos, callerId, "Owner", "Admin")
	if err != nil {
		return result, common.StringError(err)
	}

	// Refactor this into it's own function, since it's called in several places ---
	callerRole, err := GetRole(m.repos, callerId)
	if err != nil {
		return result, common.StringError(err)
	}

	memberRole, err := GetRole(m.repos, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	if callerRole == "Admin" && memberRole != "Member" {
		return result, common.StringError(serror.FORBIDDEN)
	}
	// -----------------------------------------------------------------------------

	m.repos.OrganizationMember.Activate(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	err = m.repos.DenyList.RemoveMember(memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	result, err = m.repos.OrganizationMember.GetById(ctx, memberId)
	if err != nil {
		return result, common.StringError(err)
	}

	return result, nil
}
