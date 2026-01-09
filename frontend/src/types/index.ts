// 角色相关类型
export interface Character {
  id: string;
  name: string;
  avatar?: string;
  description?: string;
  personality?: string;
  scenario?: string;
  firstMessage?: string;
  exampleDialogue?: string;
  creatorNotes?: string;
  tags?: string[];
  createdAt?: string;
  updatedAt?: string;
}

// 单条消息内容（swipe）
export interface MessageSwipe {
  id: string;                     // swipe 的唯一ID
  content: string;                // 消息内容
  timestamp: number;              // 生成时间
}

// 消息类型（一个楼层）
export interface Message {
  id: string;                     // 楼层ID
  role: 'user' | 'assistant' | 'system';  // 消息角色
  swipes: MessageSwipe[];         // 该楼层的所有消息内容
  currentSwipeIndex: number;      // 当前显示的消息索引
  isStreaming?: boolean;          // 是否正在流式输出
}

// 获取当前显示内容的辅助函数类型
export const getCurrentContent = (message: Message): string => {
  if (message.swipes.length === 0) return '';
  const index = Math.min(message.currentSwipeIndex, message.swipes.length - 1);
  return message.swipes[index]?.content || '';
};

// 获取当前时间戳的辅助函数类型
export const getCurrentTimestamp = (message: Message): number => {
  if (message.swipes.length === 0) return Date.now();
  const index = Math.min(message.currentSwipeIndex, message.swipes.length - 1);
  return message.swipes[index]?.timestamp || Date.now();
};

// 聊天会话类型
export interface ChatSession {
  id: string;
  characterId: string;
  characterName: string;
  messages: Message[];
  createdAt: string;
  updatedAt: string;
}

// 提示项类型 - 预设中的单条设定项
export interface PromptItem {
  id: string;
  identifier: string;           // 唯一标识符，如 'main', 'jailbreak', 'nsfw' 等
  name: string;                 // 显示名称
  role: 'system' | 'user' | 'assistant';  // 消息角色
  content: string;              // 提示内容
  enabled: boolean;             // 是否启用
  marker: boolean;              // 是否为标记（占位符），如 chatHistory, worldInfoBefore 等
  injection?: {                 // 注入位置设置
    position: 'before' | 'after';  // 相对位置
    depth: number;                 // 注入深度
  };
}

// 提示项排序条目
export interface PromptOrderEntry {
  identifier: string;           // 对应 PromptItem 的 identifier
  enabled: boolean;             // 该项是否启用
}

// 预设类型 - 包含多个提示项
export interface Preset {
  id: string;
  name: string;                 // 预设名称
  description?: string;         // 预设描述
  prompts: PromptItem[];        // 提示项列表
  promptOrder: PromptOrderEntry[]; // 提示项排序和启用状态
  regexRules?: RegexRule[];     // 预设绑定的正则规则
  // 生成参数
  temperature?: number;
  topP?: number;
  topK?: number;
  maxTokens?: number;
  frequencyPenalty?: number;
  presencePenalty?: number;
  // 元数据
  createdAt?: string;
  updatedAt?: string;
}

// 世界书条目类型
export interface WorldInfoEntry {
  id: string;
  worldId: string;
  keys: string[];
  secondaryKeys?: string[];
  content: string;
  comment?: string;
  enabled: boolean;
  order: number;
  probability: number;
  depth: number;
  selectiveLogic: 'and' | 'or' | 'not';
}

// 世界书类型
export interface WorldInfo {
  id: string;
  name: string;
  description?: string;
  entries: WorldInfoEntry[];
  enabled: boolean;
}

// 正则规则类型
export interface RegexRule {
  id: string;
  name: string;
  pattern: string;
  replacement: string;
  flags: string;
  scope: 'input' | 'output' | 'both';
  enabled: boolean;
  order: number;
}

// API配置类型
export interface ApiConfig {
  id: string;
  name: string;
  type: 'openai' | 'claude' | 'gemini';
  baseUrl: string;
  apiKey: string;
  model: string;
  maxTokens: number;
  temperature: number;
  topP: number;
  frequencyPenalty: number;
  presencePenalty: number;
}

// 应用设置类型
export interface AppSettings {
  theme: 'dark' | 'light';
  fontSize: number;
  language: string;
  sendOnEnter: boolean;
  streamResponse: boolean;
  autoSave: boolean;
}

// 导航菜单项类型
export interface NavMenuItem {
  key: string;
  label: string;
  icon: string;
  path: string;
}

// 系统标记类型 - 预定义的占位符标识符
export const SYSTEM_MARKERS = {
  MAIN: 'main',                     // 主提示
  NSFW: 'nsfw',                     // 辅助提示/NSFW
  JAILBREAK: 'jailbreak',           // 越狱/后置指令
  CHAT_HISTORY: 'chatHistory',      // 聊天历史
  DIALOGUE_EXAMPLES: 'dialogueExamples', // 对话示例
  WORLD_INFO_BEFORE: 'worldInfoBefore',  // 世界书（前）
  WORLD_INFO_AFTER: 'worldInfoAfter',    // 世界书（后）
  CHAR_DESCRIPTION: 'charDescription',   // 角色描述
  CHAR_PERSONALITY: 'charPersonality',   // 角色性格
  SCENARIO: 'scenario',                  // 场景
  PERSONA_DESCRIPTION: 'personaDescription', // 用户人设
  ENHANCE_DEFINITIONS: 'enhanceDefinitions', // 增强定义
} as const;

// 默认提示项模板
export const DEFAULT_PROMPT_ITEMS: PromptItem[] = [
  {
    id: 'main',
    identifier: 'main',
    name: '主提示',
    role: 'system',
    content: 'Write {{char}}\'s next reply in a fictional chat between {{char}} and {{user}}.',
    enabled: true,
    marker: false
  },
  {
    id: 'worldInfoBefore',
    identifier: 'worldInfoBefore',
    name: '世界书（前）',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'charDescription',
    identifier: 'charDescription',
    name: '角色描述',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'charPersonality',
    identifier: 'charPersonality',
    name: '角色性格',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'scenario',
    identifier: 'scenario',
    name: '场景',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'nsfw',
    identifier: 'nsfw',
    name: '辅助提示',
    role: 'system',
    content: '',
    enabled: true,
    marker: false
  },
  {
    id: 'worldInfoAfter',
    identifier: 'worldInfoAfter',
    name: '世界书（后）',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'dialogueExamples',
    identifier: 'dialogueExamples',
    name: '对话示例',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'chatHistory',
    identifier: 'chatHistory',
    name: '聊天历史',
    role: 'system',
    content: '',
    enabled: true,
    marker: true
  },
  {
    id: 'jailbreak',
    identifier: 'jailbreak',
    name: '后置指令',
    role: 'system',
    content: '',
    enabled: true,
    marker: false
  }
];

// 默认提示项排序
export const DEFAULT_PROMPT_ORDER: PromptOrderEntry[] = [
  { identifier: 'main', enabled: true },
  { identifier: 'worldInfoBefore', enabled: true },
  { identifier: 'charDescription', enabled: true },
  { identifier: 'charPersonality', enabled: true },
  { identifier: 'scenario', enabled: true },
  { identifier: 'nsfw', enabled: true },
  { identifier: 'worldInfoAfter', enabled: true },
  { identifier: 'dialogueExamples', enabled: true },
  { identifier: 'chatHistory', enabled: true },
  { identifier: 'jailbreak', enabled: true }
];
