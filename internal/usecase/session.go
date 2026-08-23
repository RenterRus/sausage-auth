package usecase

import (
	"context"
	"fmt"

	"github.com/RenterRus/sausage-profile/internal/entity"
)

func (r *profile) Logout(ctx context.Context, login, hash *string, revokeType RevokeType) error {
	switch revokeType {
	case REVOKE_ONE:
		if hash == nil || *hash == "" {
			return fmt.Errorf("Logout(one): %w", entity.ErrParametrNoFound)
		}

		if err := r.usersRepo.SetBlockRefresh(ctx, hash); err != nil {
			return fmt.Errorf("Logout.SetBlockRefresh(one): %w", err)
		}

		if err := r.usersRepo.RemoveRefreshByHash(ctx, hash); err != nil {
			return fmt.Errorf("Logout.RemoveRefreshByHash(one): %w", err)
		}
	case REVOKE_ALL:
		if login == nil || *login == "" {
			return fmt.Errorf("Logout(all): %w", entity.ErrParametrNoFound)
		}

		toBlock, err := r.usersRepo.RemoveRefreshByLogin(ctx, login)
		if err != nil {
			return fmt.Errorf("Logout.RemoveRefreshByHash(all): %w", err)
		}

		for i := range toBlock {
			if err := r.usersRepo.SetBlockRefresh(ctx, &toBlock[i]); err != nil {
				return fmt.Errorf("Logout.SetBlockRefresh(all): %w", err)
			}
		}
	}

	return nil
}

// ctx, hash -> userI, err
func (r *profile) Validation(ctx context.Context, accessHash string) (*string, error) {
	acc, err := r.jwtHashing.Decrypt(accessHash)
	if err != nil {
		return nil, fmt.Errorf("Validation.Decrypt: %w", err)
	}

	jwtDetail, err := r.jwtManager.Parse(acc)
	if err != nil {
		return nil, fmt.Errorf("Validation.Parse: %w", err)
	}

	if !jwtDetail.IsValid {
		return nil, fmt.Errorf("Validation.Parse(IsExpired): %w", entity.ErrTokenInvalid)
	}

	if jwtDetail.IsExpired {
		return nil, fmt.Errorf("Validation.Parse(IsExpired): %w", entity.ErrTokenExpired)
	}

	uuid, err := r.usersRepo.GetUUIDByLogin(ctx, &jwtDetail.UserLogin)
	if err != nil {
		return nil, fmt.Errorf("Validation.GetUUIDByLogin: %w", err)
	}

	return &uuid, nil
}

func (r *profile) LoginOTP(ctx context.Context, login, userAgent, code string) (entity.Tokens, error) {
	hash, err := r.usersRepo.OtpHash(ctx, &login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Hash: %w", err)
	}

	isValid, err := r.otpRepo.ValidateCode(code, hash)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.ValidateCode: %w", err)
	}

	if !isValid {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.isValid: %w", entity.ErrCodeInvalid)
	}

	hashs, err := r.usersRepo.DeleteOldRefresh(ctx, entity.DeleteOldRefreshParams{
		UserLogin: &login,
		UserAgent: &userAgent,
	})
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.DeleteOldRefresh: %w", err)
	}

	if len(hashs) > 0 {
		for i := range hashs {
			if err := r.usersRepo.SetBlockRefresh(ctx, &hashs[i]); err != nil {
				return entity.Tokens{}, fmt.Errorf("LoginOTP.SetBlockRefresh: %w", err)
			}
		}
	}

	access, err := r.jwtManager.GenAccess(login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenAccess: %w", err)
	}

	refresh, err := r.jwtManager.GenRefresh(login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenRefresh: %w", err)
	}

	hashRefresh, err := r.jwtHashing.Encrypt(refresh)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Encrypt(refresh): %w", err)
	}

	hashAccess, err := r.jwtHashing.Encrypt(access)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Encrypt(access): %w", err)
	}

	if err = r.usersRepo.SetRefreshHash(ctx, entity.SetRefreshHashParams{
		RefreshHash: &hashRefresh,
		UserAgent:   &userAgent,
		Login:       &login,
	}); err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.SetRefreshHash: %w", err)
	}

	if err = r.usersRepo.UpdateLastSighUp(ctx, &login); err != nil {
		fmt.Println("!!! INTO LOG. UpdateLastSighUp:", err)
	}

	return entity.Tokens{
		AccessToken:  &hashAccess,
		RefreshToken: &hashRefresh,
	}, nil
}

func (r *profile) Refresh(ctx context.Context, login, userAgent, refreshHash string) (entity.Tokens, error) {
	oldRef, err := r.usersRepo.GetRefreshToken(ctx, &login, &refreshHash)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("Refresh.GetRefreshToken: %w", err)
	}

	if oldRef.IsExpired {
		return entity.Tokens{}, fmt.Errorf("Refresh.Valid(exp): %w", entity.ErrTokenExpired)
	}

	if oldRef.Block {
		return entity.Tokens{}, fmt.Errorf("Refresh.Valid(block): %w", entity.ErrTokenBlocked)
	}

	hashs, err := r.usersRepo.DeleteOldRefresh(ctx, entity.DeleteOldRefreshParams{
		UserLogin: &login,
		UserAgent: &userAgent,
	})
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.DeleteOldRefresh: %w", err)
	}

	if len(hashs) > 0 {
		for i := range hashs {
			if err := r.usersRepo.SetBlockRefresh(ctx, &hashs[i]); err != nil {
				return entity.Tokens{}, fmt.Errorf("LoginOTP.SetBlockRefresh: %w", err)
			}
		}
	}

	access, err := r.jwtManager.GenAccess(login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenAccess: %w", err)
	}

	refresh, err := r.jwtManager.GenRefresh(login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenRefresh: %w", err)
	}

	hashRefresh, err := r.jwtHashing.Encrypt(refresh)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Encrypt(refresh): %w", err)
	}

	hashAccess, err := r.jwtHashing.Encrypt(access)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Encrypt(access): %w", err)
	}

	if err = r.usersRepo.SetRefreshHash(ctx, entity.SetRefreshHashParams{
		RefreshHash: &hashRefresh,
		UserAgent:   &userAgent,
		Login:       &login,
	}); err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.SetRefreshHash: %w", err)
	}

	if err = r.usersRepo.UpdateLastSighUp(ctx, &login); err != nil {
		fmt.Println("!!! INTO LOG. UpdateLastSighUp:", err)
	}

	return entity.Tokens{
		AccessToken:  &hashAccess,
		RefreshToken: &hashRefresh,
	}, nil
}
