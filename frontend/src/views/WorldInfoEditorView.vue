<template>
  <div class="worldinfo-editor-view">
    <n-spin :show="pageLoading" description="加载中..." style="width: 100%; height: 100%;">
      <div class="spin-content">

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
        <n-button size="small" :disabled="saving" @click="handleExport">导出</n-button>
        <n-button type="primary" size="small" :loading="saving" @click="handleSave">保存</n-button>
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

        <!-- 世界书描述 -->
        <div class="form-section">
          <label class="form-label">描述 <span class="optional">(可选)</span></label>
          <n-input
            v-model:value="worldInfoDescription"
            type="textarea"
            placeholder="输入世界书描述..."
            :autosize="{ minRows: 2, maxRows: 4 }"
          />
        </div>

        <!-- 全局开关 -->
        <div class="form-section form-section-inline">
          <label class="form-label">全局世界书</label>
          <n-switch v-model:value="worldInfoIsGlobal" />
          <span class="form-hint">全局世界书对所有角色生效</span>
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
            :key="entry.id || index"
            class="entry-card"
            :class="{ expanded: expandedIds.has(entry.id || -index - 1) }"
          >
            <!-- 词条头部 -->
            <div class="entry-header" @click="toggleEntry(entry.id || -index - 1)">
              <div class="entry-header-left">
                <n-icon class="drag-handle" :size="18" @click.stop>
                  <ReorderTwoOutline />
                </n-icon>
                <n-checkbox
                  :checked="entry.isEnabled"
                  @update:checked="(val: boolean) => entry.isEnabled = val"
                  @click.stop
                />
                <span class="entry-name">{{ entryDisplayName(entry, index) }}</span>
              </div>
              <div class="entry-header-right">
                <div v-if="!expandedIds.has(entry.id || -index - 1)" class="entry-keywords-preview">
                  <n-tag
                    v-for="kw in entry.keysList.slice(0, 3)"
                    :key="kw"
                    size="tiny"
                    :bordered="false"
                  >
                    {{ kw }}
                  </n-tag>
                  <span v-if="entry.keysList.length > 3" class="more-tag">
                    +{{ entry.keysList.length - 3 }}
                  </span>
                </div>
                <n-button quaternary circle size="tiny" @click.stop="toggleEntry(entry.id || -index - 1)">
                  <template #icon>
                    <n-icon>
                      <ChevronDownOutline v-if="!expandedIds.has(entry.id || -index - 1)" />
                      <ChevronUpOutline v-else />
                    </n-icon>
                  </template>
                </n-button>
                <n-button quaternary circle size="tiny" @click.stop="removeEntry(index)">
                  <template #icon>
                    <n-icon :size="16"><CloseOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>

            <!-- 词条展开内容 -->
            <div v-if="expandedIds.has(entry.id || -index - 1)" class="entry-body">
              <!-- 词条备注 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">备注名称</label>
                  <n-input v-model:value="entry.comment" placeholder="输入备注名称（仅管理用，不会注入）" size="small" />
                </div>
              </div>

              <!-- 关键词 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">关键词 (逗号分隔)</label>
                  <n-input
                    :value="entry.keysList.join(', ')"
                    @update:value="(val: string) => entry.keysList = val.split(',').map(k => k.trim()).filter(k => k.length > 0)"
                    placeholder="输入关键词，用逗号分隔"
                    size="small"
                  />
                </div>
              </div>

              <!-- 辅助关键词 -->
              <div class="form-row">
                <div class="form-field full">
                  <label class="field-label">辅助关键词 <span class="optional">(可选，逗号分隔)</span></label>
                  <n-input
                    :value="entry.secondaryKeys.join(', ')"
                    @update:value="(val: string) => entry.secondaryKeys = val.split(',').map(k => k.trim()).filter(k => k.length > 0)"
                    placeholder="输入辅助关键词，用逗号分隔"
                    size="small"
                  />
                </div>
              </div>

              <!-- 插入位置 / 优先级 -->
              <div class="form-row two-col">
                <div class="form-field">
                  <label class="field-label">插入位置</label>
                  <n-select
                    v-model:value="entry.position"
                    :options="positionOptions"
                    size="small"
                  />
                </div>
                <div class="form-field">
                  <label class="field-label">插入优先级</label>
                  <n-input-number
                    v-model:value="entry.insertionOrder"
                    :min="0"
                    :max="9999"
                    size="small"
                  />
                </div>
              </div>

              <!-- 深度 / 常驻 / 选择性匹配 -->
              <div class="form-row three-col">
                <div class="form-field">
                  <label class="field-label">深度</label>
                  <n-input-number
                    v-model:value="entry.depth"
                    :min="0"
                    :max="100"
                    size="small"
                    :disabled="entry.position !== EntryPosition.AtDepth"
                  />
                </div>
                <div class="form-field">
                  <label class="field-label">常驻</label>
                  <n-switch v-model:value="entry.constant" />
                </div>
                <div class="form-field">
                  <label class="field-label">选择性匹配</label>
                  <n-switch v-model:value="entry.selective" />
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
      </div>
    </div>

      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, reactive } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  NButton,
  NIcon,
  NInput,
  NInputNumber,
  NSelect,
  NCheckbox,
  NTag,
  NSwitch,
  NEmpty,
  NSpin,
  useMessage,
  useDialog
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
import { worldInfoClient } from '@/api/client';
import { useWorldInfoStore } from '@/stores/worldInfo';
import { EntryPosition, Role } from '@/gen/muse/common_pb';
import type { WorldInfoEntry } from '@/gen/muse/worldinfo_pb';

const route = useRoute();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const worldInfoStore = useWorldInfoStore();

const worldInfoId = computed(() => Number(route.params.id));
const pageLoading = ref(true);
const saving = ref(false);

// 世界书基础信息
const worldInfoName = ref('');
const worldInfoDescription = ref('');
const worldInfoIsGlobal = ref(false);

// 搜索关键词
const searchKeyword = ref('');

// 展开状态管理（使用 Set 追踪展开的条目ID）
const expandedIds = reactive(new Set<number>());

// 词条列表
const entries = ref<WorldInfoEntry[]>([]);

// 插入位置选项
const positionOptions = [
  { label: '角色描述之前', value: EntryPosition.BeforeChar },
  { label: '角色描述之后', value: EntryPosition.AfterChar },
  { label: '示例对话之前', value: EntryPosition.BeforeExample },
  { label: '示例对话之后', value: EntryPosition.AfterExample },
  { label: '按深度插入', value: EntryPosition.AtDepth }
];

// =====================
// 加载数据
// =====================

const loadWorldInfo = async () => {
  pageLoading.value = true;
  try {
    // 获取世界书基本信息
    const resp = await worldInfoClient.getWorldInfo({ id: worldInfoId.value });
    const worldInfo = resp.worldInfo;
    if (!worldInfo) {
      message.error('世界书不存在');
      router.push('/');
      return;
    }
    worldInfoName.value = worldInfo.name;
    worldInfoDescription.value = worldInfo.description ?? '';
    worldInfoIsGlobal.value = worldInfo.isGlobal;

    // 获取条目列表
    const entriesResp = await worldInfoClient.listWorldInfoEntries({
      worldInfoId: worldInfoId.value
    });
    entries.value = entriesResp.entries || [];
  } catch {
    message.error('加载世界书失败');
  } finally {
    pageLoading.value = false;
  }
};

onMounted(() => {
  loadWorldInfo();
});

// 路由参数变化时重新加载
watch(worldInfoId, (newId) => {
  if (newId) {
    expandedIds.clear();
    loadWorldInfo();
  }
});

// =====================
// 搜索过滤
// =====================

const filteredEntries = computed(() => {
  if (!searchKeyword.value.trim()) {
    return entries.value;
  }
  const keyword = searchKeyword.value.toLowerCase();
  return entries.value.filter((entry) => {
    return (
      (entry.comment || '').toLowerCase().includes(keyword) ||
      entry.keysList.some(k => k.toLowerCase().includes(keyword)) ||
      entry.content.toLowerCase().includes(keyword)
    );
  });
});

// =====================
// 辅助方法
// =====================

// 条目显示名
const entryDisplayName = (entry: WorldInfoEntry, index: number): string => {
  if (entry.comment) return entry.comment;
  if (entry.keysList.length > 0 && entry.keysList[0]) return entry.keysList[0];
  return `词条 #${index + 1}`;
};

// 简单的 Token 估算
const estimateTokens = (text: string): number => {
  if (!text) return 0;
  return Math.ceil(text.length / 4);
};

// 展开/折叠词条
const toggleEntry = (entryKey: number) => {
  if (expandedIds.has(entryKey)) {
    expandedIds.delete(entryKey);
  } else {
    expandedIds.add(entryKey);
  }
};

// =====================
// 条目操作
// =====================

// 临时ID计数器（用于新增但尚未保存的条目）
let tempIdCounter = -1;

const addEntry = () => {
  const maxSortOrder = entries.value.reduce((max, e) => Math.max(max, e.sortOrder), 0);
  const tempId = tempIdCounter--;
  const newEntry = {
    $typeName: 'muse.WorldInfoEntry',
    id: 0,
    worldInfoId: worldInfoId.value,
    uid: '',
    keysList: [] as string[],
    secondaryKeys: [] as string[],
    content: '',
    comment: '',
    isEnabled: true,
    constant: false,
    selective: false,
    insertionOrder: 100,
    position: EntryPosition.AfterChar,
    depth: 4,
    sortOrder: maxSortOrder + 1,
    createdAt: 0n,
    updatedAt: 0n,
    role: Role.System,
    _tempId: tempId
  } as WorldInfoEntry & { _tempId?: number };

  entries.value.push(newEntry);
  // 自动展开新条目（使用负数索引作为key）
  expandedIds.add(0);
};

// 删除条目
const removeEntry = (index: number) => {
  const entry = entries.value[index];
  if (!entry) return;

  const name = entry.comment || entry.keysList[0] || `词条 #${index + 1}`;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除词条"${name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      // 清理展开状态
      expandedIds.delete(entry.id || -index - 1);
      entries.value.splice(index, 1);
      message.success('词条已删除');
    }
  });
};

// =====================
// 保存操作
// =====================

const handleSave = async () => {
  if (!worldInfoName.value.trim()) {
    message.warning('请输入世界书名称');
    return;
  }

  saving.value = true;
  try {
    // 1. 更新世界书基础信息
    const updateResp = await worldInfoClient.updateWorldInfo({
      id: worldInfoId.value,
      name: worldInfoName.value,
      description: worldInfoDescription.value,
      isGlobal: worldInfoIsGlobal.value
    });
    if (updateResp.worldInfo) {
      worldInfoStore.updateWorldInfoInList(updateResp.worldInfo);
    }

    // 2. 获取远端现有条目，比较差异
    const remoteResp = await worldInfoClient.listWorldInfoEntries({
      worldInfoId: worldInfoId.value
    });
    const remoteEntries = remoteResp.entries || [];
    const remoteIds = new Set(remoteEntries.map(e => e.id));
    const localIds = new Set(entries.value.filter(e => e.id > 0).map(e => e.id));

    // 3. 删除远端已有但本地已移除的条目
    for (const remoteId of remoteIds) {
      if (!localIds.has(remoteId)) {
        await worldInfoClient.deleteWorldInfoEntry({ id: remoteId });
      }
    }

    // 4. 新增和更新条目
    for (let i = 0; i < entries.value.length; i++) {
      const entry = entries.value[i];
      if (!entry) continue;
      const entryData = {
        keysList: entry.keysList,
        secondaryKeys: entry.secondaryKeys,
        content: entry.content,
        comment: entry.comment,
        isEnabled: entry.isEnabled,
        constant: entry.constant,
        selective: entry.selective,
        insertionOrder: entry.insertionOrder,
        position: entry.position,
        depth: entry.depth,
        sortOrder: i,
        role: entry.role
      };

      if (entry.id > 0) {
        // 更新已有条目
        const resp = await worldInfoClient.updateWorldInfoEntry({
          id: entry.id,
          ...entryData
        });
        if (resp.entry) {
          entries.value[i] = resp.entry;
        }
      } else {
        // 新增条目
        const resp = await worldInfoClient.addWorldInfoEntry({
          worldInfoId: worldInfoId.value,
          ...entryData
        });
        if (resp.entry) {
          entries.value[i] = resp.entry;
        }
      }
    }

    message.success('世界书保存成功');
  } catch {
    message.error('保存失败，请重试');
  } finally {
    saving.value = false;
  }
};

// 导出
const handleExport = async () => {
  try {
    const resp = await worldInfoClient.exportWorldInfo({ id: worldInfoId.value });
    if (!resp.fileContent || resp.fileContent.length === 0) {
      message.warning('导出内容为空');
      return;
    }

    // 创建下载
    const blob = new Blob([new Uint8Array(resp.fileContent)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = resp.fileName || `${worldInfoName.value}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);

    message.success('世界书已导出');
  } catch {
    message.error('导出失败');
  }
};

// 返回对话
const backToChat = () => {
  router.push('/');
};
</script>

<style scoped>
.worldinfo-editor-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.worldinfo-editor-view :deep(.n-spin-container) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.worldinfo-editor-view :deep(.n-spin-content) {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.spin-content {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
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

.form-section-inline {
  display: flex;
  align-items: center;
  gap: 12px;
}

.form-hint {
  font-size: 12px;
  color: var(--text-tertiary);
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

.form-row.three-col {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 16px;
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
