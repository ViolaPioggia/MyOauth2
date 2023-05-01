package middleware

import (
	"net/http"
)

// CorsMiddleware ，如果做前后端分离可能会遇到cors跨域的问题，所以提供了以下函数解决有需求的情况
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. [必须]接受指定域的请求，可以使用*不加以限制，但不安全
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		// 2. [必须]设置服务器支持的所有跨域请求的方法
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,DELETE,OPTIONS")
		// 3. [可选]服务器支持的所有头信息字段，不限于浏览器在"预检"中请求的字段
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Content-Length,Token")
		// 4. [可选]设置XMLHttpRequest的响应对象能拿到的额外字段
		w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers,Token")
		// 5. [可选]是否允许后续请求携带认证信息Cookie，该值只能是true，不需要则不设置
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}
