package response

type Success struct {
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}
type Err struct {
	Err error `json:"error"`
}
type CallBack struct {
	UserId int64 `json:"user_id"`
}
