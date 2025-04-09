package public

import (
	"github.com/gin-gonic/gin"
	mysqlConf "github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model"
	"github.com/liuyang/ly/model/dbTable"
	"time"
)

type count struct{}

var Count count

var urlPathMap = map[string]string{
	"/tool/upload":  "update",
	"/tool/add":     "add",
	"/tool/comment": "comment",
}

func (co count) Url() gin.HandlerFunc {
	return func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		urlType, exist := urlPathMap[urlPath]
		if !exist {
			return
		}
		//获取用户信息
		uid, name, err := User.GetUserInfoByChan(c)
		if err != nil {
			return
		}
		userSvip := model.GetUserInfoById(int32(uid), "is_vip")
		if userSvip.IsVip == 0 {
			//所有非会员限制
		}
		go func() {
			//记录
			timeUnix := time.Now().Unix()
			cDate := time.Unix(timeUnix, 0)
			apiLogData := dbTable.APILog{
				Urlpath: urlPath,
				URLType: urlType,
				UID:     int32(uid),
				Name:    name,
				CTime:   int32(timeUnix),
				CDate:   cDate,
			}
			db := mysqlConf.GetDB()
			db.Create(&apiLogData)
		}()
	}
}
