<template>
  <div class="worldinfo-view">
    <!-- 左侧：世界书列表 -->
    <div class="worldinfo-sidebar">
      <div class="sidebar-header">
        <h3>世界书</h3>
        <n-button quaternary circle size="small" @click="createWorld">
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
            <div class="world-count">{{ world.entries.length }} 条目</div>
          </div>
          <n-switch
            v-model:value="world.enabled"
            size="small"
            @click.stop
          />
        </div>

        <n-empty v-if="filteredWorlds.length === 0" description="暂无世界书" />
      </n-scrollbar>

      <div class="sidebar-footer">
        <n-button block @click="importWorld">
          <template #icon>
            <n-icon><CloudUploadOutline /></n-icon>
          </template>
          导入世界书
        </n-button>
      </div>
    </div>

    <!-- 右侧：条目列表 -->
    <div class="worldinfo-main">
      <template v-if="selectedWorld">
        <!-- 世界书头部 -->
        <div class="main-header">
          <div class="header-info">
            <h2>{{ selectedWorld.name }}</h2>
            <p>{{ selectedWorld.description || '暂无描述' }}</p>
          </div>
          <div class="header-actions">
            <n-button @click="editWorld">
              <template #icon>
                <n-icon><CreateOutline /></n-icon>
              </template>
              编辑
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
          <TransitionGroup name="entry-list">
            <div
              v-for="entry in filteredEntries"
              :key="entry.id"
              class="entry-item"
              :class="{ 'disabled': !entry.enabled }"
            >
              <div class="entry-header">
                <n-switch v-model:value="entry.enabled" size="small" />
                <div class="entry-keys">
                  <n-tag
                    v-for="key in entry.keys.slice(0, 3)"
                    :key="key"
                    size="small"
                    type="primary"
                    :bordered="false"
                  >
                    {{ key }}
                  </n-tag>
                  <n-tag v-if="entry.keys.length > 3" size="small" :bordered="false">
                    +{{ entry.keys.length - 3 }}
                  </n-tag>
                </div>
                <div class="entry-actions">
                  <n-button quaternary circle size="tiny" @click="editEntry(entry)">
                    <template #icon>
                      <n-icon size="14"><CreateOutline /></n-icon>
                    </template>
                  </n-button>
                  <n-button quaternary circle size="tiny" @click="deleteEntry(entry)">
                    <template #icon>
                      <n-icon size="14"><TrashOutline /></n-icon>
                    </template>
                  </n-button>
                </div>
              </div>

              <div class="entry-content">
                <p>{{ entry.content }}</p>
              </div>

              <div class="entry-meta">
                <span>优先级: {{ entry.order }}</span>
                <span>深度: {{ entry.depth }}</span>
                <span>概率: {{ entry.probability }}%</span>
              </div>
            </div>
          </TransitionGroup>

          <n-empty v-if="filteredEntries.length === 0" description="暂无条目" />
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

    <!-- 条目编辑模态框 -->
    <n-modal
      v-model:show="showEntryModal"
      preset="card"
      :title="editingEntry ? '编辑条目' : '添加条目'"
      :style="{ width: '600px', maxWidth: '90vw' }"
    >
      <WorldInfoEntryEditor
        :entry="editingEntry"
        @save="handleSaveEntry"
        @cancel="showEntryModal = false"
      />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
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
  useMessage,
  useDialog
} from 'naive-ui';
import {
  AddOutline,
  SearchOutline,
  GlobeOutline,
  CloudUploadOutline,
  CreateOutline,
  TrashOutline
} from '@vicons/ionicons5';

import WorldInfoEntryEditor from '../components/worldinfo/WorldInfoEntryEditor.vue';
import type { WorldInfo, WorldInfoEntry } from '../types';

const message = useMessage();
const dialog = useDialog();

// 状态
const searchQuery = ref('');
const entrySearchQuery = ref('');
const entrySort = ref('order');
const selectedWorldId = ref<string | null>('1');
const showEntryModal = ref(false);
const editingEntry = ref<WorldInfoEntry | null>(null);

// 排序选项
const sortOptions = [
  { label: '按优先级', value: 'order' },
  { label: '按名称', value: 'name' },
  { label: '按修改时间', value: 'updated' }
];

// 模拟数据
const worlds = ref<WorldInfo[]>([
  {
    id: '1',
    name: '奇幻世界',
    description: '一个充满魔法和冒险的奇幻世界设定',
    enabled: true,
    entries: [
      {
        id: '1',
        worldId: '1',
        keys: ['魔法', '法术', 'magic'],
        content: '在这个世界中，魔法是一种可以被少数人掌握的神秘力量...',
        enabled: true,
        order: 100,
        probability: 100,
        depth: 4,
        selectiveLogic: 'or'
      },
      {
        id: '2',
        worldId: '1',
        keys: ['精灵', 'elf'],
        content: '精灵是这个世界中最古老的种族之一，他们居住在森林深处...',
        enabled: true,
        order: 90,
        probability: 100,
        depth: 4,
        selectiveLogic: 'or'
      }
    ]
  },
  {
    id: '2',
    name: '科幻设定',
    description: '未来太空时代的世界设定',
    enabled: false,
    entries: []
  }
]);

// 计算属性
const filteredWorlds = computed(() => {
  if (!searchQuery.value) return worlds.value;
  return worlds.value.filter(w =>
    w.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const selectedWorld = computed(() =>
  worlds.value.find(w => w.id === selectedWorldId.value)
);

const filteredEntries = computed(() => {
  if (!selectedWorld.value) return [];
  let entries = [...selectedWorld.value.entries];

  if (entrySearchQuery.value) {
    entries = entries.filter(e =>
      e.keys.some(k => k.toLowerCase().includes(entrySearchQuery.value.toLowerCase())) ||
      e.content.toLowerCase().includes(entrySearchQuery.value.toLowerCase())
    );
  }

  // 排序
  if (entrySort.value === 'order') {
    entries.sort((a, b) => b.order - a.order);
  }

  return entries;
});

// 方法
const selectWorld = (id: string) => {
  selectedWorldId.value = id;
};

const createWorld = () => {
  message.info('创建世界书');
};

const editWorld = () => {
  message.info('编辑世界书');
};

const importWorld = () => {
  message.info('导入世界书');
};

const addEntry = () => {
  editingEntry.value = null;
  showEntryModal.value = true;
};

const editEntry = (entry: WorldInfoEntry) => {
  editingEntry.value = entry;
  showEntryModal.value = true;
};

const deleteEntry = (entry: WorldInfoEntry) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这个条目吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      if (selectedWorld.value) {
        selectedWorld.value.entries = selectedWorld.value.entries.filter(e => e.id !== entry.id);
        message.success('条目已删除');
      }
    }
  });
};

const handleSaveEntry = (entry: WorldInfoEntry) => {
  if (!selectedWorld.value) return;

  if (editingEntry.value) {
    const index = selectedWorld.value.entries.findIndex(e => e.id === editingEntry.value!.id);
    if (index >= 0) {
      selectedWorld.value.entries[index] = entry;
    }
    message.success('条目已更新');
  } else {
    selectedWorld.value.entries.push({
      ...entry,
      id: Date.now().toString(),
      worldId: selectedWorld.value.id
    });
    message.success('条目已添加');
  }

  showEntryModal.value = false;
  editingEntry.value = null;
};
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
}

/* 侧边栏 */
.worldinfo-sidebar {
  width: 280px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-header h3 {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
}

.search-input {
  margin: 12px 16px;
  width: calc(100% - 32px);
  box-sizing: border-box;
}

.world-list {
  flex: 1;
  padding: 0 8px;
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
  padding: 16px 20px;
}

.entry-item {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 12px;
  transition: all var(--transition-fast);
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
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
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
    position: absolute;
    z-index: 10;
  }
}
</style>
