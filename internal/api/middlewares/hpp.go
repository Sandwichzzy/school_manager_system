package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuery                  bool     // 是否检查查询参数
	CheckBody                   bool     // 是否检查请求体
	CheckBodyOnlyForContentType string   // 只对特定 Content-Type 检查
	Whitelist                   []string // 允许的参数白名单
}

// 为什么需要两层返回
// 设计模式：柯里化（Currying）/函数式中间件

// 外层：配置函数（接收选项）
func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	fmt.Println("HPP Middleware...")
	// 返回一个函数，这个函数接收一个 Handler 并返回新的 Handler
	return func(next http.Handler) http.Handler {
		// 内层：实际的中间件逻辑
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("HPP Middleware being returned...")
			if options.CheckBody && r.Method == http.MethodPost && isCorrectContentType(r, options.CheckBodyOnlyForContentType) {
				//filter body params
				filterBodyParams(r, options.Whitelist)
			}
			if options.CheckQuery && r.URL.Query() != nil {
				//filter query params
				filterQueryParams(r, options.Whitelist)
			}
			// 中间件具体实现
			next.ServeHTTP(w, r)
			fmt.Println("HPP Middleware ends...")
		})
	}
}

func isCorrectContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

// 过滤 POST 表单参数
func filterBodyParams(r *http.Request, whitelist []string) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0]) //first value
			// r.Form.Set(k,v[len(v)-1]) //last value
		}
		if !isWhitelisted(k, whitelist) {
			delete(r.Form, k)
		}
	}
}

func isWhitelisted(param string, whitelist []string) bool {
	for _, item := range whitelist {
		if item == param {
			return true
		}
	}
	return false
}

// 过滤 URL 查询参数
func filterQueryParams(r *http.Request, whitelist []string) {
	query := r.URL.Query()
	for k, v := range query {
		// 如果参数有多个值 → 只取第一个值（或注释中的最后一个值）
		if len(v) > 1 {
			query.Set(k, v[0]) //first value
			// query.Set(k,v[len(v)-1]) //last value
		}
		//如果参数不在白名单中 → 删除该参数
		if !isWhitelisted(k, whitelist) {
			query.Del(k)
		}
	}
	r.URL.RawQuery = query.Encode()
}
