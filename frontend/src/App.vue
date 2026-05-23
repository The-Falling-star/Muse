<script setup lang="ts">
import { computed } from 'vue';
import { NConfigProvider, NMessageProvider } from 'naive-ui';
import { darkTheme } from 'naive-ui';
import { useThemeStore } from './stores/theme';
import { darkThemeOverrides, lightThemeOverrides } from './styles/theme';

const themeStore = useThemeStore();

// 计算当前 naive-ui 主题
const naiveTheme = computed(() => themeStore.theme === 'dark' ? darkTheme : null);

// 计算主题覆盖
const themeOverrides = computed(() =>
  themeStore.theme === 'dark' ? darkThemeOverrides : lightThemeOverrides
);
</script>

<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <router-view />
    </n-message-provider>
  </n-config-provider>
</template>

<style>
/* 移除默认的Vite样式 */
</style>
