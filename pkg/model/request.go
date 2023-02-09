package model

import "github.com/lib/pq"

type RequestPlatformCreate struct {
	PlatformName string `json:"platformName" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	Name         string `json:"name" validate:"required"`
}

type RequestPlatformUpdate struct {
	Name        *string         `json:"platformName" db:"name"`
	Description *string         `json:"description" db:"description"`
	Domains     *pq.StringArray `json:"domains" db:"domains"`
	IPAddresses *pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type RequestInviteSend struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RequestInviteUpdate struct {
	Role string `json:"role"`
	Name string `json:"name"`
}

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type RequestMemberUpdate struct {
// 	Email string `json:"email" db:"email"`
// 	Name  string `json:"name" db:"name"`

// 	// Not part of PlatformMember entity but we still want to be able to change this from the Update endpoint
// 	Role string `json:"role" db:"role"`
// }

type RequestInviteAcceptance struct {
	Id       *string `json:"id"`
	Password string  `json:"password"`
}

type RequestMemberUpdateSelf struct {
	Name        *string `json:"name" db:"name"`
	OldPassword *string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
}

type RequestMemberUpdateOther struct {
	Role string `json:"role" db:"role"`
}

type RequestApikeyUpdate struct {
	Description string `json:"description"`
}

type RequestPasswordReset struct {
	Password   string `json:"password"`
	ResetToken string `json:"resetToken"`
}
