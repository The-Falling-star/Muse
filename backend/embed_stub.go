//go:build !embed

package main

import "io/fs"

// FrontendDist 开发模式下无内嵌前端资源
// 使用 -tags embed 编译时才会内嵌前端
var FrontendDist fs.FS
