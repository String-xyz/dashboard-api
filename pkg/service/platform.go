package service

import (
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Platform interface {
	Create() (model.Platform, error)
}

type platform struct {
	repos repository.Repositories
}

func NewPlatform(repos repository.Repositories) Platform {
	return &platform{repos}
}

func (a platform) Create() (model.Platform, error) {
	return model.Platform{}, nil
}
