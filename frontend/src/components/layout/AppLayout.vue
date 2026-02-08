<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <div class="app-layout" :class="{ 'sidebar-collapsed': appStore.sidebarCollapsed }">
              <!-- 科幻背景效果 -->
              <div class="sci-fi-grid-bg"></div>
              <div class="sci-fi-grid-dots"></div>
              <div class="sci-fi-particles">
                <div
                  v-for="i in 20"
                  :key="i"
                  class="sci-fi-particle"
                  :style="{
                    left: `${Math.random() * 100}%`,
                    animationDelay: `${Math.random() * 15}s`,
                    animationDuration: `${10 + Math.random() * 10}s`
                  }"
                ></div>
              </div>
              <!-- 角落装饰 -->
              <div class="sci-fi-corner sci-fi-corner--tl"></div>
              <div class="sci-fi-corner sci-fi-corner--tr"></div>
              <div class="sci-fi-corner sci-fi-corner--bl"></div>
              <div class="sci-fi-corner sci-fi-corner--br"></div>

              <!-- 移动端遮罩 -->
              <Transition name="fade">
                <div
                  v-if="appStore.isMobile && appStore.sidebarVisible"
                  class="sidebar-overlay"
                  @click="appStore.closeSidebar"
                ></div>
              </Transition>

              <!-- 侧边栏 -->
              <Transition name="slide-sidebar">
                <aside
                  v-show="!appStore.isMobile || appStore.sidebarVisible"
                  class="app-sidebar gpu-accelerated"
                  :class="{ 'collapsed': appStore.sidebarCollapsed }"
                >
                  <AppSidebar />
                </aside>
              </Transition>

              <!-- 主内容区 -->
              <main class="app-main">
                <!-- 顶部栏 -->
                <header class="app-header">
                  <div class="header-left">
                    <n-button
                      quaternary
                      circle
                      class="menu-toggle"
                      @click="appStore.toggleSidebar"
                    >
                      <template #icon>
                        <n-icon size="22">
                          <MenuOutline />
                        </n-icon>
                      </template>
                    </n-button>

                    <div class="page-title">
                      <h1 class="title-text">{{ pageTitle }}</h1>
                      <div class="title-glow"></div>
                    </div>
                  </div>

                  <div class="header-right">
                    <!-- 主题切换 -->
                    <n-tooltip trigger="hover">
                      <template #trigger>
                        <n-button quaternary circle @click="themeStore.toggleTheme">
                          <template #icon>
                            <n-icon size="20">
                              <SunnyOutline v-if="themeStore.theme === 'dark'" />
                              <MoonOutline v-else />
                            </n-icon>
                          </template>
                        </n-button>
                      </template>
                      {{ themeStore.theme === 'dark' ? '切换到明亮模式' : '切换到暗黑模式' }}
                    </n-tooltip>

                    <!-- 用户菜单 -->
                    <n-dropdown :options="userMenuOptions" @select="handleUserMenu">
                      <n-button quaternary circle>
                        <template #icon>
                          <n-icon size="20">
                            <PersonCircleOutline />
                          </n-icon>
                        </template>
                      </n-button>
                    </n-dropdown>
                  </div>
                </header>

                <!-- 页面内容 -->
                <div class="app-content">
                  <router-view v-slot="{ Component }">
                    <Transition name="page-fade" mode="out-in">
                      <component :is="Component" class="gpu-accelerated" />
                    </Transition>
                  </router-view>
                </div>
              </main>
            </div>
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, h, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { darkTheme } from 'naive-ui';
import {
  NConfigProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  NButton,
  NIcon,
  NTooltip,
  NDropdown
} from 'naive-ui';
import {
  MenuOutline,
  SunnyOutline,
  MoonOutline,
  PersonCircleOutline,
  SettingsOutline,
  LogOutOutline
} from '@vicons/ionicons5';

import AppSidebar from './AppSidebar.vue';
import { useThemeStore } from '@/stores/theme.ts';
import { useAppStore } from '@/stores/app.ts';
import { darkThemeOverrides, lightThemeOverrides } from '@/styles/theme.ts';
import { initGlobalMessage } from '@/composables/useGlobalMessage.ts';

const route = useRoute();
const themeStore = useThemeStore();
const appStore = useAppStore();

// 计算当前naive-ui主题
const naiveTheme = computed(() => themeStore.theme === 'dark' ? darkTheme : null);

// 计算主题覆盖
const themeOverrides = computed(() =>
  themeStore.theme === 'dark' ? darkThemeOverrides : lightThemeOverrides
);

// 页面标题
const pageTitle = computed(() => {
  const title = route.meta.title as string;
  return title || 'Muse';
});

// 渲染图标函数
const renderIcon = (icon: typeof SettingsOutline) => {
  return () => h(NIcon, null, { default: () => h(icon) });
};

// 用户菜单选项
const userMenuOptions = [
  {
    label: '设置',
    key: 'settings',
    icon: renderIcon(SettingsOutline)
  },
  {
    type: 'divider',
    key: 'd1'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogOutOutline)
  }
];

const handleUserMenu = (key: string) => {
  console.log('User menu:', key);
};

// 初始化全局消息（必须在 NMessageProvider 内部调用 useMessage）
onMounted(() => {
  console.log('[AppLayout] Initializing global message...');
  initGlobalMessage();
  console.log('[AppLayout] Global message initialized!');
});
</script>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  position: relative;
  background: var(--bg-primary);
}

/* 侧边栏 */
.app-sidebar {
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  width: 260px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  z-index: 100;
  transition: width var(--transition-normal), transform var(--transition-normal);
  overflow: hidden;
}

.app-sidebar.collapsed {
  width: 72px;
}

/* 移动端侧边栏 */
@media (max-width: 768px) {
  .app-sidebar {
    transform: translateX(-100%);
    width: 280px;
    box-shadow: var(--shadow-xl);
    transition: transform var(--transition-mobile);
  }

  .app-sidebar:not(.collapsed) {
    transform: translateX(0);
  }
}

/* 侧边栏遮罩 */
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: var(--bg-overlay);
  z-index: 99;
  backdrop-filter: blur(4px);
  transition: opacity var(--transition-fast);
}

/* 主内容区 */
.app-main {
  flex: 1;
  margin-left: 260px;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  min-width: 0;
  transition: margin-left var(--transition-normal);
}

.sidebar-collapsed .app-main {
  margin-left: 72px;
}

@media (max-width: 768px) {
  .app-main {
    margin-left: 0;
    transition: none;
  }
}

/* 顶部栏 */
.app-header {
  position: sticky;
  top: 0;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  z-index: 50;
  backdrop-filter: blur(10px);
  min-width: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-toggle {
  color: var(--text-secondary);
}

.menu-toggle:hover {
  color: var(--color-primary);
}

/* 页面标题 */
.page-title {
  position: relative;
}

.title-text {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
  letter-spacing: 0.5px;
}

.title-glow {
  position: absolute;
  bottom: -4px;
  left: 0;
  width: 40px;
  height: 2px;
  background: var(--gradient-primary);
  border-radius: 2px;
  box-shadow: var(--glow-soft);
}

/* 页面内容 */
.app-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

@media (max-width: 768px) {
  .app-header {
    padding: 0 16px;
  }

  .app-content {
    padding: 16px;
  }

  .title-text {
    font-size: 18px;
  }
}

</style>
