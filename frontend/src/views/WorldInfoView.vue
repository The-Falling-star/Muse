<template>
  <div class="worldinfo-view">
    <!-- 左侧：世界书列表（移动端可折叠） -->
    <div class="worldinfo-sidebar" :class="{ 'show-mobile': showSidebar }">
      <div class="sidebar-header">
        <div class="sidebar-header-content">
          <h3>世界书</h3>
          <n-button 
            v-if="isMobile" 
            quaternary 
            circle 
            size="small" 
            class="close-sidebar-btn"
            @click="showSidebar = false"
          >
            <template #icon>
              <n-icon><CloseOutline /></n-icon>
            </template>
          </n-button>
        </div>
        <n-button quaternary circle size="small" @click="showWorldModal = true; editingWorld = null">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
        </n-button>
      </div>

      <n-input
        v-model:value="searchQuery"
        placeholder="搜索..."
        clearable
        class="search-input"
      >
        <template #prefix>
          <n-icon><SearchOutline /></n-icon>
        </template>
      </n-input>

      <n-scrollbar class="world-list">
        <n-spin :show="loading">
          <div
            v-for="world in filteredWorlds"
            :key="world.id"
            class="world-item"
            :class="{ 'active': selectedWorldId === world.id }"
            @click="selectWorld(world.id)"
          >
            <div class="world-icon">
              <n-icon><GlobeOutline /></n-icon>
            </div>
            <div class="world-info">
              <div class="world-name">{{ world.name }}</div>
              <div class="world-count">
                {{ world.isGlobal ? '全局' : '角色' }}
              </div>
            </div>
          </div>

          <n-empty v-if="filteredWorlds.length === 0 && !loading" description="暂无世界书" />
        </n-spin>
      </n-scrollbar>

      <!-- 分页控件 -->
      <div class="pagination-container" v-if="Number(totalWorlds) > pageSize">
        <n-pagination
          v-model:page="currentPage"
          :page-size="pageSize"
          :item-count="Number(totalWorlds)"
          @update:page="(page) => loadWorlds(page, pageSize)"
        />
      </div>

      <div class="sidebar-footer">
        <n-button block @click="showImportModal = true">
          <template #icon>
            <n-icon><CloudUploadOutline /></n-icon>
          </template>
          导入世界书
        </n-button>
      </div>
    </div>

    <!-- 右侧：条目列表 -->
    <div class="worldinfo-main">
      <!-- 移动端：列表切换按钮 -->
      <div v-if="isMobile" class="mobile-list-toggle">
        <n-button 
          quaternary 
          @click="showSidebar = true"
          class="list-toggle-btn"
        >
          <template #icon>
            <n-icon><ListOutline /></n-icon>
          </template>
          世界书列表
        </n-button>
      </div>
      <template v-if="selectedWorld">
        <!-- 世界书头部 -->
        <div class="main-header">
          <div class="header-info">
            <h2>{{ selectedWorld.name }}</h2>
            <p>{{ selectedWorld.description || '暂无描述' }}</p>
          </div>
          <div class="header-actions">
            <n-button @click="handleExportWorld">
              <template #icon>
                <n-icon><CloudDownloadOutline /></n-icon>
              </template>
              导出
            </n-button>
            <n-button @click="editingWorld = selectedWorld; showWorldModal = true">
              <template #icon>
                <n-icon><CreateOutline /></n-icon>
              </template>
              编辑
            </n-button>
            <n-button type="error" @click="handleDeleteWorld">
              <template #icon>
                <n-icon><TrashOutline /></n-icon>
              </template>
              删除
            </n-button>
            <n-button type="primary" @click="addEntry">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加条目
            </n-button>
          </div>
        </div>

        <!-- 条目搜索 -->
        <div class="entries-toolbar">
          <n-input
            v-model:value="entrySearchQuery"
            placeholder="搜索条目..."
            clearable
            class="entry-search"
          >
            <template #prefix>
              <n-icon><SearchOutline /></n-icon>
            </template>
          </n-input>

          <n-select
            v-model:value="entrySort"
            :options="sortOptions"
            class="entry-sort"
          />
        </div>

        <!-- 条目列表 -->
        <n-scrollbar class="entries-list">
          <n-spin :show="entriesLoading">
            <TransitionGroup name="entry-list">
              <div
                v-for="entry in filteredEntries"
                :key="entry.id"
                class="entry-item"
                :class="{ 'disabled': !entry.isEnabled }"
              >
                <div class="entry-header">
                  <n-switch
                    :value="entry.isEnabled"
                    size="small"
                    @update:value="(val: boolean) => handleToggleEntry(entry, val)"
                  />
                  <div class="entry-keys">
                    <n-tag
                      v-for="key in parseKeys(entry.keysList).slice(0, 3)"
                      :key="key"
                      size="small"
                      type="primary"
                      :bordered="false"
                    >
                      {{ key }}
                    </n-tag>
                    <n-tag v-if="parseKeys(entry.keysList).length > 3" size="small" :bordered="false">
                      +{{ parseKeys(entry.keysList).length - 3 }}
                    </n-tag>
                  </div>
                  <div class="entry-actions">
                    <n-button quaternary circle size="tiny" @click="editEntry(entry)">
                      <template #icon>
                        <n-icon size="14"><CreateOutline /></n-icon>
                      </template>
                    </n-button>
                    <n-button quaternary circle size="tiny" @click="handleDeleteEntry(entry)">
                      <template #icon>
                        <n-icon size="14"><TrashOutline /></n-icon>
                      </template>
                    </n-button>
                  </div>
                </div>

                <div class="entry-content">
                  <p 
                    :class="{ 
                      'truncated': shouldTruncateContent(entry.content) && !isEntryExpanded(entry.id),
                      'expandable': shouldTruncateContent(entry.content)
                    }"
                    @click="shouldTruncateContent(entry.content) && toggleEntryExpand(entry.id)"
                  >
                    {{ isEntryExpanded(entry.id) || !shouldTruncateContent(entry.content) 
                      ? entry.content 
                      : getTruncatedContent(entry.content) }}
                  </p>
                  <div 
                    v-if="shouldTruncateContent(entry.content)" 
                    class="expand-toggle"
                    @click="toggleEntryExpand(entry.id)"
                  >
                    {{ isEntryExpanded(entry.id) ? '收起' : '展开全部' }}
                  </div>
                </div>

                <div class="entry-meta">
                  <span>优先级: {{ entry.insertionOrder }}</span>
                  <span>深度: {{ entry.depth }}</span>
                  <span v-if="entry.constant">常驻</span>
                  <span v-if="entry.selective">选择性</span>
                </div>
              </div>
            </TransitionGroup>

            <n-empty v-if="filteredEntries.length === 0 && !entriesLoading" description="暂无条目" />
          </n-spin>
        </n-scrollbar>
      </template>

      <template v-else>
        <div class="empty-state">
          <n-icon size="64"><GlobeOutline /></n-icon>
          <h3>选择一个世界书</h3>
          <p>从左侧选择或创建一个世界书来管理条目</p>
        </div>
      </template>
    </div>

    <!-- 世界书编辑模态框 -->
    <n-modal
      v-model:show="showWorldModal"
      preset="card"
      :title="editingWorld ? '编辑世界书' : '创建世界书'"
      :style="{ width: '500px', maxWidth: '90vw' }"
    >
      <n-form ref="worldFormRef" :model="worldForm" :rules="worldRules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="worldForm.name" placeholder="输入世界书名称" />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="worldForm.description"
            type="textarea"
            placeholder="输入描述（可选）"
            :rows="3"
          />
        </n-form-item>
        <n-form-item label="类型" path="isGlobal">
          <n-switch v-model:value="worldForm.isGlobal">
            <template #checked>全局</template>
            <template #unchecked>角色</template>
          </n-switch>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showWorldModal = false">取消</n-button>
          <n-button type="primary" :loading="saving" @click="handleSaveWorld">保存</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 条目编辑模态框 -->
    <n-modal
      v-model:show="showEntryModal"
      preset="card"
      :title="editingEntry ? '编辑条目' : '添加条目'"
      :style="{ width: '600px', maxWidth: '90vw' }"
    >
      <WorldInfoEntryEditor
        :entry="editingEntry as any"
        @save="(data: any) => handleSaveEntry(data)"
        @cancel="showEntryModal = false"
      />
    </n-modal>

    <!-- 导入模态框 -->
    <n-modal
      v-model:show="showImportModal"
      preset="card"
      title="导入世界书"
      :style="{ width: '500px', maxWidth: '90vw' }"
    >
      <n-upload
        :max="1"
        accept=".json"
        :default-upload="false"
        @change="handleImportFile"
      >
        <n-upload-dragger>
          <div class="upload-content">
            <n-icon size="48" class="upload-icon"><CloudUploadOutline /></n-icon>
            <p class="upload-text">点击或拖拽文件到此处</p>
            <p class="upload-hint">支持 JSON 格式的世界书文件</p>
          </div>
        </n-upload-dragger>
      </n-upload>
    </n-modal>

    <!-- 测试按钮（仅在开发环境显示） -->
    <div v-if="__DEV__" class="test-buttons">
      <n-button @click="showSidebar = true" type="primary" size="small">
        打开侧边栏
      </n-button>
      <n-button @click="showSidebar = false" type="error" size="small">
        关闭侧边栏
      </n-button>
      <n-button @click="isMobile = !isMobile" type="warning" size="small">
        切换移动端模式: {{ isMobile ? 'ON' : 'OFF' }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import {
  NInput,
  NSelect,
  NButton,
  NIcon,
  NScrollbar,
  NSwitch,
  NTag,
  NEmpty,
  NModal,
  NForm,
  NFormItem,
  NSpace,
  NSpin,
  NUpload,
  NUploadDragger,
  NPagination,
  useMessage,
  useDialog
} from 'naive-ui';
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui';
import {
  AddOutline,
  SearchOutline,
  GlobeOutline,
  CloudUploadOutline,
  CloudDownloadOutline,
  CreateOutline,
  TrashOutline,
  CloseOutline,
  ListOutline
} from '@vicons/ionicons5';

import WorldInfoEntryEditor from '../components/worldinfo/WorldInfoEntryEditor.vue';
import { worldInfoClient } from '@/api/client';
import type { WorldInfo, WorldInfoEntry } from '@/gen/muse/muse_pb';
import { EntryPosition } from '@/gen/muse/muse_pb';

const message = useMessage();
const dialog = useDialog();

// 设备检测
const isMobile = ref(false);
const showSidebar = ref(false);

// 检测是否为移动端
const detectMobile = () => {
  isMobile.value = window.innerWidth <= 768;
};

// 监听窗口大小变化
const handleResize = () => {
  detectMobile();
  // 在桌面端自动隐藏侧边栏
  if (!isMobile.value) {
    showSidebar.value = false;
  }
};

// 在组件挂载时检测
onMounted(() => {
  detectMobile();
  window.addEventListener('resize', handleResize);
});

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});

// 提供全局属性访问
const $isMobile = isMobile;

// 状态
const loading = ref(false);
const entriesLoading = ref(false);
const saving = ref(false);
const searchQuery = ref('');
const entrySearchQuery = ref('');
const entrySort = ref('order');
const selectedWorldId = ref<number | null>(null);
const showWorldModal = ref(false);
const showEntryModal = ref(false);
const showImportModal = ref(false);
const editingWorld = ref<WorldInfo | null>(null);
const editingEntry = ref<WorldInfoEntry | null>(null);

// 展开状态管理
const expandedEntries = ref<Set<number>>(new Set());

// 分页状态
const currentPage = ref(1);
const pageSize = ref(20);
const totalWorlds = ref(0n); // 使用 bigint 类型

// 世界书列表
const worlds = ref<WorldInfo[]>([]);
// 当前选中世界书的条目
const entries = ref<WorldInfoEntry[]>([]);

// 表单
const worldFormRef = ref<FormInst | null>(null);
const worldForm = ref({
  name: '',
  description: '',
  isGlobal: false
});
const worldRules: FormRules = {
  name: [{ required: true, message: '请输入世界书名称', trigger: 'blur' }]
};

// 排序选项
const sortOptions = [
  { label: '按优先级', value: 'order' },
  { label: '按名称', value: 'name' }
];

// 计算属性
const filteredWorlds = computed(() => {
  if (!searchQuery.value) return worlds.value;
  return worlds.value.filter(w =>
    w.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const selectedWorld = computed(() =>
  worlds.value.find(w => w.id === selectedWorldId.value) || null
);

const filteredEntries = computed(() => {
  let list = [...entries.value];

  if (entrySearchQuery.value) {
    const query = entrySearchQuery.value.toLowerCase();
    list = list.filter(e =>
      e.keysList.toLowerCase().includes(query) ||
      e.content.toLowerCase().includes(query)
    );
  }

  // 排序
  if (entrySort.value === 'order') {
    list.sort((a, b) => b.insertionOrder - a.insertionOrder);
  } else if (entrySort.value === 'name') {
    list.sort((a, b) => a.keysList.localeCompare(b.keysList));
  }

  return list;
});

// 解析关键词列表
const parseKeys = (keysList: string): string[] => {
  return keysList.split(',').map(k => k.trim()).filter(k => k.length > 0);
};

// 控制条目展开/收起
const toggleEntryExpand = (entryId: number) => {
  if (expandedEntries.value.has(entryId)) {
    expandedEntries.value.delete(entryId);
  } else {
    expandedEntries.value.add(entryId);
  }
  expandedEntries.value = new Set(expandedEntries.value);
};

// 检查条目是否展开
const isEntryExpanded = (entryId: number): boolean => {
  return expandedEntries.value.has(entryId);
};

// 检查内容是否需要截断显示
const shouldTruncateContent = (content: string): boolean => {
  const lines = content.split('\n');
  return lines.length > 3;
};

// 获取截断后的内容
const getTruncatedContent = (content: string): string => {
  const lines = content.split('\n');
  if (lines.length <= 3) {
    return content;
  }
  return lines.slice(0, 3).join('\n') + '\n...';
};

// 加载世界书列表
const loadWorlds = async (page: number = 1, size: number = 20) => {
  loading.value = true;
  try {
    const response = await worldInfoClient.listWorldInfos({
      page: page,
      pageSize: size
    });
    worlds.value = response.worldInfos;
    totalWorlds.value = response.total;
    currentPage.value = response.page;
    pageSize.value = response.pageSize;
  } finally {
    loading.value = false;
  }
};

// 加载条目列表
const loadEntries = async (worldInfoId: number) => {
  entriesLoading.value = true;
  try {
    const response = await worldInfoClient.listWorldInfoEntries({ worldInfoId: worldInfoId });
    entries.value = response.entries;
  } finally {
    entriesLoading.value = false;
  }
};

// 选择世界书
const selectWorld = async (id: number) => {
  selectedWorldId.value = id;
  
  // 移动端选择后自动关闭侧边栏
  if (isMobile.value) {
    showSidebar.value = false;
  }
  
  await loadEntries(id);
};

// 保存世界书
const handleSaveWorld = async () => {
  await worldFormRef.value?.validate();
  saving.value = true;
  try {
    if (editingWorld.value) {
      // 更新
      const response = await worldInfoClient.updateWorldInfo({
        id: editingWorld.value.id,
        name: worldForm.value.name,
        description: worldForm.value.description,
        isGlobal: worldForm.value.isGlobal
      });
      if (response.worldInfo) {
        const index = worlds.value.findIndex(w => w.id === editingWorld.value!.id);
        if (index >= 0) {
          worlds.value[index] = response.worldInfo;
        }
      }
      message.success('世界书已更新');
    } else {
      // 创建
      const response = await worldInfoClient.createWorldInfo({
        name: worldForm.value.name,
        description: worldForm.value.description,
        isGlobal: worldForm.value.isGlobal
      });
      if (response.worldInfo) {
        worlds.value.push(response.worldInfo);
      }
      message.success('世界书已创建');
    }
    showWorldModal.value = false;
    resetWorldForm();
  } finally {
    saving.value = false;
  }
};

// 删除世界书
const handleDeleteWorld = () => {
  if (!selectedWorld.value) return;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除世界书"${selectedWorld.value.name}"吗？所有条目也将被删除。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await worldInfoClient.deleteWorldInfo({ id: selectedWorld.value!.id });
      worlds.value = worlds.value.filter(w => w.id !== selectedWorld.value!.id);
      selectedWorldId.value = null;
      entries.value = [];
      message.success('世界书已删除');
    }
  });
};

// 导出世界书
const handleExportWorld = async () => {
  if (!selectedWorld.value) return;
  const response = await worldInfoClient.exportWorldInfo({ id: selectedWorld.value.id });
  // 下载文件
  const blob = new Blob([new Uint8Array(response.fileContent).buffer as ArrayBuffer], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = response.fileName || `${selectedWorld.value.name}.json`;
  a.click();
  URL.revokeObjectURL(url);
  message.success('导出成功');
};

// 导入世界书
const handleImportFile = async (options: { file: UploadFileInfo; fileList: UploadFileInfo[] }) => {
  const file = options.file.file;
  if (!file) return;

  const arrayBuffer = await file.arrayBuffer();
  const fileContent = new Uint8Array(arrayBuffer);
  const response = await worldInfoClient.importWorldInfo({
    fileContent: fileContent,
    fileName: file.name
  });
  if (response.worldInfo) {
    worlds.value.push(response.worldInfo);
  }
  message.success('导入成功');
  showImportModal.value = false;
};

// 添加条目
const addEntry = () => {
  editingEntry.value = null;
  showEntryModal.value = true;
};

// 编辑条目
const editEntry = (entry: WorldInfoEntry) => {
  editingEntry.value = entry;
  showEntryModal.value = true;
};

// 保存条目
const handleSaveEntry = async (entryData: Partial<WorldInfoEntry>) => {
  if (!selectedWorld.value) return;

  saving.value = true;
  try {
    if (editingEntry.value) {
      // 更新
      const response = await worldInfoClient.updateWorldInfoEntry({
        id: editingEntry.value.id,
        uid: entryData.uid || '',
        keysList: entryData.keysList || '',
        secondaryKeys: entryData.secondaryKeys || '',
        content: entryData.content || '',
        comment: entryData.comment || '',
        isEnabled: entryData.isEnabled ?? true,
        constant: entryData.constant ?? false,
        selective: entryData.selective ?? false,
        insertionOrder: entryData.insertionOrder ?? 100,
        position: entryData.position ?? EntryPosition.BeforeChar,
        depth: entryData.depth ?? 4,
        sortOrder: entryData.sortOrder ?? 0
      });
      if (response.entry) {
        const index = entries.value.findIndex(e => e.id === editingEntry.value!.id);
        if (index >= 0) {
          entries.value[index] = response.entry;
        }
      }
      message.success('条目已更新');
    } else {
      // 添加
      const response = await worldInfoClient.addWorldInfoEntry({
        worldInfoId: selectedWorld.value.id,
        uid: entryData.uid || '',
        keysList: entryData.keysList || '',
        secondaryKeys: entryData.secondaryKeys || '',
        content: entryData.content || '',
        comment: entryData.comment || '',
        isEnabled: entryData.isEnabled ?? true,
        constant: entryData.constant ?? false,
        selective: entryData.selective ?? false,
        insertionOrder: entryData.insertionOrder ?? 100,
        position: entryData.position ?? EntryPosition.BeforeChar,
        depth: entryData.depth ?? 4,
        sortOrder: entries.value.length
      });
      if (response.entry) {
        entries.value.push(response.entry);
      }
      message.success('条目已添加');
    }
    showEntryModal.value = false;
    editingEntry.value = null;
  } finally {
    saving.value = false;
  }
};

// 删除条目
const handleDeleteEntry = (entry: WorldInfoEntry) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个条目吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await worldInfoClient.deleteWorldInfoEntry({ id: entry.id });
      entries.value = entries.value.filter(e => e.id !== entry.id);
      message.success('条目已删除');
    }
  });
};

// 切换条目启用状态
const handleToggleEntry = async (entry: WorldInfoEntry, enabled: boolean) => {
  const response = await worldInfoClient.updateWorldInfoEntry({
    id: entry.id,
    uid: entry.uid,
    keysList: entry.keysList,
    secondaryKeys: entry.secondaryKeys,
    content: entry.content,
    comment: entry.comment,
    isEnabled: enabled,
    constant: entry.constant,
    selective: entry.selective,
    insertionOrder: entry.insertionOrder,
    position: entry.position,
    depth: entry.depth,
    sortOrder: entry.sortOrder
  });
  if (response.entry) {
    const index = entries.value.findIndex(e => e.id === entry.id);
    if (index >= 0) {
      entries.value[index] = response.entry;
    }
  }
};

// 重置世界书表单
const resetWorldForm = () => {
  worldForm.value = {
    name: '',
    description: '',
    isGlobal: false
  };
  editingWorld.value = null;
};

// 监听编辑世界书
const openWorldModal = () => {
  if (editingWorld.value) {
    worldForm.value = {
      name: editingWorld.value.name,
      description: editingWorld.value.description || '',
      isGlobal: editingWorld.value.isGlobal
    };
  } else {
    resetWorldForm();
  }
};

// 初始化
onMounted(() => {
  loadWorlds(currentPage.value, pageSize.value);
});

// 监听showWorldModal变化，初始化表单
import { watch } from 'vue';
watch(showWorldModal, (val) => {
  if (val) {
    openWorldModal();
  }
});
</script>

<style scoped>
.worldinfo-view {
  display: flex;
  height: calc(100vh - 64px - 48px);
  gap: 16px;
  background: var(--bg-primary);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--border-color);
  position: relative;
}

/* 侧边栏 */
.worldinfo-sidebar {
  width: 280px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: all var(--transition-normal);
}

/* 移动端侧边栏 */
@media (max-width: 768px) {
  .worldinfo-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: 300px;
    z-index: 100;
    transform: translateX(-100%);
    box-shadow: var(--shadow-xl);
  }
  
  .worldinfo-sidebar.show-mobile {
    transform: translateX(0);
  }
  
  .worldinfo-view {
    gap: 0;
  }
}

.sidebar-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header-content h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

.close-sidebar-btn {
  margin-left: auto;
}

.search-input {
  margin: 0 16px 12px 16px;
  width: calc(100% - 32px);
  box-sizing: border-box;
}

.world-list {
  flex: 1;
  padding: 0 8px;
}

.pagination-container {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
  display: flex;
  justify-content: center;
  background: var(--bg-secondary);
}

.world-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.world-item:hover {
  background: var(--bg-card-hover);
}

.world-item.active {
  background: rgba(0, 240, 255, 0.1);
}

.world-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-primary);
  border-radius: 8px;
  color: #000;
  flex-shrink: 0;
}

.world-info {
  flex: 1;
  min-width: 0;
}

.world-name {
  font-weight: 500;
  color: var(--text-primary);
}

.world-count {
  font-size: 12px;
  color: var(--text-tertiary);
}

.sidebar-footer {
  padding: 16px;
  border-top: 1px solid var(--border-color);
}

/* 主区域 */
.worldinfo-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  position: relative;
}

/* 移动端列表切换按钮 */
.mobile-list-toggle {
  padding: 12px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.list-toggle-btn {
  width: 100%;
  justify-content: flex-start;
}

.main-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  border-bottom: 1px solid var(--border-color);
}

.header-info h2 {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 4px;
}

.header-info p {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

/* 条目工具栏 */
.entries-toolbar {
  display: flex;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.entry-search {
  flex: 1;
  max-width: 300px;
}

.entry-sort {
  width: 150px;
}

/* 条目列表 */
.entries-list {
  flex: 1;
  padding: 28px 20px 16px 20px; /* 增加顶部padding从16px到28px */
}

.entry-item {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  transition: all var(--transition-fast);
  /* 右侧间距调整：修改 margin-right 的值来控制条目右侧空白大小 */
  margin-right: 12px;
}

.entry-item:hover {
  border-color: var(--border-glow);
}

.entry-item.disabled {
  opacity: 0.5;
}

.entry-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.entry-keys {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  flex: 1;
}

.entry-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.entry-item:hover .entry-actions {
  opacity: 1;
}

.entry-content {
  margin-bottom: 12px;
}

.entry-content p {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-secondary);
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
  overflow-wrap: break-word;
  cursor: default;
}

.entry-content p.expandable {
  cursor: pointer;
}

.entry-content p.truncated {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  cursor: pointer;
}

.expand-toggle {
  color: var(--color-primary);
  font-size: 12px;
  margin-top: 8px;
  cursor: pointer;
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all var(--transition-fast);
}

.expand-toggle:hover {
  background: rgba(0, 240, 255, 0.1);
  text-decoration: underline;
}

.entry-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 空状态 */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
}

.empty-state h3 {
  font-size: 18px;
  margin: 16px 0 8px;
  color: var(--text-primary);
}

.empty-state p {
  margin: 0;
}

/* 条目列表动画 */
.entry-list-enter-active,
.entry-list-leave-active {
  transition: all 0.3s ease;
}

.entry-list-enter-from,
.entry-list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

@media (max-width: 768px) {
  .worldinfo-sidebar {
    width: 100%;
  }
  
  .worldinfo-view {
    flex-direction: column;
    height: calc(100vh - 64px - 48px);
  }
  
  .worldinfo-main {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  
  .main-header {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
    padding: 16px;
  }
  
  .header-actions {
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  
  .entries-toolbar {
    flex-direction: column;
    gap: 12px;
    padding: 12px 16px;
  }
  
  .entry-search {
    max-width: none;
  }
  
  .entry-sort {
    width: 100%;
  }
  
  .entry-header {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .entry-keys {
    order: 3;
    width: 100%;
    margin-top: 8px;
  }
  
  .entry-content p {
    font-size: 13px;
  }
  
  .entry-meta {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .entries-list {
    padding: 16px;
    flex: 1;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
  }
  
  .entry-item {
    margin-right: 0;
    padding: 12px;
  }
}

@media (max-width: 480px) {
  .worldinfo-view {
    border-radius: 0;
    border: none;
    margin: 0 -16px;
  }
  
  .main-header,
  .entries-toolbar,
  .entries-list {
    padding-left: 16px;
    padding-right: 16px;
  }
  
  .entry-item {
    padding: 12px;
    margin-bottom: 8px;
  }
  
  .entry-header {
    gap: 6px;
  }
  
  .entry-content p {
    font-size: 12px;
  }
  
  .header-actions {
    gap: 6px;
  }
  
  .header-actions .n-button {
    font-size: 12px;
    padding: 0 12px;
  }
  
  .mobile-list-toggle {
    padding: 12px 16px;
  }
}

/* 测试按钮 */
.test-buttons {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
  background: var(--bg-secondary);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

@media (max-width: 768px) {
  .test-buttons {
    bottom: 10px;
    right: 10px;
    padding: 8px;
  }
  
  .test-buttons .n-button {
    font-size: 12px;
    padding: 8px 12px;
  }
}
</style>
