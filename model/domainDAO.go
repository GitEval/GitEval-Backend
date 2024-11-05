package model

import (
	"context"
	"log"
)

type GormDomainDAO struct {
	data *Data
}

func NewGormDomainDAO(d *Data) *GormDomainDAO {
	return &GormDomainDAO{
		data: d,
	}
}

func (o *GormDomainDAO) GetDomainById(ctx context.Context, id int64) ([]string, error) {
	var lang []string
	db := o.data.Mysql.WithContext(ctx).Table(DomainTable)
	err := db.Where("user_id = ?", id).Select("domain").Find(&lang).Error
	if err != nil {
		log.Println("Error getting domain by ID")
		return nil, err
	}
	return lang, nil
}

func (o *GormDomainDAO) Create(ctx context.Context, domain []Domain) error {
	if len(domain) == 0 {
		return nil
	}
	db := o.data.DB(ctx).Table(DomainTable)
	err := db.Create(&domain).Error
	if err != nil {
		log.Println("Error creating domain")
		return err
	}
	return nil
}

func (o *GormDomainDAO) Delete(ctx context.Context, id int64) error {
	db := o.data.DB(ctx).Table(DomainTable)
	err := db.Where("user_id = ?", id).Delete(Domain{}).Error
	if err != nil {
		log.Println("Error deleting domain")
		return err
	}
	return nil
}
