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
