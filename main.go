package main // 入口包

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

/**
if r.URL.Path == "/" {
	fmt.Fprint(w, "<h1>Hello, 这里是 goblog</h1>")
} else if r.URL.Path == "/about" {
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
} else {
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
		"<p>如有疑惑，请联系我们。</p>")
}
*/

// func handlerFunc(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	if r.URL.Path == "/" {
// 		fmt.Fprint(w, "<h1>Hello, 这里是 goblog! </h1>")
// 	} else if r.URL.Path == "/about" {
// 		fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
// 			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
// 	} else {
// 		w.WriteHeader(http.StatusNot)
// 		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
// 			"<p>如有疑惑，请联系我们。</p>")
// 	}
// }

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	if r.URL.Path == "/" {
// 			fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
// 	} else {
// 			w.WriteHeader(http.StatusNotFound)
// 			fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
// 					"<p>如有疑惑，请联系我们。</p>")
// 	}
// }

// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
// 			"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
// }

// func defaultHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html; charset=utf-8")
// 	if r.URL.Path == "/" {
// 		fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog</h1>")
// 	} else {
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprint(w, "<h1>请求页面未找到 :(</h1>"+
// 			"<p>如有疑惑，请联系我们。</p>")
// 	}
// }
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}

func articlesShowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprint(w, "文章 ID："+id)
}

func articlesIndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "访问文章列表")
}

func articlesStoreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "创建新的文章")
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
	// 自定义路由  官方
	// router := http.NewServeMux()
	// 自定义路由 gorilla/mux
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET").Name("home")
	router.HandleFunc("/about", aboutHandler).Methods("GET").Name("about")

	router.HandleFunc("/articles/{id:[0-9]+}", articlesShowHandler).Methods("GET").Name("articles.show")
	router.HandleFunc("/articles", articlesIndexHandler).Methods("GET").Name("articles.index")
	router.HandleFunc("/articles", articlesStoreHandler).Methods("POST").Name("articles.store")

	// 自定义 404 页面
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// 中间件：强制内容类型为 HTML
	router.Use(forceHTMLMiddleware)

	homeURL, _ := router.Get("home").URL() // 获取指定命名的url
	fmt.Println("homeURL: ", homeURL)      // 命令行输出
	articleURL, _ := router.Get("articles.show").URL("id", "23")
	fmt.Println("articleURL: ", articleURL) // 命令行输出

	http.ListenAndServe(":3000", removeTrailingSlash(router))
}
