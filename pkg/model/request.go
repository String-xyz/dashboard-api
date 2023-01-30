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
