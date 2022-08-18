package main // 入口包

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"net/http"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

// 入口函数
func main() {
	// 初始化数据库
	bootstrap.SetupDB()
	// 自定义路由  官方
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}
