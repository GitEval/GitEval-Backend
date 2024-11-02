package model

import (
	"context"
	"github.com/GitEval/GitEval-Backend/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type contextTxKey struct{}
type Data struct {
	Mysql *gorm.DB
}

func NewData(db *gorm.DB) *Data {
	return &Data{
		Mysql: db,
	}
}
func NewDB(c *conf.DataConfig) *gorm.DB {
	newlogger := logger.New(
		//日志写入文件
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      false,
		},
	)
	db, err := gorm.Open(mysql.Open(c.Addr), &gorm.Config{Logger: newlogger})
	if err != nil {
		panic("connect mysql failed")
	}
	if err := db.AutoMigrate(&User{}, &FollowingContact{}); err != nil {
		panic(err)
	}
	return db
}
func (d *Data) InTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.Mysql.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 将tx放入到ctx中
		ctx = context.WithValue(ctx, contextTxKey{}, tx)
		return fn(ctx)
	})
}

// DB 在事务执行ORM操作的话 得需要使用这个方法获取tx！
func (d *Data) DB(ctx context.Context) *gorm.DB {
	// 从ctx中获取tx
	txKey := ctx.Value(contextTxKey{})
	tx, ok := txKey.(*gorm.DB)
	if ok {
		return tx
	}
	// Notice 如果 !ok 返回错误还是使用默认DB～这个根据实际情况来定！
	// Notice 在Data层使用事务时使用DB()方法是获取不到tx的！此时就应该用 d.MySql
	return d.Mysql
}
