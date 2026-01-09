<template>
  <div class="sidebar-container">
    <!-- Logo区域 -->
    <div class="sidebar-logo" :class="{ 'collapsed': appStore.sidebarCollapsed }">
      <div class="logo-icon">
        <div class="logo-glow"></div>
        <span class="logo-text">M</span>
      </div>
      <Transition name="fade">
        <span v-if="!appStore.sidebarCollapsed" class="logo-title">Muse</span>
      </Transition>
    </div>

    <!-- 导航菜单 -->
    <n-scrollbar class="sidebar-nav">
      <n-menu
        :value="activeKey"
        :collapsed="appStore.sidebarCollapsed"
        :collapsed-width="72"
        :collapsed-icon-size="24"
        :options="menuOptions"
        @update:value="handleMenuSelect"
      />
    </n-scrollbar>

    <!-- 底部操作区 -->
    <div class="sidebar-footer" :class="{ 'collapsed': appStore.sidebarCollapsed }">
      <n-tooltip v-if="appStore.sidebarCollapsed" trigger="hover" placement="right">
        <template #trigger>
          <n-button quaternary circle size="large" @click="appStore.toggleSidebar">
            <template #icon>
              <n-icon size="20">
                <ChevronForwardOutline />
              </n-icon>
            </template>
          </n-button>
        </template>
        展开侧边栏
      </n-tooltip>

      <n-button
        v-else
        quaternary
        class="collapse-btn"
        @click="appStore.toggleSidebar"
      >
        <template #icon>
          <n-icon size="18">
            <ChevronBackOutline />
          </n-icon>
        </template>
        收起
      </n-button>

      <!-- 版本信息 -->
      <Transition name="fade">
        <div v-if="!appStore.sidebarCollapsed" class="version-info">
          <span>v1.0.0</span>
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, h } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { NMenu, NScrollbar, NButton, NIcon, NTooltip } from 'naive-ui';
import {
  ChatbubblesOutline,
  PeopleOutline,
  BookOutline,
  DocumentTextOutline,
  CodeSlashOutline,
  SettingsOutline,
  ChevronBackOutline,
  ChevronForwardOutline
} from '@vicons/ionicons5';

import { useAppStore } from '../../stores/app';

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();

// 当前激活的菜单项
const activeKey = computed(() => route.name as string);

// 渲染图标
const renderIcon = (icon: typeof ChatbubblesOutline) => {
  return () => h(NIcon, null, { default: () => h(icon) });
};

// 菜单配置
const menuOptions = [
  {
    label: '聊天',
    key: 'Chat',
    icon: renderIcon(ChatbubblesOutline)
  },
  {
    label: '角色',
    key: 'Characters',
    icon: renderIcon(PeopleOutline)
  },
  {
    label: '世界书',
    key: 'WorldInfo',
    icon: renderIcon(BookOutline)
  },
  {
    label: '预设',
    key: 'Presets',
    icon: renderIcon(DocumentTextOutline)
  },
  {
    label: '正则',
    key: 'Regex',
    icon: renderIcon(CodeSlashOutline)
  },
  {
    type: 'divider',
    key: 'divider'
  },
  {
    label: '设置',
    key: 'Settings',
    icon: renderIcon(SettingsOutline)
  }
];

// 菜单选择处理
const handleMenuSelect = (key: string) => {
  router.push({ name: key });

  // 移动端选择后关闭侧边栏
  if (appStore.isMobile) {
    appStore.closeSidebar();
  }
};
</script>

<style scoped>
.sidebar-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 16px 12px;
}

/* Logo区域 */
.sidebar-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  margin-bottom: 24px;
  transition: all var(--transition-normal);
}

.sidebar-logo.collapsed {
  justify-content: center;
  padding: 8px 0;
}

.logo-icon {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-primary);
  border-radius: 10px;
  flex-shrink: 0;
}

.logo-glow {
  position: absolute;
  inset: -2px;
  background: var(--gradient-primary);
  border-radius: 12px;
  opacity: 0.5;
  filter: blur(8px);
  z-index: -1;
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  color: #000;
}

.logo-title {
  font-size: 24px;
  font-weight: 700;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: 1px;
}

/* 导航区域 */
.sidebar-nav {
  flex: 1;
  margin: 0 -12px;
  padding: 0 4px;
}

/* 底部区域 */
.sidebar-footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
  margin-top: 16px;
}

.sidebar-footer:not(.collapsed) {
  align-items: stretch;
}

.collapse-btn {
  justify-content: flex-start;
  width: 100%;
  color: var(--text-tertiary);
}

.collapse-btn:hover {
  color: var(--color-primary);
}

.version-info {
  text-align: center;
  font-size: 12px;
  color: var(--text-tertiary);
  padding: 4px 0;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
