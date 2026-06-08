<template>
  <aside class="right-panel">
    <!-- 面板标题栏 -->
    <div class="panel-header">
      <span class="panel-title">配置面板</span>
      <n-button quaternary circle size="small" @click="appStore.closeRightPanel">
        <template #icon>
          <n-icon :size="18">
            <Close />
          </n-icon>
        </template>
      </n-button>
    </div>

    <!-- 选项卡导航 -->
    <div class="panel-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ active: appStore.rightPanelActiveTab === tab.key }"
        @click="appStore.setRightPanelTab(tab.key)"
      >
        <n-icon :size="16">
          <component :is="tab.icon" />
        </n-icon>
        <span class="tab-label">{{ tab.label }}</span>
      </button>
    </div>

    <!-- 选项卡内容区域 -->
    <div class="panel-content">
      <!-- 模型配置 Tab -->
      <div v-if="appStore.rightPanelActiveTab === 'model'" class="tab-content">
        <div class="tab-scroll">
          <ModelConfigTab />
        </div>
      </div>

      <!-- 预设列表 Tab -->
      <div v-else-if="appStore.rightPanelActiveTab === 'preset'" class="tab-content">
        <PresetListTab />
      </div>

      <!-- 世界书列表 Tab -->
      <div v-else-if="appStore.rightPanelActiveTab === 'worldinfo'" class="tab-content">
        <WorldInfoListTab />
      </div>

      <!-- 正则列表 Tab -->
      <div v-else-if="appStore.rightPanelActiveTab === 'regex'" class="tab-content">
        <RegexListTab />
      </div>

      <!-- 用户设置 Tab -->
      <div v-else-if="appStore.rightPanelActiveTab === 'user'" class="tab-content">
        <div class="tab-scroll">
          <UserConfigTab/>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app';
import type { RightPanelTab } from '@/stores/app';
import { NButton, NIcon } from 'naive-ui';
import {
  Close,
  SettingsOutline,
  DocumentTextOutline,
  BookOutline,
  CodeSlashOutline,
} from '@vicons/ionicons5';
import {
  Bot20Regular,
} from '@vicons/fluent'
import ModelConfigTab from '@/components/panel/ModelConfigTab.vue';
import PresetListTab from '@/components/panel/PresetListTab.vue';
import WorldInfoListTab from '@/components/panel/WorldInfoListTab.vue';
import RegexListTab from '@/components/panel/RegexListTab.vue';
import UserConfigTab from "@/components/panel/UserConfigTab.vue";

const appStore = useAppStore();

interface TabItem {
  key: RightPanelTab;
  label: string;
  icon: typeof SettingsOutline;
}

const tabs: TabItem[] = [
  { key: 'model', label: '模型', icon:  Bot20Regular},
  { key: 'preset', label: '预设', icon: DocumentTextOutline },
  { key: 'worldinfo', label: '世界书', icon: BookOutline },
  { key: 'regex', label: '正则', icon: CodeSlashOutline },
  {key: 'user', label: '用户', icon: SettingsOutline},
];
</script>

<style scoped>
.right-panel {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-left: 1px solid var(--border-color);
  overflow: hidden;
  position: relative;
}

/* 科幻背景纹理 */
.right-panel::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--gradient-surface);
  pointer-events: none;
  z-index: 0;
}

.right-panel > * {
  position: relative;
  z-index: 1;
}

/* 面板标题栏 */
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.panel-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

/* 选项卡导航 */
.panel-tabs {
  display: flex;
  padding: 8px 12px 0;
  gap: 4px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 4px 10px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 8px 8px 0 0;
  transition: color 150ms, background-color 150ms;
  position: relative;
  font-size: 12px;
  line-height: 1;
}

.tab-item:hover {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.tab-item.active {
  color: var(--color-primary);
  text-shadow: 0 0 12px rgba(77, 168, 255, .3);
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 8px;
  right: 8px;
  height: 2px;
  background: var(--color-primary);
  border-radius: 2px 2px 0 0;
  box-shadow: 0 0 8px rgba(77, 168, 255, .4);
}

.tab-label {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
}

/* 选项卡内容区域 */
.panel-content {
  flex: 1;
  overflow: hidden;
}

.tab-content {
  height: 100%;
  overflow: hidden;
}

.tab-scroll {
  height: 100%;
  overflow-y: auto;
}

/* 自定义滚动条 */
.panel-content::-webkit-scrollbar {
  width: 6px;
}

.panel-content::-webkit-scrollbar-track {
  background: transparent;
}

.panel-content::-webkit-scrollbar-thumb {
  background: var(--scrollbar-color);
  border-radius: 3px;
}

.panel-content::-webkit-scrollbar-thumb:hover {
  background: var(--scrollbar-hover-color);
}
</style>
