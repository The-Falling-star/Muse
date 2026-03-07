<template>
  <div class="characters-view">
    <!-- 顶部工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索角色..."
          clearable
          class="search-input"
        >
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>

      </div>

      <div class="toolbar-right">
        <n-button-group>
          <n-button
            :type="viewMode === 'grid' ? 'primary' : 'default'"
            @click="viewMode = 'grid'"
          >
            <template #icon>
              <n-icon><GridOutline /></n-icon>
            </template>
          </n-button>
          <n-button
            :type="viewMode === 'list' ? 'primary' : 'default'"
            @click="viewMode = 'list'"
          >
            <template #icon>
              <n-icon><ListOutline /></n-icon>
            </template>
          </n-button>
        </n-button-group>

        <n-button type="primary" @click="showImportModal = true">
          <template #icon>
            <n-icon><CloudUploadOutline /></n-icon>
          </template>
          导入
        </n-button>

        <n-button type="primary" @click="showCreateModal = true">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          创建
        </n-button>
      </div>
    </div>

    <!-- 角色列表 -->
    <n-scrollbar class="characters-container">
      <!-- 网格视图 -->
      <div v-if="viewMode === 'grid'" class="characters-grid">
        <TransitionGroup name="card-list">
          <CharacterCard
            v-for="character in filteredCharacters"
            :key="character.id"
            :character="character"
            @click="selectCharacter(character)"
            @edit="editCharacter(character)"
            @delete="handleDeleteCharacter(character)"
            @chat="startChat(character)"
          />
        </TransitionGroup>
      </div>

      <!-- 列表视图 -->
      <div v-else class="characters-list">
        <n-data-table
          :columns="tableColumns"
          :data="filteredCharacters"
          :row-key="(row: Character) => row.id"
          :bordered="false"
          :single-line="false"
          striped
        />
      </div>

      <!-- 分页器 -->
      <div v-if="filteredCharacters.length > 0" class="pagination-container">
        <n-pagination
          v-model:page="page"
          v-model:page-size="pageSize"
          :item-count="total"
          :page-sizes="[10, 20, 30, 50]"
          show-size-picker
          show-quick-jumper
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        >
          <template #prefix="{ itemCount }">
            共 {{ itemCount }} 个角色
          </template>
        </n-pagination>
      </div>

      <!-- 空状态 -->
      <n-empty
        v-if="filteredCharacters.length === 0"
        description="暂无角色"
        class="empty-state"
      >
        <template #extra>
          <n-space>
            <n-button @click="showImportModal = true">导入角色卡</n-button>
            <n-button type="primary" @click="showCreateModal = true">创建新角色</n-button>
          </n-space>
        </template>
      </n-empty>
    </n-scrollbar>

    <!-- 角色详情抽屉 -->
    <n-drawer
      v-model:show="showDetailDrawer"
      :width="500"
      placement="right"
    >
      <n-drawer-content :title="selectedCharacter?.name || '角色详情'" closable>
        <CharacterDetail
          v-if="selectedCharacter"
          :character="selectedCharacter"
          @edit="handleEditSelected"
          @chat="handleChatSelected"
        />
      </n-drawer-content>
    </n-drawer>

    <!-- 创建/编辑模态框 -->
    <n-modal
      v-model:show="showCreateModal"
      preset="card"
      :title="editingCharacter ? '编辑角色' : '创建角色'"
      :style="{ width: '700px', maxWidth: '90vw' }"
      :mask-closable="false"
    >
      <CharacterEditor
        :character="editingCharacter"
        @save="handleSaveCharacter"
        @cancel="showCreateModal = false"
      />
    </n-modal>

    <!-- 导入模态框 -->
    <n-modal
      v-model:show="showImportModal"
      preset="card"
      title="导入角色卡"
      :style="{ width: '500px', maxWidth: '90vw' }"
    >
      <div class="import-modal">
        <n-upload
          multiple
          directory-dnd
          :max="10"
          accept=".png,.json"
          :default-upload="false"
          @change="handleImportFiles"
        >
          <n-upload-dragger>
            <div class="upload-content">
              <n-icon size="48" class="upload-icon"><CloudUploadOutline /></n-icon>
              <p class="upload-text">点击或拖拽文件到此处</p>
              <p class="upload-hint">支持 PNG (带嵌入数据) 或 JSON 格式的角色卡</p>
            </div>
          </n-upload-dragger>
        </n-upload>

        <n-divider>或</n-divider>

        <n-input
          v-model:value="importUrl"
          placeholder="输入角色卡链接 (CharacterHub, Chub等)"
        />

        <n-button
          type="primary"
          block
          :disabled="!importUrl"
          class="import-btn"
          @click="handleImportFromUrl"
        >
          从链接导入
        </n-button>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import {
  NInput,
  NButton,
  NButtonGroup,
  NIcon,
  NScrollbar,
  NDataTable,
  NDrawer,
  NDrawerContent,
  NModal,
  NUpload,
  NUploadDragger,
  NDivider,
  NEmpty,
  NSpace,
  NAvatar,
  NPagination,
  useMessage,
  useDialog
} from 'naive-ui';
import type { DataTableColumns, UploadFileInfo } from 'naive-ui';
import {
  SearchOutline,
  GridOutline,
  ListOutline,
  CloudUploadOutline,
  AddOutline,
  ChatbubbleOutline,
  CreateOutline,
  TrashOutline
} from '@vicons/ionicons5';

import CharacterCard from '@/components/character/CharacterCard.vue';
import CharacterDetail from '@/components/character/CharacterDetail.vue';
import CharacterEditor from '@/components/character/CharacterEditor.vue';
import { useCharacterStore } from '@/stores/character';
import { useFileStore } from '@/stores/file';
import { characterClient, chatClient } from '@/api/client';
import type { Character } from '@/gen/muse/character_pb';
import {generateSessionName} from "@/utils/common.ts";

const message = useMessage();
const dialog = useDialog();
const router = useRouter();
const characterStore = useCharacterStore();
const fileStore = useFileStore();

// 状态
const searchQuery = ref('');
const viewMode = ref<'grid' | 'list'>('grid');
const showDetailDrawer = ref(false);
const showCreateModal = ref(false);
const showImportModal = ref(false);
const selectedCharacter = ref<Character | null>(null);
const editingCharacter = ref<Character | null>(null);
const importUrl = ref('');
const loading = ref(false);

// 分页状态
const page = ref(1);
const pageSize = ref(20);
const total = computed(() => characterStore.total);

// 加载角色列表 - 直接调用client
const loadCharacters = async () => {
  loading.value = true;
  const response = await characterClient.listCharacters({
    page: page.value,
    pageSize: pageSize.value
  });
  // 更新Store缓存
  characterStore.setCharacters(response.characters, response.total);
  
  // 预加载头像
  const avatarPaths = response.characters
    .map(c => c.avatar)
    .filter((path): path is string => !!path && !path.startsWith('http') && !path.startsWith('data:'));
  if (avatarPaths.length > 0) {
    fileStore.preloadFiles(avatarPaths);
  }
  
  loading.value = false;
};

onMounted(() => {
  loadCharacters();
});

// 过滤后的角色列表
const filteredCharacters = computed(() => {
  const chars = characterStore.characters;
  if (!chars || chars.length === 0) {
    return [];
  }
  if (!searchQuery.value) {
    return chars;
  }
  const query = searchQuery.value.toLowerCase();
  return chars.filter(c =>
    c.name.toLowerCase().includes(query)
  );
});

// 表格列配置
const tableColumns: DataTableColumns<Character> = [
  {
    title: '头像',
    key: 'avatar',
    width: 60,
    render(row) {
      const avatarSrc = fileStore.getCachedUrl(row.avatar) || '';
      return h(NAvatar, {
        size: 40,
        round: true,
        src: avatarSrc,
        style: 'background: var(--gradient-primary)'
      }, { default: () => row.name.charAt(0) });
    }
  },
  {
    title: '名称',
    key: 'name',
    sorter: 'default'
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row) {
      return h('div', { style: 'display: flex; gap: 8px;' }, [
        h(NButton, {
          size: 'small',
          quaternary: true,
          circle: true,
          onClick: () => startChat(row)
        }, { icon: () => h(NIcon, null, { default: () => h(ChatbubbleOutline) }) }),
        h(NButton, {
          size: 'small',
          quaternary: true,
          circle: true,
          onClick: () => editCharacter(row)
        }, { icon: () => h(NIcon, null, { default: () => h(CreateOutline) }) }),
        h(NButton, {
          size: 'small',
          quaternary: true,
          circle: true,
          onClick: () => handleDeleteCharacter(row)
        }, { icon: () => h(NIcon, null, { default: () => h(TrashOutline) }) })
      ]);
    }
  }
];

// 方法
const selectCharacter = async (character: Character) => {
  // 先检查缓存
  const cachedDetail = characterStore.getCharacterDetail(character.id);
  if (cachedDetail) {
    selectedCharacter.value = cachedDetail;
    showDetailDrawer.value = true;
    return;
  }

  // 缓存未命中，请求后端
  loading.value = true;
  const response = await characterClient.getCharacter({ id: character.id });
  if (response.character) {
    selectedCharacter.value = response.character;
    // 缓存详情
    characterStore.cacheCharacterDetail(response.character);
  } else {
    selectedCharacter.value = character;
  }
  loading.value = false;
  showDetailDrawer.value = true;
};

const editCharacter = async (character: Character | null) => {
  if (character) {
    // 先检查缓存
    const cachedDetail = characterStore.getCharacterDetail(character.id);
    if (cachedDetail) {
      editingCharacter.value = cachedDetail;
      showCreateModal.value = true;
      showDetailDrawer.value = false;
      return;
    }

    // 缓存未命中，请求后端
    loading.value = true;
    const response = await characterClient.getCharacter({ id: character.id });
    if (response.character) {
      editingCharacter.value = response.character;
      // 缓存详情
      characterStore.cacheCharacterDetail(response.character);
    } else {
      editingCharacter.value = character;
    }
    loading.value = false;
  } else {
    editingCharacter.value = null;
  }
  showCreateModal.value = true;
  showDetailDrawer.value = false;
};

// 删除角色 - 直接调用client
const handleDeleteCharacter = (character: Character) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除角色"${character.name}"吗？此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await characterClient.deleteCharacter({ id: character.id });
      characterStore.removeCharacter(character.id);
      message.success('角色已删除');
    }
  });
};

// 开始聊天 - 直接调用client
const startChat = async (character: Character) => {
  const response = await chatClient.createChatSession({
    characterId: character.id,
    name: generateSessionName(character.name),
  });
  if (response.session) {
    router.push(`/chat/${response.session.id}`);
  }
};

// 处理从详情页触发的编辑和聊天
const handleEditSelected = () => {
  if (selectedCharacter.value) {
    editCharacter(selectedCharacter.value);
  }
};

const handleChatSelected = () => {
  if (selectedCharacter.value) {
    startChat(selectedCharacter.value);
  }
};

// 保存角色 - 直接调用client
const handleSaveCharacter = async (characterData: Partial<Character>) => {
  loading.value = true;
  if (editingCharacter.value) {
    // 编辑模式
    const response = await characterClient.updateCharacter({
      id: editingCharacter.value.id,
      name: characterData.name || '',
      avatar: characterData.avatar,
      description: characterData.description,
      firstMessage: characterData.firstMessage,
      exampleDialogue: characterData.exampleDialogue,
      creatorNotes: characterData.creatorNotes,
      version: editingCharacter.value.version
    });
    if (response.character) {
      characterStore.updateCharacterInList(response.character);
    }
    message.success('角色已更新');
  } else {
    // 创建模式
    const response = await characterClient.createCharacter({
      name: characterData.name || '',
      avatar: characterData.avatar,
      description: characterData.description,
      firstMessage: characterData.firstMessage,
      exampleDialogue: characterData.exampleDialogue,
      creatorNotes: characterData.creatorNotes
    });
    if (response.character) {
      characterStore.addCharacter(response.character);
    }
    message.success('角色已创建');
  }
  showCreateModal.value = false;
  editingCharacter.value = null;
  loading.value = false;
};

// 导入文件 - 直接调用client
const handleImportFiles = async (options: { file: UploadFileInfo; fileList: UploadFileInfo[]; event?: Event }) => {
  const files = options.fileList;
  if (files.length === 0) return;

  let successCount = 0;
  let failCount = 0;
  const importedCharacters: Character[] = [];

  for (const fileInfo of files) {
    const file = fileInfo.file;
    if (!file) continue;

    const arrayBuffer = await file.arrayBuffer();
    const fileContent = new Uint8Array(arrayBuffer);
    const response = await characterClient.importCharacter({
      fileContent: fileContent,
      fileName: file.name
    });
    if (response.character) {
      characterStore.addCharacter(response.character);
      importedCharacters.push(response.character);
      successCount++;
    } else {
      failCount++;
    }
  }

  if (successCount > 0) {
    message.success(`成功导入 ${successCount} 个角色`);
  }
  if (failCount > 0) {
    message.warning(`${failCount} 个角色导入失败`);
  }

  showImportModal.value = false;
};

const handleImportFromUrl = () => {
  message.info('暂不支持从链接导入');
  importUrl.value = '';
};

// 分页变化处理
const handlePageChange = (newPage: number) => {
  page.value = newPage;
  loadCharacters();
};

const handlePageSizeChange = (newPageSize: number) => {
  pageSize.value = newPageSize;
  page.value = 1;
  loadCharacters();
};
</script>

<style scoped>
.characters-view {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 64px - 48px);
}

/* 工具栏 */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.toolbar-left {
  display: flex;
  gap: 12px;
  flex: 1;
  min-width: 200px;
}

.search-input {
  max-width: 300px;
}

.toolbar-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-left,
  .toolbar-right {
    width: 100%;
  }

  .search-input {
    max-width: none;
    flex: 1;
  }
}

/* 角色容器 */
.characters-container {
  flex: 1;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

/* 网格视图 */
.characters-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  padding-bottom: 20px;
}

/* 卡片列表过渡 */
.card-list-enter-active,
.card-list-leave-active {
  transition: all var(--transition-normal);
}

.card-list-enter-from {
  opacity: 0;
  transform: translateY(20px) scale(0.95);
}

.card-list-leave-to {
  opacity: 0;
  transform: translateY(-20px) scale(0.95);
}

.card-list-move {
  transition: transform var(--transition-normal);
}

/* 移动端优化 */
@media (max-width: 768px) {
  .characters-grid {
    grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
    gap: 12px;
  }
  
  .card-list-enter-from {
    transform: translateX(30px);
  }
  
  .card-list-leave-to {
    transform: translateX(-30px);
  }
}

/* 列表视图 */
.characters-list {
  background: var(--bg-card);
  border-radius: 12px;
  overflow: hidden;
  transition: all var(--transition-fast);
}

.characters-list:hover {
  box-shadow: var(--shadow-md);
}

/* 分页器容器 */
.pagination-container {
  display: flex;
  justify-content: center;
  padding: 20px 0;
  margin-top: 16px;
}

/* 空状态 */
.empty-state {
  padding: 60px 20px;
}

/* 导入模态框 */
.import-modal {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
}

.upload-icon {
  color: var(--color-primary);
  margin-bottom: 16px;
}

.upload-text {
  font-size: 16px;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.upload-hint {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0;
}

.import-btn {
  margin-top: 8px;
}
</style>
