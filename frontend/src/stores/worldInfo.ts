import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type {WorldInfo, WorldInfoEntry, WorldInfoWithLen} from '@/gen/muse/worldinfo_pb';
import {worldInfoClient} from "@/api/client.ts";
import {DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE, FETCH_ALL_PAGE_SIZE} from "@/utils/constants.ts";

/**
 * 世界书 Store
 * 管理世界书列表和当前选中的世界书条目
 * CRUD操作直接使用 @/api/worldInfoApi 中的方法
 */
export const useWorldInfoStore = defineStore('worldInfo', () => {
  // =====================
  // 全局共享状态
  // =====================

  // 世界书列表（缓存）
  const worldInfos = ref<WorldInfoWithLen[]>([]);
  // 当前选中的世界书
  const selectedWorldInfo = ref<WorldInfo | null>(null);
  // 当前世界书的条目列表
  const entries = ref<WorldInfoEntry[]>([]);

  // 分页和加载状态
  const loading = ref(false);
  const loadingMore = ref(false);
  const pagination = ref({
    page: DEFAULT_PAGE_NUM,
    pageSize: DEFAULT_PAGE_SIZE,
    total: 0,
  });

  // 是否还有更多数据
  const hasMore = computed(() => {
    const total = worldInfos.value.length;
    const totalAll = Number(pagination.value.total);
    return totalAll > 0 && total < totalAll;
  });

  // =====================
  // 计算属性
  // =====================

  // 全局世界书
  const globalWorldInfos = computed(() => {
    return worldInfos.value.filter(w => w.worldInfo?.isGlobal);
  });

  // 非全局世界书
  const nonGlobalWorldInfos = computed(() => {
    return worldInfos.value.filter(w => !w.worldInfo?.isGlobal);
  });

  // 按排序顺序排列的条目
  const sortedEntries = computed(() => {
    return [...entries.value].sort((a, b) => a.sortOrder - b.sortOrder);
  });

  // =====================
  // 世界书状态管理
  // =====================
  const loadWorldInfos = async () => {
    loading.value = true;
    try {
      const rsp = await worldInfoClient.listWorldInfos({
        page: DEFAULT_PAGE_NUM,
        pageSize: DEFAULT_PAGE_SIZE,
      });
      worldInfos.value = rsp.worldInfos;
      pagination.value = {
        page: rsp.page,
        pageSize: rsp.pageSize,
        total: Number(rsp.total),
      };
    } finally {
      loading.value = false;
    }
    return worldInfos.value;
  };

  // 加载更多（无限滚动）
  const loadMore = async () => {
    if (loadingMore.value || !hasMore.value) return;
    loadingMore.value = true;
    try {
      const nextPage = pagination.value.page + 1;
      const rsp = await worldInfoClient.listWorldInfos({
        page: nextPage,
        pageSize: pagination.value.pageSize,
      });
      worldInfos.value.push(...rsp.worldInfos);
      pagination.value = {
        page: rsp.page,
        pageSize: rsp.pageSize,
        total: Number(rsp.total),
      };
    } catch (error) {
      console.error('加载更多世界书失败:', error);
    } finally {
      loadingMore.value = false;
    }
  };

  // 重置分页（重新从第一页加载）
  const resetPagination = async () => {
    pagination.value.page = DEFAULT_PAGE_NUM;
    await loadWorldInfos();
  };

  // 获取所有世界书（不分页）
  const loadAllWorldInfos = async () => {
    const rsp = await worldInfoClient.listWorldInfos({
      page: DEFAULT_PAGE_NUM,
      pageSize: FETCH_ALL_PAGE_SIZE,
    });
    worldInfos.value = rsp.worldInfos;
    pagination.value = {
      page: DEFAULT_PAGE_NUM,
      pageSize: rsp.worldInfos.length,
      total: Number(rsp.total),
    };
    return worldInfos.value;
  };

  // 添加世界书
  const addWorldInfo = (worldInfo: WorldInfo) => {
    worldInfos.value.unshift({
      $typeName: 'muse.WorldInfoWithLen',
      worldInfo: worldInfo,
      entryLength: 0,
    });
    pagination.value.total++;
  };

  // 更新世界书
  const updateWorldInfoInList = (worldInfo: WorldInfo) => {
    const index = worldInfos.value.findIndex(w => w.worldInfo?.id === worldInfo.id);
    if (index >= 0 && worldInfos.value[index]) {
      worldInfos.value[index].worldInfo = worldInfo;
    }
    if (selectedWorldInfo.value?.id === worldInfo.id) {
      selectedWorldInfo.value = worldInfo;
    }
  };

  // 移除世界书
  const removeWorldInfo = (id: number) => {
    worldInfos.value = worldInfos.value.filter(w => w.worldInfo?.id !== id);
    if (selectedWorldInfo.value?.id === id) {
      selectedWorldInfo.value = null;
      entries.value = [];
    }
    pagination.value.total = Math.max(0, pagination.value.total - 1);
  };

  // 选中世界书
  const selectWorldInfo = (worldInfo: WorldInfo | null) => {
    selectedWorldInfo.value = worldInfo;
    if (!worldInfo) {
      entries.value = [];
    }
  };

  // =====================
  // 条目状态管理
  // =====================

  // 设置条目列表
  const setEntries = (list: WorldInfoEntry[]) => {
    entries.value = list;
  };

  // 添加条目
  const addEntry = (entry: WorldInfoEntry) => {
    entries.value.push(entry);
  };

  // 更新条目
  const updateEntryInList = (entry: WorldInfoEntry) => {
    const index = entries.value.findIndex(e => e.id === entry.id);
    if (index >= 0) {
      entries.value[index] = entry;
    }
  };

  // 移除条目
  const removeEntry = (id: number) => {
    entries.value = entries.value.filter(e => e.id !== id);
  };

  // 更新本地排序
  const updateLocalEntriesOrder = (entryIds: number[]) => {
    const newOrder = new Map(entryIds.map((id, index) => [id, index]));
    entries.value.forEach(entry => {
      const newSortOrder = newOrder.get(entry.id);
      if (newSortOrder !== undefined) {
        entry.sortOrder = newSortOrder;
      }
    });
  };

  // 切换条目启用状态
  const toggleEntryEnabled = (id: number) => {
    const entry = entries.value.find(e => e.id === id);
    if (entry) {
      entry.isEnabled = !entry.isEnabled;
    }
  };

  // =====================
  // 辅助方法
  // =====================

  // 解析关键词列表
  const parseKeysList = (keysList: string[]): string[] => {
    return keysList;
  };

  // 清空缓存
  const clearCache = () => {
    worldInfos.value = [];
  };

  // 重置状态
  const reset = () => {
    worldInfos.value = [];
    selectedWorldInfo.value = null;
    entries.value = [];
  };

  return {
    // 状态
    worldInfos,
    selectedWorldInfo,
    entries,
    globalWorldInfos,
    nonGlobalWorldInfos,
    sortedEntries,
    loading,
    loadingMore,
    hasMore,

    // 世界书方法
    loadWorldInfos,
    loadMore,
    resetPagination,
    loadAllWorldInfos,
    addWorldInfo,
    updateWorldInfoInList,
    removeWorldInfo,
    selectWorldInfo,

    // 条目方法
    setEntries,
    addEntry,
    updateEntryInList,
    removeEntry,
    updateLocalEntriesOrder,
    toggleEntryEnabled,

    // 辅助方法
    parseKeysList,
    clearCache,
    reset
  };
});
