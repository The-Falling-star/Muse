<template>
  <div class="regex-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <n-input v-model:value="searchQuery" placeholder="搜索正则..." clearable class="search-input">
        <template #prefix><n-icon><SearchOutline /></n-icon></template>
      </n-input>
      <n-button type="primary" @click="createRegex">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        新建规则
      </n-button>
    </div>

    <!-- 规则列表 -->
    <n-scrollbar class="regex-container">
      <div class="regex-list">
        <TransitionGroup name="regex-list">
          <div v-for="rule in filteredRules" :key="rule.id" class="regex-card" :class="{ 'disabled': !rule.enabled }">
            <div class="regex-header">
              <n-switch v-model:value="rule.enabled" size="small" />
              <h3 class="regex-name">{{ rule.name }}</h3>
              <n-tag :type="getScopeColor(rule.scope)" size="small">{{ getScopeName(rule.scope) }}</n-tag>
              <div class="regex-actions">
                <n-button quaternary circle size="tiny" @click="editRegex(rule)">
                  <template #icon><n-icon size="14"><CreateOutline /></n-icon></template>
                </n-button>
                <n-button quaternary circle size="tiny" @click="deleteRegex(rule)">
                  <template #icon><n-icon size="14"><TrashOutline /></n-icon></template>
                </n-button>
              </div>
            </div>

            <div class="regex-body">
              <div class="regex-pattern">
                <span class="label">匹配:</span>
                <code>{{ rule.pattern }}</code>
                <n-tag size="tiny" :bordered="false">{{ rule.flags }}</n-tag>
              </div>
              <div class="regex-replacement">
                <span class="label">替换:</span>
                <code>{{ rule.replacement || '(删除匹配内容)' }}</code>
              </div>
            </div>

            <div class="regex-footer">
              <span class="regex-order"><n-icon><SwapVerticalOutline /></n-icon>{{ rule.order }}</span>
            </div>
          </div>
        </TransitionGroup>
      </div>

      <n-empty v-if="filteredRules.length === 0" description="暂无正则规则" class="empty-state">
        <template #extra><n-button type="primary" @click="createRegex">创建规则</n-button></template>
      </n-empty>
    </n-scrollbar>

    <!-- 编辑模态框 -->
    <n-modal v-model:show="showEditModal" preset="card" :title="editingRegex ? '编辑规则' : '新建规则'" :style="{ width: '600px', maxWidth: '90vw' }">
      <RegexEditor :regex="editingRegex" @save="handleSaveRegex" @cancel="showEditModal = false" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { NInput, NButton, NIcon, NScrollbar, NSwitch, NTag, NEmpty, NModal, useMessage, useDialog } from 'naive-ui';
import { SearchOutline, AddOutline, CreateOutline, TrashOutline, SwapVerticalOutline } from '@vicons/ionicons5';
import RegexEditor from '../components/regex/RegexEditor.vue';
import type { RegexRule } from '../types';

const message = useMessage();
const dialog = useDialog();

const searchQuery = ref('');
const showEditModal = ref(false);
const editingRegex = ref<RegexRule | null>(null);

const rules = ref<RegexRule[]>([
  { id: '1', name: '移除OOC内容', pattern: '\\(OOC:.*?\\)', replacement: '', flags: 'gi', scope: 'output', enabled: true, order: 100 },
  { id: '2', name: '替换用户名', pattern: '{{user}}', replacement: '你', flags: 'g', scope: 'both', enabled: true, order: 90 },
  { id: '3', name: '移除空行', pattern: '\\n{3,}', replacement: '\\n\\n', flags: 'g', scope: 'output', enabled: false, order: 80 }
]);

const filteredRules = computed(() => {
  if (!searchQuery.value) return rules.value;
  return rules.value.filter(r => r.name.toLowerCase().includes(searchQuery.value.toLowerCase()) || r.pattern.toLowerCase().includes(searchQuery.value.toLowerCase()));
});

const getScopeName = (scope: string) => ({ input: '输入', output: '输出', both: '双向' }[scope] || scope);
const getScopeColor = (scope: string): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' => ({ input: 'info', output: 'success', both: 'warning' }[scope] as 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' || 'default');

const createRegex = () => { editingRegex.value = null; showEditModal.value = true; };
const editRegex = (rule: RegexRule) => { editingRegex.value = rule; showEditModal.value = true; };

const deleteRegex = (rule: RegexRule) => {
  dialog.warning({
    title: '确认删除', content: `确定要删除规则"${rule.name}"吗？`, positiveText: '删除', negativeText: '取消',
    onPositiveClick: () => { rules.value = rules.value.filter(r => r.id !== rule.id); message.success('规则已删除'); }
  });
};

const handleSaveRegex = (rule: RegexRule) => {
  if (editingRegex.value) {
    const index = rules.value.findIndex(r => r.id === editingRegex.value!.id);
    if (index >= 0) rules.value[index] = rule;
    message.success('规则已更新');
  } else {
    rules.value.push({ ...rule, id: Date.now().toString() });
    message.success('规则已创建');
  }
  showEditModal.value = false;
  editingRegex.value = null;
};
</script>

<style scoped>
.regex-view { display: flex; flex-direction: column; height: calc(100vh - 64px - 48px); }
.toolbar { display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; }
.search-input { flex: 1; max-width: 300px; min-width: 200px; }
.regex-container { flex: 1; }
.regex-list { display: flex; flex-direction: column; gap: 12px; padding-bottom: 20px; }
.regex-card { background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 12px; padding: 16px; transition: all var(--transition-normal); }
.regex-card:hover { border-color: var(--border-glow); box-shadow: var(--glow-soft); }
.regex-card.disabled { opacity: 0.6; }
.regex-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.regex-name { flex: 1; font-size: 15px; font-weight: 600; margin: 0; color: var(--text-primary); }
.regex-actions { display: flex; gap: 4px; opacity: 0; transition: opacity var(--transition-fast); }
.regex-card:hover .regex-actions { opacity: 1; }
.regex-body { display: flex; flex-direction: column; gap: 8px; margin-bottom: 12px; }
.regex-pattern, .regex-replacement { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.regex-pattern .label, .regex-replacement .label { color: var(--text-tertiary); min-width: 40px; }
.regex-pattern code, .regex-replacement code { background: var(--bg-tertiary); padding: 4px 8px; border-radius: 4px; font-family: var(--font-mono, monospace); color: var(--color-primary); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.regex-footer { display: flex; align-items: center; font-size: 12px; color: var(--text-tertiary); }
.regex-order { display: flex; align-items: center; gap: 4px; }
.empty-state { padding: 60px 20px; }
.regex-list-enter-active, .regex-list-leave-active { transition: all 0.3s ease; }
.regex-list-enter-from, .regex-list-leave-to { opacity: 0; transform: translateX(-20px); }
</style>
