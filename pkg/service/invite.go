package service

import (
	"context"

	"github.com/String-xyz/platform-admin-api/pkg/model"
	"github.com/String-xyz/platform-admin-api/pkg/repository"
)

type Invite interface {
	Send(ctx context.Context, request model.RequestInviteSend) (model.MemberInvite, error)
	Accept(ctx context.Context, request string) (model.PlatformMember, error)
	List(ctx context.Context, request string) ([]model.MemberInvite, error)
	Resend(ctx context.Context, request string) (model.MemberInvite, error)
	Update(ctx context.Context, request model.RequestInviteUpdate, id string) (model.MemberInvite, error)
}

type invite struct {
	repos repository.Repositories
}

func NewInvite(repos repository.Repositories) Invite {
	return &invite{repos}
}

func (a invite) Send(ctx context.Context, request model.RequestInviteSend) (model.MemberInvite, error) {
	result := model.MemberInvite{}
	return result, nil
}

func (a invite) Accept(ctx context.Context, request string) (model.PlatformMember, error) {
	result := model.PlatformMember{}
	return result, nil
}

func (a invite) List(ctx context.Context, request string) ([]model.MemberInvite, error) {
	result := []model.MemberInvite{}
	return result, nil
}

func (a invite) Resend(ctx context.Context, request string) (model.MemberInvite, error) {
	result := model.MemberInvite{}
	return result, nil
}

func (a invite) Update(ctx context.Context, request model.RequestInviteUpdate, id string) (model.MemberInvite, error) {
	result := model.MemberInvite{}
	return result, nil
}
