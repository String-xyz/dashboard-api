package service

import (
	"context"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Login interface {
	Login(ctx context.Context, request model.RequestLogin) (interface{}, error) // TODO: Return JWT
}

type login struct {
	repos repository.Repositories
}

func NewLogin(repos repository.Repositories) Login {
	return &login{repos}
}

func (a login) Login(ctx context.Context, request model.RequestLogin) (interface{}, error) {
	return "{}", nil
}
