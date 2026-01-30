<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { NConfigProvider, NMessageProvider } from 'naive-ui';
import { darkTheme } from 'naive-ui';
import { useThemeStore } from './stores/theme';
import { darkThemeOverrides, lightThemeOverrides } from './styles/theme';
import AppLayout from './components/layout/AppLayout.vue';

const route = useRoute();
const themeStore = useThemeStore();

// 判断是否为登录页面
const isLoginPage = computed(() => route.name === 'Login');

// 计算当前naive-ui主题
const naiveTheme = computed(() => themeStore.theme === 'dark' ? darkTheme : null);

// 计算主题覆盖
const themeOverrides = computed(() =>
  themeStore.theme === 'dark' ? darkThemeOverrides : lightThemeOverrides
);
</script>

<template>
  <n-config-provider :theme="naiveTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <!-- 登录页面单独渲染 -->
      <router-view v-if="isLoginPage" />
      <!-- 其他页面使用AppLayout布局 -->
      <AppLayout v-else />
    </n-message-provider>
  </n-config-provider>
</template>

<style>
/* 移除默认的Vite样式 */
</style>
