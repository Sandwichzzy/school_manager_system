package middlewares

import (
	"fmt"
	"net/http"
)


func SecurityHeaders(next http.Handler) http.Handler {
	fmt.Println("SecurityHeaders Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("SecurityHeaders Middleware being returned...")
		w.Header().Set("X-DNS-Prefetch-Control", "off") //降低DNS預取的風險，防止浏览器提前解析页面中的链接域名，降低隐私泄露风险

		w.Header().Set("X-Frame-Options", "DENY") //防止网页在其他网站的ifame中加载，防止点击劫持攻击
		w.Header().Set("X-XSS-Protection", "1; mode=block") //啟用XSS過濾器,mode=block表示检测到XSS攻击时阻止页面渲染
		w.Header().Set("X-Content-Type-Options", "nosniff") //防止MIME类型嗅探攻击，强制浏览器使用响应头中的Content-Type，而不猜测文件类型

		// 		HSTS（HTTP严格传输安全）策略
		// max-age=63072000：约2年有效期
		// includeSubDomains：应用于所有子域名
		// preload：可以提交到浏览器的HSTS预加载列表
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		
		//CSP（内容安全策略），最强大的XSS防护
		//default-src 'self'：只允许从同源加载资源
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
	
		//控制Referer头的发送
		//no-referrer：完全不发送Referer头
		w.Header().Set("Referrer-Policy", "no-referrer")

		w.Header().Set("X-Powered-By", "Django") //隐藏服务器信息，防止攻击者利用已知漏洞进行攻击
		w.Header().Set("Server", "")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Permissions-Policy", "geolocation=(self), microphone=()")

		next.ServeHTTP(w, r)
		fmt.Println("SecurityHeaders Middleware ends...")
	})
}

// BASIC MIDDKLEWARE SKELETON
// func securityHeaders(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 	})
// }