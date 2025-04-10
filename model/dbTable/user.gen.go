// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbTable

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <user>
type User struct {
	ID    int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name  string    `gorm:"column:name;not null;comment:账号" json:"name"`        // 账号
	Pwd   string    `gorm:"column:pwd;not null;comment:密码" json:"pwd"`          // 密码
	IsVip int32     `gorm:"column:is_vip;not null;comment:是否vip" json:"is_vip"` // 是否vip
	CTime int32     `gorm:"column:c_time;not null;comment:添加时间" json:"c_time"`  // 添加时间
	CDate time.Time `gorm:"column:c_date;comment:添加日期" json:"c_date"`           // 添加日期
	Email string    `gorm:"column:email;not null;comment:邮箱" json:"email"`      // 邮箱
	Img   string    `gorm:"column:img;not null;comment:头像" json:"img"`          // 头像
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
