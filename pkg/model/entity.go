package model

import (
	"time"

	"github.com/lib/pq"
)

type Organization struct {
	Id          string     `json:"id,omitempty" db:"id"`
	CreatedAt   time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description" db:"description"`
}

type Platform struct {
	Id             string         `json:"id,omitempty" db:"id"`
	CreatedAt      time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt      *time.Time     `json:"deletedAt,omitempty" db:"deleted_at"`
	DeactivatedAt  *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	Name           string         `json:"name" db:"name"`
	Description    string         `json:"description" db:"description"`
	Domains        pq.StringArray `json:"domains" db:"domains" swaggertype:"array,string"`
	IPAddresses    pq.StringArray `json:"ipAddresses" db:"ip_addresses" swaggertype:"array,string"`
	OrganizationId string         `json:"organizationId" db:"organization_id"`
}

type OrganizationMember struct {
	Id            string     `json:"id,omitempty" db:"id"`
	CreatedAt     time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"-" db:"password"`
	Name          string     `json:"name" db:"name"`
}

type MemberToOrganization struct {
	OrganizationId string `json:"organizationId" db:"organization_id"`
	MemberId       string `json:"memberId" db:"member_id"`
}

type MemberRole struct {
	Id        string     `json:"id,omitempty" db:"id"`
	CreatedAt time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	Name      string     `json:"name" db:"name"`
}

type MemberToRole struct {
	MemberId string `json:"memberId" db:"member_id"`
	RoleId   string `json:"roleId" db:"role_id"`
}

type MemberInvite struct {
	Id             string     `json:"id,omitempty" db:"id"`
	CreatedAt      time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	ExpiredAt      *time.Time `json:"expiredAt,omitempty" db:"expired_at"`
	AcceptedAt     *time.Time `json:"acceptedAt,omitempty" db:"accepted_at"`
	Email          string     `json:"email" db:"email"`
	InvitedBy      *string    `json:"invitedBy,omitempty" db:"invited_by"`
	OrganizationId string     `json:"organizationId" db:"organization_id"`
	RoleId         string     `json:"-" db:"role_id"`
	Name           string     `json:"name" db:"name"`
}

type Apikey struct {
	Id             string     `json:"id,omitempty" db:"id"`
	CreatedAt      time.Time  `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
	DeletedAt      *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	Type           string     `json:"type" db:"type"`
	Data           string     `json:"data" db:"data"`
	Hint           string     `json:"hint,omitempty" db:"hint"`
	Description    *string    `json:"description,omitempty" db:"description"`
	CreatedBy      string     `json:"createdBy" db:"created_by"`
	PlatformId     *string    `json:"platformId" db:"platform_id"`
	OrganizationId string     `json:"organizationId" db:"organization_id"`
}

type Contract struct {
	Id             string         `json:"id,omitempty" db:"id"`
	CreatedAt      time.Time      `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt      time.Time      `json:"updatedAt,omitempty" db:"updated_at"`
	DeactivatedAt  *time.Time     `json:"deactivatedAt,omitempty" db:"deactivated_at"`
	DeletedAt      *time.Time     `json:"deletedAt,omitempty" db:"deleted_at"`
	Name           string         `json:"name" db:"name"`
	Address        string         `json:"address" db:"address"`
	Functions      pq.StringArray `json:"functions" db:"functions" swaggertype:"array,string"`
	Type           string         `json:"type" db:"type"`
	NetworkId      string         `json:"networkId" db:"network_id"`
	OrganizationId string         `json:"organizationId" db:"organization_id"`
	PlatformIds    pq.StringArray `json:"platformIds" db:"platform_ids" swaggertype:"array,string"`
}

type NetworkFull struct {
	Id          string     `json:"id" db:"id"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"`
	Name        string     `json:"name" db:"name"`
	NetworkId   uint64     `json:"networkId" db:"network_id"`
	ChainId     uint64     `json:"chainId" db:"chain_id"`
	GasTokenId  string     `json:"gasTokenId" db:"gas_token_id"`
	GasOracle   string     `json:"gasOracle" db:"gas_oracle"`
	RPCUrl      string     `json:"rpcUrl" db:"rpc_url"`
	ExplorerUrl string     `json:"explorerUrl" db:"explorer_url"`
}

type NetworkData struct {
	Id      string `json:"id,omitempty" db:"id"`
	Name    string `json:"name" db:"name"`
	ChainId uint64 `json:"chainId" db:"chain_id"`
}
