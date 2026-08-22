package entity

type RegisterParams struct {
	Login *string
	Hash  *string
	Link  *string
}

type GetRefreshTokenRow struct {
	RefreshHash string
	IsExpired   bool
	Block       bool
	UserAgent   *string
}

type SetRefreshHashParams struct {
	RefreshHash *string
	UserAgent   *string
	Login       *string
}
