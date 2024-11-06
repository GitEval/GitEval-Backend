package request

type SearchUser struct {
	Domain   string  `form:"domain"`
	Nation   *string `form:"nation,omitempty"`
	Page     int     `form:"page"`
	PageSize int     `form:"page_size"`
}

type GetUserInfo struct {
	UserId int64 `form:"user_id"`
}
