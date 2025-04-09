package video

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/liuyang/ly/public"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type douYin struct{}

// Gin 处理函数，接收前端请求，获取原始视频信息
func (d douYin) parseShareUrl(shareUrl string) (*public.VideoParseInfo, error) {
	// 根据请求的 URL 获取视频 ID
	videoId, err := fetchAndParseVideoID(shareUrl)
	if err != nil {
		return nil, err
	}
	fmt.Println("开始解析")
	// 根据视频 ID 解析视频信息
	videoInfo, err := parseVideoInfo(videoId)
	fmt.Println("解析完成")
	if err != nil {
		return nil, err
	}
	//if len(videoInfo.VideoUrl) > 0 {
	//	getRedirectUrl(videoInfo)
	//}
	return videoInfo, nil
}

func getRedirectUrl(videoInfo *public.VideoParseInfo) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	res2, _ := client.R().
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(videoInfo.VideoUrl)
	locationRes, _ := res2.RawResponse.Location()
	if locationRes != nil {
		(*videoInfo).VideoUrl = locationRes.String()
	}
}

// 从短链接中获取重定向后的视频 ID
func fetchAndParseVideoID(url string) (string, error) {
	client := resty.New()
	// 禁止自动重定向，以便手动获取重定向的 Location
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	res, err := client.R().
		SetHeader(public.HttpHeaderUserAgent, public.DefaultUserAgent). // 设置 User-Agent
		Get(url)
	// 如果请求失败并且不是重定向错误，返回错误
	if err != nil && !errors.Is(err, resty.ErrAutoRedirectDisabled) {
		fmt.Println(err)
		return "", fmt.Errorf("请求 URL 失败: %w", err)
	}

	// 获取重定向的目标地址
	locationRes, err := res.RawResponse.Location()
	if err != nil {
		return "", fmt.Errorf("获取重定向地址失败: %w", err)
	}

	// 从重定向地址中解析视频 ID
	videoID, err := parseVideoIDFromPath(locationRes.Path)
	if err != nil || videoID == "" {
		return "", fmt.Errorf("解析视频 ID 失败: %w", err)
	}

	return videoID, nil
}

// 解析路径中的视频 ID
func parseVideoIDFromPath(urlPath string) (string, error) {
	if urlPath == "" {
		return "", errors.New("URL 路径为空")
	}
	urlPath = strings.Trim(urlPath, "/") // 去除首尾的斜杠
	urlSplit := strings.Split(urlPath, "/")
	// 返回路径中最后一个部分作为视频 ID
	if len(urlSplit) > 0 {
		return urlSplit[len(urlSplit)-1], nil
	}
	return "", errors.New("无法从路径中解析视频 ID")
}

// 根据视频 ID 获取视频信息
func parseVideoInfo(videoId string) (*public.VideoParseInfo, error) {
	// 构造抖音视频详情页的 URL
	reqUrl := fmt.Sprintf("https://www.iesdouyin.com/share/video/%s", videoId)
	client := resty.New()
	res, err := client.R().
		SetHeader(public.HttpHeaderUserAgent, "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X)...").
		Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("请求视频详情失败: %w", err)
	}
	// 从 HTML 页面中提取视频 JSON 数据
	videoJSON, err := extractVideoJSON(res.Body())
	if err != nil {
		return nil, err
	}

	// 根据提取的 JSON 数据构造视频信息结构体
	videoInfo := constructVideoInfo(videoJSON, videoId)
	return videoInfo, nil
}

// 从 HTML 中提取包含视频信息的 JSON 数据
func extractVideoJSON(htmlBody []byte) ([]byte, error) {
	// 正则匹配 JSON 数据
	re := regexp.MustCompile(`window._ROUTER_DATA\s*=\s*(.*?)</script>`)
	match := re.FindSubmatch(htmlBody)
	// 检查匹配结果是否有效
	if len(match) < 2 {
		return nil, errors.New("无法从 HTML 中解析视频 JSON 数据")
	}
	return bytes.TrimSpace(match[1]), nil
}

// 构造视频信息结构体
func constructVideoInfo(jsonBytes []byte, videoId string) *public.VideoParseInfo {
	data := gjson.GetBytes(jsonBytes, "loaderData.video_(id)/page.videoInfoRes.item_list.0")
	if !data.Exists() {
		return nil
	}
	// 提取图片地址
	images := extractImages(data)

	// 跟换playwm 为play
	videoUrl := strings.ReplaceAll(data.Get("video.play_addr.url_list.0").String(), "playwm", "play")

	// 如果有图集，则视频地址无效，置空
	if len(images) > 0 {
		videoUrl = ""
	}

	// 构造视频信息
	videoInfo := &public.VideoParseInfo{
		Title:    data.Get("desc").String(),
		VideoUrl: videoUrl,
		MusicUrl: "",
		CoverUrl: data.Get("video.cover.url_list.0").String(),
		Images:   images,
	}
	return videoInfo
}

// 提取图集图片地址
func extractImages(data gjson.Result) []string {
	imageArray := data.Get("images").Array()
	images := make([]string, 0, len(imageArray))
	for _, item := range imageArray {
		url := item.Get("url_list.0").String()
		if url != "" {
			images = append(images, url)
		}
	}
	return images
}

// 获取真实的视频播放地址（重定向后的地址）
func resolveRedirectURL(videoInfo *public.VideoParseInfo) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	res, _ := client.R().
		SetHeader(public.HttpHeaderUserAgent, public.DefaultUserAgent).
		Get(videoInfo.VideoUrl)
	if location, _ := res.RawResponse.Location(); location != nil {
		videoInfo.VideoUrl = location.String()
	}
}

// 下载视频到本地
func downloadVideo(url, filename string) error {
	// 拼接完整的文件路径
	dir := "./downloads"
	filePath := fmt.Sprintf("%s/%s", dir, filename)
	// 检查目录是否存在，如果不存在则创建
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}
	// 发送 GET 请求下载视频
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("下载视频失败: %w", err)
	}
	defer resp.Body.Close()
	// 创建本地文件
	outFile, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer outFile.Close()
	// 写入视频数据到文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	return nil
}
