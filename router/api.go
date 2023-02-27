package router

import (
	"github.com/helloh2o/lucky"
	"github.com/kataras/iris/v12/context"
	"xvid/api"
)

const (
	OPEN    = true
	PRIVATE = false
)

func register() {
	get("/", OPEN, api.Index)
	get("/video/parse/share", OPEN, api.ParseByShare)
	post("/user/login", PRIVATE, func(c *context.Context) {})
}

func get(path string, isOpen bool, handlers ...context.Handler) {
	if isOpen {
		openApi[path] = true
	}
	lucky.Get(path, handlers...)
}
func post(path string, isOpen bool, handlers ...context.Handler) {
	if isOpen {
		openApi[path] = true
	}
	lucky.Post(path, handlers...)
}
