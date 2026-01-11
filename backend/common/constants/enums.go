// Package constants 定义项目中使用的常量和枚举类型
package constants

// MessageRole 消息角色枚举
type MessageRole int

const (
	// MessageRoleSystem 系统消息
	MessageRoleSystem MessageRole = iota + 1
	// MessageRoleUser 用户消息
	MessageRoleUser
	// MessageRoleAssistant AI助手消息
	MessageRoleAssistant
)

// String 返回消息角色的字符串表示
func (r MessageRole) String() string {
	switch r {
	case MessageRoleSystem:
		return "system"
	case MessageRoleUser:
		return "user"
	case MessageRoleAssistant:
		return "assistant"
	default:
		return "unknown"
	}
}

// PromptItemRole 提示项角色枚举
type PromptItemRole int

const (
	// PromptItemRoleSystem 系统角色
	PromptItemRoleSystem PromptItemRole = iota + 1
	// PromptItemRoleUser 用户角色
	PromptItemRoleUser
	// PromptItemRoleAssistant AI助手角色
	PromptItemRoleAssistant
)

// String 返回提示项角色的字符串表示
func (r PromptItemRole) String() string {
	switch r {
	case PromptItemRoleSystem:
		return "system"
	case PromptItemRoleUser:
		return "user"
	case PromptItemRoleAssistant:
		return "assistant"
	default:
		return "unknown"
	}
}

// InjectionPosition 注入位置枚举
type InjectionPosition int

const (
	// InjectionPositionRelative 相对位置
	InjectionPositionRelative InjectionPosition = iota
	// InjectionPositionAbsolute 绝对位置
	InjectionPositionAbsolute
)

// EntryPosition 世界书条目插入位置枚举
type EntryPosition int

const (
	// EntryPositionBeforeChar 角色定义之前
	EntryPositionBeforeChar EntryPosition = iota + 1
	// EntryPositionAfterChar 角色定义之后
	EntryPositionAfterChar
	// EntryPositionBeforeExample 示例对话之前
	EntryPositionBeforeExample
	// EntryPositionAfterExample 示例对话之后
	EntryPositionAfterExample
	// EntryPositionAtDepth 指定深度
	EntryPositionAtDepth
)

// String 返回插入位置的字符串表示
func (p EntryPosition) String() string {
	switch p {
	case EntryPositionBeforeChar:
		return "before_char"
	case EntryPositionAfterChar:
		return "after_char"
	case EntryPositionBeforeExample:
		return "before_example"
	case EntryPositionAfterExample:
		return "after_example"
	case EntryPositionAtDepth:
		return "at_depth"
	default:
		return "unknown"
	}
}

// APIProvider API提供商枚举
type APIProvider int

const (
	// APIProviderOpenAI OpenAI兼容协议
	APIProviderOpenAI APIProvider = iota + 1
	// APIProviderClaude Anthropic Claude
	APIProviderClaude
	// APIProviderGemini Google Gemini
	APIProviderGemini
)

// String 返回API提供商的字符串表示
func (p APIProvider) String() string {
	switch p {
	case APIProviderOpenAI:
		return "openai"
	case APIProviderClaude:
		return "claude"
	case APIProviderGemini:
		return "gemini"
	default:
		return "unknown"
	}
}

// Theme 主题枚举
type Theme uint

const (
	// ThemeAuto 跟随系统
	ThemeAuto Theme = iota
	// ThemeLight 亮色主题
	ThemeLight
	// ThemeDark 暗色主题
	ThemeDark
)

// String 返回主题的字符串表示
func (t Theme) String() string {
	switch t {
	case ThemeAuto:
		return "auto"
	case ThemeLight:
		return "light"
	case ThemeDark:
		return "dark"
	default:
		return "unknown"
	}
}

// RegexAffectFlags 正则规则作用范围位掩码
type RegexAffectFlags uint

const (
	// RegexAffectUserInput 作用于用户输入
	RegexAffectUserInput RegexAffectFlags = 1 << iota
	// RegexAffectAIOutput 作用于AI输出
	RegexAffectAIOutput
	// RegexAffectSlashCommand 作用于斜杠命令
	RegexAffectSlashCommand
	// RegexAffectWorldInfo 作用于世界书
	RegexAffectWorldInfo
	// RegexAffectPrompt 作用于提示词
	RegexAffectPrompt
)

// Has 检查是否包含指定的标志
func (f RegexAffectFlags) Has(flag RegexAffectFlags) bool {
	return f&flag != 0
}

// Add 添加指定的标志
func (f RegexAffectFlags) Add(flag RegexAffectFlags) RegexAffectFlags {
	return f | flag
}

// Remove 移除指定的标志
func (f RegexAffectFlags) Remove(flag RegexAffectFlags) RegexAffectFlags {
	return f &^ flag
}
