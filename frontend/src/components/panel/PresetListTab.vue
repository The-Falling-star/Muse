<template>
  <div class="preset-list-tab">
    <!-- 当前使用的预设 -->
    <div class="current-preset">
      <div class="section-label">当前使用</div>
      <n-select
        v-model:value="currentPresetId"
        :options="presetOptions"
        placeholder="选择预设"
        size="small"
      />
    </div>

    <!-- 操作按钮 -->
    <div class="action-row">
      <n-button size="small" @click="handleImport">
        <template #icon>
          <n-icon><CloudUploadOutline /></n-icon>
        </template>
        导入
      </n-button>
      <n-button size="small" type="primary" @click="handleCreate">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        新建
      </n-button>
    </div>

    <!-- 预设列表 -->
    <div class="list-container">
      <div
        v-for="preset in presets"
        :key="preset.id"
        class="list-item"
        :class="{ active: isEditing(preset.id) }"
        @click="openPresetEditor(preset.id)"
      >
        <div class="item-info">
          <div class="item-name">
            <span v-if="preset.id === currentPresetId" class="star-icon">★</span>
            {{ preset.name }}
          </div>
          <div class="item-desc">
            {{ preset.promptCount }}个Prompt项{{ preset.regexCount > 0 ? `, ${preset.regexCount}个正则` : '' }}
          </div>
        </div>
        <n-dropdown
          trigger="click"
          :options="itemMenuOptions"
          @select="(key: string) => handleMenuSelect(key, preset.id)"
          @click.stop
        >
          <n-button quaternary circle size="tiny" @click.stop>
            <template #icon>
              <n-icon :size="16"><EllipsisHorizontal /></n-icon>
            </template>
          </n-button>
        </n-dropdown>
      </div>

      <!-- 空状态 -->
      <n-empty v-if="presets.length === 0" description="暂无预设" size="small" class="empty-state" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { NSelect, NButton, NIcon, NDropdown, NEmpty } from 'naive-ui';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';

const router = useRouter();
const route = useRoute();

// TODO: 后续对接真实的预设 store
const currentPresetId = ref<string | null>('default');

interface PresetItem {
  id: string;
  name: string;
  promptCount: number;
  regexCount: number;
}

const presets = ref<PresetItem[]>([
  { id: 'default', name: 'Default', promptCount: 6, regexCount: 2 },
  { id: 'custom1', name: 'My Custom Preset', promptCount: 4, regexCount: 0 },
  { id: 'roleplay', name: 'Roleplay v2', promptCount: 8, regexCount: 3 }
]);

const presetOptions = computed(() =>
  presets.value.map((p) => ({ label: p.name, value: p.id }))
);

const itemMenuOptions = [
  { label: '设为当前使用', key: 'use' },
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

const isEditing = (presetId: string): boolean => {
  return route.name === 'PresetEditor' && route.params.id === presetId;
};

const openPresetEditor = (presetId: string) => {
  router.push(`/preset/${presetId}`);
};

const handleImport = () => {
  // TODO: 弹出文件选择器，导入预设 JSON
};

const handleCreate = () => {
  // TODO: 创建新预设并跳转编辑
};

const handleMenuSelect = (key: string, presetId: string) => {
  switch (key) {
    case 'use':
      currentPresetId.value = presetId;
      break;
    case 'copy':
      // TODO: 复制预设
      break;
    case 'export':
      // TODO: 导出预设
      break;
    case 'delete':
      // TODO: 删除预设（需确认）
      break;
  }
};
</script>

<style scoped>
.preset-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.action-row {
  display: flex;
  gap: 8px;
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 150ms;
  border: 1px solid transparent;
}

.list-item:hover {
  background: var(--bg-hover);
}

.list-item.active {
  background: var(--bg-active);
  border-color: var(--color-primary);
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.star-icon {
  color: var(--color-primary);
  margin-right: 4px;
}

.item-desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.empty-state {
  padding: 32px 0;
}
</style>
