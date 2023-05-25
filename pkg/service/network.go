package service

import (
	"context"

	"github.com/String-xyz/dashboard-api/pkg/model"
	"github.com/String-xyz/dashboard-api/pkg/repository"
	"github.com/String-xyz/go-lib/v2/common"
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

func (n network) GetAll(ctx context.Context) ([]model.NetworkData, error) {
	_, finish := Span(ctx, "service.network.GetAll")
	defer finish()

	networks, err := n.repos.Network.List(ctx, 0, 0)
	if err != nil {
		return nil, common.StringError(err)
	}

	return networks, nil
}
