package main // 入口包

import (
	"embed"
	"goblog/app/http/middlewares"
	"goblog/bootstrap"
	"goblog/config"
	c "goblog/pkg/config"
	"net/http"
)

//go:embed resources/views/articles/*
//go:embed resources/views/auth/*
//go:embed resources/views/categories/*
//go:embed resources/views/layouts/*
var tplFS embed.FS

func init() {
	// 初始化配置信息
	config.Initialize()
}

// 入口函数
func main() {
	// 初始化数据库
	bootstrap.SetupDB()

	// 初始化模板
	bootstrap.SetupTemplate(tplFS)

	// 自定义路由  官方
	router := bootstrap.SetupRoute()

	http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}
