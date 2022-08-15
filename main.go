package main // 入口包

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// 自定义路由 gorilla/mux
// var router = mux.NewRouter()
var router *mux.Router

// 定义数据库连接
var db *sql.DB

func getRouteVariable(parameterName string, r *http.Request) string {
	vars := mux.Vars(r)
	return vars[parameterName]
}

func forceHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1、设置标头
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 2、继续处理请求
		next.ServeHTTP(w, r)
	})
}

// 移除斜杠
func removeTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. 除首页以外，移除所有请求路径后面的斜杆
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}

		// 2、将请求传递下去
		next.ServeHTTP(w, r)
	})
}

// 入口函数
func main() {

	// 初始化数据库连接
	database.Initialize()
	db = database.DB

	// 初始化数据库
	bootstrap.SetupDB()
	// 自定义路由  官方
	// router := http.NewServeMux()
	router = bootstrap.SetupRoute()

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	homeURL, _ := router.Get("home").URL() // 获取指定命名的url
	fmt.Println("homeURL: ", homeURL)      // 命令行输出
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL) // 命令行输出

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
