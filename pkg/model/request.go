package model

import "github.com/lib/pq"

type RequestOrganizationCreate struct {
	OrganizationName string `json:"organizationName" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	Name             string `json:"name" validate:"required"`
}

type RequestOrganizationUpdate struct {
	Name        *string `json:"organizationName" db:"name"`
	Description *string `json:"description" db:"description"`
}

type RequestPlatformCreate struct {
	PlatformName        string `json:"platformName" validate:"required"`
	PlatformDescription string `json:"platformDescription"`
}

type RequestPlatformUpdate struct {
	Name        *string         `json:"platformName" db:"name"`
	Description *string         `json:"description" db:"description"`
	Domains     *pq.StringArray `json:"domains" db:"domains"`
	IPAddresses *pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type RequestInviteSend struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=admin member owner Admin Member Owner"`
}

type RequestInviteUpdate struct {
	Role string `json:"role"`
	Name string `json:"name"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// type RequestMemberUpdate struct {
// 	Email string `json:"email" db:"email"`
// 	Name  string `json:"name" db:"name"`

// 	// Not part of PlatformMember entity but we still want to be able to change this from the Update endpoint
// 	Role string `json:"role" db:"role"`
// }

type RequestInviteAcceptance struct {
	Password string `json:"password" validate:"required,min=8,max=100"`
	Token    string `json:"token" validate:"required,base64"`
}

type RequestMemberUpdateSelf struct {
	Name        *string `json:"name" db:"name"`
	OldPassword *string `json:"oldPassword" validate:"required_with=NewPassword,min=8,max=100"`
	NewPassword *string `json:"newPassword" validate:"required_with=OldPassword,min=8,max=100"`
}

type RequestMemberUpdateOther struct {
	Role string `json:"role" db:"role"`
}

type RequestTransferOwnership struct {
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RequestApikeyUpdate struct {
	Description string `json:"description"`
}

type RequestPasswordResetEmail struct {
	Email string `query:"email" validate:"required,email"`
}

type RequestPasswordReset struct {
	Password   string `json:"password" validate:"required,min=8,max=100"`
	ResetToken string `json:"resetToken" validate:"required"`
}

type RequestContractCreate struct {
	Name      string         `json:"name" db:"name"`
	Address   string         `json:"address" db:"address"`
	Functions pq.StringArray `json:"functions" db:"functions"`
	NetworkId string         `json:"networkId" db:"network_id"`
}

type RequestContractUpdate struct {
	Name      *string         `json:"name" db:"name"`
	Address   *string         `json:"address" db:"address"`
	Functions *pq.StringArray `json:"functions" db:"functions"`
	NetworkId *string         `json:"networkId" db:"network_id"`
}
