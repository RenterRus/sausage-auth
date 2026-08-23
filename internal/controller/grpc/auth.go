package grpc

import (
	"context"

	"github.com/AlekSi/pointer"
	proto "github.com/RenterRus/sausage-profile/docs/proto/v1"
	"github.com/RenterRus/sausage-profile/internal/entity"
	"github.com/RenterRus/sausage-profile/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RevokeSession implements authpb.AuthServiceServer.
func (m *Manager) RevokeSession(ctx context.Context, req *proto.RevokeSessionRequest) (*proto.RevokeSessionResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "RevokeSession: %w", entity.ErrParametrNoFound)
	}

	switch req.Mode.(type) {
	case *proto.RevokeSessionRequest_Current_:
		if err := m.profile.Logout(ctx, pointer.To(req.GetCurrent().GetHash()), usecase.REVOKE_ONE); err != nil {
			return &proto.RevokeSessionResponse{
				Status: entity.STATUS_FAILED,
			}, status.Errorf(codes.Internal, "RevokeSession.Current: %w", err)
		}
	case *proto.RevokeSessionRequest_All_:
		if err := m.profile.Logout(ctx, pointer.To(req.GetAll().GetLogin()), usecase.REVOKE_ALL); err != nil {
			return &proto.RevokeSessionResponse{
				Status: entity.STATUS_FAILED,
			}, status.Errorf(codes.Internal, "RevokeSession.All: %w", err)
		}
	}

	return &proto.RevokeSessionResponse{
		Status: entity.STATUS_OK,
	}, nil
}

func (m *Manager) ValidateToken(ctx context.Context, req *proto.ValidateTokenRequest) (*proto.ValidateTokenResponse, error) {
	if req == nil || req.GetAccess() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "ValidateToken: %w", entity.ErrParametrNoFound)
	}

	uuid, err := m.profile.Validation(ctx, req.GetAccess())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ValidateToken.Validation: %w", err)
	}

	return &proto.ValidateTokenResponse{
		Uuid: *uuid,
	}, nil
}

func (m *Manager) LoginOTP(ctx context.Context, req *proto.LoginOTPRequest) (*proto.LoginOTPResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "LoginOTP: %w", entity.ErrParametrNoFound)
	}

	tokens, err := m.profile.LoginOTP(ctx, usecase.LoginRequest{
		Login:     req.GetLogin(),
		UserAgent: req.GetUserAgent(),
		Code:      req.GetCode(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LoginOTP.LoginOTP: %w", err)
	}

	return &proto.LoginOTPResponse{
		Tokens: &proto.Tokens{
			Access:  *tokens.AccessToken,
			Refresh: *tokens.RefreshToken,
		},
	}, nil
}

func (m *Manager) RefreshToken(ctx context.Context, req *proto.RefreshRequest) (*proto.RefreshResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "LoginOTP: %w", entity.ErrParametrNoFound)
	}

	tokens, err := m.profile.Refresh(ctx, usecase.RefreshRequest{
		Login:        req.GetLogin(),
		UserAgent:    req.GetUserAgent(),
		RefreshToken: req.GetRefreshToken(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "LoginOTP.LoginOTP: %w", err)
	}

	return &proto.RefreshResponse{
		Tokens: &proto.Tokens{
			Access:  *tokens.AccessToken,
			Refresh: *tokens.RefreshToken,
		},
	}, nil
}
