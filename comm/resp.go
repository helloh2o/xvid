package comm

import v12ctx "github.com/kataras/iris/v12/context"

// WriteHttpResponse 返回HTTP 数据
func WriteHttpResponse(ctx *v12ctx.Context, code int, err string, data interface{}) {
	resp := make(map[string]interface{})
	resp["code"] = code
	resp["error"] = err
	if data != nil {
		resp["data"] = data
	} else {
		resp["data"] = struct{}{}
	}
	ctx.JSON(resp)
}
