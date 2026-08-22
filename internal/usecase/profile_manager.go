package usecase

import (
	"context"
	"fmt"

	"github.com/RenterRus/sausage-profile/internal/entity"
	"github.com/RenterRus/sausage-profile/internal/repo/psql"
	"github.com/RenterRus/sausage-profile/internal/usecase/hashing"
	"github.com/RenterRus/sausage-profile/internal/usecase/jwt"
	"github.com/RenterRus/sausage-profile/internal/usecase/otp"
)

type register struct {
	otpRepo    otp.OTP
	jwtHashing hashing.Hashing
	usersRepo  psql.UsersRepo
	jwtManager jwt.JWT
}

func NewRegisterManager(otpRepo otp.OTP, hash hashing.Hashing, jwtManager jwt.JWT, usersRepo psql.UsersRepo) Register {
	return &register{
		otpRepo:    otpRepo,
		jwtHashing: hash,
		usersRepo:  usersRepo,
		jwtManager: jwtManager,
	}
}

func (r *register) Registration(ctx context.Context, login string) (string, error) {
	userLogin, err := r.usersRepo.IsExist(ctx, &login)
	if err != nil {
		return "", fmt.Errorf("Registration.IsExists: %w", err)
	}

	if userLogin {
		return "", fmt.Errorf("Registration.IsExists(exists): %w", entity.ErrAlreadyExists)
	}

	hash, url, err := r.otpRepo.GenerateHash(login)
	if err != nil {
		return "", nil
	}

	if err := r.usersRepo.Register(ctx, entity.RegisterParams{
		Login: &login,
		Hash:  &hash,
		Link:  &url,
	}); err != nil {
		return "", fmt.Errorf("Registration.Register: %w", err)
	}

	return url, nil
}

func (r *register) Confirmed(ctx context.Context, login, code string) error {
	hash, err := r.usersRepo.Hash(ctx, &login)
	if err != nil {
		return fmt.Errorf("Confirmed.Hash: %w", err)
	}

	isValid, err := r.otpRepo.ValidateCode(code, hash)
	if err != nil {
		return fmt.Errorf("Confirmed.ValidateCode: %w", err)
	}

	if !isValid {
		return fmt.Errorf("Confirmed.ValidateCode(invalid): %w", entity.ErrCodeInvalid)
	}

	if err = r.usersRepo.Confirmed(ctx, &login); err != nil {
		return fmt.Errorf("Confirmed.Confirmed: %w", err)
	}

	return nil

}
