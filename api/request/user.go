package request

type GetUserInfo struct {
	UserID int64 `form:"user_id"`
}

type GetRanking struct {
	UserID int64 `form:"user_id"`
}

type GetEvaluation struct {
	UserID int64 `form:"user_id"`
}

type GetNation struct {
	UserID int64 `form:"user_id"`
}

type GetDomain struct {
	UserID int64 `form:"user_id"`
}

type SearchUser struct {
	Domain   string `form:"domain"`
	Nation   string `form:"nation"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
