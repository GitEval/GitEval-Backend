package request

type SearchUser struct {
	Domain   string `form:"domain"`
	Nation   string `form:"nation"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}
