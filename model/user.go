package model

const (
	UserTable    = "users"
	ContactTable = "contacts"
)

// User 模型
type User struct {
	ID                int64   `gorm:"column:id;primaryKey" `
	LoginName         string  `gorm:"column:login_name" json:"login_name"`                   //用户的登录名
	Name              *string `gorm:"column:name" json:"name"`                               //真实姓名
	Location          *string `gorm:"column:location" json:"location"`                       //地区
	Email             string  `gorm:"column:email" json:"email"`                             //邮箱
	Following         int     `gorm:"column:following" json:"following"`                     //关注数
	Followers         int     `gorm:"column:followers" json:"followers"`                     //粉丝数
	Blog              *string `gorm:"column:blog" json:"blog"`                               //博客连接
	Bio               *string `gorm:"column:bio" json:"Bio"`                                 //用户的个人简介
	PublicRepos       int     `gorm:"column:public_repos" json:"public_repos"`               //用户公开的仓库的数量
	TotalPrivateRepos int     `gorm:"column:total_private_repos" json:"total_private_repos"` //用户的私有仓库总数
	Company           *string `gorm:"column:company" json:"company"`                         //用户所属的公司
	AvatarURL         string  `gorm:"column:avatar_url" json:"avatar_url"`                   //用户头像的 URL
	Collaborators     int     `gorm:"column:collaborators" json:"collaborators"`             //协作者的数量
	Nationality       *string `gorm:"column:nationality" json:"nationality"`                 //国籍
	Score             float64 `gorm:"column:score" json:"score"`                             //评分
}

type FollowingContact struct {
	//subject 关注 object
	Subject int64 `gorm:"column:subject;index:idx_contact" json:"subject"` //主体
	Object  int64 `gorm:"column:object;index:idx_contact" json:"object"`   //被关注的客体
}

func (u *User) TableName() string {
	return UserTable
}
func (c *FollowingContact) TableName() string {
	return ContactTable
}
