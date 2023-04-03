package repository

type Repositories struct {
	Platform         Platform
	PlatformMember   PlatformMember
	MemberToRole     MemberToRole
	MemberToPlatform MemberToPlatform
	MemberRole       MemberRole
	MemberInvite     MemberInvite
	Apikey           Apikey
	DenyList         DenyList
	Contract         Contract
	Network          Network
}
