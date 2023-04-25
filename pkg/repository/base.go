package repository

type Repositories struct {
	Organization         Organization
	Platform             Platform
	OrganizationMember   OrganizationMember
	MemberToRole         MemberToRole
	MemberToOrganization MemberToOrganization
	MemberRole           MemberRole
	MemberInvite         MemberInvite
	Apikey               Apikey
	DenyList             DenyList
	Contract             Contract
	Network              Network
}
