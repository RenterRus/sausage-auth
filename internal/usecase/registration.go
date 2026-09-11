package usecase

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/sausage-auth/internal/entity"
	"github.com/RenterRus/sausage-auth/internal/repo/inmem"
)

func otpCacheKey(uuid string) string {
	return fmt.Sprintf("URL_%s", uuid)
}

func (r *profile) UrlOTP(ctx context.Context, accessHash string) (string, error) {
	uuid, err := r.Validation(ctx, accessHash)
	if err != nil {
		return "", fmt.Errorf("UrlOTP.Validation: %w", err)
	}

	if pointer.Get(uuid) == "" {
		return "", fmt.Errorf("UrlOTP.Login: %w", entity.ErrAuthFailed)
	}

	url, isValid, err := r.accCache.Get(ctx, otpCacheKey(*uuid))
	if err == nil && isValid {
		return url, nil
	}

	login, err := r.usersRepo.LoginByUUID(ctx, uuid)
	if pointer.Get(uuid) == "" {
		return "", fmt.Errorf("UrlOTP.LoginByUUID: %w", err)
	}

	url, err = r.otpRepo.GenerateUrl(login)
	if err != nil {
		return "", fmt.Errorf("UrlOTP.GenerateURL: %w", err)
	}

	r.accCache.Set(ctx, inmem.AccessCacheRequest{
		Key:   otpCacheKey(*uuid),
		Value: url,
	})

	return url, err
}

func (r *profile) Registration(ctx context.Context, login string) (string, error) {
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
	}); err != nil {
		return "", fmt.Errorf("Registration.Register: %w", err)
	}

	defer func() {
		uuid, err := r.usersRepo.GetUUIDByLogin(ctx, &login)
		if err != nil {
			return
		}

		r.accCache.Set(ctx, inmem.AccessCacheRequest{
			Key:   otpCacheKey(uuid),
			Value: url,
		})
	}()

	return url, nil
}

func (r *profile) Confirmed(ctx context.Context, login, code string) error {
	hash, err := r.usersRepo.OtpHash(ctx, &login)
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
