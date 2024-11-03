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

// CreateUsers 这里更新的数据不包括国籍和评价
func (o *GormUserDAO) CreateUsers(ctx context.Context, users []User) error {
	db := o.data.DB(ctx).Table(UserTable)

	// 定义要更新的字段（除 Nationality 和 Evaluation 外的所有字段）,有点弱智但是刚刚好
	updateFields := []string{
		"login_name", "name", "location", "email", "following", "followers",
		"blog", "bio", "public_repos", "total_private_repos", "company",
		"avatar_url", "collaborators", "score",
	}

	// 设置冲突时更新指定字段
	err := db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns(updateFields),
	}).Create(&users).Error

	// 错误处理
	if err != nil {
		log.Println("Error creating users with selective update")
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

func (o *GormUserDAO) SaveUser(ctx context.Context, user User) error {
	db := o.data.DB(ctx).Table(UserTable)
	err := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error
	if err != nil {
		log.Println("Error saving user")
		return err
	}
	return nil
}
