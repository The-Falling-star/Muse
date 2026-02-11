<template>
  <header class="top-bar" :class="{ 'top-bar--mobile': appStore.isMobile }">
    <!-- 左侧区域 -->
    <div class="top-bar__left">
      <!-- 汉堡菜单 -->
      <n-button quaternary circle size="small" class="top-bar__icon-btn" @click="appStore.toggleLeftSidebar">
        <template #icon>
          <n-icon size="20">
            <MenuOutline />
          </n-icon>
        </template>
      </n-button>

      <!-- Logo -->
      <div class="top-bar__logo" @click="handleLogoClick">
        <span class="top-bar__logo-icon">M</span>
        <span v-if="!appStore.isMobile" class="top-bar__logo-text">Muse</span>
      </div>

      <!-- 返回对话按钮（仅编辑视图时显示） -->
      <template v-if="isEditorView">
        <n-button text class="top-bar__back-btn" @click="backToChat">
          <template #icon>
            <n-icon size="18">
              <ArrowBackOutline />
            </n-icon>
          </template>
          返回对话
        </n-button>

        <span class="top-bar__view-title">{{ viewTitle }}</span>
      </template>
    </div>

    <!-- 中间区域 - 模型选择器（仅对话视图时显示） -->
    <div v-if="!isEditorView" class="top-bar__center">
      <n-button
        class="top-bar__model-selector"
        quaternary
        @click="handleModelSelectorClick"
      >
        <template #icon>
          <n-icon size="16">
            <SparklesOutline />
          </n-icon>
        </template>
        <span class="top-bar__model-name">{{ activeModelName }}</span>
        <n-icon size="14" class="top-bar__model-arrow">
          <ChevronDownOutline />
        </n-icon>
      </n-button>
    </div>

    <!-- 右侧区域 -->
    <div class="top-bar__right">
      <!-- 右侧面板切换按钮 -->
      <n-tooltip trigger="hover" :disabled="appStore.isMobile">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="small"
            class="top-bar__icon-btn"
            :class="{ 'top-bar__icon-btn--active': appStore.rightPanelVisible }"
            @click="appStore.toggleRightPanel"
          >
            <template #icon>
              <n-icon size="20">
                <OptionsOutline />
              </n-icon>
            </template>
          </n-button>
        </template>
        {{ appStore.rightPanelVisible ? '关闭配置面板' : '打开配置面板' }}
      </n-tooltip>

      <!-- 设置按钮 -->
      <n-tooltip trigger="hover" :disabled="appStore.isMobile">
        <template #trigger>
          <n-button quaternary circle size="small" class="top-bar__icon-btn" @click="appStore.openSettings">
            <template #icon>
              <n-icon size="20">
                <SettingsOutline />
              </n-icon>
            </template>
          </n-button>
        </template>
        设置
      </n-tooltip>

      <!-- 用户头像/菜单 -->
      <n-dropdown :options="userMenuOptions" @select="handleUserMenu" trigger="click">
        <n-button quaternary circle size="small" class="top-bar__icon-btn">
          <template #icon>
            <n-icon size="20">
              <PersonCircleOutline />
            </n-icon>
          </template>
        </n-button>
      </n-dropdown>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed, h } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  NButton,
  NIcon,
  NTooltip,
  NDropdown
} from 'naive-ui';
import {
  MenuOutline,
  ArrowBackOutline,
  SparklesOutline,
  ChevronDownOutline,
  OptionsOutline,
  SettingsOutline,
  PersonCircleOutline,
  SunnyOutline,
  MoonOutline,
  LogOutOutline
} from '@vicons/ionicons5';

import { useAppStore } from '@/stores/app';
import { useThemeStore } from '@/stores/theme';
import { useUserStore } from '@/stores/user';

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();
const themeStore = useThemeStore();
const userStore = useUserStore();

// 是否为编辑视图
const isEditorView = computed(() => {
  const name = route.name as string;
  return ['PresetEditor', 'WorldInfoEditor', 'RegexEditor'].includes(name);
});

// 编辑视图标题
const viewTitle = computed(() => {
  const title = route.meta.title as string;
  return title || '';
});

// 当前模型名称
const activeModelName = computed(() => {
  const config = userStore.activeApiConfig;
  if (!config) {
    return '选择模型';
  }
  return config.model || config.name || '选择模型';
});

// 点击Logo回到对话
const handleLogoClick = () => {
  if (isEditorView.value) {
    router.push('/');
  }
};

// 返回对话
const backToChat = () => {
  router.push('/');
};

// 模型选择器点击
const handleModelSelectorClick = () => {
  appStore.openRightPanel('model');
};

// 渲染图标
const renderIcon = (icon: typeof SunnyOutline) => {
  return () => h(NIcon, null, { default: () => h(icon) });
};

// 用户菜单选项
const userMenuOptions = computed(() => [
  {
    label: themeStore.theme === 'dark' ? '切换到明亮模式' : '切换到暗黑模式',
    key: 'toggleTheme',
    icon: renderIcon(themeStore.theme === 'dark' ? SunnyOutline : MoonOutline)
  },
  {
    type: 'divider' as const,
    key: 'd1'
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogOutOutline)
  }
]);

// 用户菜单操作
const handleUserMenu = (key: string) => {
  if (key === 'toggleTheme') {
    themeStore.toggleTheme();
    return;
  }
  if (key === 'logout') {
    userStore.logout();
    router.push('/login');
  }
};
</script>

<style scoped>
.top-bar {
  position: sticky;
  top: 0;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  z-index: 100;
  backdrop-filter: blur(12px);
  flex-shrink: 0;
  background-image: var(--gradient-surface);
}

.top-bar--mobile {
  height: 48px;
  padding: 0 12px;
}

/* 左侧区域 */
.top-bar__left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex-shrink: 0;
}

/* Logo */
.top-bar__logo {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  user-select: none;
  flex-shrink: 0;
}

.top-bar__logo-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-primary);
  border-radius: 8px;
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  box-shadow: var(--glow-primary-sm);
  transition: box-shadow var(--transition-fast);
}

.top-bar__logo:hover .top-bar__logo-icon {
  box-shadow: var(--glow-primary);
}

.top-bar__logo-text {
  font-size: 18px;
  font-weight: 600;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: .5px;
}

/* 返回按钮 */
.top-bar__back-btn {
  font-size: 14px;
  color: var(--color-primary);
  margin-left: 8px;
}

/* 视图标题 */
.top-bar__view-title {
  font-size: 14px;
  color: var(--text-secondary);
  margin-left: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 中间区域 */
.top-bar__center {
  flex: 1;
  display: flex;
  justify-content: center;
  min-width: 0;
}

/* 模型选择器 */
.top-bar__model-selector {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid var(--border-light);
  font-size: 14px;
  color: var(--text-primary);
  max-width: 280px;
  transition: all var(--transition-fast);
  background: var(--gradient-glow);
}

.top-bar__model-selector:hover {
  border-color: var(--color-primary);
  background: var(--bg-hover);
  box-shadow: var(--glow-primary-sm);
}

.top-bar__model-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.top-bar__model-arrow {
  margin-left: 4px;
  flex-shrink: 0;
}

/* 右侧区域 */
.top-bar__right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

/* 图标按钮通用 */
.top-bar__icon-btn {
  color: var(--text-secondary);
  transition: color var(--transition-fast);
}

.top-bar__icon-btn:hover {
  color: var(--text-primary);
}

.top-bar__icon-btn--active {
  color: var(--color-primary);
  text-shadow: 0 0 10px rgba(77, 168, 255, .4);
}

/* 移动端适配 */
@media (max-width: 767px) {
  .top-bar__model-selector {
    max-width: 160px;
    padding: 4px 12px;
    font-size: 13px;
  }

  .top-bar__view-title {
    max-width: 120px;
    font-size: 13px;
  }
}
</style>
