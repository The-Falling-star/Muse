package sillytavern

// CharacterCard 角色卡数据结构（兼容V2/V3格式）
// 字段说明参考 SillyTavern Character Card V2/V3 规范
type CharacterCard struct {
	// 基本信息
	Name        string `json:"name"`        // 角色名称
	Description string `json:"description"` // 角色描述/背景故事
	Personality string `json:"personality"` // 性格特征
	Scenario    string `json:"scenario"`    // 场景设定
	FirstMes    string `json:"first_mes"`   // 第一条消息/开场白
	MesExample  string `json:"mes_example"` // 示例对话

	// 扩展信息
	CreatorNotes       string     `json:"creator_notes"`             // 创作者备注
	SystemPrompt       string     `json:"system_prompt"`             // 系统提示词
	PostHistoryInstr   string     `json:"post_history_instructions"` // 历史记录后指令
	AlternateGreetings []string   `json:"alternate_greetings"`       // 备选开场白
	Tags               []string   `json:"tags"`                      // 标签
	Creator            string     `json:"creator"`                   // 创作者
	CharacterVersion   string     `json:"character_version"`         // 角色版本
	CharacterBook      *WorldBook `json:"character_book,omitempty"`  // 关联的世界书

	// V3扩展字段
	Assets []Asset `json:"assets,omitempty"` // 资源文件列表

	// 深度提示（可在对话特定深度插入）
	DepthPrompt *DepthPrompt `json:"depth_prompt,omitempty"`
}

// WorldBook 世界书/知识库结构
type WorldBook struct {
	Name              string      `json:"name"`               // 世界书名称
	Description       string      `json:"description"`        // 描述
	ScanDepth         int         `json:"scan_depth"`         // 扫描深度
	TokenBudget       int         `json:"token_budget"`       // Token预算
	RecursiveScanning bool        `json:"recursive_scanning"` // 递归扫描
	Entries           []BookEntry `json:"entries"`            // 条目列表
}

// BookEntry 世界书条目
type BookEntry struct {
	Keys           []string `json:"keys"`            // 触发关键词
	SecondaryKeys  []string `json:"secondary_keys"`  // 次级关键词
	Content        string   `json:"content"`         // 条目内容
	Enabled        bool     `json:"enabled"`         // 是否启用
	InsertionOrder int      `json:"insertion_order"` // 插入顺序
	CaseSensitive  bool     `json:"case_sensitive"`  // 大小写敏感
	Name           string   `json:"name"`            // 条目名称
	Priority       int      `json:"priority"`        // 优先级
	Comment        string   `json:"comment"`         // 备注
	Selective      bool     `json:"selective"`       // 选择性触发
	Constant       bool     `json:"constant"`        // 常驻条目
	Position       string   `json:"position"`        // 插入位置
	Depth          int      `json:"depth"`           // 深度（当position为at_depth时）
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
	Temperature       float64 `json:"temperature"`        // 温度
	FrequencyPenalty  float64 `json:"frequency_penalty"`  // 频率惩罚
	PresencePenalty   float64 `json:"presence_penalty"`   // 存在惩罚
	TopP              float64 `json:"top_p"`              // Top P
	TopK              int     `json:"top_k"`              // Top K
	TopA              float64 `json:"top_a"`              // Top A
	MinP              float64 `json:"min_p"`              // Min P
	RepetitionPenalty float64 `json:"repetition_penalty"` // 重复惩罚
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
}

// PresetPromptItem 预设中的提示词项
type PresetPromptItem struct {
	Identifier        string   `json:"identifier"`                   // 唯一标识符
	Name              string   `json:"name"`                         // 显示名称
	Content           string   `json:"content,omitempty"`            // 提示词内容
	Role              string   `json:"role,omitempty"`               // 角色：system, user, assistant
	SystemPrompt      bool     `json:"system_prompt,omitempty"`      // 是否为系统提示词
	Marker            bool     `json:"marker,omitempty"`             // 是否为标记（占位符）
	InjectionPosition int      `json:"injection_position,omitempty"` // 注入位置
	InjectionDepth    int      `json:"injection_depth,omitempty"`    // 注入深度
	InjectionOrder    int      `json:"injection_order,omitempty"`    // 注入顺序
	InjectionTrigger  []string `json:"injection_trigger,omitempty"`  // 注入触发器
	ForbidOverrides   bool     `json:"forbid_overrides,omitempty"`   // 禁止覆盖
	Enabled           bool     `json:"enabled,omitempty"`            // 是否启用
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
