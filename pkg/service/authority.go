package service

import (
	"os"

	"github.com/String-xyz/go-lib/common"
	serrors "github.com/String-xyz/go-lib/stringerror"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

// In: Handle to repos, PlatformMember ID of caller, allowable roles (["Member", "Admin", "Owner"])
// TODO: Refactor to use enums instead of strings
func RequireAuthority(repos repository.Repositories, callerId string, rolesAllowed ...string) error {
	ownerId := os.Getenv("MEMBER_ROLE_OWNER_ID")
	adminId := os.Getenv("MEMBER_ROLE_ADMIN_ID")
	memberId := os.Getenv("MEMBER_ROLE_MEMBER_ID")

	// Convert rolesAllowed from Name to UUID
	var roleIdsAllowed []string
	for _, stringRole := range rolesAllowed {
		if stringRole == "Member" || stringRole == "member" {
			roleIdsAllowed = append(roleIdsAllowed, memberId)
		} else if stringRole == "Admin" || stringRole == "admin" {
			roleIdsAllowed = append(roleIdsAllowed, adminId)
		} else if stringRole == "Owner" || stringRole == "owner" {
			roleIdsAllowed = append(roleIdsAllowed, ownerId)
		}
	}

	// Get the callers role (they can only have one)
	callerRole, err := repos.MemberToRole.GetByMember(callerId)
	if err != nil {
		return common.StringError(err)
	}

	// Is the callers role allowed?
	for _, allowed := range roleIdsAllowed {
		if callerRole.RoleID == allowed {
			return nil
		}
	}

	return common.StringError(serrors.FORBIDDEN)
}

func GetRole(repos repository.Repositories, memberId string) (string, error) {
	role := "Not Found"
	ownerRoleId := os.Getenv("MEMBER_ROLE_OWNER_ID")
	adminRoleId := os.Getenv("MEMBER_ROLE_ADMIN_ID")
	memberRoleId := os.Getenv("MEMBER_ROLE_MEMBER_ID")

	memberToRole, err := repos.MemberToRole.GetByMember(memberId)
	if err != nil {
		return role, common.StringError(err)
	}

	// Ordered by statistical distribution
	if memberToRole.RoleID == memberRoleId {
		role = "Member"
	} else if memberToRole.RoleID == adminRoleId {
		role = "Admin"
	} else if memberToRole.RoleID == ownerRoleId {
		role = "Owner"
	}

	return role, nil
}

func GetRoleId(roleName string) string {
	role := "Not Found"
	ownerRoleId := os.Getenv("MEMBER_ROLE_OWNER_ID")
	adminRoleId := os.Getenv("MEMBER_ROLE_ADMIN_ID")
	memberRoleId := os.Getenv("MEMBER_ROLE_MEMBER_ID")

	if roleName == "Member" || roleName == "member" {
		role = memberRoleId
	} else if roleName == "Admin" || roleName == "admin" {
		role = adminRoleId
	} else if roleName == "Owner" || roleName == "owner" {
		role = ownerRoleId
	}
	return role
}
