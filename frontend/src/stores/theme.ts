import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

export type ThemeMode = 'dark' | 'light';

export const useThemeStore = defineStore('theme', () => {
  // 从localStorage读取主题设置，默认为dark
  const getInitialTheme = (): ThemeMode => {
    const saved = localStorage.getItem('muse-theme');
    if (saved === 'light' || saved === 'dark') {
      return saved;
    }
    // 检测系统主题偏好
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
      return 'light';
    }
    return 'dark';
  };

  const theme = ref<ThemeMode>(getInitialTheme());

  // 切换主题
  const toggleTheme = () => {
    theme.value = theme.value === 'dark' ? 'light' : 'dark';
  };

  // 设置主题
  const setTheme = (newTheme: ThemeMode) => {
    theme.value = newTheme;
  };

  // 监听主题变化，保存到localStorage并更新document
  watch(theme, (newTheme) => {
    localStorage.setItem('muse-theme', newTheme);
    document.documentElement.setAttribute('data-theme', newTheme);
  }, { immediate: true });

  return {
    theme,
    toggleTheme,
    setTheme
  };
});
