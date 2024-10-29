package model

type User struct {
	LoginName         string  `json:"login_name"`          //用户的登录名
	Name              *string `json:"name"`                //真实姓名
	Location          *string `json:"location"`            //地区
	Email             string  `json:"email"`               //邮箱
	Following         int     `json:"following"`           //关注数
	Followers         int     `json:"followers"`           //粉丝数
	Blog              *string `json:"blog"`                //博客连接
	Bio               *string `json:"Bio"`                 //用户的个人简介
	PublicRepos       int     `json:"public_repos"`        //用户公开的仓库的数量
	TotalPrivateRepos int     `json:"total_private_repos"` //用户的私有仓库总数
	Company           *string `json:"company"`             //用户所属的公司
	AvatarURL         string  `json:"avatar_url"`          //用户头像的 URL
	Collaborators     int     `json:"collaborators"`       //协作者的数量
	Nationality       *string `json:"nationality"`         //国籍
	Score             float64 `json:"score"`               //评分
}

type FollowingContact struct {
	//subject 关注 object
	Subject string `json:"subject"` //主体
	Object  string `json:"object"`  //被关注的客体
}
