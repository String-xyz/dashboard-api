package model

import (
	"time"

	"github.com/lib/pq"
)

type Platform struct {
	ID            string         `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ActivatedAt   *time.Time     `json:"activatedAt,omitempty" db:"activated_at"`
	Name          string         `json:"name" db:"name"`
	Description   string         `json:"description" db:"description"`
	Domains       pq.StringArray `json:"domains" db:"domains"`
	IPAddresses   pq.StringArray `json:"ipAddresses" db:"ip_addresses"`
}

type PlatformMember struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"-" db:"password"`
	Name          string     `json:"name" db:"name"`
}

type MemberToPlatform struct {
	PlatformId string `json:"platformId" db:"platform_id"`
	MemberID   string `json:"memberId" db:"member_id"`
}

type MemberRole struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string     `json:"name" db:"name"`
}

type MemberToRole struct {
	MemberID string `json:"memberId" db:"member_id"`
	RoleID   string `json:"roleId" db:"role_id"`
}

type MemberInvite struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ExpiredAt     *time.Time `json:"expiredAt,omitempty" db:"expired_at"`
	AcceptedAt    *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	Email         string     `json:"email" db:"email"`
	InvitedBy     *string    `json:"invitedBy,omitempty" db:"invited_by"`
	PlatformId    string     `json:"platformId" db:"platform_id"`
	RoleID        string     `json:"roleId" db:"role_id"`
	Name          string     `json:"name" db:"name"`
}

type Apikey struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Public        string     `json:"public" db:"public"` // an unhased public key
	Secret        *string    `json:"secret" db:"secret"` // a hashed secret key
	Description   *string    `json:"description" db:"description"`
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
