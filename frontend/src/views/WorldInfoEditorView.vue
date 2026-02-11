<template>
  <div class="worldinfo-editor-view">
    <!-- 编辑器顶部栏 -->
    <div class="editor-top-bar">
      <n-button quaternary @click="backToChat">
        <template #icon>
          <n-icon><ArrowBackOutline /></n-icon>
        </template>
        返回对话
      </n-button>
      <span class="editor-title">编辑世界书: {{ worldInfoName }}</span>
      <div class="editor-actions">
        <n-button size="small" @click="handleExport">导出</n-button>
        <n-button type="primary" size="small" @click="handleSave">保存</n-button>
      </div>
    </div>

    <!-- 编辑器内容 -->
    <div class="editor-scroll-area">
      <div class="editor-content">
        <!-- 世界书名称 -->
        <div class="form-section">
          <label class="form-label">世界书名称</label>
          <n-input v-model:value="worldInfoName" placeholder="输入世界书名称" />
        </div>

        <!-- 搜索 + 添加词条 -->
        <div class="search-action-row">
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索词条..."
            clearable
            size="small"
            class="search-input"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>
          <n-button size="small" type="primary" @click="addEntry">
            <template #icon>
              <n-icon><AddOutline /></n-icon>
            </template>
            添加词条
          </n-button>
        </div>

        <!-- 词条列表 -->
        <div class="section-divider">
          <span>词条列表 ({{ filteredEntries.length }})</span>
        </div>

        <div class="entries-list">
          <div
            v-for="(entry, index) in filteredEntries"
            :key="entry.id"
            class="entry-card"
            :class="{ expanded: entry.expanded }"
          >
            <!-- 词条头部 -->
            <div class="entry-header" @click="toggleEntry(entry.id)">
              <div class="entry-header-left">
                <n-icon class="drag-handle" :size="18" @click.stop>
                  <ReorderTwoOutline />
                </n-icon>
                <n-checkbox
                  :checked="entry.enabled"
                  @update:checked="(val: boolean) => entry.enabled = val"
                  @click.stop
                />
                <span class="entry-name">{{ entry.name || `词条 #${index + 1}` }}</span>
              </div>
              <div class="entry-header-right">
                <div v-if="!entry.expanded" class="entry-keywords-preview">
                  <n-tag
                    v-for="kw in entry.keywords.slice(0, 3)"
                    :key="kw"
                    size="tiny"
                    :bordered="false"
                  >
                    {{ kw }}
                  </n-tag>
                  <span v-if="entry.keywords.length > 3" class="more-tag">+{{ entry.keywords.length - 3 }}</span>
                </div>
                <n-button quaternary circle size="tiny" @click.stop="toggleEntry(entry.id)">
                  <template #icon>
                    <n-icon>
                      <ChevronDownOutline v-if="!entry.expanded" />
                      <ChevronUpOutline v-else />
                    </n-icon>
                  </template>
                </n-button>
                <n-button quaternary circle size="tiny" @click.stop="removeEntry(entry.id)">
                  <template #icon>
                    <n-icon :size="16"><CloseOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>

            <!-- 词条展开内容 -->
            <div v-if="entry.expanded" class="entry-body">
              <!-- 词条名称 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">词条名称</label>
                  <n-input v-model:value="entry.name" placeholder="输入词条名称" size="small" />
                </div>
              </div>

              <!-- 关键词 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">关键词</label>
                  <n-dynamic-tags v-model:value="entry.keywords" />
                </div>
              </div>

              <!-- 辅助关键词 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">辅助关键词 <span class="optional">(可选)</span></label>
                  <n-dynamic-tags v-model:value="entry.secondaryKeywords" />
                </div>
              </div>

              <!-- 触发方式 / 插入位置 -->
              <div class="form-row two-col">
                <div class="form-field">
                  <label class="field-label">触发方式</label>
                  <n-select
                    v-model:value="entry.triggerMode"
                    :options="triggerModeOptions"
                    size="small"
                  />
                </div>
                <div class="form-field">
                  <label class="field-label">插入位置</label>
                  <n-select
                    v-model:value="entry.insertPosition"
                    :options="insertPositionOptions"
                    size="small"
                  />
                </div>
              </div>

              <!-- 优先级 / 状态 -->
              <div class="form-row two-col">
                <div class="form-field">
                  <label class="field-label">优先级</label>
                  <n-input-number
                    v-model:value="entry.priority"
                    :min="0"
                    :max="1000"
                    size="small"
                  />
                </div>
                <div class="form-field">
                  <label class="field-label">状态</label>
                  <n-select
                    v-model:value="entry.status"
                    :options="statusOptions"
                    size="small"
                  />
                </div>
              </div>

              <!-- 内容编辑器 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">内容</label>
                  <n-input
                    v-model:value="entry.content"
                    type="textarea"
                    placeholder="输入词条内容，支持 {{char}}、{{user}} 等宏变量..."
                    :autosize="{ minRows: 5, maxRows: 20 }"
                    class="content-editor"
                  />
                  <div class="token-count">Token 数: {{ estimateTokens(entry.content) }}</div>
                </div>
              </div>
            </div>
          </div>

          <!-- 空状态 -->
          <n-empty
            v-if="filteredEntries.length === 0 && entries.length > 0"
            description="没有匹配的词条"
            size="small"
            class="empty-state"
          />
          <n-empty
            v-if="entries.length === 0"
            description="暂无词条，点击上方按钮添加"
            size="small"
            class="empty-state"
          />
        </div>

        <!-- 底部添加按钮 -->
        <n-button
          v-if="entries.length > 0"
          dashed
          block
          size="small"
          class="add-entry-bottom"
          @click="addEntry"
        >
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          添加词条
        </n-button>

        <!-- 全局设置 -->
        <div class="section-divider">
          <span>全局设置</span>
        </div>

        <div class="global-settings">
          <div class="form-row two-col">
            <div class="form-field">
              <label class="field-label">递归扫描</label>
              <n-switch v-model:value="globalSettings.recursiveScan" />
            </div>
            <div class="form-field">
              <label class="field-label">扫描深度</label>
              <n-input-number
                v-model:value="globalSettings.scanDepth"
                :min="1"
                :max="10"
                size="small"
                :disabled="!globalSettings.recursiveScan"
              />
            </div>
          </div>
          <div class="form-row two-col">
            <div class="form-field">
              <label class="field-label">Token 预算</label>
              <n-input-number
                v-model:value="globalSettings.tokenBudget"
                :min="0"
                :max="65536"
                :step="256"
                size="small"
              />
            </div>
            <div class="form-field" />
          </div>
        </div>
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
  NInputNumber,
  NSelect,
  NCheckbox,
  NDynamicTags,
  NTag,
  NSwitch,
  NEmpty
} from 'naive-ui';
import {
  ArrowBackOutline,
  SearchOutline,
  AddOutline,
  ReorderTwoOutline,
  ChevronDownOutline,
  ChevronUpOutline,
  CloseOutline
} from '@vicons/ionicons5';

const route = useRoute();
const router = useRouter();

// 世界书名称
const worldInfoName = ref(route.params.id as string || '未命名世界书');

// 搜索关键词
const searchKeyword = ref('');

// 词条数据结构
interface WorldInfoEntry {
  id: string;
  name: string;
  keywords: string[];
  secondaryKeywords: string[];
  triggerMode: string;
  insertPosition: string;
  priority: number;
  status: string;
  content: string;
  enabled: boolean;
  expanded: boolean;
}

// 词条列表（TODO: 后续对接真实 store）
const entries = ref<WorldInfoEntry[]>([
  {
    id: '1',
    name: '角色背景',
    keywords: ['background', '背景', '人物设定'],
    secondaryKeywords: [],
    triggerMode: 'keyword',
    insertPosition: 'after_char_desc',
    priority: 100,
    status: 'enabled',
    content: '{{char}}的背景故事:\n{{char}}出生于一个小镇...',
    enabled: true,
    expanded: true
  },
  {
    id: '2',
    name: '世界地理',
    keywords: ['geography', '地理', '地图'],
    secondaryKeywords: ['位置', '城市'],
    triggerMode: 'keyword',
    insertPosition: 'before_system',
    priority: 80,
    status: 'enabled',
    content: '世界地理设定...',
    enabled: true,
    expanded: false
  },
  {
    id: '3',
    name: '魔法体系',
    keywords: ['magic', '魔法', '技能'],
    secondaryKeywords: [],
    triggerMode: 'keyword',
    insertPosition: 'after_char_desc',
    priority: 60,
    status: 'enabled',
    content: '魔法体系说明...',
    enabled: true,
    expanded: false
  }
]);

// 全局设置
const globalSettings = ref({
  recursiveScan: true,
  scanDepth: 2,
  tokenBudget: 2048
});

// 下拉选项
const triggerModeOptions = [
  { label: '关键词匹配', value: 'keyword' },
  { label: '常驻', value: 'constant' },
  { label: '角色专属', value: 'character_only' },
  { label: '禁用', value: 'disabled' }
];

const insertPositionOptions = [
  { label: '系统提示前', value: 'before_system' },
  { label: '系统提示后', value: 'after_system' },
  { label: '角色描述前', value: 'before_char_desc' },
  { label: '角色描述后', value: 'after_char_desc' },
  { label: '对话历史前', value: 'before_chat' },
  { label: '对话历史后', value: 'after_chat' }
];

const statusOptions = [
  { label: '启用', value: 'enabled' },
  { label: '禁用', value: 'disabled' }
];

// 搜索过滤
const filteredEntries = computed(() => {
  if (!searchKeyword.value.trim()) {
    return entries.value;
  }
  const keyword = searchKeyword.value.toLowerCase();
  return entries.value.filter((entry) => {
    return (
      entry.name.toLowerCase().includes(keyword) ||
      entry.keywords.some((kw) => kw.toLowerCase().includes(keyword)) ||
      entry.content.toLowerCase().includes(keyword)
    );
  });
});

// 简单的 Token 估算（每4个字符约1个 token）
const estimateTokens = (text: string): number => {
  if (!text) return 0;
  return Math.ceil(text.length / 4);
};

// 展开/折叠词条
const toggleEntry = (entryId: string) => {
  const entry = entries.value.find((e) => e.id === entryId);
  if (entry) {
    entry.expanded = !entry.expanded;
  }
};

// 添加词条
let nextId = 100;
const addEntry = () => {
  const newEntry: WorldInfoEntry = {
    id: String(nextId++),
    name: '',
    keywords: [],
    secondaryKeywords: [],
    triggerMode: 'keyword',
    insertPosition: 'after_char_desc',
    priority: 50,
    status: 'enabled',
    content: '',
    enabled: true,
    expanded: true
  };
  entries.value.push(newEntry);
};

// 删除词条
const removeEntry = (entryId: string) => {
  const idx = entries.value.findIndex((e) => e.id === entryId);
  if (idx !== -1) {
    entries.value.splice(idx, 1);
  }
};

// 返回对话
const backToChat = () => {
  router.push('/');
};

// 保存
const handleSave = () => {
  // TODO: 调用API保存世界书
};

// 导出
const handleExport = () => {
  // TODO: 导出世界书为JSON
};
</script>

<style scoped>
.worldinfo-editor-view {
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
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.editor-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.editor-scroll-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.editor-content {
  max-width: 960px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 表单部分 */
.form-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

/* 搜索+添加行 */
.search-action-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  flex: 1;
}

/* 分割线标题 */
.section-divider {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 500;
}

.section-divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border-color);
}

/* 词条列表 */
.entries-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 词条卡片 */
.entry-card {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  transition: border-color 150ms;
}

.entry-card:hover {
  border-color: var(--text-tertiary);
}

.entry-card.expanded {
  border-color: var(--color-primary);
}

/* 词条头部 */
.entry-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  min-height: 48px;
  transition: background-color 150ms;
}

.entry-header:hover {
  background: var(--bg-hover);
}

.entry-header-left {
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

.entry-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.entry-header-right {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.entry-keywords-preview {
  display: flex;
  gap: 4px;
  align-items: center;
}

.more-tag {
  font-size: 11px;
  color: var(--text-tertiary);
}

/* 词条内容区 */
.entry-body {
  padding: 16px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row.two-col .form-field {
  flex: 1;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-field.full {
  flex: 1;
}

.field-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
}

.optional {
  font-weight: 400;
  color: var(--text-tertiary);
}

.content-editor :deep(textarea) {
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
}

.token-count {
  font-size: 11px;
  color: var(--text-tertiary);
  text-align: right;
  margin-top: 4px;
}

/* 底部添加按钮 */
.add-entry-bottom {
  margin-top: 4px;
}

/* 全局设置 */
.global-settings {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px solid var(--border-color);
}

/* 空状态 */
.empty-state {
  padding: 32px 0;
}

/* 响应式适配 */
@media (max-width: 767px) {
  .editor-top-bar {
    padding: 10px 16px;
  }

  .editor-scroll-area {
    padding: 16px 12px;
  }

  .form-row.two-col {
    flex-direction: column;
    gap: 14px;
  }

  .search-action-row {
    flex-direction: column;
    gap: 8px;
  }

  .entry-keywords-preview {
    display: none;
  }
}
</style>
