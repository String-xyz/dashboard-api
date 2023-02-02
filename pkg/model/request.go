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
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RequestInviteUpdate struct {
	Status string `json:"status"`
	Role   string `json:"role"`
}

type RequestLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestMemberUpdate struct {
	Email string `json:"email" db:"email"`
	Name  string `json:"name" db:"name"`
}

type RequestApikeyUpdate struct {
	TodoDefineParameters string `json:"undefined"`
}

type RequestPasswordReset struct {
	Password   string `json:"password"`
	ResetToken string `json:"resetToken"`
}
