package home

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/model"
	"github.com/liuyang/ly/public"
	"math/rand"
	"time"
)

func UserInfo(c *gin.Context) {
	uId, _, err := public.User.GetUserInfoByChan(c)

	if err != nil {
		c.JSON(400, gin.H{"code": 3001, "msg": err.Error()})
		return
	}
	uid := int32(uId)
	//通过uid查询用户信息
	userInfo := model.GetUserInfoById(uid, "c_time,is_vip,name,id,img")
	times := time.Unix(int64(userInfo.CTime), 0)
	date := times.Format("2006-01-02 15:04:05")
	data, _ := json.Marshal(userInfo)
	userMap := make(map[string]interface{})
	json.Unmarshal(data, &userMap)
	userMap["date"] = date
	c.JSON(200, gin.H{"code": 200, "msg": "success", "data": userMap})
}

func Img(c *gin.Context) {

}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	for i := range b {
		b[i] = letters[b[i]%byte(len(letters))]
	}
	return string(b)
}
