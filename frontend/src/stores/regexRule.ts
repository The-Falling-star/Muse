import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { regexRuleClient } from '@/api/client';
import type { RegexRule, RegexAffectFlags } from '@/gen/muse/regex_pb';

export const useRegexRuleStore = defineStore('regexRule', () => {
  // =====================
  // 状态
  // =====================

  // 正则规则列表
  const rules = ref<RegexRule[]>([]);
  // 加载状态
  const loading = ref(false);
  // 搜索关键词
  const searchQuery = ref('');
  // 当前关联的预设ID（0表示全局规则）
  const currentPresetId = ref(0);

  // =====================
  // 计算属性
  // =====================

  // 过滤后的规则列表
  const filteredRules = computed(() => {
    if (!searchQuery.value) return rules.value;
    const query = searchQuery.value.toLowerCase();
    return rules.value.filter(r =>
      r.name.toLowerCase().includes(query) ||
      r.findPattern.toLowerCase().includes(query)
    );
  });

  // 按排序顺序排列的规则
  const sortedRules = computed(() => {
    return [...rules.value].sort((a, b) => a.sortOrder - b.sortOrder);
  });

  // =====================
  // 方法
  // =====================

  // 获取正则规则列表
  const fetchRules = async (presetId: number = 0) => {
    currentPresetId.value = presetId;
    const response = await regexRuleClient.listRegexRules({});
    rules.value = response.rules;
    return response.rules;
  };

  // 添加正则规则
  const addRule = async (data: {
    presetId: number;
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
    const response = await regexRuleClient.addRegexRule(data);
    if (response.rule) {
      rules.value.push(response.rule);
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
      const index = rules.value.findIndex(r => r.id === id);
      if (index >= 0) {
        rules.value[index] = response.rule;
      }
    }
    return response.rule;
  };

  // 删除正则规则
  const deleteRule = async (id: number) => {
    await regexRuleClient.deleteRegexRule({ id });
    rules.value = rules.value.filter(r => r.id !== id);
  };

  // 更新正则规则排序
  const updateRulesOrder = async (presetId: number, ruleIds: number[]) => {
    await regexRuleClient.updateRegexRulesOrder({ presetId, ruleIds });
    // 更新本地排序
    const newOrder = new Map(ruleIds.map((id, index) => [id, index]));
    rules.value.forEach(rule => {
      const newSortOrder = newOrder.get(rule.id);
      if (newSortOrder !== undefined) {
        rule.sortOrder = newSortOrder;
      }
    });
  };

  // 导入正则规则（全局规则）
  const importRules = async (fileContent: Uint8Array, fileName: string) => {
    await regexRuleClient.importRegexRules({ fileContent, fileName });
    // 重新加载规则列表
    await fetchRules(currentPresetId.value);
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
    const rule = rules.value.find(r => r.id === id);
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

  // =====================
  // 辅助方法
  // =====================

  // 设置搜索关键词
  const setSearchQuery = (query: string) => {
    searchQuery.value = query;
  };

  // 重置状态
  const reset = () => {
    rules.value = [];
    searchQuery.value = '';
    currentPresetId.value = 0;
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
    rules,
    loading,
    searchQuery,
    currentPresetId,
    filteredRules,
    sortedRules,

    // 方法
    fetchRules,
    addRule,
    updateRule,
    deleteRule,
    updateRulesOrder,
    importRules,
    exportRules,
    toggleRuleEnabled,

    // 辅助方法
    setSearchQuery,
    reset,
    createDefaultAffectFlags
  };
});
