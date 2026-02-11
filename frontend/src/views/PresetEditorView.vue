<template>
  <div class="preset-editor-view">
    <!-- 编辑器顶部栏 -->
    <div class="editor-top-bar">
      <n-button quaternary @click="backToChat">
        <template #icon>
          <n-icon><ArrowBackOutline /></n-icon>
        </template>
        返回对话
      </n-button>
      <span class="editor-title">编辑预设: {{ presetName }}</span>
      <div class="editor-actions">
        <n-button size="small" @click="handleSaveAs">另存为</n-button>
        <n-button type="primary" size="small" @click="handleSave">保存</n-button>
      </div>
    </div>

    <!-- 编辑器内容区域 -->
    <div class="editor-body">
      <div class="editor-content">
        <!-- 预设名称 -->
        <div class="field-group">
          <label class="field-label">预设名称</label>
          <n-input
            v-model:value="presetName"
            placeholder="输入预设名称"
            size="large"
          />
        </div>

        <!-- Prompt 排列 -->
        <div class="section-header">
          <span class="section-title">Prompt 排列</span>
          <span class="section-desc">拖拽调整顺序，展开编辑内容</span>
        </div>

        <div class="prompt-list">
          <div
            v-for="(prompt, index) in prompts"
            :key="prompt.id"
            class="prompt-card"
            :class="{ expanded: prompt.expanded }"
          >
            <!-- 折叠态头部 -->
            <div class="prompt-header" @click="togglePrompt(index)">
              <div class="prompt-header-left">
                <n-icon class="drag-handle" :size="18" @click.stop>
                  <ReorderTwoOutline />
                </n-icon>
                <n-checkbox
                  :checked="prompt.enabled"
                  @update:checked="(val: boolean) => prompt.enabled = val"
                  @click.stop
                />
                <span class="prompt-name">{{ prompt.name || '未命名 Prompt' }}</span>
                <n-tag :type="roleTagType(prompt.role)" size="small" round>
                  {{ roleLabel(prompt.role) }}
                </n-tag>
              </div>
              <div class="prompt-header-right">
                <n-button
                  quaternary
                  circle
                  size="tiny"
                  @click.stop="togglePrompt(index)"
                >
                  <template #icon>
                    <n-icon :size="16">
                      <ChevronDownOutline v-if="!prompt.expanded" />
                      <ChevronUpOutline v-else />
                    </n-icon>
                  </template>
                </n-button>
                <n-button
                  quaternary
                  circle
                  size="tiny"
                  @click.stop="removePrompt(index)"
                >
                  <template #icon>
                    <n-icon :size="16"><CloseOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>

            <!-- 展开态内容 -->
            <div v-if="prompt.expanded" class="prompt-body">
              <div class="prompt-fields">
                <div class="prompt-field-row">
                  <div class="prompt-field">
                    <label class="field-label-sm">名称</label>
                    <n-input
                      v-model:value="prompt.name"
                      placeholder="Prompt 名称"
                      size="small"
                    />
                  </div>
                  <div class="prompt-field">
                    <label class="field-label-sm">角色</label>
                    <n-select
                      v-model:value="prompt.role"
                      :options="roleOptions"
                      size="small"
                    />
                  </div>
                </div>
              </div>
              <div class="prompt-editor">
                <n-input
                  v-model:value="prompt.content"
                  type="textarea"
                  placeholder="输入 Prompt 内容，支持 {{char}}、{{user}} 等宏变量..."
                  :autosize="{ minRows: 6, maxRows: 20 }"
                  class="prompt-textarea"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 添加 Prompt 按钮 -->
        <n-button dashed block @click="addPrompt" class="add-prompt-btn">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          添加 Prompt 项
        </n-button>

        <!-- 预设内正则规则 -->
        <div class="section-header" style="margin-top: 32px;">
          <span class="section-title">预设内正则规则</span>
          <span class="section-desc">此预设关联的正则替换规则</span>
        </div>

        <div class="regex-list">
          <div
            v-for="(rule, index) in regexRules"
            :key="rule.id"
            class="regex-item"
          >
            <n-checkbox
              :checked="rule.enabled"
              @update:checked="(val: boolean) => rule.enabled = val"
            />
            <span class="regex-name">{{ rule.name }}</span>
            <code class="regex-preview">{{ rule.find }}</code>
            <span class="regex-arrow">→</span>
            <span class="regex-replace">{{ rule.replace || '(空)' }}</span>
            <n-button
              quaternary
              circle
              size="tiny"
              @click="removeRegex(index)"
            >
              <template #icon>
                <n-icon :size="14"><CloseOutline /></n-icon>
              </template>
            </n-button>
          </div>

          <n-empty
            v-if="regexRules.length === 0"
            description="暂无关联正则规则"
            size="small"
            class="empty-regex"
          />
        </div>

        <n-button dashed block size="small" @click="addRegex" class="add-regex-btn">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          添加正则
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  NButton,
  NIcon,
  NInput,
  NSelect,
  NCheckbox,
  NTag,
  NEmpty
} from 'naive-ui';
import {
  ArrowBackOutline,
  AddOutline,
  CloseOutline,
  ReorderTwoOutline,
  ChevronDownOutline,
  ChevronUpOutline
} from '@vicons/ionicons5';

const route = useRoute();
const router = useRouter();

const presetId = computed(() => route.params.id as string);

// TODO: 后续对接真实的预设 store，根据 presetId 加载数据
const presetName = ref('Default');

// Prompt 角色选项
const roleOptions = [
  { label: '系统 (System)', value: 'system' },
  { label: '用户 (User)', value: 'user' },
  { label: '助手 (Assistant)', value: 'assistant' }
];

// Prompt 项接口
interface PromptItem {
  id: string;
  name: string;
  role: string;
  content: string;
  enabled: boolean;
  expanded: boolean;
}

// 示例数据
const prompts = ref<PromptItem[]>([
  {
    id: '1',
    name: 'System Prompt',
    role: 'system',
    content: 'You are {{char}}, a creative writing assistant...',
    enabled: true,
    expanded: true
  },
  {
    id: '2',
    name: 'Character Description',
    role: 'system',
    content: '{{char}}的详细描述和设定...',
    enabled: true,
    expanded: false
  },
  {
    id: '3',
    name: 'Chat History',
    role: 'system',
    content: '',
    enabled: true,
    expanded: false
  },
  {
    id: '4',
    name: 'User Prompt',
    role: 'user',
    content: '',
    enabled: true,
    expanded: false
  }
]);

// 正则规则接口
interface RegexRule {
  id: string;
  name: string;
  find: string;
  replace: string;
  enabled: boolean;
}

const regexRules = ref<RegexRule[]>([
  { id: 'r1', name: '删除OOC', find: '\\(OOC:.*?\\)', replace: '', enabled: true },
  { id: 'r2', name: '格式化思考', find: '<think>.*?</think>', replace: '', enabled: true }
]);

// 角色标签类型映射
const roleTagType = (role: string) => {
  const map: Record<string, 'info' | 'success' | 'warning'> = {
    system: 'info',
    user: 'success',
    assistant: 'warning'
  };
  return map[role] || 'info';
};

// 角色标签文字
const roleLabel = (role: string) => {
  const map: Record<string, string> = {
    system: '系统',
    user: '用户',
    assistant: '助手'
  };
  return map[role] || role;
};

// 展开/折叠 Prompt
const togglePrompt = (index: number) => {
  prompts.value[index].expanded = !prompts.value[index].expanded;
};

// 添加 Prompt
const addPrompt = () => {
  prompts.value.push({
    id: `new-${Date.now()}`,
    name: '',
    role: 'system',
    content: '',
    enabled: true,
    expanded: true
  });
};

// 删除 Prompt
const removePrompt = (index: number) => {
  prompts.value.splice(index, 1);
};

// 添加正则
const addRegex = () => {
  regexRules.value.push({
    id: `regex-${Date.now()}`,
    name: '新规则',
    find: '',
    replace: '',
    enabled: true
  });
};

// 删除正则
const removeRegex = (index: number) => {
  regexRules.value.splice(index, 1);
};

// 返回对话
const backToChat = () => {
  router.push('/');
};

// 保存
const handleSave = () => {
  // TODO: 调用API保存预设
  console.log('保存预设:', presetId.value, presetName.value);
};

// 另存为
const handleSaveAs = () => {
  // TODO: 弹出对话框输入新名称，复制并保存
  console.log('另存为预设');
};
</script>

<style scoped>
.preset-editor-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.editor-top-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.editor-title {
  flex: 1;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
}

.editor-actions {
  display: flex;
  gap: 8px;
}

.editor-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.editor-content {
  max-width: 960px;
  margin: 0 auto;
  width: 100%;
}

.field-group {
  margin-bottom: 24px;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.field-label-sm {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.section-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* Prompt 列表 */
.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.prompt-card {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  overflow: hidden;
  transition: border-color 150ms;
}

.prompt-card:hover {
  border-color: var(--color-primary-light, var(--border-color));
}

.prompt-card.expanded {
  border-color: var(--color-primary);
}

.prompt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  min-height: 48px;
  transition: background-color 100ms;
}

.prompt-header:hover {
  background: var(--bg-hover);
}

.prompt-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.drag-handle {
  cursor: grab;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.drag-handle:active {
  cursor: grabbing;
}

.prompt-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.prompt-header-right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.prompt-body {
  padding: 0 16px 16px;
  border-top: 1px solid var(--border-color);
}

.prompt-fields {
  padding-top: 12px;
  margin-bottom: 12px;
}

.prompt-field-row {
  display: grid;
  grid-template-columns: 1fr 200px;
  gap: 12px;
}

.prompt-textarea :deep(textarea) {
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
}

/* 添加 Prompt 按钮 */
.add-prompt-btn {
  margin-bottom: 8px;
}

/* 正则规则列表 */
.regex-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.regex-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
}

.regex-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
}

.regex-preview {
  font-family: var(--font-mono);
  font-size: 11px;
  background: var(--bg-tertiary);
  padding: 2px 6px;
  border-radius: 3px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}

.regex-arrow {
  color: var(--text-tertiary);
  flex-shrink: 0;
  font-size: 12px;
}

.regex-replace {
  font-size: 12px;
  color: var(--text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-regex {
  padding: 16px 0;
}

.add-regex-btn {
  margin-bottom: 16px;
}

/* 响应式 */
@media (max-width: 767px) {
  .editor-top-bar {
    padding: 8px 16px;
  }

  .editor-body {
    padding: 16px 12px;
  }

  .prompt-field-row {
    grid-template-columns: 1fr;
  }

  .editor-title {
    font-size: 14px;
  }
}
</style>
