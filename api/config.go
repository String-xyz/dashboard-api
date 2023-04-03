package api

import (
	"github.com/String-xyz/go-lib/database"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"github.com/String-xyz/platform-admin-api/pkg/service"
	"github.com/String-xyz/platform-admin-api/pkg/store"
)

func newRepos(config APIConfig) repository.Repositories {
	return repository.Repositories{
		Platform:         repository.NewPlatform(config.DB),
		PlatformMember:   repository.NewPlatformMember(config.DB),
		MemberToRole:     repository.NewMemberToRole(config.DB),
		MemberToPlatform: repository.NewMemberToPlatform(config.DB),
		MemberRole:       repository.NewMemberRole(config.DB),
		MemberInvite:     repository.NewMemberInvite(config.DB),
		Apikey:           repository.NewApikey(config.DB),
		DenyList:         repository.NewDenyList(config.Redis),
		Contract:         repository.NewContract(config.DB),
		Network:          repository.NewNetwork(config.DB),
	}
}

func newRedis() database.RedisStore {
	return store.NewRedis()
}

func newServices(config APIConfig, repos repository.Repositories, redis database.RedisStore) service.Services {
	plat := service.NewPlatform(repos, redis) // redis here is passed into Invite
	member := service.NewMember(repos)
	auth := service.NewAuth(repos, redis)
	login := service.NewLogin(repos, auth)
	invite := service.NewInvite(repos, redis)
	apikey := service.NewApikey(repos)
	contract := service.NewContract(repos)
	network := service.NewNetwork(repos)

	return service.Services{Platform: plat, Member: member, Auth: auth, Login: login, Invite: invite, Apikey: apikey, Contract: contract, Network: network}
}
