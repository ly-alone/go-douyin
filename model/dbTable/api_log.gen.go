// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dbTable

import (
	"time"
)

const TableNameAPILog = "api_log"

// APILog mapped from table <api_log>
type APILog struct {
	ID      int32     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Urlpath string    `gorm:"column:urlpath;not null;comment:请求url地址" json:"urlpath"` // 请求url地址
	URLType string    `gorm:"column:url_type;not null;comment:接口" json:"url_type"`    // 接口
	UID     int32     `gorm:"column:uid;not null;comment:uid" json:"uid"`             // uid
	Name    string    `gorm:"column:name;not null;comment:用户名" json:"name"`           // 用户名
	CTime   int32     `gorm:"column:c_time;not null;comment:请求时间" json:"c_time"`      // 请求时间
	CDate   time.Time `gorm:"column:c_date;not null;comment:请求日期" json:"c_date"`      // 请求日期
}

// TableName APILog's table name
func (*APILog) TableName() string {
	return TableNameAPILog
}
