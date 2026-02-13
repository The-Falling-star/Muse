<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <div class="app-layout">
              <!-- ===== PC端：推挤式左侧边栏 ===== -->
              <aside
                v-if="appStore.isDesktop"
                class="left-sidebar-wrapper"
                :class="{ collapsed: !appStore.leftSidebarVisible }"
              >
                <div class="sidebar-inner">
                  <LeftSidebar />
                </div>
              </aside>

              <!-- ===== 平板/移动端：Drawer式左侧边栏 ===== -->
              <n-drawer
                v-if="!appStore.isDesktop"
                v-model:show="appStore.leftSidebarVisible"
                :width="appStore.isMobile ? '80%' : '60%'"
                placement="left"
                :trap-focus="false"
                :block-scroll="true"
                :native-scrollbar="true"
              >
                <n-drawer-content
                  :body-content-style="{ padding: 0, height: '100%' }"
                  :content-style="{ height: '100%' }"
                >
                  <LeftSidebar />
                </n-drawer-content>
              </n-drawer>

              <!-- ===== 中间主区域 ===== -->
              <main class="main-content">
                <TopBar />
                <div class="main-view">
                  <router-view v-slot="{ Component, route: viewRoute }">
                    <transition name="fade" mode="out-in">
                      <component :is="Component" :key="viewRoute.fullPath" />
                    </transition>
                  </router-view>
                </div>
              </main>

              <!-- ===== PC端：推挤式右侧面板 ===== -->
              <aside
                v-if="appStore.isDesktop"
                class="right-panel-wrapper"
                :class="{ collapsed: !appStore.rightPanelVisible }"
              >
                <div class="panel-inner">
                  <RightPanel />
                </div>
              </aside>

              <!-- ===== 平板/移动端：Drawer式右侧面板 ===== -->
              <n-drawer
                v-if="!appStore.isDesktop"
                v-model:show="appStore.rightPanelVisible"
                :width="appStore.isMobile ? '80%' : '60%'"
                placement="right"
                :trap-focus="false"
                :block-scroll="true"
                :native-scrollbar="true"
              >
                <n-drawer-content
                  :body-content-style="{ padding: 0, height: '100%' }"
                  :content-style="{ height: '100%' }"
                >
                  <RightPanel />
                </n-drawer-content>
              </n-drawer>

              <!-- ===== 设置抽屉 ===== -->
              <SettingsDrawer
                v-model:visible="appStore.settingsDrawerVisible"
              />
            </div>
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue';
import { darkTheme } from 'naive-ui';
import {
  NConfigProvider,
  NLoadingBarProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  NDrawer,
  NDrawerContent
} from 'naive-ui';

import TopBar from './TopBar.vue';
import LeftSidebar from './LeftSidebar.vue';
import RightPanel from './RightPanel.vue';
import SettingsDrawer from './SettingsDrawer.vue';
import { useThemeStore } from '@/stores/theme';
import { useAppStore } from '@/stores/app';
import { darkThemeOverrides, lightThemeOverrides } from '@/styles/theme';
import { initGlobalMessage } from '@/composables/useGlobalMessage';

const themeStore = useThemeStore();
const appStore = useAppStore();

// 计算当前 naive-ui 主题
const naiveTheme = computed(() => themeStore.theme === 'dark' ? darkTheme : null);

// 计算主题覆盖
const themeOverrides = computed(() =>
  themeStore.theme === 'dark' ? darkThemeOverrides : lightThemeOverrides
);

// 初始化全局消息
onMounted(() => {
  initGlobalMessage();
});
</script>

<style scoped>
/* ====== 三栏布局框架 ====== */
.app-layout {
  display: flex;
  height: 100vh;
  max-height: 100vh;
  width: 100%;
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
}

/* ====== 左侧边栏（PC推挤式） ====== */
.left-sidebar-wrapper {
  width: 300px;
  min-width: 300px;
  height: 100vh;
  position: sticky;
  top: 0;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  z-index: 30;
  overflow: hidden;
  transition: width 300ms cubic-bezier(0.4, 0, 0.2, 1),
              min-width 300ms cubic-bezier(0.4, 0, 0.2, 1),
              border-color 300ms cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 300ms cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 1px 0 12px rgba(77, 168, 255, .03);
}

.left-sidebar-wrapper.collapsed {
  width: 0;
  min-width: 0;
  border-right-color: transparent;
  box-shadow: none;
}

.sidebar-inner {
  width: 300px;
  height: 100%;
}

/* ====== 中间主区域 ====== */
.main-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.main-view {
  flex: 1;
  overflow: hidden;
  min-height: 0;
}

/* ====== 右侧面板（PC推挤式） ====== */
.right-panel-wrapper {
  width: 300px;
  min-width: 300px;
  height: 100vh;
  position: sticky;
  top: 0;
  background: var(--bg-secondary);
  border-left: 1px solid var(--border-color);
  z-index: 30;
  overflow: hidden;
  transition: width 300ms cubic-bezier(0.4, 0, 0.2, 1),
              min-width 300ms cubic-bezier(0.4, 0, 0.2, 1),
              border-color 300ms cubic-bezier(0.4, 0, 0.2, 1),
              box-shadow 300ms cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: -1px 0 12px rgba(77, 168, 255, .03);
}

.right-panel-wrapper.collapsed {
  width: 0;
  min-width: 0;
  border-left-color: transparent;
  box-shadow: none;
}

.panel-inner {
  width: 300px;
  height: 100%;
}

/* ====== 视图切换过渡动画 ====== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 200ms ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ====== 响应式 ====== */

/* 平板端 */
@media (min-width: 768px) and (max-width: 1023px) {
  .main-content {
    width: 100%;
  }
}

/* 移动端 */
@media (max-width: 767px) {
  .main-content {
    width: 100%;
  }
}
</style>
