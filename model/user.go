package model

import (
	"gorm.io/gorm"
	"time"
)

// User 模型
type User struct {
	ID                      int64           `json:"id,omitempty"`
	Login                   string          `json:"login,omitempty"`
	NodeID                  string          `json:"node_id,omitempty"`
	AvatarURL               string          `json:"avatar_url,omitempty"`
	HTMLURL                 string          `json:"html_url,omitempty"`
	GravatarID              string          `json:"gravatar_id,omitempty"`
	Name                    string          `json:"name,omitempty"`
	Company                 string          `json:"company,omitempty"`
	Blog                    string          `json:"blog,omitempty"`
	Location                string          `json:"location,omitempty"`
	Email                   string          `json:"email,omitempty"`
	Hireable                bool            `json:"hireable,omitempty"`
	Bio                     string          `json:"bio,omitempty"`
	TwitterUsername         string          `json:"twitter_username,omitempty"`
	PublicRepos             int             `json:"public_repos,omitempty"`
	PublicGists             int             `json:"public_gists,omitempty"`
	Followers               int             `json:"followers,omitempty"`
	Following               int             `json:"following,omitempty"`
	CreatedAt               time.Time       `json:"created_at,omitempty"`
	UpdatedAt               time.Time       `json:"updated_at,omitempty"`
	SuspendedAt             time.Time       `json:"suspended_at,omitempty"`
	Type                    string          `json:"type,omitempty"`
	SiteAdmin               bool            `json:"site_admin,omitempty"`
	TotalPrivateRepos       int             `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos       int             `json:"owned_private_repos,omitempty"`
	PrivateGists            int             `json:"private_gists,omitempty"`
	DiskUsage               int             `json:"disk_usage,omitempty"`
	Collaborators           int             `json:"collaborators,omitempty"`
	TwoFactorAuthentication bool            `json:"two_factor_authentication,omitempty"`
	LdapDn                  string          `json:"ldap_dn,omitempty"`
	URL                     string          `json:"url,omitempty"`
	EventsURL               string          `json:"events_url,omitempty"`
	FollowingURL            string          `json:"following_url,omitempty"`
	FollowersURL            string          `json:"followers_url,omitempty"`
	GistsURL                string          `json:"gists_url,omitempty"`
	OrganizationsURL        string          `json:"organizations_url,omitempty"`
	ReceivedEventsURL       string          `json:"received_events_url,omitempty"`
	ReposURL                string          `json:"repos_url,omitempty"`
	StarredURL              string          `json:"starred_url,omitempty"`
	SubscriptionsURL        string          `json:"subscriptions_url,omitempty"`
	Permissions             map[string]bool `json:"permissions,omitempty"`
	RoleName                string          `json:"role_name,omitempty"`
}

// UserDAO 接口
type UserDAO interface {
	CreateUser(user *User) error
	GetUserByID(id int64) (*User, error)
	UpdateUser(user *User) error
	DeleteUser(id int64) error
}

// GormUserDAO 实现了 UserDAO 接口
type GormUserDAO struct {
	db *gorm.DB
}

// NewGormUserDAO 构造函数
func NewGormUserDAO(db *gorm.DB) *GormUserDAO {
	return &GormUserDAO{db: db}
}

// CreateUser 创建用户
func (r *GormUserDAO) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

// GetUserByID 通过 ID 获取用户
func (r *GormUserDAO) GetUserByID(id int64) (*User, error) {
	var user User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (r *GormUserDAO) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

// DeleteUser 删除用户
func (r *GormUserDAO) DeleteUser(id int64) error {
	return r.db.Delete(&User{}, id).Error
}
