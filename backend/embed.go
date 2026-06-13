//go:build embed

package main

import "embed"

// FrontendDist 嵌入前端构建产物
// 仅在使用 -tags embed 编译时启用
// 使用 compile.sh 脚本可自动完成前端构建和后端编译
//
//go:embed all:static
var FrontendDist embed.FS
