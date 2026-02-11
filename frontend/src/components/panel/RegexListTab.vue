<template>
  <div class="regex-list-tab">
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
        新建规则
      </n-button>
    </div>

    <!-- 正则列表 -->
    <div class="list-container">
      <div
        v-for="rule in regexRules"
        :key="rule.id"
        class="list-item"
        :class="{ active: isEditing(rule.id) }"
        @click="openRegexEditor(rule.id)"
      >
        <n-checkbox
          :checked="rule.enabled"
          @update:checked="(val: boolean) => toggleEnabled(rule.id, val)"
          @click.stop
        />
        <div class="item-info">
          <div class="item-name">{{ rule.name }}</div>
          <div class="item-desc">
            <code class="regex-preview">{{ rule.findPreview }}</code>
            <span class="arrow">→</span>
            <span class="replace-preview">{{ rule.replacePreview || '(空)' }}</span>
          </div>
        </div>
        <n-dropdown
          trigger="click"
          :options="itemMenuOptions"
          @select="(key: string) => handleMenuSelect(key, rule.id)"
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
      <n-empty v-if="regexRules.length === 0" description="暂无正则规则" size="small" class="empty-state" />
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

// TODO: 后续对接真实的正则 store
interface RegexRuleItem {
  id: string;
  name: string;
  findPreview: string;
  replacePreview: string;
  enabled: boolean;
}

const regexRules = ref<RegexRuleItem[]>([
  { id: 'ooc', name: '删除OOC标记', findPreview: '/\\(OOC:.*?\\)/g', replacePreview: '', enabled: true },
  { id: 'think', name: '格式化思考', findPreview: '/<think>.*?/s', replacePreview: '...', enabled: true },
  { id: 'old', name: '旧规则', findPreview: '/old/', replacePreview: 'new', enabled: false }
]);

const itemMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

const isEditing = (ruleId: string): boolean => {
  return route.name === 'RegexEditor' && route.params.id === ruleId;
};

const openRegexEditor = (ruleId: string) => {
  router.push(`/regex/${ruleId}`);
};

const toggleEnabled = (ruleId: string, enabled: boolean) => {
  const item = regexRules.value.find((r) => r.id === ruleId);
  if (item) {
    item.enabled = enabled;
  }
};

const handleImport = () => {
  // TODO: 导入正则规则
};

const handleCreate = () => {
  // TODO: 创建新规则并跳转编辑
};

const handleMenuSelect = (key: string, ruleId: string) => {
  switch (key) {
    case 'copy':
      // TODO: 复制规则
      break;
    case 'export':
      // TODO: 导出规则
      break;
    case 'delete':
      // TODO: 删除规则（需确认）
      break;
  }
};
</script>

<style scoped>
.regex-list-tab {
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
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  overflow: hidden;
}

.regex-preview {
  font-family: var(--font-mono);
  font-size: 10px;
  background: var(--bg-tertiary);
  padding: 1px 4px;
  border-radius: 3px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 120px;
}

.arrow {
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.replace-preview {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-tertiary);
  font-size: 11px;
}

.empty-state {
  padding: 32px 0;
}
</style>
