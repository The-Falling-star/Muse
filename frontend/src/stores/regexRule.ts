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
  return { global: global, preset, character };
};

// 在所有存储中查找规则（用于更新/删除/切换）
const findRuleInLists = (id: number, global: RegexRule[], preset: Map<number, RegexRule[]>, character: Map<number, RegexRule[]>): RegexRule | null => {
  const inGlobal = global.find(r => r.id === id);
  if (inGlobal) return inGlobal;

  for (const [, rules] of preset) {
    const found = rules.find(r => r.id === id);
    if (found) return found;
  }

  for (const [, rules] of character) {
    const found = rules.find(r => r.id === id);
    if (found) return found;
  }

  return null;
};

export const useRegexRuleStore = defineStore('regexRule', () => {
  // =====================
  // 状态：三种类型的正则规则
  // =====================
  const globalRegex = ref<RegexRule[]>([]);
  const presetRegex = ref<Map<number, RegexRule[]>>(new Map());
  const charRegex = ref<Map<number, RegexRule[]>>(new Map());
  const loading = ref(false);

  // =====================
  // 方法
  // =====================

  // 全量拉取正则规则（全局 + 指定预设 + 指定角色）
  const loadRegexs = async (presetId: number = 0, characterId: number = 0) => {
    if (presetRegex.value.has(presetId)) {
      console.info("已存在全局正则: ", presetId)
      presetId = 0;
    }
    if (charRegex.value.has(characterId)) {
      console.info("已存在角色正则: ", characterId)
      characterId = 0;
    }
    if (presetId === 0 && characterId === 0) {
      return;
    }
    console.info("拉取预设正则: ", presetId, "角色正则: ", characterId)
    loading.value = true;
    try {
      const rsp = await regexRuleClient.listRegexRules({ presetId, characterId });
      const { global, preset, character } = distributeRules(rsp.rules);
      globalRegex.value = global;
      if (presetId !== 0) {
        presetRegex.value.set(presetId, preset);
      }
      if (characterId !== 0) {
        charRegex.value.set(characterId, character);
      }
    } finally {
      loading.value = false;
    }
  };

  // 设置预设的正则规则
  const setPresetRegex = (presetId: number, regex: RegexRule[]) => {
    presetRegex.value.set(presetId, regex)
  };

  // 设置角色的正则规则
  const setCharRegex = (characterId: number, regex: RegexRule[]) => {
    charRegex.value.set(characterId, regex)
  };

  // 添加正则规则
  const addRegex = async (data: {
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
      if (type === REGEX_RULE_TYPE.GLOBAL) {
        globalRegex.value.push(response.rule);
      } else if (type === REGEX_RULE_TYPE.PRESET) {
        const key = response.rule.presetId;
        if (!presetRegex.value.has(key)) presetRegex.value.set(key, []);
        presetRegex.value.get(key)!.push(response.rule);
      } else {
        const key = response.rule.characterId;
        if (!charRegex.value.has(key)) charRegex.value.set(key, []);
        charRegex.value.get(key)!.push(response.rule);
      }
    }
    return response.rule;
  };

  // 更新正则规则
  const updateRegex = async (
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
      let idx = globalRegex.value.findIndex(r => r.id === id);
      if (idx >= 0) {
        globalRegex.value.splice(idx, 1);
      } else {
        for (const [, rules] of presetRegex.value) {
          idx = rules.findIndex(r => r.id === id);
          if (idx >= 0) { rules.splice(idx, 1); break; }
        }
        if (idx < 0) {
          for (const [, rules] of charRegex.value) {
            idx = rules.findIndex(r => r.id === id);
            if (idx >= 0) { rules.splice(idx, 1); break; }
          }
        }
      }
      const type = classifyRule(response.rule);
      if (type === REGEX_RULE_TYPE.GLOBAL) {
        globalRegex.value.push(response.rule);
      } else if (type === REGEX_RULE_TYPE.PRESET) {
        const key = response.rule.presetId;
        if (!presetRegex.value.has(key)) presetRegex.value.set(key, []);
        presetRegex.value.get(key)!.push(response.rule);
      } else {
        const key = response.rule.characterId;
        if (!charRegex.value.has(key)) charRegex.value.set(key, []);
        charRegex.value.get(key)!.push(response.rule);
      }
    }
    return response.rule;
  };

  // 删除正则规则
  const deleteRegex = async (id: number) => {
    await regexRuleClient.deleteRegexRule({ id });
    let idx = globalRegex.value.findIndex(r => r.id === id);
    if (idx >= 0) {
      globalRegex.value.splice(idx, 1);
      return;
    }
    for (const [, rules] of presetRegex.value) {
      idx = rules.findIndex(r => r.id === id);
      if (idx >= 0) { rules.splice(idx, 1); return; }
    }
    for (const [, rules] of charRegex.value) {
      idx = rules.findIndex(r => r.id === id);
      if (idx >= 0) { rules.splice(idx, 1); return; }
    }
  };

  // 更新正则规则排序
  const updateRegexOrder = async (presetId: number, ruleIds: number[]) => {
    await regexRuleClient.updateRegexRulesOrder({ presetId, ruleIds });
    const newOrder = new Map(ruleIds.map((id, index) => [id, index]));

    // 更新全局规则排序
    for (const rule of globalRegex.value) {
      const newSortOrder = newOrder.get(rule.id);
      if (newSortOrder !== undefined) {
        rule.sortOrder = newSortOrder;
      }
    }

    // 更新预设规则排序
    for (const rules of presetRegex.value.values()) {
      for (const rule of rules) {
        const newSortOrder = newOrder.get(rule.id);
        if (newSortOrder !== undefined) {
          rule.sortOrder = newSortOrder;
        }
      }
    }

    // 更新角色规则排序
    for (const rules of charRegex.value.values()) {
      for (const rule of rules) {
        const newSortOrder = newOrder.get(rule.id);
        if (newSortOrder !== undefined) {
          rule.sortOrder = newSortOrder;
        }
      }
    }
  };

  // 导入正则规则（全局规则）
  const importRegex = async (fileContent: Uint8Array, fileName: string) => {
    await regexRuleClient.importRegexRules({ fileContent, fileName });
    await loadRegexs();
  };

  // 导出正则规则
  const exportRegex = async (presetId: number) => {
    const response = await regexRuleClient.exportRegexRules({ presetId });
    return {
      fileContent: response.fileContent,
      fileName: response.fileName
    };
  };

  // 切换规则启用状态
  const toggleRegexEnabled = async (id: number) => {
    const rule = findRuleInLists(id, globalRegex.value, presetRegex.value, charRegex.value);
    if (!rule) return;

    return updateRegex(id, {
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
    globalRegex.value = [];
    presetRegex.value = new Map();
    charRegex.value = new Map();
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
    globalRegex,
    presetRegex,
    charRegex,
    loading,

    // 方法
    loadRegexs,
    setPresetRegex,
    setCharRegex,
    addRegex,
    updateRegex,
    deleteRegex,
    updateRegexOrder,
    importRegex,
    exportRegex,
    toggleRegexEnabled,

    // 辅助
    reset,
    createDefaultAffectFlags
  };
});
