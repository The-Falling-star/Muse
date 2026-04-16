package constant

const (
	// DefaultPageNum 默认页码
	DefaultPageNum = 1
	// DefaultPageSize 默认页大小
	DefaultPageSize = 20
	// MaxPageSize 最大页码
	MaxPageSize = 100
)

const (
	// DefaultTemperature 默认温度
	DefaultTemperature = 1.0
	// DefaultTopK 默认TopK
	DefaultTopK = 10
	// DefaultTopP 默认TopP
	DefaultTopP = 1.0
	// DefaultMaxTokens 默认最大令牌数
	DefaultMaxTokens = 2048
	// DefaultCandidateCount 默认候选词数
	DefaultCandidateCount = 1
)

// CtxKey 上下文键
type CtxKey string

const (
	// TransactionKey 获取ctx中的事务键
	TransactionKey CtxKey = "transaction"
	// UserIDKey 获取ctx中的用户ID键
	UserIDKey CtxKey = "userId"
)

// ModuleType 模块类型
type ModuleType int

const (
	// UserInput 用户输入
	UserInput ModuleType = iota
	// AIOutput AI输出
	AIOutput
	// Preset 预设
	Preset
	// WorldInfo 世界信息
	WorldInfo
)

// STInjectPos 酒馆提示词注入位置
type STInjectPos int

const (
	// STInjectPosRelative 相对注入
	STInjectPosRelative STInjectPos = iota
	// STInjectPosAbsolute 绝对注入
	STInjectPosAbsolute
)

// NormalizePagination 处理分页参数，返回规范化后的 page 和 pageSize
func NormalizePagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = DefaultPageNum
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return page, pageSize
}

type Macro string

const (
	// 基础/用户/角色宏
	User1           = "{{user}}"
	User2           = "<user>"
	Char1           = "{{char}}"
	Char2           = "<char>"
	Char3           = "<bot>"
	Char4           = "<BOT>"
	LastMessage     = "{{lastMessage}}"
	LastUserMessage = "{{lastUserMessage}}"
	LastCharMessage = "{{lastCharMessage}}"
	Time            = "{{time}}"
	Date            = "{{date}}"
	Input           = "{{input}}"

	// 新增宏（按表格顺序）
	Pipe               = "{{pipe}}"
	Newline            = "{{newline}}"
	Trim               = "{{trim}}"
	Noop               = "{{noop}}"
	CharPrompt         = "{{charPrompt}}"
	CharJailbreak      = "{{charJailbreak}}"
	Group              = "{{group}}"
	CharIfNotGroup     = "{{charIfNotGroup}}"
	GroupNotMuted      = "{{groupNotMuted}}"
	Description        = "{{description}}"
	Scenario           = "{{scenario}}"
	Personality        = "{{personality}}"
	Persona            = "{{persona}}"
	MesExamples        = "{{mesExamples}}"
	CharVersion        = "{{char_version}}"
	Model              = "{{model}}"
	LastMessageId      = "{{lastMessageId}}"
	FirstIncludedMsgId = "{{firstIncludedMessageId}}"
	CurrentSwipeId     = "{{currentSwipeId}}"
	LastSwipeId        = "{{lastSwipeId}}"
	LastGenerationType = "{{lastGenerationType}}"
	Original           = "{{original}}"
	TimeDiff           = "{{timeDiff}}" // 使用时需带参数，如 {{timeDiff::time1::time2}}
	Weekday            = "{{weekday}}"
	IsoTime            = "{{isotime}}"
	IsoDate            = "{{isodate}}"
	DateTimeFormat     = "{{datetimeformat}}" // 使用时带格式，如 {{datetimeformat DD.MM.YYYY HH:mm}}
	IdleDuration       = "{{idle_duration}}"
	Random             = "{{random}}" // 支持 :(args) 或 ::arg1::arg2 语法
	Pick               = "{{pick}}"
	Roll               = "{{roll}}"
	Bias               = "{{bias}}"   // 使用时带引号文本，如 {{bias "text"}}
	Note               = "{{//}}"     // 注释宏，AI 不可见
	Banned             = "{{banned}}" // 使用时带引号文本，如 {{banned "text"}}
	Reverse            = "{{reverse}}"

	// 指令模式和上下文模板宏
	ExampleSeparator                = "{{exampleSeparator}}"
	ChatStart                       = "{{chatStart}}"
	InstructSystemPrompt            = "{{instructSystemPrompt}}"
	InstructSystemPromptPrefix      = "{{instructSystemPromptPrefix}}"
	InstructSystemPromptSuffix      = "{{instructSystemPromptSuffix}}"
	InstructUserPrefix              = "{{instructUserPrefix}}"
	InstructAssistantPrefix         = "{{instructAssistantPrefix}}"
	InstructSystemPrefix            = "{{instructSystemPrefix}}"
	InstructUserSuffix              = "{{instructUserSuffix}}"
	InstructAssistantSuffix         = "{{instructAssistantSuffix}}"
	InstructSystemSuffix            = "{{instructSystemSuffix}}"
	InstructFirstAssistantPrefix    = "{{instructFirstAssistantPrefix}}"
	InstructLastAssistantPrefix     = "{{instructLastAssistantPrefix}}"
	InstructFirstUserPrefix         = "{{instructFirstUserPrefix}}"
	InstructLastUserPrefix          = "{{instructLastUserPrefix}}"
	InstructSystemInstructionPrefix = "{{instructSystemInstructionPrefix}}"
	InstructUserFiller              = "{{instructUserFiller}}"
	InstructStop                    = "{{instructStop}}"
	MaxPrompt                       = "{{maxPrompt}}"
	SystemPrompt                    = "{{systemPrompt}}"
	DefaultSystemPrompt             = "{{defaultSystemPrompt}}"

	// 聊天变量宏（本地变量）
	Getvar = "{{getvar}}" // 使用时带 name，如 {{getvar::name}}
	Setvar = "{{setvar}}" // 使用时带 name::value，如 {{setvar::name::value}}
	Addvar = "{{addvar}}" // 使用时带 name::increment
	Incvar = "{{incvar}}" // 使用时带 name
	Decvar = "{{decvar}}" // 使用时带 name
	// 全局变量
	Getglobalvar = "{{getglobalvar}}"
	Setglobalvar = "{{setglobalvar}}"
	Addglobalvar = "{{addglobalvar}}"
	Incglobalvar = "{{incglobalvar}}"
	Decglobalvar = "{{decglobalvar}}"
	Var          = "{{var}}" // 作用域变量，支持 ::name 或 ::name::index

	// 扩展特定宏
	Summary = "{{summary}}"
)

// Value 返回所有已定义的宏常量（完整列表）
func (m Macro) Value() []Macro {
	return []Macro{
		User1, User2, Char1, Char2, Char3, Char4,
		LastMessage, LastUserMessage, LastCharMessage,
		Time, Date, Input,
		Pipe, Newline, Trim, Noop,
		CharPrompt, CharJailbreak,
		Group, CharIfNotGroup, GroupNotMuted,
		Description, Scenario, Personality, Persona, MesExamples,
		CharVersion, Model, LastMessageId, FirstIncludedMsgId,
		CurrentSwipeId, LastSwipeId, LastGenerationType, Original,
		TimeDiff, Weekday, IsoTime, IsoDate, DateTimeFormat, IdleDuration,
		Random, Pick, Roll, Bias, Note, Banned, Reverse,
		ExampleSeparator, ChatStart,
		InstructSystemPrompt, InstructSystemPromptPrefix, InstructSystemPromptSuffix,
		InstructUserPrefix, InstructAssistantPrefix, InstructSystemPrefix,
		InstructUserSuffix, InstructAssistantSuffix, InstructSystemSuffix,
		InstructFirstAssistantPrefix, InstructLastAssistantPrefix,
		InstructFirstUserPrefix, InstructLastUserPrefix,
		InstructSystemInstructionPrefix, InstructUserFiller, InstructStop,
		MaxPrompt, SystemPrompt, DefaultSystemPrompt,
		Getvar, Setvar, Addvar, Incvar, Decvar,
		Getglobalvar, Setglobalvar, Addglobalvar, Incglobalvar, Decglobalvar, Var,
		Summary,
	}
}
