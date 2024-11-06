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
	if len(users) == 0 {
		return nil
	}

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

func (o *GormUserDAO) SearchUser(ctx context.Context, nation *string, domain string, page int, pageSize int) (users []User, err error) {
	db := o.data.Mysql.WithContext(ctx)

	// 基础查询：联合查询用户和域名，按域名和分数排序,这里使用|进行划分
	query := db.Select("DISTINCT users.*").
		Joins("JOIN domain ON domain.user_id = users.id").
		Where("SUBSTRING_INDEX(domain.domain, '|', 1) = ?", domain).
		Order("users.score DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	// 如果 nation 不为 nil，添加国家筛选条件
	if nation != nil {
		query = query.Where("SUBSTRING_INDEX(users.nationality, '|', 1) = ?", *nation)
	}

	// 执行查询
	err = query.Find(&users).Error
	if err != nil {
		log.Println("Error getting followers users:", err)
		return nil, err
	}

	return users, nil
}
