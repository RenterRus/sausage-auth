package grpc

import (
	proto "github.com/RenterRus/sausage-auth/docs/proto/v1"
	"github.com/RenterRus/sausage-auth/internal/usecase"
)

type Manager struct {
	profile usecase.Register

	proto.UnimplementedAuthServiceServer
}

func NewManager(register usecase.Register) proto.AuthServiceServer {
	return &Manager{
		profile: register,
	}
}
