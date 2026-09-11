package grpc

import (
	"context"

	proto "github.com/RenterRus/sausage-auth/docs/proto/v1"
	"github.com/RenterRus/sausage-auth/internal/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *Manager) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	if req == nil || req.Login == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Register: %s", entity.ErrParametrNoFound.Error())
	}

	url, err := t.profile.Registration(ctx, req.GetLogin())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Register.Registration: %s", err.Error())
	}

	return &proto.RegisterResponse{
		Url: url,
	}, nil
}

func (t *Manager) Confirm(ctx context.Context, req *proto.AcceptRequest) (*proto.AcceptResponse, error) {
	if req == nil || req.Login == "" || req.OtpCode == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Confirm: %s", entity.ErrParametrNoFound.Error())
	}

	if err := t.profile.Confirmed(ctx, req.GetLogin(), req.GetOtpCode()); err != nil {
		return &proto.AcceptResponse{
			Status: entity.STATUS_FAILED,
		}, status.Errorf(codes.InvalidArgument, "Confirm.Confirmed: %s", err.Error())
	}

	return &proto.AcceptResponse{
		Status: entity.STATUS_OK,
	}, nil
}

func (t *Manager) UrlOTP(ctx context.Context, req *proto.UrlOTPRequest) (*proto.UrlOTPResponse, error) {
	if req == nil || req.Access == "" {
		return nil, status.Errorf(codes.InvalidArgument, "UrlOTP: %s", entity.ErrParametrNoFound.Error())
	}

	url, err := t.profile.UrlOTP(ctx, req.GetAccess())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "UrlOTP.UrlOTP: %s", err.Error())
	}

	return &proto.UrlOTPResponse{
		Url: url,
	}, nil
}
