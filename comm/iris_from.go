package comm

import (
	v12ctx "github.com/kataras/iris/v12/context"
	"strconv"
)

// IrisValueInt 获取iris int值
func IrisValueInt(context *v12ctx.Context, key string, defaultVal int) int {
	valStr := context.FormValue(key)
	if val, err := strconv.ParseInt(valStr, 10, 64); err != nil {
		return defaultVal
	} else {
		return int(val)
	}
}

func FormValueInt(context *v12ctx.Context, key string, defaultVal int) int {
	return IrisValueInt(context, key, defaultVal)
}

// FormValueFloat 获取iris float值
func FormValueFloat(context *v12ctx.Context, key string, defaultVal float64) float64 {
	valStr := context.FormValue(key)
	if val, err := strconv.ParseFloat(valStr, 10); err != nil {
		return defaultVal
	} else {
		return val
	}
}

// FormValueBool 获取iris bool值
func FormValueBool(context *v12ctx.Context, key string) bool {
	valStr := context.FormValue(key)
	if val, err := strconv.ParseBool(valStr); err != nil {
		return false
	} else {
		return val
	}
}
