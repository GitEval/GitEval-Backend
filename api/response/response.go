package response

import "github.com/GitEval/GitEval-Backend/model"

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

type User struct {
	U      model.User `json:"user"`
	Domain []string   `json:"domain"`
}
