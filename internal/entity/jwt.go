package entity

type BaseJWT struct {
	UserLogin string
	IsExpired bool
	IsValid   bool
}

type DeleteOldRefreshParams struct {
	UserLogin *string
	UserAgent *string
}
