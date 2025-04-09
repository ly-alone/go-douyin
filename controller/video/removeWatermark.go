package video

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liuyang/ly/public"
	"io"
	"net/http"
	"strings"
)

// 请求参数结构体，包含需要解析的 URL
type Param struct {
	Url string `form:"videoUrl" json:"videoUrl" binding:"required"`
}

func RemoveWatermark(c *gin.Context) {
	var param Param
	if err := c.ShouldBindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	videoShareUrl, _ := public.RegexpMatchUrlFromString(param.Url)
	res, err := ParseVideoShareUrl(videoShareUrl)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{"data": res, "code": http.StatusOK, "msg": "success"})
}

// ParseVideoShareUrl 根据视频分享链接解析视频信息: 分享链接需是正常http链接
func ParseVideoShareUrl(shareUrl string) (*public.VideoParseInfo, error) {
	// 根据分享url判断source
	source := ""
	for itemSource, itemSourceInfo := range videoSourceInfoMapping {
		for _, itemUrlDomain := range itemSourceInfo.VideoShareUrlDomain {
			if strings.Contains(shareUrl, itemUrlDomain) {
				source = itemSource
				break
			}
		}
		if len(source) > 0 {
			break
		}
	}

	// 没有找到对应source
	if len(source) <= 0 {
		return nil, fmt.Errorf("share url [%s] not have source config", shareUrl)
	}

	// 没有对应的视频链接解析方法
	urlParser := videoSourceInfoMapping[source].VideoShareUrlParser
	if urlParser == nil {
		return nil, fmt.Errorf("source %s has no video share url parser", source)
	}
	return urlParser.parseShareUrl(shareUrl)
}

// ParseVideoId 根据视频id解析视频信息
func ParseVideoId(source, videoId string) (*public.VideoParseInfo, error) {
	if len(videoId) <= 0 || len(source) <= 0 {
		return nil, errors.New("video id or source is empty")
	}

	idParser := videoSourceInfoMapping[source].VideoIdParser
	if idParser == nil {
		return nil, fmt.Errorf("source %s has no video id parser", source)
	}

	return idParser.parseVideoID(videoId)
}

// 代理视频请求
func VideoProxy(c *gin.Context) {
	videoURL := c.DefaultQuery("url", "") // 获取 URL 参数
	if videoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Missing video URL",
		})
		return
	}

	// 向外部视频服务器发起请求
	resp, err := http.Get(videoURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to fetch video: %v", err),
		})
		return
	}
	defer resp.Body.Close()

	// 设置返回的内容类型为视频流
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline; filename=video.mp4")

	// 将外部视频流直接写入响应体
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to stream video: %v", err),
		})
		return
	}
}
