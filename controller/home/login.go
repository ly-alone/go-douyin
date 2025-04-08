package home

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/model"
	"github.com/liuyang/ly/model/dbTable"
	"github.com/liuyang/ly/public"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var loginRequest LoginRequest
	// 绑定请求体的 JSON 数据到 loginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4010, "msg": "参数错误"})
		return
	}

	// 查询账号是否存在
	res := model.GetUserInfoByName(loginRequest.Username)
	if res == (dbTable.User{}) {
		// 返回 401 而不是 200，因为这是身份验证失败
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4010, "msg": "账号密码错误"})
		return
	}

	// 验证密码是否匹配
	hash := md5.New()
	hash.Write([]byte(loginRequest.Password))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)

	if hashString != res.Pwd {
		// 密码不匹配，返回 401
		c.JSON(http.StatusUnauthorized, gin.H{"code": 4010, "msg": "账号密码错误"})
		return
	}

	// 创建 JWT 用户信息（不要传递密码）
	jwtUser := make(map[string]interface{})
	jwtUser["user_name"] = loginRequest.Username
	jwtUser["user_id"] = res.ID

	// 将 jwtUser 转换为 JSON 字符串
	jsonData, err := json.Marshal(jwtUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2001, "msg": "系统错误，请稍后再试"})
		return
	}

	// 生成 JWT token
	token, err := public.CreateJwt(res.Name, string(jsonData))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 2001, "msg": "请稍后再试"})
		return
	}
	login := make(map[string]string)
	login["token"] = token
	// 返回生成的 token
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": login, "msg": "成功"})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func Register(c *gin.Context) {
	var Request RegisterRequest
	if err := c.ShouldBindJSON(&Request); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1001, "msg": "参数错误"})
		return
	}
	//验证用户名是否存在
	IsData := model.GetUserInfoByName(Request.Username)
	if IsData != (dbTable.User{}) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 1003, "msg": "用户名已经存在"})
		return
	}
	//加密密码
	hash := md5.New()
	hash.Write([]byte(Request.Password))
	hashBytes := hash.Sum(nil)
	pwd := hex.EncodeToString(hashBytes)
	_, err := model.AddUserData(Request.Username, pwd, Request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 2001, "msg": "添加失败"})
		return
	}
	c.JSON(200, gin.H{"code": 200, "msg": "success"})
	return
}
