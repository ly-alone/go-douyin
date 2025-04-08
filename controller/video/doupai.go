package video

import (
	"errors"
	"fmt"
	"github.com/liuyang/ly/public"
	"net/url"

	"github.com/tidwall/gjson"

	"github.com/go-resty/resty/v2"
)

type douPai struct {
}

func (d douPai) parseShareUrl(shareUrl string) (*public.VideoParseInfo, error) {
	urlInfo, err := url.Parse(shareUrl)
	if err != nil {
		return nil, errors.New("parse share url fail")
	}
	if len(urlInfo.Query()["id"]) <= 0 {
		return nil, errors.New("can not parse video id from share url")
	}

	return d.parseVideoID(urlInfo.Query()["id"][0])
}

func (d douPai) parseVideoID(videoId string) (*public.VideoParseInfo, error) {
	reqUrl := fmt.Sprintf("https://v2.doupai.cc/topic/%s.json", videoId)
	headers := map[string]string{
		HttpHeaderUserAgent: DefaultUserAgent,
	}

	client := resty.New()
	res, err := client.R().
		SetHeaders(headers).
		Get(reqUrl)
	if err != nil {
		return nil, err
	}

	data := gjson.GetBytes(res.Body(), "data")
	parseInfo := &public.VideoParseInfo{
		Title:    data.Get("name").String(),
		VideoUrl: data.Get("videoUrl").String(),
		CoverUrl: data.Get("imageUrl").String(),
	}
	parseInfo.Author.Uid = data.Get("userId.id").String()
	parseInfo.Author.Name = data.Get("userId.name").String()
	parseInfo.Author.Avatar = data.Get("userId.avatar").String()

	return parseInfo, nil
}
