package entity

type BaseJWT struct {
	UserLogin string
	IsExpired bool
	IsValid   bool
}

type DeleteOldRefreshParams struct {
	Refresh   string
	UserAgent string
}

type RemoveOldRefreshByLoginUA struct {
	Login     string
	UserAgent string
}

type OldRefreshResponse struct {
	Refresh string
	Login   string
}
