package model

import (
	"context"
	"gorm.io/gorm/clause"
	"log"
)

type GormContactDAO struct {
	data *Data
}

func NewGormContactDAO(data *Data) *GormContactDAO {
	return &GormContactDAO{
		data: data,
	}
}
func (g GormContactDAO) GetCountOfFollowing(ctx context.Context, id int64) (cnt int64, err error) {
	db := g.data.Mysql.WithContext(ctx).Table(ContactTable)
	err = db.Where("subject = ?", id).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (g GormContactDAO) GetCountOfFollowers(ctx context.Context, id int64) (cnt int64, err error) {
	db := g.data.Mysql.WithContext(ctx).Table(ContactTable)
	err = db.Where("object = ?", id).Count(&cnt).Error
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (g GormContactDAO) CreateContacts(ctx context.Context, contacts []FollowingContact) error {
	db := g.data.DB(ctx).Table(ContactTable)
	err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&contacts).Error
	if err != nil {
		log.Println("Error creating contacts")
		return err
	}
	return nil
}
