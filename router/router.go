package router

import (
	"github.com/helloh2o/lucky"
	"github.com/helloh2o/lucky/log"
	"github.com/helloh2o/lucky/utils"
	v12ctx "github.com/kataras/iris/v12/context"
	"math"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"xvid/comm"
	"xvid/config"
	"xvid/constans"
)

func InitRouter(auth func(string) bool) {
	defer register()
	// auth x-token
	lucky.Iris().Use(func(ctx *v12ctx.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("context panic err:%v", string(debug.Stack()))
				panic(r)
			}
		}()
		xt := getCtxValue(ctx, "_xt") // sign token
		ts := getCtxValue(ctx, "_ts") // timestamp
		rs := getCtxValue(ctx, "_rs") // random str
		xc := getCtxValue(ctx, "_xc") // channel
		isOpen := isOpen(ctx.Request().URL.Path)
		switch {
		case isOpen:
		default:
			if xt != utils.GetMD5Hash(xc+config.Get().Signature+rs+ts) {
				log.Error("error sign:%s, ts:%s, path:%s", xt, ts, ctx.Request().URL.Path)
				comm.WriteHttpResponse(ctx, constans.ErrForbidden, constans.ErrTxtParamErr, nil)
				return
			}
		}
		if xt != "" {
			// 非白名单
			if !isOpen {
				// 维护中
				if config.Get().Maintenance && comm.GetRemoteAddr(ctx) != config.Get().MaintenanceIp {
					comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ErrServerFix, nil)
					return
				}
			}
			ts, _ := strconv.ParseInt(ts, 10, 64)
			now := time.Now().Unix()
			// 时间差值是否过大
			if math.Abs(float64(now-ts)) > 86400 {
				comm.WriteHttpResponse(ctx, constans.ErrForbidden, constans.ErrVerifyTime, nil)
				return
			}
			// 防止业务攻击，
			lk := xt + ":" + ctx.Request().URL.Path
			release, lock, _ := utils.RDLockOpTimeout(lk, time.Second*10)
			defer release()
			if !lock {
				log.Error("lock key:%s attack from ip:%s, client:%s, channel:%s", lk, comm.GetRemoteAddr(ctx), xc, ctx.FormValue("channel"))
				comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.EMPTY, nil)
				return
			}
		}
		// 拦截可能的SQL 注入
		if !isSqlSave(ctx) {
			comm.WriteHttpResponse(ctx, constans.ErrNormal, "X="+constans.ErrTxtParamErr, nil)
			return
		}
		// 拦截权限
		if auth != nil && !auth(getCtxValue(ctx, "token")) {
			comm.WriteHttpResponse(ctx, constans.ErrForbidden, constans.ErrNoAuth, nil)
			return
		}
		// 捕获逻辑异常
		defer func() {
			if r := recover(); r != nil {
				comm.WriteHttpResponse(ctx, constans.ErrNormal, constans.ServerPanic, nil)
				panic(r)
			}
		}()
		// 计算执行逻辑事件
		begin := time.Now().UnixNano() / int64(time.Millisecond)
		ctx.Next()
		costs := time.Now().UnixNano()/int64(time.Millisecond) - begin
		if ctx.Request().Method == "POST" {
			log.Debug("===> execute_api:%s, costs %dms <===", ctx.Request().URL.Path, costs)
		}
	})
}
func getCtxValue(ctx *v12ctx.Context, key string) string {
	v := ctx.Request().Header.Get(key)
	if v == "" {
		v = ctx.FormValue(key)
	}
	return v
}

// 遍历参数是否有sql异常
func isSqlSave(ctx *v12ctx.Context) bool {
	kw := []string{
		"dump(", "md5(", "delay'", "expr ", "select ", "set(", "cast(", "concat(", "${", "extractvalue(",
		"waitfor/", "delete ", "update ", "drop ", "truncate ", "convert(", "count(", "group by", "and 1=", "0<>(",
		"union ", "from ",
	}
	for fk, v := range ctx.Request().Form {
		for _, k := range kw {
			if strings.Contains(strings.ToLower(v[0]), k) {
				log.Error("form key:%s, error value:%s, sql kw:%s", fk, v[0], k)
				return false
			}
		}
	}
	return true
}
