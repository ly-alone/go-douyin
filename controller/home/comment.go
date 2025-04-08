package home

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/conf"
	"github.com/liuyang/ly/model/dbTable"
	"github.com/tidwall/gjson"
	"net/http"
	"time"
)

type CommentParam struct {
	Content string `json:"feedback"`
}

func Comment(c *gin.Context) {
	user := c.Value("user")
	userJson, err := json.Marshal(user)
	if err != nil {
		c.JSON(400, gin.H{"code": 2003, "msg": "系统错误"}) // 如果 content 为空，则返回 400 错误
		return
	}
	userinfo := gjson.GetBytes(userJson, "userinfo")
	userinfoStr := userinfo.String()
	user_id := gjson.Get(userinfoStr, "user_id")
	uid := user_id.Int()
	var commentParam CommentParam
	// 获取 'content' 字段的值，如果没有则为空字符串
	if err := c.ShouldBindJSON(&commentParam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if commentParam.Content == "" {
		c.JSON(400, gin.H{"error": "Content is required"}) // 如果 content 为空，则返回 400 错误
		return
	}

	db := conf.GetDB()
	times := int32(time.Now().Unix())
	data := &dbTable.Comment{
		Content:   commentParam.Content,
		CreatedTs: times,
		UID:       int32(uid),
	}

	res := db.Create(data) // 插入数据
	if res.Error != nil {
		c.JSON(500, gin.H{"error": res.Error.Error()}) // 如果插入失败，则返回错误信息
		return
	}

	//newID := data.ID // GORM 会自动将插入记录的 ID 填充到 data 结构体中
	c.JSON(200, gin.H{"code": 200, "data": data, "msg": "success"})
	return
}

func ShowComment(c *gin.Context) {
	page, _ := c.Get("page")
	pageNum, _ := c.Get("page_num")
	if page == nil || pageNum == nil {
		page = 1
		pageNum = 20
	}
	offsets, _ := page.(int)
	limits, _ := pageNum.(int)
	db := conf.GetDB()
	var data []map[string]interface{}
	db.Debug().Table("comment").Offset(offsets).Limit(limits).Scan(&data)
	var newData []map[string]interface{}
	uids := make([]int32, 0)
	for _, v := range data {
		uid, ok := v["uid"].(int32)
		if ok {
			uids = append(uids, uid)
		}
		createdTs, ok := v["created_ts"].(int64)
		if ok == false {
			createdTs = int64(0)
		}
		times := time.Unix(createdTs, 0)
		date := times.Format("2006-01-02 15:04:05")
		v["created_ts"] = date
		newData = append(newData, v)
	}
	var userInfo []map[string]interface{}
	db.Debug().Table("user").Select("name,id").Where("id in (?)", uids).Scan(&userInfo)
	for k, v := range newData {
		v["user_info"] = map[string]interface{}{
			"name": "",
			"id":   v["id"],
		}
		vId, _ := v["id"].(int32)
		for _, val := range userInfo {
			valId, _ := val["id"].(uint32)
			if int(vId) == int(valId) {
				v["user_info"] = map[string]interface{}{
					"name": val["name"],
					"id":   val["id"],
				}
				newData[k] = v
			}
		}
	}
	c.JSON(200, gin.H{"code": 200, "data": newData, "msg": "success"})
	return
}
