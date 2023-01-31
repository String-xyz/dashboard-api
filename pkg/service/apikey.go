package service

import (
	"context"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Apikey interface {
	Create(ctx context.Context, request string) (model.Apikey, error)
	GetAll(ctx context.Context, request string) ([]model.Apikey, error)
	Get(ctx context.Context, request string, id string) (model.Apikey, error)
	Deactivate(ctx context.Context, request string, id string) (model.Apikey, error)
	Update(ctx context.Context, request model.RequestApikeyUpdate, id string) (model.Apikey, error)
}

type apikey struct {
	repos repository.Repositories
}

func NewApikey(repos repository.Repositories) Apikey {
	return &apikey{repos}
}

func (a apikey) Create(ctx context.Context, request string) (model.Apikey, error) {
	return model.Apikey{}, nil
}

func (a apikey) GetAll(ctx context.Context, request string) ([]model.Apikey, error) {
	return nil, nil
}

func (a apikey) Get(ctx context.Context, request string, id string) (model.Apikey, error) {
	return model.Apikey{}, nil
}

func (a apikey) Deactivate(ctx context.Context, request string, id string) (model.Apikey, error) {
	return model.Apikey{}, nil
}

func (a apikey) Update(ctx context.Context, request model.RequestApikeyUpdate, id string) (model.Apikey, error) {
	return model.Apikey{}, nil
}
