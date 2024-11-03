package model

const (
	DomainTable = "domain"
)

type Domain struct {
	UserID int64  `gorm:"index;column:user_id" json:"user_id"`
	Domain string `gorm:"index;column:domain" json:"domain"`
}

func (d *Domain) TableName() string {
	return DomainTable
}
