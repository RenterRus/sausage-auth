package usecase

import (
	"github.com/RenterRus/sausage-profile/internal/repo/psql"
	"github.com/RenterRus/sausage-profile/internal/usecase/hashing"
	"github.com/RenterRus/sausage-profile/internal/usecase/jwt"
	"github.com/RenterRus/sausage-profile/internal/usecase/otp"
)

type profile struct {
	otpRepo    otp.OTP
	jwtHashing hashing.Hashing
	usersRepo  psql.UsersRepo
	jwtManager jwt.JWT
}

type ProfileConf struct {
	OtpRepo    otp.OTP
	Hash       hashing.Hashing
	JwtManager jwt.JWT
	UsersRepo  psql.UsersRepo
}

func NewProfileManager(conf ProfileConf) Register {
	return &profile{
		otpRepo:    conf.OtpRepo,
		jwtHashing: conf.Hash,
		usersRepo:  conf.UsersRepo,
		jwtManager: conf.JwtManager,
	}
}
