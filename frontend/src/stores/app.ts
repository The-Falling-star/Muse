import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useMediaQuery } from '@vueuse/core';

export const useAppStore = defineStore('app', () => {
  // 侧边栏状态
  const sidebarCollapsed = ref(false);
  const sidebarVisible = ref(true);
  
  // 响应式检测
  const isMobile = useMediaQuery('(max-width: 768px)');
  const isTablet = useMediaQuery('(max-width: 1024px)');
  
  // 当前活跃的角色ID
  const activeCharacterId = ref<string | null>(null);
  
  // 当前活跃的聊天会话ID
  const activeChatId = ref<string | null>(null);
  
  // 全局加载状态
  const isLoading = ref(false);
  
  // 计算属性
  const isCompactMode = computed(() => isMobile.value || sidebarCollapsed.value);
  
  // 切换侧边栏
  const toggleSidebar = () => {
    if (isMobile.value) {
      sidebarVisible.value = !sidebarVisible.value;
    } else {
      sidebarCollapsed.value = !sidebarCollapsed.value;
    }
  };
  
  // 关闭侧边栏（移动端）
  const closeSidebar = () => {
    if (isMobile.value) {
      sidebarVisible.value = false;
    }
  };
  
  // 打开侧边栏
  const openSidebar = () => {
    sidebarVisible.value = true;
  };
  
  // 设置加载状态
  const setLoading = (loading: boolean) => {
    isLoading.value = loading;
  };
  
  // 设置当前角色
  const setActiveCharacter = (characterId: string | null) => {
    activeCharacterId.value = characterId;
  };
  
  // 设置当前聊天
  const setActiveChat = (chatId: string | null) => {
    activeChatId.value = chatId;
  };

  return {
    sidebarCollapsed,
    sidebarVisible,
    isMobile,
    isTablet,
    isCompactMode,
    activeCharacterId,
    activeChatId,
    isLoading,
    toggleSidebar,
    closeSidebar,
    openSidebar,
    setLoading,
    setActiveCharacter,
    setActiveChat
  };
});
