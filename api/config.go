package api

import (
	"github.com/String-xyz/platform-admin-api/pkg/repository"
	"github.com/String-xyz/platform-admin-api/pkg/service"
)

func newRepos(config APIConfig) repository.Repositories {
	return repository.Repositories{
		Platform: repository.NewPlatform(config.DB),
	}
}

func newServices(config APIConfig, repos repository.Repositories) service.Services {
	plat := service.NewPlatform(repos)
	return service.Services{Platform: plat}
}
