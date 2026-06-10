import { defineStore } from 'pinia';
import { ref } from 'vue';
import { regexRuleClient } from '@/api/client';
import { REGEX_RULE_TYPE } from '@/utils/constants';
import type { RegexRuleType } from '@/utils/constants';
import type { RegexRule, RegexAffectFlags } from '@/gen/muse/regex_pb';

// 根据规则字段分类到对应类型
const classifyRule = (rule: RegexRule): RegexRuleType => {
  if (rule.characterId > 0) return REGEX_RULE_TYPE.CHARACTER;
  if (rule.presetId > 0) return REGEX_RULE_TYPE.PRESET;
  return REGEX_RULE_TYPE.GLOBAL;
};

// 将规则列表分发到三个分类数组
const distributeRules = (rules: RegexRule[]) => {
  const global: RegexRule[] = [];
  const preset: RegexRule[] = [];
  const character: RegexRule[] = [];
  for (const rule of rules) {
    switch (classifyRule(rule)) {
      case REGEX_RULE_TYPE.GLOBAL:
        global.push(rule);
        break;
      case REGEX_RULE_TYPE.PRESET:
        preset.push(rule);
        break;
      case REGEX_RULE_TYPE.CHARACTER:
        character.push(rule);
        break;
    }
  }
  return { global, preset, character };
};

// 在三个列表中查找规则（用于更新/删除/切换）
const findRuleInLists = (id: number, lists: RegexRule[][]): RegexRule | null => {
  for (const list of lists) {
    const found = list.find(r => r.id === id);
    if (found) return found;
  }
  return null;
};

export const useRegexRuleStore = defineStore('regexRule', () => {
  // =====================
  // 状态：三种类型的正则规则
  // =====================
  const globalRules = ref<RegexRule[]>([]);
  const presetRules = ref<RegexRule[]>([]);
  const characterRules = ref<RegexRule[]>([]);
  const loading = ref(false);
  const loaded = ref(false);

  // =====================
  // 方法
  // =====================

  // 全量拉取正则规则（全局 + 指定预设 + 指定角色）
  const loadRules = async (presetId: number = 0, characterId: number = 0) => {
    if (loaded.value) {
      return;
    }
    loading.value = true;
    try {
      const rsp = await regexRuleClient.listRegexRules({ presetId, characterId });
      const { global, preset, character } = distributeRules(rsp.rules);
      globalRules.value = global;
      presetRules.value = preset;
      characterRules.value = character;
      loaded.value = true;
    } finally {
      loading.value = false;
    }
  };

  // 添加正则规则
  const addRule = async (data: {
    presetId: number;
    characterId?: number;
    name: string;
    findPattern: string;
    replacePattern?: string;
    isEnabled: boolean;
    runOnEdit: boolean;
    substituteRegex: boolean;
    minDepth?: number;
    maxDepth?: number;
    affectFlags?: RegexAffectFlags;
    sortOrder: number;
  }) => {
    const response = await regexRuleClient.addRegexRule({
      presetId: data.presetId,
      characterId: data.characterId ?? 0,
      name: data.name,
      findPattern: data.findPattern,
      replacePattern: data.replacePattern,
      isEnabled: data.isEnabled,
      runOnEdit: data.runOnEdit,
      substituteRegex: data.substituteRegex,
      minDepth: data.minDepth,
      maxDepth: data.maxDepth,
      affectFlags: data.affectFlags,
      sortOrder: data.sortOrder,
    });
    if (response.rule) {
      const type = classifyRule(response.rule);
      if (type === REGEX_RULE_TYPE.GLOBAL) globalRules.value.push(response.rule);
      else if (type === REGEX_RULE_TYPE.PRESET) presetRules.value.push(response.rule);
      else characterRules.value.push(response.rule);
    }
    return response.rule;
  };

  // 更新正则规则
  const updateRule = async (
    id: number,
    data: {
      name: string;
      findPattern: string;
      replacePattern?: string;
      isEnabled: boolean;
      runOnEdit: boolean;
      substituteRegex: boolean;
      minDepth?: number;
      maxDepth?: number;
      affectFlags?: RegexAffectFlags;
      sortOrder: number;
    }
  ) => {
    const response = await regexRuleClient.updateRegexRule({ id, ...data });
    if (response.rule) {
      // 先从所有列表中移除旧项，再按新类型添加
      const lists = [globalRules.value, presetRules.value, characterRules.value];
      for (const list of lists) {
        const index = list.findIndex(r => r.id === id);
        if (index >= 0) {
          list.splice(index, 1);
          break;
        }
      }
      const type = classifyRule(response.rule);
      if (type === REGEX_RULE_TYPE.GLOBAL) globalRules.value.push(response.rule);
      else if (type === REGEX_RULE_TYPE.PRESET) presetRules.value.push(response.rule);
      else characterRules.value.push(response.rule);
    }
    return response.rule;
  };

  // 删除正则规则
  const deleteRule = async (id: number) => {
    await regexRuleClient.deleteRegexRule({ id });
    const lists = [globalRules.value, presetRules.value, characterRules.value];
    for (const list of lists) {
      const index = list.findIndex(r => r.id === id);
      if (index >= 0) {
        list.splice(index, 1);
        return;
      }
    }
  };

  // 更新正则规则排序
  const updateRulesOrder = async (presetId: number, ruleIds: number[]) => {
    await regexRuleClient.updateRegexRulesOrder({ presetId, ruleIds });
    const newOrder = new Map(ruleIds.map((id, index) => [id, index]));
    const lists = [globalRules.value, presetRules.value, characterRules.value];
    for (const list of lists) {
      for (const rule of list) {
        const newSortOrder = newOrder.get(rule.id);
        if (newSortOrder !== undefined) {
          rule.sortOrder = newSortOrder;
        }
      }
    }
  };

  // 导入正则规则（全局规则）
  const importRules = async (fileContent: Uint8Array, fileName: string) => {
    await regexRuleClient.importRegexRules({ fileContent, fileName });
    await loadRules();
  };

  // 导出正则规则
  const exportRules = async (presetId: number) => {
    const response = await regexRuleClient.exportRegexRules({ presetId });
    return {
      fileContent: response.fileContent,
      fileName: response.fileName
    };
  };

  // 切换规则启用状态
  const toggleRuleEnabled = async (id: number) => {
    const rule = findRuleInLists(id, [globalRules.value, presetRules.value, characterRules.value]);
    if (!rule) return;

    return updateRule(id, {
      name: rule.name,
      findPattern: rule.findPattern,
      replacePattern: rule.replacePattern,
      isEnabled: !rule.isEnabled,
      runOnEdit: rule.runOnEdit,
      substituteRegex: rule.substituteRegex,
      minDepth: rule.minDepth,
      maxDepth: rule.maxDepth,
      affectFlags: rule.affectFlags,
      sortOrder: rule.sortOrder
    });
  };

  // 重置状态
  const reset = () => {
    globalRules.value = [];
    presetRules.value = [];
    characterRules.value = [];
    loaded.value = false;
  };

  // 创建默认的AffectFlags
  const createDefaultAffectFlags = (): RegexAffectFlags => ({
    userInput: false,
    aiOutput: true,
    slashCommand: false,
    worldInfo: false,
    prompt: false,
    $typeName: 'muse.RegexAffectFlags'
  });

  return {
    // 状态
    globalRules,
    presetRules,
    characterRules,
    loading,
    loaded,

    // 方法
    loadRules,
    addRule,
    updateRule,
    deleteRule,
    updateRulesOrder,
    importRules,
    exportRules,
    toggleRuleEnabled,

    // 辅助
    reset,
    createDefaultAffectFlags
  };
});
