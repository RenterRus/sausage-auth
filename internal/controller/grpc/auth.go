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
