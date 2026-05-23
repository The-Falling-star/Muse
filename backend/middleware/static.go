// Package middleware 静态文件服务中间件
package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// StaticFileHandler 创建静态文件服务处理器
func StaticFileHandler(frontendDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API请求直接返回404，让主路由器处理
		if strings.HasPrefix(r.URL.Path, "/muse.") {
			http.NotFound(w, r)
			return
		}

		// 构建文件路径
		path := filepath.Join(frontendDir, r.URL.Path)

		// 检查文件是否存在
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// SPA路由支持 - 不存在的文件都返回index.html
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		}

		// 存在的文件直接提供服务
		http.FileServer(http.Dir(frontendDir)).ServeHTTP(w, r)
	})
}

// ServeStaticFiles 条件性地添加静态文件服务到路由器
func ServeStaticFiles(mux *http.ServeMux, enabled bool, frontendDir string) {
	if !enabled {
		log.Info("静态文件服务已禁用")
		return
	}

	// 检查前端目录是否存在
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		log.Warnf("前端目录不存在: %s，请先构建前端项目", frontendDir)
		return
	}

	// 添加静态文件服务（放在所有API路由之后）
	mux.Handle("/", StaticFileHandler(frontendDir))
	log.Infof("静态文件服务已启用，目录: %s", frontendDir)
}
