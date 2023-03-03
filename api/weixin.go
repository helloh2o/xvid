package api

import (
	"fmt"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/utils"
	v12ctx "github.com/kataras/iris/v12/context"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"time"
	"xvid/comm"
	"xvid/constans"
	"xvid/entity"
)

const (
	LaoBaiAppId     = "wx79439c9ad7d9b76f"
	LaoBaiAppSecret = "235b6cca94343634bcb58f914c3be3a3"
)

func WeiXinLogin(ctx *v12ctx.Context) {
	appid := LaoBaiAppId
	secret := LaoBaiAppSecret
	code := ctx.FormValue("code")
	channel := ctx.FormValue("channel")
	if code == constans.EMPTY {
		comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ErrTxtParamErr, nil)
		return
	}
	release := utils.RDLockOpWait(code)
	defer release()
	if limited, _ := utils.Limiter.IsV2Limited(ctx.Path()+code, time.Second, 1); limited {
		comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ErrTxtTooFast, nil)
		return
	}
	finalErr := "微信登录获取Token失败"
	// APP
	//api := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appid, secret, code)
	// 微信小程序
	api := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, code)
	resp, err := http.Get(api)
	if err != nil {
		log.Error("WX login access_token error %v", err)
	} else {
		if bodyData, err := ioutil.ReadAll(resp.Body); err != nil {
			log.Error("WX login access_token error %v", err)
		} else {
			log.Release("raw body:%s", string(bodyData))
			unionId := gjson.GetBytes(bodyData, "unionid").String()
			openId := gjson.GetBytes(bodyData, "openid").String()
			if openId != constans.EMPTY && unionId == constans.EMPTY {
				unionId = openId
			}
			if unionId == constans.EMPTY && openId == constans.EMPTY {
				log.Error("not get open id or union id: raw=> %s", string(bodyData))
			} else {
				Login(ctx, entity.LoginParam{
					UserUid: unionId,
					OpenId:  openId,
					Channel: channel,
				})
				return
			}
		}
	}
	comm.WriteHttpResponse(ctx, constans.ErrNormal, finalErr, nil)
}
