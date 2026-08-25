package usecase

import (
	"github.com/RenterRus/sausage-auth/internal/repo/inmem"
	"github.com/RenterRus/sausage-auth/internal/repo/psql"
	"github.com/RenterRus/sausage-auth/internal/usecase/hashing"
	"github.com/RenterRus/sausage-auth/internal/usecase/jwt"
	"github.com/RenterRus/sausage-auth/internal/usecase/otp"
)

type profile struct {
	otpRepo    otp.OTP
	jwtHashing hashing.Hashing
	usersRepo  psql.UsersRepo
	jwtManager jwt.JWT

	accCache inmem.AccessCache
}

type ProfileConf struct {
	OtpRepo    otp.OTP
	Hash       hashing.Hashing
	JwtManager jwt.JWT
	UsersRepo  psql.UsersRepo

	AccCache inmem.AccessCache
}

func NewProfileManager(conf ProfileConf) Register {
	return &profile{
		otpRepo:    conf.OtpRepo,
		jwtHashing: conf.Hash,
		usersRepo:  conf.UsersRepo,
		jwtManager: conf.JwtManager,
		accCache:   conf.AccCache,
	}
}
