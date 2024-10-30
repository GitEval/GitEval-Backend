package model

import "gorm.io/gorm"

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
