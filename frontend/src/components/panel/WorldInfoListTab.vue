<template>
  <div class="worldinfo-list-tab">
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

    <!-- 世界书列表 -->
    <div class="list-container">
      <div
        v-for="wi in worldInfos"
        :key="wi.id"
        class="list-item"
        :class="{ active: isEditing(wi.id) }"
        @click="openWorldInfoEditor(wi.id)"
      >
        <n-checkbox
          :checked="wi.enabled"
          @update:checked="(val: boolean) => toggleEnabled(wi.id, val)"
          @click.stop
        />
        <div class="item-info">
          <div class="item-name">{{ wi.name }}</div>
          <div class="item-desc">{{ wi.entryCount }}个词条</div>
        </div>
        <n-dropdown
          trigger="click"
          :options="itemMenuOptions"
          @select="(key: string) => handleMenuSelect(key, wi.id)"
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
      <n-empty v-if="worldInfos.length === 0" description="暂无世界书" size="small" class="empty-state" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { NButton, NIcon, NCheckbox, NDropdown, NEmpty } from 'naive-ui';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';

const router = useRouter();
const route = useRoute();

// TODO: 后续对接真实的世界书 store
interface WorldInfoItem {
  id: string;
  name: string;
  entryCount: number;
  enabled: boolean;
}

const worldInfos = ref<WorldInfoItem[]>([
  { id: 'main', name: '主世界设定', entryCount: 12, enabled: true },
  { id: 'char', name: '角色专属世界书', entryCount: 5, enabled: true },
  { id: 'backup', name: '备用设定', entryCount: 8, enabled: false }
]);

const itemMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

const isEditing = (wiId: string): boolean => {
  return route.name === 'WorldInfoEditor' && route.params.id === wiId;
};

const openWorldInfoEditor = (wiId: string) => {
  router.push(`/worldinfo/${wiId}`);
};

const toggleEnabled = (wiId: string, enabled: boolean) => {
  const item = worldInfos.value.find((w) => w.id === wiId);
  if (item) {
    item.enabled = enabled;
  }
};

const handleImport = () => {
  // TODO: 弹出文件选择器，导入世界书 JSON
};

const handleCreate = () => {
  // TODO: 创建新世界书并跳转编辑
};

const handleMenuSelect = (key: string, wiId: string) => {
  switch (key) {
    case 'copy':
      // TODO: 复制世界书
      break;
    case 'export':
      // TODO: 导出世界书
      break;
    case 'delete':
      // TODO: 删除世界书（需确认）
      break;
  }
};
</script>

<style scoped>
.worldinfo-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
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
  gap: 10px;
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

.item-desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.empty-state {
  padding: 32px 0;
}
</style>
