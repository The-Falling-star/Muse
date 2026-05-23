package convert

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	log "github.com/sirupsen/logrus"
)

type JSRegex struct {
	Engine   *regexp2.Regexp
	IsGlobal bool
}

// ConvertToRegexp2 智能转换器：完美兼容 JS风格正则 和 Go原生正则
func ConvertToRegexp2(regexStr string) (*JSRegex, error) {
	pattern := regexStr
	flags := ""
	isJSFormat := false

	// 1. 尝试识别 JS 格式 (/pattern/flags)
	if strings.HasPrefix(regexStr, "/") {
		lastSlashIdx := strings.LastIndex(regexStr, "/")
		if lastSlashIdx > 0 { // 确保有超过一个 '/'
			// 提取可能是修饰符的部分
			possibleFlags := regexStr[lastSlashIdx+1:]

			// 验证 possibleFlags 是否只包含合法的 JS 修饰符
			if isValidJSFlags(possibleFlags) {
				pattern = regexStr[1:lastSlashIdx]
				flags = possibleFlags
				isJSFormat = true
			}
			// 如果验证不通过，说明这不是 JS 包装格式，回退为普通正则 (isJSFormat 保持 false)
		}
	}

	// 2. 解析修饰符 (只有确认为 JS 格式时才执行)
	options := regexp2.ECMAScript
	isGlobal := false

	if isJSFormat {
		for _, flag := range flags {
			switch flag {
			case 'i':
				options |= regexp2.IgnoreCase
			case 'm':
				options |= regexp2.Multiline
			case 's':
				options |= regexp2.Singleline
			case 'g':
				isGlobal = true
				// u (unicode), y (sticky), d 等高级修饰符在此忽略或支持
			}
		}
	}

	// 3. 编译正则
	log.Debugf("js正则: %s 转换为go正则: %s", regexStr, pattern)
	re, err := regexp2.Compile(pattern, regexp2.RegexOptions(options))
	if err != nil {
		return nil, fmt.Errorf("正则编译失败: %w", err)
	}

	return &JSRegex{
		Engine:   re,
		IsGlobal: isGlobal,
	}, nil
}

// isValidJSFlags 检查字符串是否仅由合法的 JS 正则修饰符组成
func isValidJSFlags(flags string) bool {
	// JS 支持的合法修饰符
	validChars := "gimsuy"
	for _, char := range flags {
		if !strings.ContainsRune(validChars, char) {
			return false // 发现非法字符，绝对不是 JS 修饰符
		}
	}
	return true
}
