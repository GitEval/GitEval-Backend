package model

import (
	"fmt"
	"github.com/google/go-github/v50/github"
	"gorm.io/gorm"
)

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
	Nationality       string  `gorm:"column:nationality" json:"nationality"`                 //国籍
	Score             float64 `gorm:"column:score;index" json:"score"`                       //评分
	Evaluation        string  `gorm:"column:evaluation" json:"evaluation"`                   //评估
}

type FollowingContact struct {
	//subject 关注 object
	ID      string `gorm:"column:id;primaryKey" json:"id"`
	Subject int64  `gorm:"column:subject;index:idx_contact" json:"subject"` //主体
	Object  int64  `gorm:"column:object;index:idx_contact" json:"object"`   //被关注的客体
}

type Leaderboard struct {
	UserID    int64   `json:"user_id"`
	UserName  string  `json:"user_name"`
	AvatarURL string  `json:"avatar_url"`
	Score     float64 `json:"score"`
}

func (u *User) TableName() string {
	return UserTable
}
func (c *FollowingContact) TableName() string {
	return ContactTable
}
func (c *FollowingContact) GenerateID() {
	c.ID = fmt.Sprintf("%d->%d", c.Subject, c.Object)
}
func (c *FollowingContact) BeforeCreate(tx *gorm.DB) (err error) {
	c.GenerateID()
	return nil
}

func TransformUser(userInfo *github.User) User {
	return User{
		ID:                userInfo.GetID(),
		AvatarURL:         userInfo.GetAvatarURL(),
		Name:              userInfo.Name,
		Company:           userInfo.Company,
		Blog:              userInfo.Blog,
		Location:          userInfo.Location,
		Email:             userInfo.GetEmail(),
		Bio:               userInfo.Bio,
		PublicRepos:       userInfo.GetPublicRepos(),
		Followers:         userInfo.GetFollowers(),
		Following:         userInfo.GetFollowing(),
		TotalPrivateRepos: userInfo.GetTotalPrivateRepos(),
		Collaborators:     userInfo.GetCollaborators(),
	}
}
func TransformUsers(users []*github.User) []User {
	var transformedUsers []User
	for _, user := range users {
		transformedUsers = append(transformedUsers, TransformUser(user))
	}
	return transformedUsers
}
