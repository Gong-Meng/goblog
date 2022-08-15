package main // 入口包

import (
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/pkg/logger"
	"net/http"
)

// 入口函数
func main() {
	// 初始化数据库
	bootstrap.SetupDB()
	// 自定义路由  官方
	router := bootstrap.SetupRoute()

	err := http.ListenAndServe(":3000", middlewares.RemoveTrailingSlash(router))
	logger.LogError(err)
}
