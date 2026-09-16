package grpc

import (
	"context"
	"errors"

	"github.com/AlekSi/pointer"
	proto "github.com/RenterRus/sausage-auth/docs/proto/v1"
	"github.com/RenterRus/sausage-auth/internal/entity"
	"github.com/RenterRus/sausage-auth/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RevokeSession implements authpb.AuthServiceServer.
func (m *Manager) RevokeSession(ctx context.Context, req *proto.RevokeSessionRequest) (*proto.RevokeSessionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "RevokeSession: %s", entity.ErrParametrNoFound.Error())
	}

	switch req.Mode.(type) {
	case *proto.RevokeSessionRequest_Current_:
		if err := m.profile.Logout(ctx, pointer.To(req.GetCurrent().GetHash()), usecase.REVOKE_ONE); err != nil {
			return &proto.RevokeSessionResponse{
				Status: entity.STATUS_FAILED,
			}, status.Errorf(codes.Internal, "RevokeSession.Current: %s", err.Error())
		}
	case *proto.RevokeSessionRequest_All_:
		if err := m.profile.Logout(ctx, pointer.To(req.GetAll().GetLogin()), usecase.REVOKE_ALL); err != nil {
			return &proto.RevokeSessionResponse{
				Status: entity.STATUS_FAILED,
			}, status.Errorf(codes.Internal, "RevokeSession.All: %s", err.Error())
		}
	}

	return &proto.RevokeSessionResponse{
		Status: entity.STATUS_OK,
	}, nil
}

func (m *Manager) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	if req == nil || req.GetAccess() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ValidateToken: %s", entity.ErrParametrNoFound.Error())
	}

	uuid, err := m.profile.Validation(ctx, req.GetAccess())
	st := codes.Internal
	if errors.Is(err, entity.ErrTokenExpired) {
		st = codes.DeadlineExceeded
	}

	if err != nil {
		return nil, status.Errorf(st, "ValidateToken.Validation: %s", err.Error())
	}

	return &proto.ValidateTokenResponse{
		Uuid: *uuid,
	}, nil
}

func (m *Manager) LoginOTP(ctx context.Context, req *proto.LoginOTPRequest) (*proto.LoginOTPResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "LoginOTP: %s", entity.ErrParametrNoFound.Error())
	}

	tokens, err := m.profile.LoginOTP(ctx, usecase.LoginRequest{
		Login:     req.GetLogin(),
		UserAgent: req.GetUserAgent(),
		Code:      req.GetCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LoginOTP.LoginOTP: %s", err.Error())
	}

	return &proto.LoginOTPResponse{
		Tokens: &proto.Tokens{
			Access:  tokens.AccessToken,
			Refresh: tokens.RefreshToken,
		},
		Uuid: tokens.UUID,
	}, nil
}

func (m *Manager) RefreshToken(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "LoginOTP: %s", entity.ErrParametrNoFound.Error())
	}

	tokens, err := m.profile.Refresh(ctx, usecase.RefreshRequest{
		UserAgent:    req.GetUserAgent(),
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LoginOTP.LoginOTP: %s", err.Error())
	}

	return &proto.RefreshResponse{
		Tokens: &proto.Tokens{
			Access:  tokens.AccessToken,
			Refresh: tokens.RefreshToken,
		},
		Uuid: tokens.UUID,
	}, nil
}
