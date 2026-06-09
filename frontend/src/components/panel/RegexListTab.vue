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
    <n-spin :show="loading" description="加载中..." style="min-height: 60px;">
      <n-infinite-scroll class="list-container" @load="loadMore" :distance="100">
        <div
          v-for="rule in sortedRules"
          :key="rule.id"
          class="list-item"
          :class="{ active: isEditing(rule.id) }"
          @click="openRegexEditor(rule.id)"
        >
          <n-checkbox
            :checked="rule.isEnabled"
            @update:checked="() => toggleEnabled(rule.id)"
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

        <!-- 加载更多指示 -->
        <div v-if="loadingMore" class="loading-more">
          <n-spin size="small" />
        </div>

        <!-- 空状态 -->
        <n-empty v-if="rules.length === 0 && !loading" description="暂无正则规则" size="small" class="empty-state" />
      </n-infinite-scroll>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { NButton, NIcon, NCheckbox, NDropdown, NEmpty, NSpin, useDialog, useMessage } from 'naive-ui';
import { NInfiniteScroll } from 'naive-ui/es/infinite-scroll';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';
import { useRegexRuleStore } from '@/stores/regexRule';

const router = useRouter();
const route = useRoute();
const dialog = useDialog();
const message = useMessage();
const regexRuleStore = useRegexRuleStore();

const rules = computed(() => regexRuleStore.rules || []);
const sortedRules = computed(() => regexRuleStore.sortedRules || []);
const loading = computed(() => regexRuleStore.loading);
const loadingMore = computed(() => regexRuleStore.loadingMore);
const hasMore = computed(() => regexRuleStore.hasMore);

const loadMore = () => {
  if (!hasMore.value || loadingMore.value) return Promise.resolve();
  return regexRuleStore.loadMore();
};

// 菜单选项
const itemMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

onMounted(async () => {
  await regexRuleStore.loadRules();
});

// 判断是否正在编辑
const isEditing = (ruleId: number): boolean => {
  return route.name === 'RegexEditor' && Number(route.params.id) === ruleId;
};

// 打开正则编辑器
const openRegexEditor = (ruleId: number) => {
  router.push(`/regex/${ruleId}`);
};

// 切换启用/禁用
const toggleEnabled = async (ruleId: number) => {
  try {
    await regexRuleStore.toggleRuleEnabled(ruleId);
  } catch (e) {
    console.error('切换规则状态失败:', e);
  }
};

// 导入正则规则
const handleImport = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = '.json';
  input.onchange = async (event: Event) => {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (!file) return;
    try {
      const buffer = await file.arrayBuffer();
      await regexRuleStore.importRules(new Uint8Array(buffer), file.name);
      message.success('导入成功');
    } catch (e) {
      console.error('导入正则规则失败:', e);
    }
  };
  input.click();
};

// 新建规则
const handleCreate = async () => {
  try {
    const defaultFlags = regexRuleStore.createDefaultAffectFlags();
    const newRule = await regexRuleStore.addRule({
      presetId: 0,
      name: '新规则',
      findPattern: '',
      replacePattern: '',
      isEnabled: true,
      runOnEdit: false,
      substituteRegex: false,
      affectFlags: defaultFlags,
      sortOrder: regexRuleStore.rules!.length
    });
    if (newRule) {
      message.success('规则已创建');
      router.push(`/regex/${newRule.id}`);
    }
  } catch (e) {
    console.error('创建规则失败:', e);
  }
};

// 菜单操作
const handleMenuSelect = (key: string, ruleId: number) => {
  switch (key) {
    case 'copy':
      handleCopy(ruleId);
      break;
    case 'export':
      handleExport();
      break;
    case 'delete':
      handleDelete(ruleId);
      break;
  }
};

// 复制规则
const handleCopy = async (ruleId: number) => {
  const rule = regexRuleStore.rules!.find(r => r.id === ruleId);
  if (!rule) return;
  try {
    const newRule = await regexRuleStore.addRule({
      presetId: 0,
      name: `${rule.name} (副本)`,
      findPattern: rule.findPattern,
      replacePattern: rule.replacePattern,
      isEnabled: rule.isEnabled,
      runOnEdit: rule.runOnEdit,
      substituteRegex: rule.substituteRegex,
      minDepth: rule.minDepth,
      maxDepth: rule.maxDepth,
      affectFlags: rule.affectFlags,
      sortOrder: regexRuleStore.rules!.length
    });
    if (newRule) {
      message.success('规则已复制');
    }
  } catch (e) {
    console.error('复制规则失败:', e);
  }
};

// 导出规则
const handleExport = async () => {
  try {
    const result = await regexRuleStore.exportRules(0);
    const blob = new Blob([new Uint8Array(result.fileContent)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = result.fileName || 'regex-rules.json';
    a.click();
    URL.revokeObjectURL(url);
    message.success('导出成功');
  } catch (e) {
    console.error('导出正则规则失败:', e);
  }
};

// 删除规则
const handleDelete = (ruleId: number) => {
  const rule = regexRuleStore.rules!.find(r => r.id === ruleId);
  if (!rule) return;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除正则规则「${rule.name}」吗？此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await regexRuleStore.deleteRule(ruleId);
        message.success('规则已删除');
        // 如果当前正在编辑该规则，跳转回聊天页
        if (isEditing(ruleId)) {
          router.push('/');
        }
      } catch (e) {
        console.error('删除规则失败:', e);
      }
    }
  });
};
</script>

<style scoped>
.regex-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  min-height: 0;
}

.action-row {
  display: flex;
  gap: 8px;
}

.list-container {
  flex: 1;
  min-height: 0;
  overflow: auto;
}

/* 穿透 n-spin 内部结构，保证 flex 布局链条完整 */
.regex-list-tab :deep(.n-spin-container),
.regex-list-tab :deep(.n-spin-content) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.loading-more {
  display: flex;
  justify-content: center;
  padding: 12px 0;
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
