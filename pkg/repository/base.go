package repository

import "errors"

var ErrNotFound = errors.New("not found")

type Repositories struct {
	Platform         Platform
	PlatformMember   PlatformMember
	MemberToRole     MemberToRole
	MemberToPlatform MemberToPlatform
	MemberRole       MemberRole
	MemberInvite     MemberInvite
	Apikey           Apikey
}
