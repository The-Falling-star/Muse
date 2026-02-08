import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { Character } from '@/gen/muse/muse_pb';

/**
 * 角色 Store
 * 只管理全局共享的状态：角色列表缓存、当前选中角色
 * CRUD操作直接使用 @/api/characterApi 中的方法
 */
export const useCharacterStore = defineStore('character', () => {
  // =====================
  // 全局共享状态
  // =====================

  // 角色列表（缓存）
  const characters = ref<Character[]>([]);
  // 当前选中的角色（用于聊天页面等跨页面共享）
  const selectedCharacter = ref<Character | null>(null);
  // 总数
  const total = ref(0);
  // 角色详情缓存（key: characterId, value: 完整的角色信息）
  const characterDetailsCache = ref<Map<number, Character>>(new Map());

  // =====================
  // 计算属性
  // =====================

  // 是否有缓存数据
  const hasCached = computed(() => characters.value.length > 0);

  // =====================
  // 状态管理方法
  // =====================

  // 设置角色列表（由View调用API后更新）
  const setCharacters = (list: Character[], totalCount: number) => {
    characters.value = list;
    total.value = totalCount;
  };

  // 添加角色到列表头部（创建/导入后调用）
  const addCharacter = (character: Character) => {
    characters.value.unshift(character);
    total.value += 1;
  };

  // 更新列表中的角色
  const updateCharacterInList = (character: Character) => {
    const index = characters.value.findIndex(c => c.id === character.id);
    if (index >= 0) {
      characters.value[index] = character;
    }
    if (selectedCharacter.value?.id === character.id) {
      selectedCharacter.value = character;
    }
    // 更新详情缓存
    characterDetailsCache.value.set(character.id, character);
  };

  // 从列表中移除角色
  const removeCharacter = (id: number) => {
    characters.value = characters.value.filter(c => c.id !== id);
    total.value -= 1;
    if (selectedCharacter.value?.id === id) {
      selectedCharacter.value = null;
    }
    // 清除详情缓存
    characterDetailsCache.value.delete(id);
  };

  // 选中角色
  const selectCharacter = (character: Character | null) => {
    selectedCharacter.value = character;
  };

  // 清空缓存
  const clearCache = () => {
    characters.value = [];
    total.value = 0;
    characterDetailsCache.value.clear();
  };

  // 重置全部状态
  const reset = () => {
    clearCache();
    selectedCharacter.value = null;
  };

  // 获取角色详情缓存
  const getCharacterDetail = (id: number): Character | undefined => {
    return characterDetailsCache.value.get(id);
  };

  // 缓存角色详情
  const cacheCharacterDetail = (character: Character) => {
    characterDetailsCache.value.set(character.id, character);
  };

  // 检查是否有角色详情缓存
  const hasCharacterDetail = (id: number): boolean => {
    return characterDetailsCache.value.has(id);
  };

  return {
    // 状态
    characters,
    selectedCharacter,
    total,
    hasCached,

    // 方法
    setCharacters,
    addCharacter,
    updateCharacterInList,
    removeCharacter,
    selectCharacter,
    clearCache,
    reset,
    getCharacterDetail,
    cacheCharacterDetail,
    hasCharacterDetail
  };
});
