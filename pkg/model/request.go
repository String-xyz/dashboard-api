package model

type RequestPlatformCreate struct {
	PlatformName string `json:"platformName"`
	Email        string `json:"email"`
	Name         string `json:"name"`
}

type RequestPlatformUpdate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Domains     []string `json:"domains"`
	IPAddresses []string `json:"ipAddresses"`
}

type RequestInviteSend struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Role    string  `json:"role"`
	Invitee *string `json:"invitee"`
}

type RequestInviteUpdate struct {
	Status string `json:"status"`
	Role   string `json:"role"`
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

type RequestMemberUpdateSelf struct {
	Name string `json:"name" db:"name"`
	// Password string `json:"password" db:"password"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
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
