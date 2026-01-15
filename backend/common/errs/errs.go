package errs

import (
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
)

// New 创建错误
func New(code int32, msg string) *connect.Error {
	return connect.NewError(connect.Code(code), fmt.Errorf(msg))
}

// Newf 格式化创建错误
func Newf(code int32, format string, args ...any) *connect.Error {
	return connect.NewError(connect.Code(code), fmt.Errorf(format, args...))
}

// NewStandard 创建标准错误
func NewStandard(code connect.Code, msg string) *connect.Error {
	return connect.NewError(code, fmt.Errorf(msg))
}

// NewStandardf 格式化创建标准错误
func NewStandardf(code connect.Code, format string, args ...any) *connect.Error {
	return connect.NewError(code, fmt.Errorf(format, args...))
}
