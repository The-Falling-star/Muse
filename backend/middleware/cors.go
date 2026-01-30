package middleware

import (
	"net/http"
)

// CORS 跨域中间件
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 响应头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, connect-protocol-version, connect-timeout-ms, connect-grpc-web-protocol")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400") // 24小时

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 调用下一个处理器
		next.ServeHTTP(w, r)
	})
}
