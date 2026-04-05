package errs

import (
	"errors"
	"fmt"

	"connectrpc.com/connect"
)

const (
	// General
	DataConflict  = "数据已被修改，请刷新后重试"
	EmptySortData = "排序数据不能为空"

	// Character
	InvalidCharacterID = "无效的角色ID"
	EmptyCharacterName = "角色名称不能为空"
	CharacterNotFound  = "角色不存在"

	// Chat
	InvalidSessionID = "无效的会话ID"
	EmptySessionName = "会话名称不能为空"
	SessionNotFound  = "会话不存在"

	// Preset
	InvalidPresetID     = "无效的预设ID"
	EmptyPresetName     = "预设名称不能为空"
	PresetNotFound      = "预设不存在"
	InvalidPromptItemID = "无效的提示项ID"
	EmptyPromptItemName = "提示项名称不能为空"
	PromptItemNotFound  = "提示项不存在"

	// Regex Rule
	InvalidRegexRuleID = "无效的规则ID"
	EmptyRegexRuleName = "规则名称不能为空"
	RegexRuleNotFound  = "规则不存在"
	EmptyFindPattern   = "查找模式不能为空"
	EmptyAffectFlags   = "影响标志不能为空"

	// World Info
	InvalidWorldInfoID      = "无效的世界书ID"
	EmptyWorldInfoName      = "世界书名称不能为空"
	WorldInfoNotFound       = "世界书不存在"
	InvalidWorldInfoEntryID = "无效的条目ID"
	EmptyKeysList           = "关键词列表不能为空"
	EmptyContent            = "内容不能为空"
	WorldInfoEntryNotFound  = "条目不存在"
	EmptyFileContent        = "文件内容不能为空"
	InvalidWorldInfoFile    = "无效的世界书文件格式"

	// User
	InvalidUserID      = "无效的用户ID"
	EmptyUsername      = "用户名不能为空"
	EmptyPassword      = "密码不能为空"
	UserNotFound       = "用户不存在"
	UserAlreadyExists  = "用户已存在"
	InvalidPassword    = "密码错误"
	InvalidOldPassword = "原密码错误"

	// Persona
	InvalidPersonaID = "无效的人设ID"
	EmptyPersonaName = "人设名称不能为空"
	PersonaNotFound  = "人设不存在"

	// API Config
	InvalidAPIConfigID = "无效的API配置ID"
	APIConfigNotFound  = "API配置不存在"
	EmptyAPIKey        = "API Key不能为空"
	InvalidAPIProvider = "无效的API提供商"
)

// New 创建错误
func New(code int32, msg string) error {
	return connect.NewError(connect.Code(code), fmt.Errorf(msg))
}

// Newf 格式化创建错误
func Newf(code int32, format string, args ...any) error {
	return connect.NewError(connect.Code(code), fmt.Errorf(format, args...))
}

// NewStandard 创建标准错误
func NewStandard(code connect.Code, msg string) error {
	return connect.NewError(code, fmt.Errorf(msg))
}

// NewStandardf 格式化创建标准错误
func NewStandardf(code connect.Code, format string, args ...any) error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}

// Code 获取错误码
func Code(err error) connect.Code {
	var e *connect.Error
	if errors.As(err, &e) {
		return e.Code()
	}
	return connect.CodeUnknown
}
