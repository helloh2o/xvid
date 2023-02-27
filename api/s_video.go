package api

import (
	"fmt"
	"github.com/helloh2o/lucky/log"
	v12ctx "github.com/kataras/iris/v12/context"
	"github.com/wujunwei928/parse-video/parser"
	"regexp"
	"xvid/comm"
	"xvid/constans"
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
		comm.WriteHttpResponse(ctx, constans.Success, constans.EMPTY, parseRes)
	}
}
