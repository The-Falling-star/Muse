import { defineStore } from 'pinia';
import { ref, computed, watch } from 'vue';
import { useMediaQuery } from '@vueuse/core';

// 右侧面板Tab类型
export type RightPanelTab = 'model' | 'preset' | 'worldinfo' | 'regex' | 'user';

// 左侧边栏视图类型
export type SidebarView = 'sessions' | 'characters';

export const useAppStore = defineStore('app', () => {
  // ====== 响应式检测 ======
  const isMobile = useMediaQuery('(max-width: 767px)');
  const isTablet = useMediaQuery('(min-width: 768px) and (max-width: 1023px)');
  const isDesktop = useMediaQuery('(min-width: 1024px)');

  // ====== 左侧边栏状态 ======
  const leftSidebarVisible = ref(true);
  const sidebarView = ref<SidebarView>('sessions');

  // ====== 右侧面板状态 ======
  const rightPanelVisible = ref(false);
  const rightPanelActiveTab = ref<RightPanelTab>('model');

  // ====== 设置抽屉状态 ======
  const settingsDrawerVisible = ref(false);

  // ====== 当前活跃的角色ID ======
  const activeCharacterId = ref<string | null>(null);

  // ====== 当前活跃的聊天会话ID ======
  const activeChatId = ref<string | null>(null);

  // ====== 全局加载状态 ======
  const isLoading = ref(false);

  // ====== 计算属性 ======
  const isCompactMode = computed(() => isMobile.value || isTablet.value);

  // ====== 左侧边栏操作 ======
  const toggleLeftSidebar = () => {
    leftSidebarVisible.value = !leftSidebarVisible.value;
  };

  const closeLeftSidebar = () => {
    leftSidebarVisible.value = false;
  };

  const openLeftSidebar = () => {
    leftSidebarVisible.value = true;
  };

  // ====== 右侧面板操作 ======
  const toggleRightPanel = () => {
    rightPanelVisible.value = !rightPanelVisible.value;
  };

  const openRightPanel = (tab?: RightPanelTab) => {
    if (tab) {
      rightPanelActiveTab.value = tab;
    }
    rightPanelVisible.value = true;
  };

  const closeRightPanel = () => {
    rightPanelVisible.value = false;
  };

  const setRightPanelTab = (tab: RightPanelTab) => {
    rightPanelActiveTab.value = tab;
  };

  // ====== 左侧边栏视图切换 ======
  const setSidebarView = (view: SidebarView) => {
    sidebarView.value = view;
  };

  // ====== 设置抽屉操作 ======
  const openSettings = () => {
    settingsDrawerVisible.value = true;
  };

  const closeSettings = () => {
    settingsDrawerVisible.value = false;
  };

  // ====== 通用状态操作 ======
  const setLoading = (loading: boolean) => {
    isLoading.value = loading;
  };

  const setActiveCharacter = (characterId: string | null) => {
    activeCharacterId.value = characterId;
  };

  const setActiveChat = (chatId: string | null) => {
    activeChatId.value = chatId;
  };

  // ====== 响应式自适应 ======
  // 平板和移动端默认隐藏左侧边栏
  watch(isDesktop, (desktop) => {
    if (!desktop) {
      leftSidebarVisible.value = false;
    } else {
      leftSidebarVisible.value = true;
    }
  }, { immediate: true });

  // 移动端进入编辑视图时自动关闭右侧面板
  watch(isMobile, (mobile) => {
    if (mobile && rightPanelVisible.value) {
      rightPanelVisible.value = false;
    }
  });

  return {
    // 响应式
    isMobile,
    isTablet,
    isDesktop,
    isCompactMode,

    // 左侧边栏
    leftSidebarVisible,
    sidebarView,
    toggleLeftSidebar,
    closeLeftSidebar,
    openLeftSidebar,
    setSidebarView,

    // 右侧面板
    rightPanelVisible,
    rightPanelActiveTab,
    toggleRightPanel,
    openRightPanel,
    closeRightPanel,
    setRightPanelTab,

    // 设置抽屉
    settingsDrawerVisible,
    openSettings,
    closeSettings,

    // 通用状态
    activeCharacterId,
    activeChatId,
    isLoading,
    setLoading,
    setActiveCharacter,
    setActiveChat
  };
});
