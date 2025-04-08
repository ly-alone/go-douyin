package video

import (
	"errors"
	"github.com/liuyang/ly/public"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
)

type sixRoom struct {
}

func (s sixRoom) parseShareUrl(shareUrl string) (*public.VideoParseInfo, error) {
	urlInfo, err := url.Parse(shareUrl)
	if err != nil {
		return nil, errors.New("parse share url fail")
	}
	var videoId string
	if strings.Contains(shareUrl, "watchMini.php?vid=") {
		if len(urlInfo.Query()["vid"]) <= 0 {
			return nil, errors.New("can not parse video id from share url")
		}
		videoId = urlInfo.Query()["vid"][0]
	} else {
		videoId = strings.ReplaceAll(urlInfo.Path, "/v/", "")
	}
	return s.parseVideoID(videoId)
}

func (s sixRoom) parseVideoID(videoId string) (*public.VideoParseInfo, error) {
	reqUrl := "https://v.6.cn/coop/mobile/index.php?padapi=minivideo-watchVideo.php&av=3.0&encpass=&logiuid=&isnew=1&from=0&vid=" + videoId
	client := resty.New()
	videoRes, err := client.R().
		SetHeader(HttpHeaderReferer, "https://m.6.cn/v/"+videoId).
		SetHeader(HttpHeaderUserAgent, DefaultUserAgent).
		Get(reqUrl)
	if err != nil {
		return nil, err
	}

	data := gjson.GetBytes(videoRes.Body(), "content")
	parseInfo := &public.VideoParseInfo{
		Title:    data.Get("title").String(),
		VideoUrl: data.Get("playurl").String(),
		CoverUrl: data.Get("picurl").String(),
	}
	parseInfo.Author.Name = data.Get("alias").String()
	parseInfo.Author.Avatar = data.Get("picuser").String()

	return parseInfo, nil
}
