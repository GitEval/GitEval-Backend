package model

import (
	"context"
	"gorm.io/gorm/clause"
	"log"
)

// GormUserDAO 实现了 UserDAO 接口
type GormUserDAO struct {
	data *Data
}

// NewGormUserDAO 构造函数
func NewGormUserDAO(data *Data) *GormUserDAO {
	return &GormUserDAO{
		data: data,
	}
}
func (o *GormUserDAO) CreateUsers(ctx context.Context, users []User) error {
	db := o.data.DB(ctx).Table(UserTable)
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&users).Error
	if err != nil {
		log.Println("Error creating user")
		return err
	}
	return nil
}
func (o *GormUserDAO) GetUserByID(ctx context.Context, id int64) (u User, err error) {
	db := o.data.Mysql.WithContext(ctx).Table(UserTable)
	err = db.Where("id = ?", id).First(&u).Error
	if err != nil {
		log.Println("Error getting user by ID")
		return User{}, err
	}
	return u, nil
}
func (o *GormUserDAO) GetFollowingUsersJoinContact(ctx context.Context, id int64) (users []User, err error) {
	db := o.data.Mysql.WithContext(ctx)
	err = db.Joins("JOIN contacts ON users.id = contacts.subject AND users.id = ?", id).Find(&users).Error
	if err != nil {
		log.Println("Error getting following users")
		return nil, err
	}
	return users, nil
}

func (o *GormUserDAO) GetFollowersUsersJoinContact(ctx context.Context, id int64) (users []User, err error) {
	db := o.data.Mysql.WithContext(ctx)
	err = db.Joins("JOIN contacts ON users.id = contacts.object AND users.id =?", id).Find(&users).Error
	if err != nil {
		log.Println("Error getting followers users")
		return nil, err
	}
	return users, nil
}
