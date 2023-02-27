package api

import v12ctx "github.com/kataras/iris/v12/context"

func Index(ctx *v12ctx.Context)  {
	ctx.WriteString("welcome :) ")
}
