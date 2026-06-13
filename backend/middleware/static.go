// Package middleware 静态文件服务中间件
package middleware

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// embeddedFS 全局内嵌文件系统实例，由 RegisterEmbeddedFS 设置
var embeddedFS fs.FS

// RegisterEmbeddedFS 注册内嵌文件系统
// 在 main 包中调用，将 embed.FS 传入
func RegisterEmbeddedFS(fsys fs.FS) {
	embeddedFS = fsys
}

// hasEmbeddedFrontend 检查内嵌文件系统中是否包含前端资源
func hasEmbeddedFrontend() bool {
	if embeddedFS == nil {
		return false
	}
	entries, err := fs.ReadDir(embeddedFS, "static")
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// ServeStaticFiles 条件性地添加静态文件服务到路由器
// 优先使用内嵌前端资源，如果没有则降级到磁盘路径
func ServeStaticFiles(mux *http.ServeMux, enabled bool, frontendDir string) {
	if !enabled {
		log.Info("静态文件服务已禁用")
		return
	}

	// 优先使用内嵌前端资源
	if hasEmbeddedFrontend() {
		sub, err := fs.Sub(embeddedFS, "static")
		if err != nil {
			log.Errorf("创建内嵌文件子系统失败: %v", err)
			return
		}
		mux.Handle("/", spaHandler{fileServer: http.FileServer(http.FS(sub)), staticFS: sub})
		log.Info("静态文件服务已启用（内嵌模式）")
		return
	}

	// 降级到磁盘路径
	if _, err := os.Stat(frontendDir); os.IsNotExist(err) {
		log.Warnf("前端目录不存在: %s，请先构建前端项目", frontendDir)
		return
	}

	mux.Handle("/", diskSPAHandler(frontendDir))
	log.Infof("静态文件服务已启用，目录: %s", frontendDir)
}

// spaHandler 基于 embed.FS 的 SPA 路由处理器
type spaHandler struct {
	fileServer http.Handler
	staticFS   fs.FS
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")

	// API 路径跳过
	if strings.HasPrefix(r.URL.Path, "/muse.") {
		http.NotFound(w, r)
		return
	}

	// 检查内嵌文件系统中是否存在该文件
	if path != "" {
		if _, err := fs.Stat(h.staticFS, path); err == nil {
			h.fileServer.ServeHTTP(w, r)
			return
		}
	}

	// SPA 路由支持：不存在的路径返回 index.html
	r.URL.Path = "/"
	h.fileServer.ServeHTTP(w, r)
}

// diskSPAHandler 基于磁盘的 SPA 路由处理器
func diskSPAHandler(frontendDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// API 路径跳过
		if strings.HasPrefix(r.URL.Path, "/muse.") {
			http.NotFound(w, r)
			return
		}

		// 构建文件路径
		path := filepath.Join(frontendDir, r.URL.Path)

		// 检查文件是否存在
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// SPA 路由支持 - 不存在的文件都返回 index.html
			http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
			return
		}

		// 存在的文件直接提供服务
		http.FileServer(http.Dir(frontendDir)).ServeHTTP(w, r)
	})
}
