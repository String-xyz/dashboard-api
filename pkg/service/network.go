package service

import (
	"context"

	"github.com/String-xyz/go-lib/common"
	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Network interface {
	GetAll(ctx context.Context) ([]model.NetworkData, error)
}

type network struct {
	repos repository.Repositories
}

func NewNetwork(repos repository.Repositories) Network {
	return &network{repos}
}

func (a network) GetAll(ctx context.Context) ([]model.NetworkData, error) {
	networks, err := a.repos.Network.List(ctx, 0, 0)
	if err != nil {
		return nil, common.StringError(err)
	}

	return networks, nil
}
