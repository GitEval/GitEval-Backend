package request

type GetUserInfo struct {
	UserID int64 `json:"user_id"`
}
type GetRanking struct {
	UserID int64 `json:"user_id"`
}

type GetEvaluation struct {
	UserID int64 `json:"user_id"`
}
