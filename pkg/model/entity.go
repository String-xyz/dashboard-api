package model

import (
	"time"

	"github.com/lib/pq"
)

type Organization struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ActivatedAt   *time.Time `json:"activatedAt,omitempty" db:"activated_at"`
	Name          string     `json:"name" db:"name"`
	Description   string     `json:"description" db:"description"`
}

type Platform struct {
	Id            string         `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Domains       pq.StringArray `json:"domains" db:"domains"`
	IPAddresses   pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type OrganizationMember struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"-" db:"password"`
	Name          string     `json:"name" db:"name"`
}

type MemberToOrganization struct {
	OrganizationId string `json:"organizationId" db:"organization_id"`
	MemberId       string `json:"memberId" db:"member_id"`
}

type MemberRole struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string     `json:"name" db:"name"`
}

type MemberToRole struct {
	MemberId string `json:"memberId" db:"member_id"`
	RoleId   string `json:"roleId" db:"role_id"`
}

type MemberInvite struct {
	Id             string     `json:"id,omitempty" db:"id"`
	CreatedAt      time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt  *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ExpiredAt      *time.Time `json:"expiredAt,omitempty" db:"expired_at"`
	AcceptedAt     *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	Email          string     `json:"email" db:"email"`
	InvitedBy      *string    `json:"invitedBy,omitempty" db:"invited_by"`
	OrganizationId string     `json:"organizationId" db:"organization_id"`
	RoleId         string     `json:"roleId" db:"role_id"`
	Name           string     `json:"name" db:"name"`
}

type Apikey struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Data          string     `json:"data" db:"data"`
	Hint          string     `json:"hint,omitempty" db:"hint"`
	Description   *string    `json:"description,omitempty" db:"description"`
	CreatedBy     string     `json:"createdBy" db:"created_by"`
	PlatformId    string     `json:"platformId" db:"platform_id"`
}
type Contract struct {
	Id            string         `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string         `json:"name" db:"name"`
	Address       string         `json:"address" db:"address"`
	Functions     pq.StringArray `json:"functions" db:"functions"`
	NetworkId     string         `json:"networkId" db:"network_id"`
	PlatformId    string         `json:"platformId" db:"platform_id"`
}

type NetworkData struct {
	Id   string `json:"id,omitempty" db:"id"`
	Name string `json:"name" db:"name"`
}
