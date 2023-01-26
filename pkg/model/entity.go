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

type Member struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"password" db:"password"`
}

type PlatformMember struct {
	PlatformID string `json:"platformId" db:"platform_id"`
	MemberID   string `json:"memberId" db:"member_id"`
}

type Role struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name          string     `json:"name" db:"name"`
}

type MemberRole struct {
	MemberID string `json:"memberId" db:"member_id"`
	RoleID   string `json:"roleId" db:"role_id"`
}

type Invite struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	ExpiredAt     *time.Time `json:"expiredAt,omitempty" db:"expired_at"`
	AcceptedAt    *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	Email         string     `json:"email" db:"email"`
	InvitedBy     string     `json:"invitedBy" db:"invited_by"`
	PlatformID    string     `json:"platformId" db:"platform_id"`
}

type Apikey struct {
	ID            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Type          string     `json:"type" db:"type"`
	Data          string     `json:"data" db:"data"`
	Description   string     `json:"description" db:"description"`
	CreatedBy     string     `json:"createdBy" db:"created_by"`
	PlatformID    string     `json:"platformId" db:"platform_id"`
}
