package sillytavern

import "github.com/ling/muse/common/constant"

// CharacterCard 角色卡数据结构（兼容V2/V3格式）
// 字段说明参考 SillyTavern Character Card V2/V3 规范
type CharacterCard struct {
	// V3 规范字段
	Spec        string        `json:"spec"`         // 规范名称（如 "chara_card_v3"）
	SpecVersion string        `json:"spec_version"` // 规范版本（如 "3.0"）
	Data        CharacterData `json:"data"`         // 扩展数据

	// 基本信息（兼容 V2，部分字段同时存在于顶层）
	Name        string   `json:"name"`        // 角色名称
	Description string   `json:"description"` // 角色描述/背景故事
	Personality string   `json:"personality"` // 性格特征
	Scenario    string   `json:"scenario"`    // 场景设定
	FirstMes    string   `json:"first_mes"`   // 第一条消息/开场白
	MesExample  string   `json:"mes_example"` // 示例对话
	Tags        []string `json:"tags"`        // 标签

	// 额外的顶层字段
	Fav            bool   `json:"fav"`            // 是否收藏
	Talkativeness  string `json:"talkativeness"`  // 话多程度
	CreateDate     string `json:"create_date"`    // 创建日期
	CreatorComment string `json:"creatorcomment"` // 创作者评论
	Avatar         string `json:"avatar"`         // 头像
}

// CharacterData 扩展数据信息
type CharacterData struct {
	// 基本信息
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 角色描述/背景故事
	Personality string `json:"personality"` // 性格特征
	Scenario    string `json:"scenario"`    // 场景设定
	FirstMes    string `json:"first_mes"`   // 第一条消息/开场白
	MesExample  string `json:"mes_example"` // 示例对话

	CreatorNotes       string        `json:"creator_notes"`             // 创作者备注
	SystemPrompt       string        `json:"system_prompt"`             // 系统提示词,v2
	PostHistoryInstr   string        `json:"post_history_instructions"` // 历史记录后指令,v2
	AlternateGreetings []string      `json:"alternate_greetings"`       // 备选开场白
	GroupOnlyGreetings []string      `json:"group_only_greetings"`      // 仅群聊使用开场白
	Tags               []string      `json:"tags"`                      // 标签
	Creator            string        `json:"creator"`                   // 创作者
	CharacterVersion   string        `json:"character_version"`         // 角色版本
	CharacterBook      CharacterBook `json:"character_book,omitempty"`  // 关联的世界书（角色卡内嵌格式）

	// 扩展字段（包含角色内嵌的正则脚本, 世界书等）
	Extensions DataExtensions `json:"extensions,omitempty"`
}

// DataExtensions data字段中的扩展数据
type DataExtensions struct {
	Fav           bool          `json:"fav"`                     // 是否收藏
	Talkativeness string        `json:"talkativeness"`           // 话多程度
	World         string        `json:"world"`                   // 世界名称
	DepthPrompt   DepthPrompt   `json:"depth_prompt"`            // 深度提示
	RegexScripts  []RegexScript `json:"regex_scripts,omitempty"` // 角色内嵌的正则脚本（Scoped Scripts）
}

// WorldBook 世界书/知识库结构（SillyTavern 独立导出格式）
// 用于导入导出 SillyTavern 世界书 JSON 文件
// 注意：独立导出的世界书文件使用 map 格式（以 UID 为 key）
type WorldBook struct {
	Name              string               `json:"name,omitempty"`         // 世界书名称
	Description       string               `json:"description,omitempty"`  // 描述
	ScanDepth         int                  `json:"scan_depth,omitempty"`   // 扫描深度
	TokenBudget       int                  `json:"token_budget,omitempty"` // Token预算
	RecursiveScanning bool                 `json:"recursive_scanning"`     // 递归扫描
	Entries           map[string]BookEntry `json:"entries"`                // 条目列表（以 UID 为 key）
}

// BookEntry 世界书条目（SillyTavern 独立导出格式）
// 字段说明参考 SillyTavern newWorldInfoEntryDefinition
type BookEntry struct {
	// 基础字段
	UID            int      `json:"uid"`                      // 唯一标识符
	Key            []string `json:"key"`                      // 触发关键词（主键）
	KeySecondary   []string `json:"keysecondary,omitempty"`   // 次级关键词
	Content        string   `json:"content"`                  // 条目内容
	Comment        string   `json:"comment,omitempty"`        // 备注/名称
	Disable        bool     `json:"disable"`                  // 是否禁用
	Constant       bool     `json:"constant"`                 // 常驻条目
	Selective      bool     `json:"selective"`                // 选择性触发
	SelectiveLogic int      `json:"selectiveLogic,omitempty"` // 选择性逻辑（0=AND_ANY, 1=NOT_ALL, 2=NOT_ANY, 3=AND_ALL）

	// 位置和排序
	Order    int `json:"order"`    // 插入顺序
	Position int `json:"position"` // 插入位置（0=before_char, 1=after_char, 2=before_desc, 3=after_desc, 4=at_depth）
	Depth    int `json:"depth"`    // 深度（当 position=4 时使用）

	// 扩展字段
	AddMemo        bool   `json:"addMemo,omitempty"`        // 添加备注到内容
	Group          string `json:"group,omitempty"`          // 分组
	GroupOverride  bool   `json:"groupOverride,omitempty"`  // 分组覆盖
	GroupWeight    int    `json:"groupWeight,omitempty"`    // 分组权重
	Probability    int    `json:"probability,omitempty"`    // 触发概率
	UseProbability bool   `json:"useProbability,omitempty"` // 使用概率
	DisplayIndex   int    `json:"displayIndex,omitempty"`   // 显示索引

	// 递归控制
	ExcludeRecursion    bool `json:"excludeRecursion,omitempty"`    // 排除递归
	PreventRecursion    bool `json:"preventRecursion,omitempty"`    // 阻止递归
	DelayUntilRecursion int  `json:"delayUntilRecursion,omitempty"` // 延迟到递归

	// 匹配选项
	CaseSensitive   bool `json:"caseSensitive,omitempty"`   // 大小写敏感
	MatchWholeWords bool `json:"matchWholeWords,omitempty"` // 匹配整词

	// 高级匹配
	MatchPersonaDescription   bool `json:"matchPersonaDescription,omitempty"`
	MatchCharacterDescription bool `json:"matchCharacterDescription,omitempty"`
	MatchCharacterPersonality bool `json:"matchCharacterPersonality,omitempty"`
	MatchCharacterDepthPrompt bool `json:"matchCharacterDepthPrompt,omitempty"`
	MatchScenario             bool `json:"matchScenario,omitempty"`
	MatchCreatorNotes         bool `json:"matchCreatorNotes,omitempty"`
}

// CharacterBook 角色卡内嵌的世界书结构（数组格式）
// 符合 Character Card V2/V3 规范
type CharacterBook struct {
	Name    string               `json:"name,omitempty"` // 世界书名称
	Entries []CharacterBookEntry `json:"entries"`        // 条目列表（数组格式）
}

// CharacterBookEntry 角色卡内嵌世界书的条目（数组格式）
// 符合 Character Card V2/V3 规范中的 CharacterBookEntry 定义
type CharacterBookEntry struct {
	ID             int                          `json:"id,omitempty"`             // 条目ID
	Keys           []string                     `json:"keys"`                     // 触发关键词
	SecondaryKeys  []string                     `json:"secondary_keys,omitempty"` // 次级关键词
	Content        string                       `json:"content"`                  // 条目内容
	Comment        string                       `json:"comment,omitempty"`        // 备注
	Enabled        bool                         `json:"enabled"`                  // 是否启用
	InsertionOrder int                          `json:"insertion_order"`          // 插入顺序
	Selective      bool                         `json:"selective,omitempty"`      // 选择性触发
	Constant       bool                         `json:"constant,omitempty"`       // 常驻条目
	Position       string                       `json:"position,omitempty"`       // 插入位置,v2,只支持 "before_char", "after_char"
	UseRegex       bool                         `json:"use_regex,omitempty"`      // 使用正则表达式
	Extensions     CharacterBookEntryExtensions `json:"extensions,omitempty"`     // 条目扩展字段
}

// CharacterBookEntryExtensions 世界书条目扩展字段
type CharacterBookEntryExtensions struct {
	// 插入位置（0=before, 1=after, 2=ANTop, 3=ANBottom, 4=atDepth, 5=EMTop, 6=EMBottom, 7=outlet）
	Position                  int      `json:"position"`
	ExcludeRecursion          bool     `json:"exclude_recursion"`            // 排除递归
	DisplayIndex              int      `json:"display_index"`                // 显示索引
	Probability               int      `json:"probability"`                  // 触发概率
	UseProbability            bool     `json:"use_probability"`              // 使用概率
	Depth                     int      `json:"depth"`                        // 深度
	SelectiveLogic            int      `json:"selective_logic"`              // 选择性逻辑（0=AND_ANY, 1=NOT_ALL, 2=NOT_ANY, 3=AND_ALL）
	Group                     string   `json:"group"`                        // 分组
	GroupOverride             bool     `json:"group_override"`               // 分组覆盖
	GroupWeight               int      `json:"group_weight"`                 // 分组权重
	PreventRecursion          bool     `json:"prevent_recursion"`            // 阻止递归
	DelayUntilRecursion       bool     `json:"delay_until_recursion"`        // 延迟到递归
	ScanDepth                 *int     `json:"scan_depth"`                   // 扫描深度（可为null）
	MatchWholeWords           *bool    `json:"match_whole_words"`            // 匹配整词（可为null）
	UseGroupScoring           bool     `json:"use_group_scoring"`            // 使用分组评分
	CaseSensitive             *bool    `json:"case_sensitive"`               // 大小写敏感（可为null）
	MatchPersonaDescription   bool     `json:"match_persona_description"`    // 匹配角色描述
	MatchCharacterDescription bool     `json:"match_character_description"`  // 匹配角色描述
	MatchCharacterPersonality bool     `json:"match_character_personality"`  // 匹配角色性格
	MatchCharacterDepthPrompt bool     `json:"match_character_depth_prompt"` // 匹配深度提示
	MatchScenario             bool     `json:"match_scenario"`               // 匹配场景
	MatchCreatorNotes         bool     `json:"match_creator_notes"`          // 匹配创作者备注
	Triggers                  []string `json:"triggers"`                     // 触发器
	IgnoreBudget              bool     `json:"ignore_budget"`                // 忽略预算
	AutomationID              string   `json:"automation_id"`                // 自动化ID
	Role                      int      `json:"role"`                         // 角色, 当 position=atDepth 时使用
	Vectorized                bool     `json:"vectorized"`                   // 向量化
	Sticky                    int      `json:"sticky"`                       // 粘性
	Cooldown                  int      `json:"cooldown"`                     // 冷却时间
	Delay                     int      `json:"delay"`                        // 延迟
	OutletName                string   `json:"outlet_name"`                  // 出口名称
}

// Asset 资源文件（V3格式）
type Asset struct {
	Type string `json:"type"` // 资源类型：icon, background, user_icon等
	URI  string `json:"uri"`  // 资源URI（可以是data URI或外部链接）
	Name string `json:"name"` // 资源名称
}

// DepthPrompt 深度提示配置
type DepthPrompt struct {
	Prompt string `json:"prompt"` // 提示内容
	Depth  int    `json:"depth"`  // 插入深度
	Role   string `json:"role"`   // 角色：system, user, assistant
}

// OpenAIPreset SillyTavern 导出的 OpenAI 预设 JSON 文件结构
// 字段说明参考 SillyTavern settingsToUpdate 定义
type OpenAIPreset struct {
	// ============ API 来源和模型 ============
	ChatCompletionSource string `json:"chat_completion_source"` // API 来源: openai/claude/google 等
	OpenAIModel          string `json:"openai_model"`           // OpenAI 模型
	ClaudeModel          string `json:"claude_model"`           // Claude 模型
	GoogleModel          string `json:"google_model"`           // Google 模型
	VertexAIModel        string `json:"vertexai_model"`         // Vertex AI 模型
	OpenRouterModel      string `json:"openrouter_model"`       // OpenRouter 模型
	AI21Model            string `json:"ai21_model"`             // AI21 模型
	MistralAIModel       string `json:"mistralai_model"`        // Mistral AI 模型
	CohereModel          string `json:"cohere_model"`           // Cohere 模型
	PerplexityModel      string `json:"perplexity_model"`       // Perplexity 模型
	GroqModel            string `json:"groq_model"`             // Groq 模型
	ChutesModel          string `json:"chutes_model"`           // Chutes 模型
	SiliconFlowModel     string `json:"siliconflow_model"`      // SiliconFlow 模型
	ElectronHubModel     string `json:"electronhub_model"`      // ElectronHub 模型
	NanoGPTModel         string `json:"nanogpt_model"`          // NanoGPT 模型
	DeepSeekModel        string `json:"deepseek_model"`         // DeepSeek 模型
	AIMLAPIModel         string `json:"aimlapi_model"`          // AIMLAPI 模型
	XAIModel             string `json:"xai_model"`              // xAI 模型
	PollinationsModel    string `json:"pollinations_model"`     // Pollinations 模型
	MoonshotModel        string `json:"moonshot_model"`         // Moonshot 模型
	FireworksModel       string `json:"fireworks_model"`        // Fireworks 模型
	CometAPIModel        string `json:"cometapi_model"`         // CometAPI 模型
	ZAIModel             string `json:"zai_model"`              // ZAI 模型
	CustomModel          string `json:"custom_model"`           // 自定义模型 ID

	// ============ 基础生成参数 ============
	Temperature       float32 `json:"temperature"`        // 温度
	FrequencyPenalty  float32 `json:"frequency_penalty"`  // 频率惩罚
	PresencePenalty   float32 `json:"presence_penalty"`   // 存在惩罚
	TopP              float32 `json:"top_p"`              // Top P
	TopK              int     `json:"top_k"`              // Top K
	TopA              float32 `json:"top_a"`              // Top A
	MinP              float32 `json:"min_p"`              // Min P
	RepetitionPenalty float32 `json:"repetition_penalty"` // 重复惩罚
	Seed              int     `json:"seed"`               // 随机种子 (-1 表示随机)
	N                 int     `json:"n"`                  // 生成数量

	// ============ 上下文和 Token ============
	OpenAIMaxContext   int  `json:"openai_max_context"`   // 最大上下文长度
	OpenAIMaxTokens    int  `json:"openai_max_tokens"`    // 最大生成 Token 数
	MaxContextUnlocked bool `json:"max_context_unlocked"` // 解锁最大上下文

	// ============ 自定义 API ============
	CustomURL                  string `json:"custom_url"`                    // 自定义 API URL
	CustomIncludeBody          string `json:"custom_include_body"`           // 自定义请求体包含字段
	CustomExcludeBody          string `json:"custom_exclude_body"`           // 自定义请求体排除字段
	CustomIncludeHeaders       string `json:"custom_include_headers"`        // 自定义请求头
	CustomPromptPostProcessing string `json:"custom_prompt_post_processing"` // 提示词后处理

	// ============ 代理设置（敏感字段）============
	ReverseProxy  string `json:"reverse_proxy"`  // 反向代理地址
	ProxyPassword string `json:"proxy_password"` // 代理密码

	// ============ 提示词模板 ============
	NamesBehavior        int    `json:"names_behavior"`          // 名称行为
	SendIfEmpty          string `json:"send_if_empty"`           // 空消息时发送的内容
	ImpersonationPrompt  string `json:"impersonation_prompt"`    // 扮演提示词
	NewChatPrompt        string `json:"new_chat_prompt"`         // 新对话提示词
	NewGroupChatPrompt   string `json:"new_group_chat_prompt"`   // 新群聊提示词
	NewExampleChatPrompt string `json:"new_example_chat_prompt"` // 新示例对话提示词
	ContinueNudgePrompt  string `json:"continue_nudge_prompt"`   // 继续提示词
	GroupNudgePrompt     string `json:"group_nudge_prompt"`      // 群组提示词
	WIFormat             string `json:"wi_format"`               // 世界书格式
	ScenarioFormat       string `json:"scenario_format"`         // 场景格式
	PersonalityFormat    string `json:"personality_format"`      // 性格格式

	// ============ 提示词管理器 ============
	Prompts            []PresetPromptItem `json:"prompts"`              // 提示词数组
	PromptOrder        []PromptOrderItem  `json:"prompt_order"`         // 提示词顺序
	BiasPresetSelected string             `json:"bias_preset_selected"` // 选中的 Logit Bias 预设

	// ============ Claude 专用 ============
	AssistantPrefill       string `json:"assistant_prefill"`       // 助手预填充
	AssistantImpersonation string `json:"assistant_impersonation"` // 助手扮演
	UseSysprompt           bool   `json:"use_sysprompt"`           // 使用系统提示词

	// ============ 功能开关 ============
	StreamOpenAI         bool   `json:"stream_openai"`          // 流式输出
	SquashSystemMessages bool   `json:"squash_system_messages"` // 合并系统消息
	MediaInlining        bool   `json:"media_inlining"`         // 媒体内联
	InlineImageQuality   string `json:"inline_image_quality"`   // 内联图片质量
	ContinuePrefill      bool   `json:"continue_prefill"`       // 继续预填充
	ContinuePostfix      string `json:"continue_postfix"`       // 继续后缀
	FunctionCalling      bool   `json:"function_calling"`       // 函数调用
	ShowThoughts         bool   `json:"show_thoughts"`          // 显示思考过程
	ReasoningEffort      string `json:"reasoning_effort"`       // 推理努力程度
	Verbosity            string `json:"verbosity"`              // 详细程度
	EnableWebSearch      bool   `json:"enable_web_search"`      // 启用网络搜索
	BypassStatusCheck    bool   `json:"bypass_status_check"`    // 跳过状态检查
	ShowExternalModels   bool   `json:"show_external_models"`   // 显示外部模型

	// ============ 扩展字段 ============
	Extensions PresetExtensions `json:"extensions,omitempty"` // 扩展字段，包含预设内嵌的正则脚本等
}

// PresetExtensions 预设扩展字段
// 用于存储预设内嵌的扩展数据，如正则脚本
type PresetExtensions struct {
	RegexScripts []RegexScript `json:"regex_scripts,omitempty"` // 预设内嵌的正则脚本
}

// PresetPromptItem 预设中的提示词项
type PresetPromptItem struct {
	Identifier        string               `json:"identifier"`                   // 唯一标识符
	Name              string               `json:"name"`                         // 显示名称
	Content           string               `json:"content,omitempty"`            // 提示词内容
	Role              string               `json:"role,omitempty"`               // 角色：system, user, assistant
	SystemPrompt      bool                 `json:"system_prompt,omitempty"`      // 是否为系统提示词
	Marker            bool                 `json:"marker,omitempty"`             // 是否为标记（占位符）
	InjectionPosition constant.STInjectPos `json:"injection_position,omitempty"` // 注入位置
	InjectionDepth    int                  `json:"injection_depth,omitempty"`    // 注入深度
	InjectionOrder    int                  `json:"injection_order,omitempty"`    // 注入顺序
	InjectionTrigger  []string             `json:"injection_trigger,omitempty"`  // 注入触发器
	ForbidOverrides   bool                 `json:"forbid_overrides,omitempty"`   // 禁止覆盖
	Enabled           bool                 `json:"enabled,omitempty"`            // 是否启用
}

// PromptOrderItem 提示词顺序项
type PromptOrderItem struct {
	CharacterID int                     `json:"character_id"` // 角色ID（100000=默认，100001=群聊）
	Order       []PromptOrderIdentifier `json:"order"`        // 顺序列表
}

// PromptOrderIdentifier 提示词顺序中的标识符
type PromptOrderIdentifier struct {
	Identifier string `json:"identifier"` // 提示词标识符
	Enabled    bool   `json:"enabled"`    // 是否启用
}

// RegexScript SillyTavern 正则脚本数据结构
// 用于解析 SillyTavern 导出的正则规则 JSON 文件
type RegexScript struct {
	ID              string   `json:"id"`              // UUID 标识符
	ScriptName      string   `json:"scriptName"`      // 脚本名称
	FindRegex       string   `json:"findRegex"`       // 查找正则表达式
	ReplaceString   string   `json:"replaceString"`   // 替换字符串
	TrimStrings     []string `json:"trimStrings"`     // 裁剪字符串列表
	Placement       []int    `json:"placement"`       // 作用位置（0=用户输入, 1=AI输出, 2=斜杠命令, 3=世界书, 4=提示词）
	Disabled        bool     `json:"disabled"`        // 是否禁用
	MarkdownOnly    bool     `json:"markdownOnly"`    // 仅应用于 Markdown 显示
	PromptOnly      bool     `json:"promptOnly"`      // 仅应用于提示词
	RunOnEdit       bool     `json:"runOnEdit"`       // 编辑时运行
	SubstituteRegex int      `json:"substituteRegex"` // 替换正则模式（0=无, 1={{user}}, 2={{char}}, 3=全部）
	MinDepth        int      `json:"minDepth"`        // 最小深度
	MaxDepth        int      `json:"maxDepth"`        // 最大深度
}
