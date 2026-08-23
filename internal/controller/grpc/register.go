package grpc

import (
	"context"
	"fmt"

	proto "github.com/RenterRus/sausage-profile/docs/proto/v1"
	"github.com/RenterRus/sausage-profile/internal/entity"
)

func (t *Manager) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req == nil || req.Login == "" {
		return nil, fmt.Errorf("Register: %w", entity.ErrParametrNoFound)
	}

	url, err := t.profile.Registration(ctx, req.GetLogin())
	if err != nil {
		return nil, fmt.Errorf("Registration: %w", err)
	}

	return &proto.RegisterResponse{
		Url: url,
	}, nil
}

func (t *Manager) Confirm(ctx context.Context, req *proto.AcceptRequest) (*proto.AcceptResponse, error) {
	if req == nil || req.Login == "" || req.OtpCode == "" {
		return nil, fmt.Errorf("Confirm: %w", entity.ErrParametrNoFound)
	}

	if err := t.profile.Confirmed(ctx, req.GetLogin(), req.GetOtpCode()); err != nil {
		return &proto.AcceptResponse{
			Status: entity.STATUS_FAILED,
		}, fmt.Errorf("Confirm.Confirmed: %w", err)
	}

	return &proto.AcceptResponse{
		Status: entity.STATUS_OK,
	}, nil
}
