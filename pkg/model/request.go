package model

type RequestPlatformCreate struct {
	PlatformName string `json:"platformName"`
	Email        string `json:"email"`
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
	TodoDefineParameters string `json:"undefined"`
}

type RequestApikeyUpdate struct {
	TodoDefineParameters string `json:"undefined"`
}
