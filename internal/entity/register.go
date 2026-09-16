package entity

type RegisterParams struct {
	Login string
	Hash  string
	URL   string
}

type GetRefreshTokenRow struct {
	RefreshHash string
	IsExpired   bool
	Block       bool
	UserAgent   *string
	Login       string
}

type SetRefreshHashParams struct {
	RefreshHash *string
	UserAgent   *string
	Login       *string
}
