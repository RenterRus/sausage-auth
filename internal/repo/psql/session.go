package psql

import (
	"context"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/RenterRus/sausage-auth/internal/entity"
	"github.com/RenterRus/sausage-auth/internal/repo/psql/db"
	"github.com/samber/lo"
)

func (u *UserRepo) GetRefreshToken(ctx context.Context, req GetRefreshReq) (entity.GetRefreshTokenRow, error) {
	if req.Login == nil || *req.Login == "" {
		return entity.GetRefreshTokenRow{}, fmt.Errorf("GetRefreshToken: %w", entity.ErrParametrNoFound)
	}

	resp, err := u.Queries.GetRefreshToken(ctx, db.GetRefreshTokenParams{
		UserLogin: req.Login,
		Hash:      req.Hash,
		UserAgent: req.UserAgent,
	})
	if err != nil {
		return entity.GetRefreshTokenRow{}, fmt.Errorf("GetRefreshToken.GetRefreshToken: %w", err)
	}

	return entity.GetRefreshTokenRow{
		RefreshHash: resp.RefreshHash,
		IsExpired:   resp.IsExpired,
		Block:       resp.Block,
		UserAgent:   resp.UserAgent,
		Login:       resp.UserLogin,
	}, nil
}

func (u *UserRepo) SetBlockRefresh(ctx context.Context, refreshHash *string) error {
	if refreshHash == nil || *refreshHash == "" {
		return fmt.Errorf("SetBlockRefresh: %w", entity.ErrParametrNoFound)
	}

	if err := u.Queries.SetBlockRefresh(ctx, refreshHash); err != nil {
		return fmt.Errorf("SetBlockRefresh.SetBlockRefresh: %w", err)
	}

	return nil
}

func (u *UserRepo) SetRefreshHash(ctx context.Context, arg entity.SetRefreshHashParams) error {
	if err := u.Queries.SetRefreshHash(ctx, db.SetRefreshHashParams{
		RefreshHash: arg.RefreshHash,
		UserLogin:   arg.Login,
		UserAgent:   arg.UserAgent,
	}); err != nil {
		return fmt.Errorf("SetRefreshHash.SetRefreshHash: %w", err)
	}

	return nil
}

func (u *UserRepo) GetUUIDByLogin(ctx context.Context, login *string) (string, error) {
	if login == nil || *login == "" {
		return "", fmt.Errorf("GetUUIDByLogin: %w", entity.ErrParametrNoFound)
	}

	uuid, err := u.Queries.GetUUIDByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("GetUUIDByLogin.GetUUIDByLogin: %w", err)
	}

	return uuid, nil
}

func (u *UserRepo) RemoveRefreshByHash(ctx context.Context, refreshHash *string) error {
	if refreshHash == nil || *refreshHash == "" {
		return fmt.Errorf("RemoveRefreshByHash: %w", entity.ErrParametrNoFound)
	}

	if err := u.Queries.RemoveRefreshByHash(ctx, refreshHash); err != nil {
		return fmt.Errorf("RemoveRefreshByHash.RemoveRefreshByHash: %w", err)
	}

	return nil
}

func (u *UserRepo) RemoveRefreshByLogin(ctx context.Context, userLogin *string) ([]string, error) {
	if userLogin == nil || *userLogin == "" {
		return nil, fmt.Errorf("RemoveRefreshByLogin: %w", entity.ErrParametrNoFound)
	}

	refreshs, err := u.Queries.RemoveRefreshByLogin(ctx, userLogin)
	if err != nil {
		return nil, fmt.Errorf("RemoveRefreshByLogin.RemoveRefreshByLogin: %w", err)
	}

	return refreshs, nil
}

func (u *UserRepo) RemoveOldRefreshByLoginUA(ctx context.Context, arg entity.RemoveOldRefreshByLoginUA) ([]entity.OldRefreshResponse, error) {
	if arg.UserAgent == "" || arg.Login == "" {
		return nil, fmt.Errorf("DeleteOldRefresh(validation): %w", entity.ErrParametrNoFound)
	}

	hashs, err := u.Queries.RemoveOldRefreshByLoginUA(ctx, db.RemoveOldRefreshByLoginUAParams{
		UserLogin: &arg.Login,
		UserAgent: &arg.UserAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("DeleteOldRefresh: %w", err)
	}

	return lo.Map(hashs, func(item db.RemoveOldRefreshByLoginUARow, _ int) entity.OldRefreshResponse {
		return entity.OldRefreshResponse{
			Refresh: item.RefreshHash,
			Login:   item.UserLogin,
		}
	}), nil
}

func (u *UserRepo) DeleteOldRefresh(ctx context.Context, arg entity.DeleteOldRefreshParams) ([]entity.OldRefreshResponse, error) {
	if arg.UserAgent == "" || arg.Refresh == "" {
		return nil, fmt.Errorf("DeleteOldRefresh(validation): %w", entity.ErrParametrNoFound)
	}

	hashs, err := u.Queries.RemoveOldRefresh(ctx, db.RemoveOldRefreshParams{
		RefreshHash: &arg.Refresh,
		UserAgent:   &arg.UserAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("DeleteOldRefresh: %w", err)
	}

	return lo.Map(hashs, func(item db.RemoveOldRefreshRow, _ int) entity.OldRefreshResponse {
		return entity.OldRefreshResponse{
			Refresh: item.RefreshHash,
			Login:   item.UserLogin,
		}
	}), nil
}

func (u *UserRepo) UpdateLastSighUp(ctx context.Context, userLogin *string) error {
	if pointer.Get(userLogin) == "" {
		return fmt.Errorf("UpdateLastSighUp(validation): %w", entity.ErrParametrNoFound)
	}

	if err := u.Queries.UpdateLastSighUp(ctx, userLogin); err != nil {
		return fmt.Errorf("UpdateLastSighUp: %w", err)
	}

	return nil
}
