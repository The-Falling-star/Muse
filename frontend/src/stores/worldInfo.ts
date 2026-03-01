import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { WorldInfo, WorldInfoEntry } from '@/gen/muse/worldinfo_pb';

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
  const worldInfos = ref<WorldInfo[]>([]);
  // 当前选中的世界书
  const selectedWorldInfo = ref<WorldInfo | null>(null);
  // 当前世界书的条目列表
  const entries = ref<WorldInfoEntry[]>([]);

  // =====================
  // 计算属性
  // =====================

  // 全局世界书
  const globalWorldInfos = computed(() => {
    return worldInfos.value.filter(w => w.isGlobal);
  });

  // 非全局世界书
  const nonGlobalWorldInfos = computed(() => {
    return worldInfos.value.filter(w => !w.isGlobal);
  });

  // 按排序顺序排列的条目
  const sortedEntries = computed(() => {
    return [...entries.value].sort((a, b) => a.sortOrder - b.sortOrder);
  });

  // 是否有缓存
  const hasCached = computed(() => worldInfos.value.length > 0);

  // =====================
  // 世界书状态管理
  // =====================

  // 设置世界书列表
  const setWorldInfos = (list: WorldInfo[]) => {
    worldInfos.value = list;
  };

  // 添加世界书
  const addWorldInfo = (worldInfo: WorldInfo) => {
    worldInfos.value.push(worldInfo);
  };

  // 更新世界书
  const updateWorldInfoInList = (worldInfo: WorldInfo) => {
    const index = worldInfos.value.findIndex(w => w.id === worldInfo.id);
    if (index >= 0) {
      worldInfos.value[index] = worldInfo;
    }
    if (selectedWorldInfo.value?.id === worldInfo.id) {
      selectedWorldInfo.value = worldInfo;
    }
  };

  // 移除世界书
  const removeWorldInfo = (id: number) => {
    worldInfos.value = worldInfos.value.filter(w => w.id !== id);
    if (selectedWorldInfo.value?.id === id) {
      selectedWorldInfo.value = null;
      entries.value = [];
    }
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
    hasCached,

    // 世界书方法
    setWorldInfos,
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
