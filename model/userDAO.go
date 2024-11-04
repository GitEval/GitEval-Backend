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
	err = db.Select("DISTINCT users.*").
		Joins("JOIN contacts ON contacts.object = users.id").
		Where("contacts.subject = ?", id).
		Find(&users).Error
	if err != nil {
		log.Println("Error getting followers users")
		return nil, err
	}
	return users, nil
}

func (o *GormUserDAO) GetFollowersUsersJoinContact(ctx context.Context, id int64) (users []User, err error) {
	db := o.data.Mysql.WithContext(ctx)
	err = db.Select("DISTINCT users.*").
		Joins("JOIN contacts ON contacts.subject = users.id").
		Where("contacts.object = ?", id).
		Find(&users).Error
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

func (o *GormUserDAO) SearchUser(ctx context.Context, nation, domain string, page int, pageSize int) (users []User, err error) {
	db := o.data.Mysql.WithContext(ctx)

	// 这里用 LIKE 来查找 nation 字段在 '|' 之前与传入的 nation 相匹配,domain也使用相同的方法进行模糊匹配
	err = db.Select("DISTINCT users.*").
		Joins("JOIN contacts ON contacts.user_id = users.id").
		Where("SUBSTRING_INDEX(contacts.domain, '|', 1) = ?", domain).
		Where("SUBSTRING_INDEX(users.nation, '|', 1) = ?", nation).
		Order("users.score DESC"). // 按照 score 字段从高到低排序
		Offset((page - 1) * pageSize). // 分页
		Limit(pageSize). // 设置每页大小
		Find(&users).Error

	if err != nil {
		log.Println("Error getting followers users:", err)
		return nil, err
	}

	return users, nil
}
