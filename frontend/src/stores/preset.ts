import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { presetClient } from '@/api/client';
import type { PromptItem, PresetWithAll, PresetWithPromptLen } from '@/gen/muse/preset_pb';
import type { Role, InjectionPosition, PromptItemIdentifier } from '@/gen/muse/common_pb';
import { DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE, FETCH_ALL_PAGE_SIZE } from '@/utils/constants';

export const usePresetStore = defineStore('preset', () => {
  // =====================
  // 状态
  // =====================

  // 预设列表
  const presets = ref<PresetWithPromptLen[]>([]);
  // 当前选中的预设
  const selectedPreset = ref<PresetWithAll | null>(null);
  // 当前预设的提示项列表
  const promptItems = ref<PromptItem[]>([]);
  // 加载状态
  const loading = ref(false);
  const loadingMore = ref(false);
  // 搜索关键词
  const searchQuery = ref('');
  // 分页状态
  const pagination = ref({
    page: DEFAULT_PAGE_NUM,
    pageSize: DEFAULT_PAGE_SIZE,
    total: 0,
  });

  // 是否还有更多数据
  const hasMore = computed(() => {
    const total = presets.value.length;
    const totalAll = Number(pagination.value.total);
    return totalAll > 0 && total < totalAll;
  });

  // =====================
  // 计算属性
  // =====================

  // 过滤后的预设列表
  const filteredPresets = computed(() => {
    if (!searchQuery.value) return presets.value;
    const query = searchQuery.value.toLowerCase();
    return presets.value.filter(p => p.preset?.name.toLowerCase().includes(query));
  });

  // 按排序顺序排列的提示项（后端已按顺序返回）
  const sortedPromptItems = computed(() => {
    return [...promptItems.value];
  });

  // =====================
  // 预设相关方法
  // =====================

  // 获取预设列表（替换）
  const fetchPresets = async (page?: number, pageSize?: number) => {
    const requestPage = page ?? pagination.value.page;
    const requestPageSize = pageSize ?? pagination.value.pageSize;

    const response = await presetClient.listPresets({
      page: requestPage,
      pageSize: requestPageSize,
    });

    presets.value = response.presets;
    pagination.value = {
      page: response.page,
      pageSize: response.pageSize,
      total: Number(response.total),
    };

    return response.presets;
  };

  // 获取所有预设（不分页，用于需要全部数据的场景）
  const fetchAllPresets = async () => {
    const response = await presetClient.listPresets({
      page: DEFAULT_PAGE_NUM,
      pageSize: FETCH_ALL_PAGE_SIZE,
    });
    presets.value = response.presets;
    pagination.value = {
      page: DEFAULT_PAGE_NUM,
      pageSize: response.presets.length,
      total: Number(response.total),
    };
    return response.presets;
  };

  // 加载更多（无限滚动用）
  const loadMore = async () => {
    if (loadingMore.value || !hasMore.value) return;
    loadingMore.value = true;
    try {
      const nextPage = pagination.value.page + 1;
      const response = await presetClient.listPresets({
        page: nextPage,
        pageSize: pagination.value.pageSize,
      });
      presets.value.push(...response.presets);
      pagination.value = {
        page: response.page,
        pageSize: response.pageSize,
        total: Number(response.total),
      };
    } catch (error) {
      console.error('加载更多预设失败:', error);
    } finally {
      loadingMore.value = false;
    }
  };

  // 重置分页（重新从第一页加载）
  const resetPagination = async () => {
    pagination.value.page = DEFAULT_PAGE_NUM;
    await fetchPresets(DEFAULT_PAGE_NUM, pagination.value.pageSize);
  };

  // 获取单个预设
  const fetchPreset = async (id: number) => {
    const response = await presetClient.getPreset({ id });
    if (response.preset) {
      const index = presets.value.findIndex(p => p.preset?.id === id);
      if (index >= 0) {
        // 更新列表中的预设信息（保持prompt_len）
        presets.value[index] = {
          $typeName: 'muse.PresetWithPromptLen',
          preset: response.preset.preset,
          promptLen: response.preset.promptItems?.length || 0
        };
      }
      // 同步更新提示项
      if (response.preset.promptItems) {
        promptItems.value = response.preset.promptItems;
      }
    }
    return response.preset;
  };

  // 创建预设
  const createPreset = async (data: {
    name: string;
    temperature: number;
    topP: number;
    topK: number;
    maxTokens: number;
    frequencyPenalty: number;
    presencePenalty: number;
    promptItems?: Array<{
      identifier: PromptItemIdentifier;
      name: string;
      content?: string;
      role: Role;
      isEnabled: boolean;
      injectionPosition: InjectionPosition;
      injectionDepth: number;
      forbidOverrides: boolean;
      sortOrder: number;
    }>;
  }) => {
    await presetClient.createPreset(data);
    // 创建成功后重新获取第一页列表
    pagination.value.page = DEFAULT_PAGE_NUM;
    await fetchPresets(DEFAULT_PAGE_NUM, pagination.value.pageSize);
  };

  // 更新预设
  const updatePreset = async (
    id: number,
    data: {
      name: string;
      temperature: number;
      topP: number;
      topK: number;
      maxTokens: number;
      frequencyPenalty: number;
      presencePenalty: number;
      version: bigint;
    }
  ) => {
    await presetClient.updatePreset({ id, ...data });
    // 更新成功后重新获取该预设的完整信息
    const updatedPreset = await fetchPreset(id);
    return updatedPreset;
  };

  // 删除预设
  const deletePreset = async (id: number) => {
    await presetClient.deletePreset({ id });
    presets.value = presets.value.filter(p => p.preset?.id !== id);
    if (selectedPreset.value?.preset?.id === id) {
      selectedPreset.value = null;
      promptItems.value = [];
    }
  };

  // 设置活跃预设
  const setActivePreset = async (presetId: number) => {
    await presetClient.setActivePreset({ presetId });
  };

  // 导入预设
  const importPreset = async (fileContent: Uint8Array, fileName: string) => {
    const response = await presetClient.importPreset({ fileContent, fileName });
    if (response.preset) {
      // 添加到列表
      presets.value.push({
        $typeName: 'muse.PresetWithPromptLen',
        preset: response.preset.preset,
        promptLen: response.preset.promptItems?.length || 0
      });
    }
    return response.preset;
  };

  // 导出预设
  const exportPreset = async (id: number) => {
    const response = await presetClient.exportPreset({ id });
    return {
      fileContent: response.fileContent,
      fileName: response.fileName
    };
  };

  // 选中预设
  const selectPreset = async (preset: PresetWithPromptLen | null) => {
    if (preset && preset.preset) {
      // 加载完整的预设信息
      const fullPreset = await fetchPreset(preset.preset.id);
      if (fullPreset) {
        selectedPreset.value = fullPreset;
      }
    } else {
      selectedPreset.value = null;
      promptItems.value = [];
    }
  };

  // =====================
  // 提示项相关方法
  // =====================

  // 获取预设的提示项列表
  const fetchPromptItems = async (presetId: number) => {
    const response = await presetClient.listPromptItems({ presetId });
    promptItems.value = response.items;
    return response.items;
  };

  // 添加提示项
  const addPromptItem = async (
    presetId: number,
    data: {
      identifier: PromptItemIdentifier;
      name: string;
      content?: string;
      role: Role;
      isEnabled: boolean;
      injectionPosition: InjectionPosition;
      injectionDepth: number;
      forbidOverrides: boolean;
      sortOrder: number;
    }
  ) => {
    const response = await presetClient.addPromptItem({ presetId, ...data });
    if (response.item) {
      promptItems.value.push(response.item);
    }
    return response.item;
  };

  // 更新提示项
  const updatePromptItem = async (
    id: number,
    data: {
      identifier: PromptItemIdentifier;
      name: string;
      content?: string;
      role: Role;
      isEnabled: boolean;
      injectionPosition: InjectionPosition;
      injectionDepth: number;
      forbidOverrides: boolean;
      sortOrder: number;
    }
  ) => {
    const response = await presetClient.updatePromptItem({ id, ...data });
    if (response.item) {
      const index = promptItems.value.findIndex(item => item.id === id);
      if (index >= 0) {
        promptItems.value[index] = response.item;
      }
    }
    return response.item;
  };

  // 删除提示项
  const deletePromptItem = async (id: number) => {
    await presetClient.deletePromptItem({ id });
    promptItems.value = promptItems.value.filter(item => item.id !== id);
  };

  // 更新提示项排序
  const updatePromptItemsOrder = async (presetId: number) => {
    await presetClient.updatePromptItemsOrder({ presetId });
    // 重新获取提示项列表以更新本地状态
    await fetchPromptItems(presetId);
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
    presets.value = [];
    selectedPreset.value = null;
    promptItems.value = [];
    searchQuery.value = '';
    pagination.value = {
      page: DEFAULT_PAGE_NUM,
      pageSize: DEFAULT_PAGE_SIZE,
      total: 0
    };
  };

  return {
    // 状态
    presets,
    selectedPreset,
    promptItems,
    sortedPromptItems,
    loading,
    loadingMore,
    searchQuery,
    filteredPresets,
    hasMore,

    // 预设方法
    fetchPresets,
    fetchAllPresets,
    loadMore,
    resetPagination,
    fetchPreset,
    createPreset,
    updatePreset,
    deletePreset,
    setActivePreset,
    importPreset,
    exportPreset,
    selectPreset,

    // 提示项方法
    fetchPromptItems,
    addPromptItem,
    updatePromptItem,
    deletePromptItem,
    updatePromptItemsOrder,

    // 辅助方法
    setSearchQuery,
    reset
  };
});
