<template>
  <aside class="left-sidebar">
    <!-- 会话列表视图 -->
    <template v-if="appStore.sidebarView === 'sessions'">
      <!-- 新建对话按钮 -->
      <div class="sidebar-top">
        <n-button class="new-chat-btn" block quaternary @click="handleNewChat">
          <template #icon>
            <n-icon size="18"><AddOutline /></n-icon>
          </template>
          新建对话
        </n-button>
      </div>

      <!-- 搜索框 -->
      <div class="sidebar-search">
        <n-input
          v-model:value="searchQuery"
          placeholder="搜索对话..."
          clearable
          size="small"
        >
          <template #prefix>
            <n-icon size="16"><SearchOutline /></n-icon>
          </template>
        </n-input>
      </div>

      <!-- 会话列表 -->
      <n-scrollbar class="session-list-scroll">
        <div class="session-list">
          <template v-if="filteredGroupedSessions.length > 0">
            <div
              v-for="group in filteredGroupedSessions"
              :key="group.label"
              class="session-group"
            >
              <div class="group-label">{{ group.label }}</div>
              <div
                v-for="session in group.sessions"
                :key="session.id"
                class="session-item"
                :class="{ active: chatStore.activeSessionId === session.id }"
                @click="handleSelectSession(session)"
                @contextmenu.prevent="handleSessionContextMenu($event, session)"
              >
                <!-- 角色头像 -->
                <n-avatar
                  v-if="session.character?.avatar"
                  :src="session.character.avatar"
                  :size="24"
                  round
                  class="session-avatar"
                />
                <div v-else class="session-avatar-placeholder">
                  <n-icon size="14"><ChatbubblesOutline /></n-icon>
                </div>

                <!-- 会话标题 -->
                <span class="session-title">{{ session.name || '新对话' }}</span>

                <!-- 更多操作按钮 -->
                <n-button
                  quaternary
                  circle
                  size="tiny"
                  class="session-more-btn"
                  @click.stop="handleSessionContextMenu($event, session)"
                >
                  <template #icon>
                    <n-icon size="14"><EllipsisHorizontalOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </template>

          <!-- 空状态 -->
          <div v-else class="empty-state">
            <n-icon size="32" color="var(--text-tertiary)"><ChatbubblesOutline /></n-icon>
            <span>{{ searchQuery ? '没有找到匹配的对话' : '暂无对话记录' }}</span>
          </div>
        </div>
      </n-scrollbar>

      <!-- 底部快捷入口 -->
      <div class="sidebar-bottom">
        <div class="bottom-entry" @click="appStore.setSidebarView('characters')">
          <n-icon size="18"><PeopleOutline /></n-icon>
          <span>角色管理</span>
          <n-icon size="14" class="entry-arrow"><ChevronForwardOutline /></n-icon>
        </div>
        <div class="bottom-entry" @click="appStore.openSettings()">
          <n-icon size="18"><SettingsOutline /></n-icon>
          <span>设置</span>
        </div>
      </div>
    </template>

    <!-- 角色管理视图 -->
    <template v-else-if="appStore.sidebarView === 'characters'">
      <!-- 顶部返回栏 -->
      <div class="sidebar-top character-top">
        <n-button quaternary size="small" @click="appStore.setSidebarView('sessions')">
          <template #icon>
            <n-icon size="16"><ArrowBackOutline /></n-icon>
          </template>
          返回
        </n-button>
        <span class="panel-title">角色管理</span>
      </div>

      <!-- 角色搜索和导入 -->
      <div class="sidebar-search">
        <n-input
          v-model:value="characterSearchQuery"
          placeholder="搜索角色..."
          clearable
          size="small"
        >
          <template #prefix>
            <n-icon size="16"><SearchOutline /></n-icon>
          </template>
        </n-input>
      </div>

      <div class="character-actions">
        <n-button size="small" quaternary block @click="handleImportCharacter">
          <template #icon>
            <n-icon size="16"><CloudUploadOutline /></n-icon>
          </template>
          导入角色
        </n-button>
      </div>

      <!-- 角色列表 -->
      <n-scrollbar class="character-list-scroll">
        <div class="character-list">
          <div
            v-for="char in filteredCharacters"
            :key="char.id"
            class="character-item"
            @click="handleSelectCharacter(char)"
            @contextmenu.prevent="handleCharacterContextMenu($event, char)"
          >
            <n-avatar
              v-if="char.avatar"
              :src="char.avatar"
              :size="40"
              round
              class="character-avatar"
            />
            <div v-else class="character-avatar-default">
              <n-icon size="20"><PersonOutline /></n-icon>
            </div>
            <div class="character-info">
              <span class="character-name">{{ char.name || '未命名角色' }}</span>
              <span class="character-desc">{{ char.description || '暂无描述' }}</span>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredCharacters.length === 0" class="empty-state">
            <n-icon size="32" color="var(--text-tertiary)"><PeopleOutline /></n-icon>
            <span>{{ characterSearchQuery ? '没有找到匹配的角色' : '暂无角色' }}</span>
          </div>
        </div>
      </n-scrollbar>
    </template>

    <!-- 会话右键菜单 -->
    <n-dropdown
      trigger="manual"
      :show="showSessionMenu"
      :options="sessionMenuOptions"
      :x="menuX"
      :y="menuY"
      placement="bottom-start"
      @select="handleSessionMenuSelect"
      @clickoutside="showSessionMenu = false"
    />

    <!-- 角色右键菜单 -->
    <n-dropdown
      trigger="manual"
      :show="showCharacterMenu"
      :options="characterMenuOptions"
      :x="menuX"
      :y="menuY"
      placement="bottom-start"
      @select="handleCharacterMenuSelect"
      @clickoutside="showCharacterMenu = false"
    />
  </aside>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import {
  NButton, NIcon, NInput, NScrollbar, NAvatar, NDropdown
} from 'naive-ui';
import {
  AddOutline,
  SearchOutline,
  ChatbubblesOutline,
  EllipsisHorizontalOutline,
  PeopleOutline,
  SettingsOutline,
  ChevronForwardOutline,
  ArrowBackOutline,
  CloudUploadOutline,
  PersonOutline
} from '@vicons/ionicons5';

import { useAppStore } from '@/stores/app';
import { useChatStore } from '@/stores/chat';
import { useCharacterStore } from '@/stores/character';
import type { ChatSession, Character } from '@/gen/muse/muse_pb';

const router = useRouter();
const appStore = useAppStore();
const chatStore = useChatStore();
const characterStore = useCharacterStore();

// ====== 搜索状态 ======
const searchQuery = ref('');
const characterSearchQuery = ref('');

// ====== 右键菜单状态 ======
const showSessionMenu = ref(false);
const showCharacterMenu = ref(false);
const menuX = ref(0);
const menuY = ref(0);
const contextSession = ref<ChatSession | null>(null);
const contextCharacter = ref<Character | null>(null);

// ====== 会话按时间分组 ======
interface SessionGroup {
  label: string;
  sessions: ChatSession[];
}

const groupSessionsByTime = (sessions: ChatSession[]): SessionGroup[] => {
  const now = Date.now();
  const todayStart = new Date();
  todayStart.setHours(0, 0, 0, 0);
  const yesterdayStart = new Date(todayStart);
  yesterdayStart.setDate(yesterdayStart.getDate() - 1);
  const weekStart = new Date(todayStart);
  weekStart.setDate(weekStart.getDate() - 7);

  const groups: Record<string, ChatSession[]> = {
    '今天': [],
    '昨天': [],
    '过去 7 天': [],
    '更早': []
  };

  for (const session of sessions) {
    const updatedAt = Number(session.updatedAt);
    // updatedAt 可能是秒级或毫秒级时间戳
    const ts = updatedAt > 1e12 ? updatedAt : updatedAt * 1000;

    if (ts >= todayStart.getTime()) {
      groups['今天'].push(session);
    } else if (ts >= yesterdayStart.getTime()) {
      groups['昨天'].push(session);
    } else if (ts >= weekStart.getTime()) {
      groups['过去 7 天'].push(session);
    } else {
      groups['更早'].push(session);
    }
  }

  return Object.entries(groups)
    .filter(([, sessions]) => sessions.length > 0)
    .map(([label, sessions]) => ({ label, sessions }));
};

// ====== 计算属性 ======
const filteredGroupedSessions = computed(() => {
  let sessions = chatStore.sessions;
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    sessions = sessions.filter(s =>
      (s.name || '').toLowerCase().includes(query)
    );
  }
  return groupSessionsByTime(sessions);
});

const filteredCharacters = computed(() => {
  let chars = characterStore.characters;
  if (characterSearchQuery.value) {
    const query = characterSearchQuery.value.toLowerCase();
    chars = chars.filter(c =>
      (c.name || '').toLowerCase().includes(query) ||
      (c.description || '').toLowerCase().includes(query)
    );
  }
  return chars;
});

// ====== 会话操作 ======
const handleNewChat = () => {
  // 跳转到对话视图，由对话视图处理新建逻辑
  router.push('/');
  // 移动端自动关闭侧边栏
  if (appStore.isCompactMode) {
    appStore.closeLeftSidebar();
  }
};

const handleSelectSession = (session: ChatSession) => {
  chatStore.setActiveSession(session);
  router.push('/');
  if (appStore.isCompactMode) {
    appStore.closeLeftSidebar();
  }
};

const handleSessionContextMenu = (e: MouseEvent, session: ChatSession) => {
  contextSession.value = session;
  menuX.value = e.clientX;
  menuY.value = e.clientY;
  showSessionMenu.value = true;
  showCharacterMenu.value = false;
};

const sessionMenuOptions = [
  { label: '重命名', key: 'rename' },
  { label: '删除', key: 'delete' },
  { label: '导出', key: 'export' }
];

const handleSessionMenuSelect = (key: string) => {
  showSessionMenu.value = false;
  if (!contextSession.value) return;

  switch (key) {
    case 'rename':
      // TODO: 实现重命名逻辑
      break;
    case 'delete':
      // TODO: 实现删除逻辑
      break;
    case 'export':
      // TODO: 实现导出逻辑
      break;
  }
};

// ====== 角色操作 ======
const handleSelectCharacter = (char: Character) => {
  characterStore.selectCharacter(char);
  // 以该角色创建新会话或显示角色详情
  // TODO: 实现角色选择后的行为
  if (appStore.isCompactMode) {
    appStore.closeLeftSidebar();
  }
};

const handleImportCharacter = () => {
  // TODO: 实现角色导入逻辑（文件选择器）
};

const handleCharacterContextMenu = (e: MouseEvent, char: Character) => {
  contextCharacter.value = char;
  menuX.value = e.clientX;
  menuY.value = e.clientY;
  showCharacterMenu.value = true;
  showSessionMenu.value = false;
};

const characterMenuOptions = [
  { label: '编辑', key: 'edit' },
  { label: '删除', key: 'delete' },
  { label: '导出', key: 'export' }
];

const handleCharacterMenuSelect = (key: string) => {
  showCharacterMenu.value = false;
  if (!contextCharacter.value) return;

  switch (key) {
    case 'edit':
      // TODO: 实现编辑逻辑
      break;
    case 'delete':
      // TODO: 实现删除逻辑
      break;
    case 'export':
      // TODO: 实现导出逻辑
      break;
  }
};
</script>

<style scoped>
.left-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
  overflow: hidden;
}

/* ====== 顶部区域 ====== */
.sidebar-top {
  padding: 12px 12px 0;
  flex-shrink: 0;
}

.character-top {
  display: flex;
  align-items: center;
  gap: 8px;
}

.panel-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.new-chat-btn {
  justify-content: flex-start;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  height: 40px;
  font-weight: 500;
  color: var(--text-primary);
  transition: background-color 150ms ease;
}

.new-chat-btn:hover {
  background: var(--bg-hover);
}

/* ====== 搜索框 ====== */
.sidebar-search {
  padding: 12px 12px 8px;
  flex-shrink: 0;
}

/* ====== 角色操作区 ====== */
.character-actions {
  padding: 0 12px 8px;
  flex-shrink: 0;
}

/* ====== 会话列表 ====== */
.session-list-scroll {
  flex: 1;
  min-height: 0;
}

.session-list {
  padding: 0 8px 8px;
}

.session-group {
  margin-bottom: 4px;
}

.group-label {
  padding: 8px 8px 4px;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-tertiary);
  user-select: none;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 100ms ease;
  position: relative;
}

.session-item:hover {
  background: var(--bg-hover);
}

.session-item.active {
  background: var(--bg-active);
}

.session-avatar {
  flex-shrink: 0;
}

.session-avatar-placeholder {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--text-tertiary);
}

.session-title {
  flex: 1;
  font-size: 13px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.4;
}

.session-more-btn {
  opacity: 0;
  flex-shrink: 0;
  transition: opacity 100ms ease;
}

.session-item:hover .session-more-btn {
  opacity: 1;
}

/* ====== 角色列表 ====== */
.character-list-scroll {
  flex: 1;
  min-height: 0;
}

.character-list {
  padding: 0 8px 8px;
}

.character-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 100ms ease;
}

.character-item:hover {
  background: var(--bg-hover);
}

.character-avatar {
  flex-shrink: 0;
}

.character-avatar-default {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: var(--text-tertiary);
}

.character-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.character-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.character-desc {
  font-size: 12px;
  color: var(--text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ====== 底部快捷入口 ====== */
.sidebar-bottom {
  flex-shrink: 0;
  border-top: 1px solid var(--border-color);
  padding: 8px;
}

.bottom-entry {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: background-color 100ms ease;
  user-select: none;
}

.bottom-entry:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.bottom-entry span {
  flex: 1;
}

.entry-arrow {
  color: var(--text-tertiary);
}

/* ====== 空状态 ====== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 16px;
  color: var(--text-tertiary);
  font-size: 13px;
}

/* ====== 移动端触摸优化 ====== */
@media (max-width: 767px) {
  .session-item {
    padding: 10px 12px;
    min-height: 48px;
  }

  .session-more-btn {
    opacity: 1;
  }

  .character-item {
    min-height: 56px;
  }

  .bottom-entry {
    min-height: 44px;
  }
}
</style>
