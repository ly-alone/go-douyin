package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"strings"
	"time"
)

var secretKey = []byte("ly-xcd-shdbz")

func CreateJwt(username string, userinfo string) (string, error) {
	auth := jwt.MapClaims{}
	auth["exp"] = time.Now().Add(time.Hour * 30).Unix() // JWT的过期时间
	auth["iat"] = time.Now().Unix()                     // 表示JWT的签发时间
	auth["nbf"] = time.Now().Unix()                     // 表示JWT在何时之前不可被接受处理的时间。这可以用于防止JWT在预定的生效时间之前被使用。
	auth["iss"] = "pescms-rent"                         // JWT的实体或服务的标识
	auth["sub"] = "pescms-rent"                         // JWT所面向的用户或实体的标识
	auth["username"] = username
	auth["userinfo"] = userinfo
	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth)

	// 签名 token
	signedToken, err := token.SignedString(secretKey)
	fmt.Println(signedToken, string(signedToken))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return signedToken, nil
}

func ValidateJwt(c *gin.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		signedToken := c.GetHeader("Authorization")
		if signedToken == "" {
			c.JSON(401, gin.H{"error": "请先登录"})
			c.Abort()
			return
		}
		if strings.Contains(signedToken, "Bearer") == true {
			data := strings.Split(signedToken, " ")
			signedToken = data[1]
		}
		parsedToken, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
			// 确认签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// 返回签名密钥
			return secretKey, nil
		})
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}
		// 确认 token 有效
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
			// 将 claims 存入上下文，后续处理可以访问
			c.Set("user", claims)
		} else {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
