package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
)

// http 相关
const (
	HttpHeaderUserAgent   = "User-Agent" //http header
	HttpHeaderReferer     = "Referer"
	HttpHeaderContentType = "Content-Type"
	HttpHeaderCookie      = "Cookie"
	DefaultUserAgent      = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"
)

// VideoParseInfo 视频解析信息
type VideoParseInfo struct {
	Author struct {
		Uid    string `json:"uid"`    // 作者id
		Name   string `json:"name"`   // 作者名称
		Avatar string `json:"avatar"` // 作者头像
	} `json:"author"`
	Title    string   `json:"title"`     // 描述
	VideoUrl string   `json:"video_url"` // 视频播放地址
	MusicUrl string   `json:"music_url"` // 音乐播放地址
	CoverUrl string   `json:"cover_url"` // 视频封面地址
	Images   []string `json:"images"`    // 图集图片地址列表
}

type ReturnJson struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func RegexpMatchUrlFromString(str string) (string, error) {
	urlReg, err := regexp.Compile(`https?://[\w.-]+[\w/-]*[\w.-:]*\??[\w=&:\-+%.]*/*`)
	if err != nil {
		return "", fmt.Errorf("match url regexp compile error: %s", err.Error())
	}

	findStr := urlReg.FindString(str)
	if len(findStr) <= 0 {
		return "", fmt.Errorf("str not have url")
	}

	return findStr, nil
}

func Page() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取分页参数，如果没有传递则使用默认值
		page := c.DefaultQuery("page", "1")           // 默认从第1页开始
		pageSize := c.DefaultQuery("page_size", "10") // 默认每页10条记录

		// 将查询参数转为整数
		pages, _ := strconv.Atoi(page)
		pageSizes, _ := strconv.Atoi(pageSize)

		// 计算偏移量（offset）
		offset := (pages - 1) * pageSizes

		// 将分页参数存储到上下文中
		c.Set("offset", offset)
		c.Set("limit", pageSizes)

		// 继续处理下一个中间件
		c.Next()
	}
}
