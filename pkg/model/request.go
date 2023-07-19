package model

import "github.com/lib/pq"

type RequestOrganizationCreate struct {
	OrganizationName string `json:"organizationName" validate:"required,max=100"`
	Email            string `json:"email" validate:"required,email"`
	Name             string `json:"name" validate:"required,max=100"`
}

type RequestOrganizationUpdate struct {
	Name        *string `json:"organizationName" db:"name" validate:"omitempty,max=100"`
	Description *string `json:"description" db:"description" validate:"omitempty,max=144"`
}

type RequestPlatformCreate struct {
	PlatformName        string `json:"platformName" validate:"required"`
	PlatformDescription string `json:"platformDescription"`
}

type RequestPlatformUpdate struct {
	Name        *string         `json:"platformName" db:"name"`
	Description *string         `json:"description" db:"description"`
	Domains     *pq.StringArray `json:"domains" db:"domains" swaggertype:"array,string"`
	IPAddresses *pq.StringArray `json:"ipAddresses" db:"ip_addresses" swaggertype:"array,string"`
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
	OldPassword *string `json:"oldPassword" validate:"required_with_all=OldPassword,NewPassword,min=8,max=100"`
	NewPassword *string `json:"newPassword" validate:"required_with_all=NewPassword,OldPassword,min=8,max=100"`
}

type RequestMemberUpdateOther struct {
	Role string `json:"role" db:"role"`
}

type RequestTransferOwnership struct {
	Password string `json:"password" validate:"required,min=8,max=100"`
}

type RequestAPIKeyCreate struct {
	PlatformId string `json:"platformId" validate:"required,uuid"`
}

type RequestApikeyUpdate struct {
	Description string `json:"description" validate:"omitempty,max=144"`
}

type RequestPasswordResetEmail struct {
	Email string `query:"email" validate:"required,email"`
}

type RequestPasswordReset struct {
	Password   string `json:"password" validate:"required,min=8,max=100"`
	ResetToken string `json:"resetToken" validate:"required"`
}

type RequestContractCreate struct {
	Name           string         `json:"name" db:"name" validate:"required,max=100"`
	Address        string         `json:"address" db:"address" validate:"required,eth_addr"`
	Functions      pq.StringArray `json:"functions" db:"functions" validate:"required" swaggertype:"array,string" `
	Type           string         `json:"type" db:"type" validate:"required,oneof=NFT TOKEN NFT_AND_TOKEN"`
	NetworkId      string         `json:"networkId" db:"network_id" validate:"required,uuid"`
	PlatformIds    pq.StringArray `json:"platformIds" db:"platform_ids" validate:"required,min=1" swaggertype:"array,string"`
	OrganizationId string         `json:"-" db:"organization_id"`
}

type RequestContractUpdate struct {
	Name        *string        `json:"name" validate:"omitempty,max=100"`
	Address     *string        `json:"address" validate:"omitempty,eth_addr"`
	Functions   pq.StringArray `json:"functions" swaggertype:"array,string"`
	Type        *string        `json:"type" validate:"omitempty,oneof=NFT TOKEN NFT_AND_TOKEN"`
	NetworkId   *string        `json:"networkId" validate:"omitempty,uuid"`
	PlatformIds pq.StringArray `json:"platformIds" validate:"omitempty,min=1" swaggertype:"array,string"`
}
