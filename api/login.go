package api

import (
	"encoding/base64"
	"fmt"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/utils"
	"github.com/helloh2o/lucky/xdb"
	v12ctx "github.com/kataras/iris/v12/context"
	"strings"
	"time"
	"xvid/comm"
	"xvid/constans"
	"xvid/entity"
	"xvid/queues"
)

const tt = "t_t@"

func LoginCtx(ctx *v12ctx.Context) {
	Login(ctx, entity.LoginParam{
		UserUid: ctx.FormValue("user_uid"),
		Channel: ctx.FormValue("channel"),
		Device:  ctx.FormValue("device"),
		Version: comm.IrisValueInt(ctx, "version", 0),
	})
}

func Login(ctx *v12ctx.Context, params entity.LoginParam) {
	if limited, _ := utils.Limiter.IsV2Limited(ctx.Path()+params.UserUid, time.Second, 1); limited {
		comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ErrTxtTooFast, nil)
		return
	}
	// 生成同样的token
	userToken := tt + utils.GetMD5Hash(base64.StdEncoding.EncodeToString([]byte(params.UserUid)))
	appUser, done := comm.GetAppUser(userToken)
	defer done()
	if appUser == nil {
		var temp entity.AppUser
		// find in db
		if err := xdb.QqsDB().Where("user_uid=?", params.UserUid).First(&temp).Error; err != nil {
			appUser = &entity.AppUser{
				UserUid:     params.UserUid,
				OpenId:      params.OpenId,
				Channel:     params.Channel,
				Nickname:    params.Nickname,
				Avatar:      params.Avatar,
				LeftTimesWT: 5,
				RegTime:     time.Now(),
				OpTime:      time.Now(),
				Token:       tt + base64.StdEncoding.EncodeToString([]byte(params.UserUid)),
			}
			if err = xdb.QqsDB().Create(appUser).Error; err != nil {
				comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ErrBusyTxt, nil)
				return
			}
			if strings.Contains(params.Channel, "wx_min") {
				appUser.Nickname = fmt.Sprintf("微信用户:%d", 1000+appUser.Id)
			}
			log.Release("create new user:%s", params.UserUid)
		} else {
			appUser = &temp
			comm.SetAppUserCache(appUser.Token, appUser)
		}
	} else {
		appUser.AfterLogic = func() {
			queues.AppUserQ.PushToQueueWait(appUser.Token)
		}
	}
	comm.WriteHttpResponse(ctx, constans.Success, constans.EMPTY, appUser)
}
