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
	Token string `json:"token"`
}

type User struct {
	U      model.User `json:"user"`
	Domain []string   `json:"domain"`
}

type Ranking struct {
	Leaderboard []model.Leaderboard `json:"leaderboard"`
}

type EvaluationResp struct {
	Evaluation string `json:"evaluation"`
}

type NationResp struct {
	Nation string `json:"nation"`
}

type DomainResp struct {
	Domain []string `json:"domain"`
}
type SearchResp struct {
	Users []model.User `json:"users"`
}
