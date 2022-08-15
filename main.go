package main // 入口包

import (
	"database/sql"
	"fmt"
	"goblog/bootstrap"
	"goblog/pkg/database"
	"goblog/pkg/logger"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

// 自定义路由 gorilla/mux
// var router = mux.NewRouter()
var router *mux.Router

// 定义数据库连接
var db *sql.DB

func saveArticleToDB(title string, body string) (int64, error) {
	// 变量初始化
	var (
		id   int64
		err  error
		rs   sql.Result
		stmt *sql.Stmt
	)

	// 1. 获取一个 perpare 声明语句
	stmt, err = db.Prepare("INSERT INTO articles (title, body) VALUES(?,?)")
	// 列行的错误检测
	if err != nil {
		return 0, err
	}

	// 2. 在此函数运行结束后关闭此语句，防止sql占用连接
	defer stmt.Close()

	// 3. 执行请求，传参进入绑定的容器
	rs, err = stmt.Exec(title, body)

	if err != nil {
		return 0, err
	}

	// 4. 插入成功的话，绘返回自增 ID
	if id, err = rs.LastInsertId(); id > 0 {
		return id, nil
	}

	return 0, err
}

// Article 对应一条文章数据
type Article struct {
	Title, Body string
	ID          int64
}

// Delete 方法用以从数据库中删除单条记录
func (a Article) Delete() (rowsAffected int64, err error) {
	rs, err := db.Exec("DELETE FROM articles WHERE id = " + strconv.FormatInt(a.ID, 10))

	if err != nil {
		return 0, err
	}

	// 删除成功，跳转到文章详情页
	if n, _ := rs.RowsAffected(); n > 0 {
		return n, nil
	}

	return 0, nil
}

// RouteName2URL 通过路由名称来获取 URL
func RouteName2URL(routeName string, pairs ...string) string {
	url, err := router.Get(routeName).URL(pairs...)
	if err != nil {
		logger.LogError(err)
		return ""
	}

	return url.String()
}

// Int64ToString 将 int64 转换为 string
func Int64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

func getArticleByID(id string) (Article, error) {
	article := Article{}
	query := "SELECT * FROM articles WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&article.ID, &article.Title, &article.Body)
	return article, err
}

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

func articlesDeleteHandler(w http.ResponseWriter, r *http.Request) {

	// 1. 获取URL参数
	id := getRouteVariable("id", r)

	// 2. 读取对应的文章数据
	article, err := getArticleByID(id)

	// 3. 如果出现错误
	if err != nil {
		if err == sql.ErrNoRows {
			// 3.1 数据未找到
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "404 文章未找到")
		} else {
			// 3.2 数据库错误
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		}
	} else {
		// 4. 未出现错误，执行删除操作
		rowsAffected, err := article.Delete()

		// 4.1 发生错误
		if err != nil {
			// 应该是sql报错
			logger.LogError(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 服务器内部错误")
		} else {
			// 4.2 未发生错误
			if rowsAffected > 0 {
				// 重定向到文章列表页
				indexURL, _ := router.Get("articles.index").URL()
				http.Redirect(w, r, indexURL.String(), http.StatusFound)
			} else {
				// Edge case
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprint(w, "404 文章未找到")
			}
		}

	}
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

	router.HandleFunc("/articles/{id:[0-9]+}/delete", articlesDeleteHandler).Methods("POST").Name("articles.delete")

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	homeURL, _ := router.Get("home").URL() // 获取指定命名的url
	fmt.Println("homeURL: ", homeURL)      // 命令行输出
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL) // 命令行输出

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
