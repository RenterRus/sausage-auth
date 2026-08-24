package usecase

import (
	"context"
	"fmt"

	"github.com/RenterRus/sausage-profile/internal/entity"
	"github.com/RenterRus/sausage-profile/internal/repo/inmem"
)

func (r *profile) Logout(ctx context.Context, sign *string, revokeType RevokeType) error {
	if sign == nil || *sign == "" {
		return fmt.Errorf("Logout: %w", entity.ErrParametrNoFound)
	}

	switch revokeType {
	case REVOKE_ONE:
		if err := r.usersRepo.SetBlockRefresh(ctx, sign); err != nil {
			return fmt.Errorf("Logout.SetBlockRefresh(one): %w", err)
		}

		if err := r.usersRepo.RemoveRefreshByHash(ctx, sign); err != nil {
			return fmt.Errorf("Logout.RemoveRefreshByHash(one): %w", err)
		}
	case REVOKE_ALL:
		toBlock, err := r.usersRepo.RemoveRefreshByLogin(ctx, sign)
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

	uuid, isValid, err := r.accCache.Get(ctx, accessHash)
	if err == nil && isValid {
		return &uuid, nil
	}

	uuid, err = r.usersRepo.GetUUIDByLogin(ctx, &jwtDetail.UserLogin)
	if err != nil {
		return nil, fmt.Errorf("Validation.GetUUIDByLogin: %w", err)
	}

	r.accCache.Set(ctx, inmem.AccessCacheRequest{
		Access:   accessHash,
		UserUUID: uuid,
	})

	return &uuid, nil
}

func (r *profile) LoginOTP(ctx context.Context, req LoginRequest) (entity.Tokens, error) {
	hash, err := r.usersRepo.OtpHash(ctx, &req.Login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.Hash: %w", err)
	}

	isValid, err := r.otpRepo.ValidateCode(req.Code, hash)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.ValidateCode: %w", err)
	}

	if !isValid {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.isValid: %w", entity.ErrCodeInvalid)
	}

	hashs, err := r.usersRepo.DeleteOldRefresh(ctx, entity.DeleteOldRefreshParams{
		UserLogin: &req.Login,
		UserAgent: &req.UserAgent,
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

	access, err := r.jwtManager.GenAccess(req.Login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenAccess: %w", err)
	}

	refresh, err := r.jwtManager.GenRefresh(req.Login)
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
		UserAgent:   &req.UserAgent,
		Login:       &req.Login,
	}); err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.SetRefreshHash: %w", err)
	}

	if err = r.usersRepo.UpdateLastSighUp(ctx, &req.Login); err != nil {
		fmt.Println("!!! INTO LOG. UpdateLastSighUp:", err)
	}

	uuid, err := r.usersRepo.GetUUIDByLogin(ctx, &req.Login)
	if err == nil {
		r.accCache.Set(ctx, inmem.AccessCacheRequest{
			Access:   hashAccess,
			UserUUID: uuid,
		})
	}

	return entity.Tokens{
		AccessToken:  &hashAccess,
		RefreshToken: &hashRefresh,
	}, nil
}

func (r *profile) Refresh(ctx context.Context, req RefreshRequest) (entity.Tokens, error) {
	oldRef, err := r.usersRepo.GetRefreshToken(ctx, &req.Login, &req.RefreshToken)
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
		UserLogin: &req.Login,
		UserAgent: &req.UserAgent,
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

	access, err := r.jwtManager.GenAccess(req.Login)
	if err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.GenAccess: %w", err)
	}

	refresh, err := r.jwtManager.GenRefresh(req.Login)
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
		UserAgent:   &req.UserAgent,
		Login:       &req.Login,
	}); err != nil {
		return entity.Tokens{}, fmt.Errorf("LoginOTP.SetRefreshHash: %w", err)
	}

	if err = r.usersRepo.UpdateLastSighUp(ctx, &req.Login); err != nil {
		fmt.Println("!!! INTO LOG. UpdateLastSighUp:", err)
	}

	uuid, err := r.usersRepo.GetUUIDByLogin(ctx, &req.Login)
	if err == nil {
		r.accCache.Set(ctx, inmem.AccessCacheRequest{
			Access:   hashAccess,
			UserUUID: uuid,
		})
	}

	return entity.Tokens{
		AccessToken:  &hashAccess,
		RefreshToken: &hashRefresh,
	}, nil
}
