package service

import (
	"context"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type MemberCreateResponse struct {
	JWT  JWT                  `json:"authToken"`
	User model.PlatformMember `json:"member"`
}

type Member interface {
	GetAll(ctx context.Context, request string) ([]model.PlatformMember, error)
	Get(ctx context.Context, request string, id string) (model.PlatformMember, error)
	Update(ctx context.Context, request model.RequestMemberUpdate, id string) (model.PlatformMember, error)
}

type member struct {
	repos repository.Repositories
}

func NewMember(repos repository.Repositories) Member {
	return &member{repos}
}

func (a member) GetAll(ctx context.Context, request string) ([]model.PlatformMember, error) {
	return nil, nil
}

func (a member) Get(ctx context.Context, request string, id string) (model.PlatformMember, error) {
	return model.PlatformMember{}, nil
}

func (a member) Update(ctx context.Context, request model.RequestMemberUpdate, id string) (model.PlatformMember, error) {
	return model.PlatformMember{}, nil
}
