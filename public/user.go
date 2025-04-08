package public

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

type user struct{}

var User user

func (u user) GetUserInfoByChan(c context.Context) (int64, string, error) {
	// 从上下文中获取用户信息
	userChan := c.Value("user")
	userJson, err := json.Marshal(userChan)
	if err != nil {
		return 0, "", errors.New("user not exist")
	}
	userinfo := gjson.GetBytes(userJson, "userinfo")
	userinfoStr := userinfo.String()
	// 使用 gjson 解析 JSON 字符串
	userId := gjson.Get(userinfoStr, "user_id").Int()
	userName := gjson.Get(userinfoStr, "user_name").String()
	// 校验解析结果
	if userId == 0 || userName == "" {
		return 0, "", fmt.Errorf("invalid user_id or user_name")
	}

	return userId, userName, nil
}
