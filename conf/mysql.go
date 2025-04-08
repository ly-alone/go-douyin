package conf

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDB 返回数据库连接实例
func GetDB() *gorm.DB {
	return dbInstance
}

func MysqlDb() {
	dsn := "root:2LCqvSOJ6m0Utgg6@tcp(localhost:3307)/ly?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	once.Do(func() {
		dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		panic(err)
	}
}
