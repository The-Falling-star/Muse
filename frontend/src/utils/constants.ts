/**
 * 分页相关常量
 */
export const DEFAULT_PAGE_NUM = 1;
export const DEFAULT_PAGE_SIZE = 20;
export const MAX_PAGE_SIZE = 100;

// 获取所有数据时的分页大小
export const FETCH_ALL_PAGE_SIZE = 1000;

// 宏
export const MACRO_USER = ['{{user}}', '{{User}}', '<user>']
export const MACRO_CHAR = ['{{char}}', '{{Char}}', '<char>', '<BOT>']

/**
 * 正则规则类型
 * 根据 preset_id / character_id 组合决定：
 *   - 全局正则: preset_id=0 且 character_id=0
 *   - 预设正则: preset_id>0 且 character_id=0
 *   - 角色正则: character_id>0
 */
export const REGEX_RULE_TYPE = {
  GLOBAL: 'global',
  PRESET: 'preset',
  CHARACTER: 'character',
} as const;

export type RegexRuleType = typeof REGEX_RULE_TYPE[keyof typeof REGEX_RULE_TYPE];

export const INPUT_EXTRA_OPTIONS_KEY= {
  IMPORT_CHAT_HISTORY: 'import_chat_history',
}