<template>
  <aside class="left-sidebar">
    <!-- 会话列表视图 -->
    <template v-if="appStore.sidebarView === 'sessions'">
      <!-- 新建对话按钮 -->
      <div class="sidebar-top">
        <n-button class="new-chat-btn" block quaternary @click="handleNewChat">
          <template #icon>
            <n-icon size="18">
              <AddOutline/>
            </n-icon>
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
            <n-icon size="16">
              <SearchOutline/>
            </n-icon>
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
                    :src="getAvatarUrlSync(session.character.avatar)"
                    :size="24"
                    round
                    class="session-avatar"
                />
                <div v-else class="session-avatar-placeholder">
                  <n-icon size="14">
                    <ChatbubblesOutline/>
                  </n-icon>
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
                    <n-icon size="14">
                      <EllipsisHorizontalOutline/>
                    </n-icon>
                  </template>
                </n-button>
              </div>
            </div>
          </template>

          <!-- 空状态 -->
          <div v-else class="empty-state">
            <n-icon size="32" color="var(--text-tertiary)">
              <ChatbubblesOutline/>
            </n-icon>
            <span>{{ searchQuery ? '没有找到匹配的对话' : '暂无对话记录' }}</span>
          </div>
        </div>
      </n-scrollbar>

      <!-- 底部快捷入口 -->
      <div class="sidebar-bottom">
        <div class="bottom-entry" @click="appStore.setSidebarView('characters')">
          <n-icon size="18">
            <PeopleOutline/>
          </n-icon>
          <span>角色管理</span>
          <n-icon size="14" class="entry-arrow">
            <ChevronForwardOutline/>
          </n-icon>
        </div>
        <div class="bottom-entry" @click="appStore.openSettings()">
          <n-icon size="18">
            <SettingsOutline/>
          </n-icon>
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
            <n-icon size="16">
              <ArrowBackOutline/>
            </n-icon>
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
            <n-icon size="16">
              <SearchOutline/>
            </n-icon>
          </template>
        </n-input>
      </div>

      <div class="character-actions">
        <n-button size="small" quaternary block @click="handleImportCharacter">
          <template #icon>
            <n-icon size="16">
              <CloudUploadOutline/>
            </n-icon>
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
                :src="getAvatarUrlSync(char.avatar)"
                :size="40"
                round
                class="character-avatar"
            />
            <div v-else class="character-avatar-default">
              <n-icon size="20">
                <PersonOutline/>
              </n-icon>
            </div>
            <div class="character-info">
              <span class="character-name">{{ char.name || '未命名角色' }}</span>
              <span class="character-desc">{{ char.description || '暂无描述' }}</span>
            </div>
          </div>

          <!-- 空状态 -->
          <div v-if="filteredCharacters.length === 0" class="empty-state">
            <n-icon size="32" color="var(--text-tertiary)">
              <PeopleOutline/>
            </n-icon>
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
// ====== 切换到角色视图时加载角色数据 ======
import {computed, h, onMounted, ref, watch} from 'vue';
import {useRouter} from 'vue-router';
import {NAvatar, NButton, NDropdown, NIcon, NInput, NScrollbar, useDialog, useMessage} from 'naive-ui';
import {
  AddOutline,
  ArrowBackOutline,
  ChatbubblesOutline,
  ChevronForwardOutline,
  CloudUploadOutline,
  EllipsisHorizontalOutline,
  PeopleOutline,
  PersonOutline,
  SearchOutline,
  SettingsOutline
} from '@vicons/ionicons5';

import {useAppStore} from '@/stores/app';
import {useChatStore} from '@/stores/chat';
import {useCharacterStore} from '@/stores/character';
import {characterClient, chatClient} from '@/api/client';
import type {ChatSession} from '@/gen/muse/chat_pb';
import type {Character} from '@/gen/muse/character_pb';
import {DEFAULT_PAGE_NUM, DEFAULT_PAGE_SIZE} from "@/utils/constants.ts";
import {getAvatarUrlSync} from '@/utils/common';

const router = useRouter();
const appStore = useAppStore();
const chatStore = useChatStore();
const characterStore = useCharacterStore();
const message = useMessage();
const dialog = useDialog();

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

// ====== 加载数据 ======
const loadSessions = async () => {
  try {
    const response = await chatClient.listChatSessions({page: DEFAULT_PAGE_NUM, pageSize: DEFAULT_PAGE_SIZE});
    chatStore.setSessions(response.sessions);
  } catch {
    // 错误由拦截器统一处理
  }
};

const loadCharacters = async () => {
  if (characterStore.hasCached) return;
  try {
    const response = await characterClient.listCharacters({page: DEFAULT_PAGE_NUM, pageSize: DEFAULT_PAGE_SIZE});
    characterStore.setCharacters(response.characters, response.total);
  } catch {
    // 错误由拦截器统一处理
  }
};

onMounted(() => {
  loadSessions();
});

// ====== 会话按时间分组 ======
interface SessionGroup {
  label: string;
  sessions: ChatSession[];
}

const groupSessionsByTime = (sessions: ChatSession[]): SessionGroup[] => {
  const todayStart = new Date();
  todayStart.setHours(0, 0, 0, 0);
  const yesterdayStart = new Date(todayStart);
  yesterdayStart.setDate(yesterdayStart.getDate() - 1);
  const weekStart = new Date(todayStart);
  weekStart.setDate(weekStart.getDate() - 7);

  const groups = new Map<string, ChatSession[]>([
    ['今天', []],
    ['昨天', []],
    ['过去 7 天', []],
    ['更早', []]
  ]);

  for (const session of sessions) {
    const updatedAt = Number(session.updatedAt ?? 0);
    const ts = updatedAt > 1e12 ? updatedAt : updatedAt * 1000;

    if (ts >= todayStart.getTime()) {
      groups.get('今天')!.push(session);
    } else if (ts >= yesterdayStart.getTime()) {
      groups.get('昨天')!.push(session);
    } else if (ts >= weekStart.getTime()) {
      groups.get('过去 7 天')!.push(session);
    } else {
      groups.get('更早')!.push(session);
    }
  }

  return Array.from(groups.entries())
      .filter(([, sessions]) => sessions.length > 0)
      .map(([label, sessions]) => ({label, sessions}));
};

// ====== 计算属性 ======
const filteredGroupedSessions = computed(() => {
  let sessions = [...chatStore.sessions.values()];
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
function handleNewChat() {
  router.push('/');
  chatStore.setActiveSession(null)
  if (appStore.isCompactMode) {
    appStore.closeLeftSidebar();
  }
}

const handleSelectSession = async (session: ChatSession) => {
  if (!session.messages?.length) {
    const sessionWithMsg = await chatClient.getChatSession({id: session.id});
    if (!sessionWithMsg) {
      console.warn("获取到的会话为null")
    } else {
      session = sessionWithMsg.session!
    }
  }
  chatStore.setActiveSession(session);
  // 更新会话访问时间
  try {
    await chatClient.updateSessionTime({sessionId: session.id})
  } catch (e) {
    console.error("更新会话访问时间失败: ", e);
  }
  await router.push('/');
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
  {label: '重命名', key: 'rename'},
  {label: '删除', key: 'delete'}
];

const handleSessionMenuSelect = (key: string) => {
  showSessionMenu.value = false;
  if (!contextSession.value) return;
  const session = contextSession.value;

  switch (key) {
    case 'rename': {
      const newName = ref(session.name || '');
      dialog.create({
        title: '重命名对话',
        content: () => {
          return h(NInput, {
            value: newName.value,
            'onUpdate:value': (v: string) => {
              newName.value = v;
            },
            placeholder: '输入新的对话名称',
            autofocus: true
          });
        },
        positiveText: '确定',
        negativeText: '取消',
        onPositiveClick: async () => {
          if (!newName.value.trim()) {
            message.warning('名称不能为空');
            return false;
          }
          try {
            const response = await chatClient.updateChatSession({
              id: session.id,
              name: newName.value.trim()
            });
            if (response.session) {
              chatStore.updateSessionInList(response.session);
            }
            message.success('已重命名');
          } catch {
            // 错误由拦截器统一处理
          }
        }
      });
      break;
    }
    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除对话"${session.name || '新对话'}"吗？`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: async () => {
          try {
            await chatClient.deleteChatSession({id: session.id});
            chatStore.removeSession(session.id);
            message.success('对话已删除');
          } catch {
            // 错误由拦截器统一处理
          }
        }
      });
      break;
  }
};

// ====== 角色操作 ======
const handleSelectCharacter = async (char: Character) => {
  // 获取角色最新的会话
  const response = await chatClient.getCharLatestSession({
    characterId: char.id,
  });
  if (response.session) {
    chatStore.addSession(response.session);
    chatStore.setActiveSession(response.session);
    try {
      await chatClient.updateSessionTime({sessionId: response.session.id})
    } catch (e) {
      console.error("更新会话访问时间失败: ", e)
    }
    // 切换回会话视图并跳转到聊天页
    appStore.setSidebarView('sessions');
    router.push('/');
  }
  if (appStore.isCompactMode) {
    appStore.closeLeftSidebar();
  }
};

const handleImportCharacter = () => {
  // 创建隐藏的文件选择器
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = '.png,.json';
  input.multiple = true;
  input.onchange = async (e) => {
    const files = (e.target as HTMLInputElement).files;
    if (!files || files.length === 0) return;

    let successCount = 0;
    let failCount = 0;

    for (const file of Array.from(files)) {
      try {
        const arrayBuffer = await file.arrayBuffer();
        const fileContent = new Uint8Array(arrayBuffer);
        const response = await characterClient.importCharacter({
          fileContent,
          fileName: file.name
        });
        if (response.character) {
          characterStore.addCharacter(response.character);
          successCount++;
        } else {
          failCount++;
        }
      } catch {
        failCount++;
      }
    }

    if (successCount > 0) {
      message.success(`成功导入 ${successCount} 个角色`);
    }
    if (failCount > 0) {
      message.warning(`${failCount} 个角色导入失败`);
    }
  };
  input.click();
};

const handleCharacterContextMenu = (e: MouseEvent, char: Character) => {
  contextCharacter.value = char;
  menuX.value = e.clientX;
  menuY.value = e.clientY;
  showCharacterMenu.value = true;
  showSessionMenu.value = false;
};

const characterMenuOptions = [
  {label: '开始聊天', key: 'chat'},
  {label: '编辑', key: 'edit'},
  {type: 'divider', key: 'd1'},
  {label: '删除', key: 'delete'}
];

const handleCharacterMenuSelect = (key: string) => {
  showCharacterMenu.value = false;
  if (!contextCharacter.value) return;
  const char = contextCharacter.value;

  switch (key) {
    case 'chat':
      handleSelectCharacter(char);
      break;
    case 'edit':
      // 跳转到角色管理页面并选择该角色
      router.push('/characters');
      break;
    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除角色"${char.name}"吗？此操作不可撤销。`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: async () => {
          try {
            await characterClient.deleteCharacter({id: char.id});
            characterStore.removeCharacter(char.id);
            message.success('角色已删除');
          } catch {
            // 错误由拦截器统一处理
          }
        }
      });
      break;
  }
};

watch(() => appStore.sidebarView, (view) => {
  if (view === 'characters') {
    loadCharacters();
  }
});
</script>

<style scoped>
.left-sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: var(--bg-secondary);
  overflow: hidden;
  position: relative;
}

/* 侧边栏科幻背景纹理 */
.left-sidebar::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--gradient-surface);
  pointer-events: none;
  z-index: 0;
}

.left-sidebar > * {
  position: relative;
  z-index: 1;
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
  border: 1px solid var(--border-light);
  border-radius: 8px;
  height: 40px;
  font-weight: 500;
  color: var(--text-primary);
  transition: all 200ms ease;
  background: var(--gradient-glow);
}

.new-chat-btn:hover {
  background: var(--bg-hover);
  border-color: var(--color-primary);
  box-shadow: var(--glow-primary-sm);
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
  border-left: 2px solid var(--color-primary);
  box-shadow: inset 0 0 20px rgba(77, 168, 255, .04);
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
  background: var(--gradient-glow);
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
  gap: 12px;
  padding: 60px 24px;
  color: var(--text-tertiary);
  font-size: 13px;
}

.empty-state .n-icon {
  opacity: .5;
  animation: glow-pulse 3s ease-in-out infinite;
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: .4;
    filter: drop-shadow(0 0 4px rgba(77, 168, 255, .1));
  }
  50% {
    opacity: .7;
    filter: drop-shadow(0 0 8px rgba(77, 168, 255, .3));
  }
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
