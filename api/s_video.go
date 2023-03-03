package api

import (
	"context"
	"fmt"
	"github.com/helloh2o/lucky/cache"
	"github.com/helloh2o/lucky/log"
	v12ctx "github.com/kataras/iris/v12/context"
	"github.com/tidwall/gjson"
	"github.com/wujunwei928/parse-video/parser"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"xvid/comm"
	"xvid/constans"
)

const (
	baseHost = "http://localhost:8080"
)

// ParseByShare 短视频服务
func ParseByShare(ctx *v12ctx.Context) {
	urlReg := regexp.MustCompile(`http[s]?:\/\/[\w.-]+[\w\/-]*[\w.-]*\??[\w=&:\-\+\%]*`)
	shareInfo := ctx.FormValue("url")
	videoShareUrl := urlReg.FindString(shareInfo)
	parseRes, err := parser.ParseVideoShareUrl(videoShareUrl)
	if err != nil {
		log.Error("parse video error")
		comm.WriteHttpResponse(ctx, constans.ErrNormal, fmt.Sprintf("解析失败:%s", err.Error()), nil)
	} else {
		if parseRes.VideoUrl == constans.EMPTY && strings.Contains(videoShareUrl, "douyin.com") {
			// try another way
			client := http.DefaultClient
			data := "share_link=" + url.QueryEscape(shareInfo)
			if resp1, err1 := client.Post("https://dy.gyh.im/douyin", "application/x-www-form-urlencoded", strings.NewReader(data)); err1 != nil {
				log.Error("req dy.gyh.im error:%v", err1)
				comm.WriteHttpResponse(ctx, constans.ErrNormal, fmt.Sprintf("解析失败:%s", err1.Error()), nil)
			} else {
				defer resp1.Body.Close()
				if respData, err := ioutil.ReadAll(resp1.Body); err != nil {
					log.Error("read data dy.gyh.im error:%v", err)
					comm.WriteHttpResponse(ctx, constans.ErrNormal, fmt.Sprintf("解析失败:%s", err.Error()), nil)
				} else {
					log.Release("dy.gyh.im resp raw:%s", respData)
					parseRes.Author.Name = gjson.GetBytes(respData, "author.nickname").String()
					parseRes.Author.Uid = gjson.GetBytes(respData, "author.uid").String()
					parseRes.Author.Avatar = gjson.GetBytes(respData, "author.avatar.url_list.0").String()
					parseRes.VideoUrl = gjson.GetBytes(respData, "video.play_addr.url_list").String()
					parseRes.Title = gjson.GetBytes(respData, "desc").String()
					parseRes.CoverUrl = gjson.GetBytes(respData, "music.cover_medium.url_list.0").String()
					parseRes.MusicUrl = gjson.GetBytes(respData, "video.cover.url_list.0").String()
					// 302
					if strings.Contains(parseRes.VideoUrl, "video_id=") {
						// get location
						if resp2, err2 := http.Get(parseRes.VideoUrl); err2 == nil {
							defer resp2.Body.Close()
							if resp2.StatusCode == 200 {
								parseRes.VideoUrl = resp2.Request.URL.String()
							}
						}
					}
				}
			}
		}
		if parseRes.VideoUrl == constans.EMPTY {
			comm.WriteHttpResponse(ctx, constans.ErrNormal, fmt.Sprintf("解析失败:%s", "请稍后重试"), nil)
		} else {
			proxyUrl := fmt.Sprintf("%s/video/proxy/url?redirect=%s", baseHost, url.QueryEscape(parseRes.VideoUrl))
			dk := time.Now().Format("2006-01-02") + ":parse"
			cache.RedisC.IncrBy(context.Background(), dk, 1)
			cache.RedisC.Expire(context.Background(), dk, time.Hour*24)
			comm.WriteHttpResponse(ctx, constans.Success, proxyUrl, parseRes)
		}
	}

}

func ProxyVid(ctx *v12ctx.Context) {
	target := ctx.FormValue("redirect")
	target, _ = url.QueryUnescape(target)
	if resp, err := http.Get(target); err == nil {
		defer resp.Body.Close()
		io.Copy(ctx.ResponseWriter(), resp.Body)
		log.Release("redirect ok ...")
	} else {
		ctx.StatusCode(500)
	}
}
