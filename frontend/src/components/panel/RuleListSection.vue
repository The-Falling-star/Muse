<template>
  <div class="rule-list-section">
    <n-spin :show="loading" description="加载中..." class="spin-container">
      <div class="list-container">
        <div
          v-for="rule in rules"
          :key="rule.id"
          class="list-item"
          :class="{ active: isEditing(rule.id) }"
          @click="emit('edit', rule.id)"
        >
          <n-checkbox
            :checked="rule.isEnabled"
            @update:checked="() => emit('toggle', rule.id)"
            @click.stop
          />
          <div class="item-info">
            <div class="item-name">{{ rule.name }}</div>
            <div class="item-desc">
              <code class="regex-preview">{{ rule.findPattern }}</code>
              <span class="arrow">→</span>
              <span class="replace-preview">{{ rule.replacePattern || '(空)' }}</span>
            </div>
          </div>
          <n-dropdown
            trigger="click"
            :options="itemMenuOptions"
            @select="(key: string) => emit('menuAction', key, rule.id)"
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
        <n-empty
          v-if="rules.length === 0 && !loading"
          :description="emptyText"
          size="small"
          class="empty-state"
        />
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { NButton, NIcon, NCheckbox, NDropdown, NEmpty, NSpin } from 'naive-ui';
import { EllipsisHorizontal } from '@vicons/ionicons5';
import type { RegexRule } from '@/gen/muse/regex_pb';

const route = useRoute();

defineProps<{
  rules: RegexRule[];
  loading: boolean;
  emptyText: string;
}>();

const emit = defineEmits<{
  toggle: [ruleId: number];
  edit: [ruleId: number];
  menuAction: [key: string, ruleId: number];
}>();

const itemMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
] as const;

const isEditing = (ruleId: number): boolean =>
  route.name === 'RegexEditor' && Number(route.params.id) === ruleId;
</script>

<style scoped>
.rule-list-section {
  height: 100%;
  min-height: 0;
}

.spin-container {
  height: 100%;
}

.spin-container :deep(.n-spin-container),
.spin-container :deep(.n-spin-content) {
  height: 100%;
}

.list-container {
  height: 100%;
  overflow: auto;
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
