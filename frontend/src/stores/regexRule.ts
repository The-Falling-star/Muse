import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { regexRuleClient } from '@/api/client';
import type { RegexRule, RegexAffectFlags } from '@/gen/muse/regex_pb';
import { DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE } from '@/utils/constants';

export const useRegexRuleStore = defineStore('regexRule', () => {
  // =====================
  // 状态
  // =====================

  // 正则规则列表
  const rules = ref<RegexRule[]>([]);
  // 加载状态
  const loading = ref(false);
  const loadingMore = ref(false);
  // 搜索关键词
  const searchQuery = ref('');
  // 当前关联的预设ID（0表示全局规则）
  const currentPresetId = ref(0);
  // 分页状态
  const pagination = ref({
    page: DEFAULT_PAGE_NUM,
    pageSize: DEFAULT_PAGE_SIZE,
    total: 0,
  });

  // 是否还有更多数据
  const hasMore = computed(() => {
    const total = rules.value.length;
    const totalAll = Number(pagination.value.total);
    return totalAll > 0 && total < totalAll;
  });

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

  // 获取正则规则列表（初始加载，从第一页开始）
  const fetchRules = async (presetId: number = 0) => {
    currentPresetId.value = presetId;
    loading.value = true;
    try {
      const response = await regexRuleClient.listRegexRules({
        page: DEFAULT_PAGE_NUM,
        pageSize: DEFAULT_PAGE_SIZE,
      });
      rules.value = response.rules;
      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: Number(response.total),
      };
    } finally {
      loading.value = false;
    }
    return rules.value;
  };

  // 加载更多（无限滚动）
  const loadMore = async () => {
    if (loadingMore.value || !hasMore.value) return;
    loadingMore.value = true;
    try {
      const nextPage = pagination.value.page + 1;
      const response = await regexRuleClient.listRegexRules({
        page: nextPage,
        pageSize: pagination.value.pageSize,
      });
      rules.value.push(...response.rules);
      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: Number(response.total),
      };
    } catch (error) {
      console.error('加载更多正则规则失败:', error);
    } finally {
      loadingMore.value = false;
    }
  };

  // 重置分页（重新从第一页加载）
  const resetPagination = async () => {
    pagination.value.page = DEFAULT_PAGE_NUM;
    await fetchRules(currentPresetId.value);
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
    // 重新加载第一页
    pagination.value.page = DEFAULT_PAGE_NUM;
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
    loadingMore,
    searchQuery,
    currentPresetId,
    filteredRules,
    sortedRules,
    hasMore,

    // 方法
    fetchRules,
    loadMore,
    resetPagination,
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
