<template>
  <div class="chat-view">
    <!-- 左侧：会话列表 -->
    <div class="chat-sessions" :class="{ 'show-mobile': showSessionsPanel }">
      <div class="sessions-header">
        <h3>会话列表</h3>
        <n-button quaternary circle size="small" @click="createNewSession">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
        </n-button>
      </div>

      <n-scrollbar class="sessions-list">
        <div
          v-for="session in sessions"
          :key="session.id"
          class="session-item"
          :class="{ 'active': activeSessionId === session.id }"
          @click="selectSession(session.id)"
        >
          <n-avatar
            :size="40"
            round
            :src="session.character?.avatar"
            class="session-avatar"
          >
            {{ getSessionDisplayName(session).charAt(0) }}
          </n-avatar>
          <div class="session-info">
            <div class="session-name">{{ getSessionDisplayName(session) }}</div>
            <div class="session-preview">{{ getSessionPreview(session) }}</div>
          </div>
          <div class="session-meta">
            <span class="session-time">{{ formatSessionTime(session.updatedAt) }}</span>
          </div>
        </div>

        <n-empty v-if="sessions.length === 0" description="暂无会话" class="sessions-empty">
          <template #extra>
            <n-button size="small" type="primary" @click="createNewSession">
              开始新对话
            </n-button>
          </template>
        </n-empty>
      </n-scrollbar>
    </div>

    <!-- 右侧：聊天区域 -->
    <div class="chat-main">
      <!-- 聊天头部 -->
      <div class="chat-header">
        <n-button
          quaternary
          circle
          class="hide-desktop mobile-back-btn"
          @click="showSessionsPanel = !showSessionsPanel"
        >
          <template #icon>
            <n-icon><ListOutline /></n-icon>
          </template>
        </n-button>

        <div v-if="currentCharacter" class="chat-character">
          <n-avatar :size="36" round :src="currentCharacter.avatar" class="character-avatar">
            {{ currentCharacter.name.charAt(0) }}
          </n-avatar>
          <div class="character-info">
            <span class="character-name">{{ currentCharacter.name }}</span>
            <span class="character-status">
              <span class="status-dot"></span>
              在线
            </span>
          </div>
        </div>

        <div class="chat-actions">
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="small">
                <template #icon>
                  <n-icon><RefreshOutline /></n-icon>
                </template>
              </n-button>
            </template>
            重新生成
          </n-tooltip>

          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="small">
                <template #icon>
                  <n-icon><SettingsOutline /></n-icon>
                </template>
              </n-button>
            </template>
            会话设置
          </n-tooltip>

<n-dropdown :options="moreOptions" trigger="click" @select="handleMoreOptionSelect">
            <n-button quaternary circle size="small">
              <template #icon>
                <n-icon><EllipsisVerticalOutline /></n-icon>
              </template>
            </n-button>
          </n-dropdown>

          <!-- 隐藏的文件上传组件 -->
          <n-upload
            ref="uploadRef"
            :show-file-list="false"
            accept=".json"
            :custom-request="handleImportFile"
            style="display: none;"
          />
        </div>
      </div>

      <!-- 消息列表 -->
      <div class="chat-messages">
        <!-- 欢迎信息 -->
        <div v-if="messages.length === 0" class="welcome-message">
          <div class="welcome-icon">
            <n-icon size="48"><SparklesOutline /></n-icon>
          </div>
          <h2>开始新的对话</h2>
          <p>选择一个角色或开始输入来创建对话</p>
        </div>

        <!-- 虚拟滚动消息列表 -->
        <n-virtual-list
          v-else
          ref="virtualListRef"
          style="height: 100%"
          :items="messagesWithTyping"
          :item-size="120"
          item-resizable
          key-field="id"
          :padding-top="20"
          :padding-bottom="20"
          :intersection-observer-options="{ rootMargin: '100px 0px 100px 0px' }"
        >
          <template #default="{ item }">
            <div :key="item.id" class="message-item-wrapper gpu-accelerated">
            <!-- 正在输入指示器 -->
            <div v-if="item.id === 'typing-indicator'" class="typing-indicator">
              <div class="typing-dots">
                <span></span>
                <span></span>
                <span></span>
              </div>
              <span class="typing-text">{{ currentCharacter?.name }} 正在输入...</span>
            </div>
            <!-- 消息项 -->
            <MessageItem
              v-else
              :message="item"
              :character="currentCharacter"
              :persona="currentPersona"
              @edit="handleEditMessage"
              @delete="handleDeleteMessage"
              @delete-swipe="handleDeleteSwipe"
              @regenerate="handleRegenerateMessage"
              @swipe-change="handleSwipeChange"
              @branch="handleBranch"
              @duplicate-swipe="handleDuplicateSwipe"
            />
            </div>
          </template>
        </n-virtual-list>
      </div>

      <!-- 输入区域 -->
      <div class="chat-input-area">
        <MessageInput
          v-model="inputMessage"
          :disabled="isTyping"
          @send="handleSendMessage"
          @stop="handleStopGeneration"
          @persona-change="(p: any) => handlePersonaChange(p)"
        />
      </div>
    </div>

    <!-- 移动端会话面板遮罩 -->
    <Transition name="fade">
      <div
        v-if="showSessionsPanel"
        class="sessions-overlay hide-desktop"
        @click="showSessionsPanel = false"
      ></div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  NScrollbar,
  NAvatar,
  NButton,
  NIcon,
  NTooltip,
  NDropdown,
  NEmpty,
  NVirtualList,
  NUpload,
  useMessage,
  useDialog
} from 'naive-ui';
import type { UploadFileInfo } from 'naive-ui';
import type { VirtualListInst } from 'naive-ui';
import {
  AddOutline,
  ListOutline,
  RefreshOutline,
  SettingsOutline,
  EllipsisVerticalOutline,
  SparklesOutline
} from '@vicons/ionicons5';

import MessageItem from '../components/chat/MessageItem.vue';
import MessageInput from '../components/chat/MessageInput.vue';
import { useChatStore } from '@/stores/chat';
import { useUserStore } from '@/stores/user';
import { chatClient } from '@/api/client';
import type { ChatSession, Message as PbMessage, Character, Persona } from '@/gen/muse/muse_pb';

// 本地Message类型适配
interface LocalMessage {
  id: string;
  role: 'user' | 'assistant' | 'system';
  swipes: Array<{
    id: string;
    content: string;
    timestamp: number;
  }>;
  currentSwipeIndex: number;
}

const route = useRoute();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const chatStore = useChatStore();
const userStore = useUserStore();

// 响应式状态
const showSessionsPanel = ref(false);
const inputMessage = ref('');
const virtualListRef = ref<VirtualListInst | null>(null);
const uploadRef = ref<InstanceType<typeof NUpload> | null>(null);
const loading = ref(false);

// 从Store获取数据
const sessions = computed(() => chatStore.sessions);
const activeSessionId = computed(() => chatStore.activeSessionId);
const isTyping = computed(() => chatStore.isStreaming);

// 当前角色
const currentCharacter = computed<Character | null>(() => {
  return chatStore.activeCharacter ?? null;
});

// 当前人设
const currentPersona = computed<Persona | null>(() => {
  return userStore.activePersona ?? null;
});

// 将PbMessage转换为本地格式
const convertToLocalMessage = (msg: PbMessage): LocalMessage => {
  const roleMap: Record<number, 'user' | 'assistant' | 'system'> = {
    1: 'system',
    2: 'user',
    3: 'assistant'
  };
  return {
    id: msg.id.toString(),
    role: roleMap[msg.role] || 'assistant',
    swipes: msg.swipes.map(s => ({
      id: s.id.toString(),
      content: s.content,
      timestamp: Number(s.createdAt)
    })),
    currentSwipeIndex: msg.activeSwipeIndex
  };
};

// 会话辅助函数
const getSessionDisplayName = (session: ChatSession): string => {
  return session.name || session.character?.name || '未命名会话';
};

const getSessionPreview = (session: ChatSession): string => {
  if (session.messages && session.messages.length > 0) {
    const lastMsg = session.messages[session.messages.length - 1];
    if (lastMsg && lastMsg.swipes && lastMsg.swipes.length > 0) {
      const swipe = lastMsg.swipes[lastMsg.activeSwipeIndex];
      const content = swipe?.content || '';
      return content.length > 50 ? content.slice(0, 50) + '...' : content;
    }
  }
  return '暂无消息';
};

const formatSessionTime = (timestamp: bigint): string => {
  const date = new Date(Number(timestamp));
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`;
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`;

  return date.toLocaleDateString();
};

// 消息列表（转换为本地格式）
const messages = computed<LocalMessage[]>(() => {
  return chatStore.messages.map(convertToLocalMessage);
});

// 计算属性：消息列表（包含输入指示器）
const messagesWithTyping = computed(() => {
  const items = [...messages.value];
  if (isTyping.value) {
    items.push({
      id: 'typing-indicator',
      role: 'assistant',
      swipes: [],
      currentSwipeIndex: 0
    });
  }
  return items;
});

// 加载会话列表
const loadSessions = async (characterId?: number) => {
  loading.value = true;
  const response = await chatClient.listChatSessions({
    characterId: characterId,
    page: 1,
    pageSize: 50
  });
  chatStore.setSessions(response.sessions);
  loading.value = false;
};

// 加载会话详情（包含消息）
const loadSession = async (sessionId: number) => {
  loading.value = true;
  const response = await chatClient.getChatSession({
    id: sessionId,
    includeMessages: true
  });
  if (response.session) {
    chatStore.setActiveSession(response.session);
  }
  loading.value = false;
};

// 初始化
onMounted(async () => {
  await loadSessions();

  // 如果URL中有sessionId参数，加载该会话
  const sessionId = route.params.sessionId;
  if (sessionId) {
    await loadSession(Number(sessionId));
  }
});

// 监听路由参数变化
watch(() => route.params.sessionId, async (newId) => {
  if (newId) {
    await loadSession(Number(newId));
  } else {
    chatStore.setActiveSession(null);
  }
});

// 滚动到底部
const scrollToBottom = () => {
  if (virtualListRef.value) {
    virtualListRef.value.scrollTo({ position: 'bottom', behavior: 'smooth' });
  }
};

// 更多操作选项
const moreOptions = [
  { label: '导入聊天记录', key: 'import', icon: () => null },
  { label: '导出会话', key: 'export', icon: () => null },
  { label: '分享', key: 'share', icon: () => null },
  { type: 'divider', key: 'd1' },
  { label: '清空会话', key: 'clear', icon: () => null },
  { label: '删除会话', key: 'delete', icon: () => null }
];

// 处理更多选项点击
const handleMoreOptionSelect = (key: string) => {
  switch (key) {
    case 'import':
      triggerImportChat();
      break;
    case 'export':
      handleExportChat();
      break;
    case 'clear':
      handleClearChat();
      break;
    case 'delete':
      handleDeleteChat();
      break;
  }
};

// 触发导入聊天记录
const triggerImportChat = () => {
  uploadRef.value?.openOpenFileDialog();
};

// 处理文件上传
const handleImportFile = (options: { file: UploadFileInfo }) => {
  const { file } = options;
  if (!file.file) return false;

  const reader = new FileReader();
reader.onload = (e) => {
    try {
      const content = e.target?.result as string;
      JSON.parse(content); // 验证JSON格式
      // TODO: 调用后端导入接口
      message.info('导入功能暂未实现');
    } catch (err) {
      message.error('解析文件失败，请确保是有效的JSON文件');
    }
  };
  reader.readAsText(file.file);
  return false;
};

// 导出聊天记录
const handleExportChat = async () => {
  if (!chatStore.activeSessionId) {
    message.warning('请先选择一个会话');
    return;
  }

  // TODO: 后端实现导出接口后启用
  message.info('导出功能暂未实现');
};

// 清空会话
const handleClearChat = () => {
  dialog.warning({
    title: '确认清空',
    content: '确定要清空当前会话的所有消息吗？此操作不可撤销。',
    positiveText: '清空',
    negativeText: '取消',
    onPositiveClick: async () => {
      if (!chatStore.activeSessionId) return;
      // TODO: 后端实现清空接口后启用
      message.info('清空功能暂未实现');
    }
  });
};

// 删除会话
const handleDeleteChat = () => {
  if (!chatStore.activeSession) {
    message.warning('请先选择一个会话');
    return;
  }

  dialog.warning({
    title: '确认删除',
    content: '确定要删除当前会话吗？此操作不可撤销。',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      if (!chatStore.activeSessionId) return;
      try {
        await chatClient.deleteChatSession({ id: chatStore.activeSessionId });
        chatStore.removeSession(chatStore.activeSessionId);
        router.push('/chat');
        message.success('会话已删除');
      } catch (e) {
        message.error('删除失败');
      }
    }
  });
};

// 创建新会话
const createNewSession = () => {
  // 跳转到角色选择页面或显示角色选择弹窗
  router.push('/characters');
};

// 选择会话
const selectSession = async (id: number) => {
  router.push(`/chat/${id}`);
  showSessionsPanel.value = false;
};

// 发送消息
const handleSendMessage = async (content: string) => {
  if (!content.trim() || !chatStore.activeSessionId) return;

  // 创建临时用户消息
  const tempUserMsg = chatStore.createTempUserMessage(content);
  chatStore.addMessage(tempUserMsg);
  inputMessage.value = '';

  // 创建临时AI消息占位
  const tempAiMsg = chatStore.createTempAiMessage();
  chatStore.addMessage(tempAiMsg);

  // 滚动到底部
  setTimeout(scrollToBottom, 50);

  // 开始流式生成
  chatStore.startStreaming();
  const signal = chatStore.createAbortController();

  try {
    const stream = chatClient.sendMessage(
      { sessionId: chatStore.activeSessionId, content: content },
      { signal: signal }
    );

    // 存储每个候选回复的内容
    const swipeContents: Map<number, string> = new Map();

    for await (const response of stream) {
      const index = response.index;

      // 累加内容
      const currentContent = swipeContents.get(index) || '';
      const newContent = currentContent + response.content;
      swipeContents.set(index, newContent);

      // 更新临时AI消息的swipes
      const aiMsgIndex = chatStore.messages.findIndex(m => m.id === tempAiMsg.id);
      if (aiMsgIndex >= 0) {
        const aiMsg = chatStore.messages[aiMsgIndex];
        if (aiMsg) {
          const swipe = aiMsg.swipes.find(s => s.sortOrder === index);
          if (swipe) {
            swipe.content = newContent;
          } else {
            aiMsg.swipes.push({
              id: tempAiMsg.id + index,
              messageId: tempAiMsg.id,
              content: newContent,
              sortOrder: index,
              createdAt: BigInt(Date.now()),
              $typeName: 'muse.MessageSwipe'
            } as any);
          }
        }
      }
    }

    // 流式完成后，重新加载会话获取真实的消息ID
    await loadSession(chatStore.activeSessionId);

  } catch (error) {
    if ((error as Error).name === 'AbortError') {
      message.info('消息生成已取消');
    } else {
      message.error('发送消息失败');
      console.error('发送消息失败:', error);
    }
  } finally {
    chatStore.stopStreaming();
    setTimeout(scrollToBottom, 50);
  }
};

// 停止生成
const handleStopGeneration = () => {
  chatStore.abortGeneration();
};

// 编辑消息内容
const handleEditMessage = async (messageId: number, swipeId: number, content: string) => {
  try {
    await chatClient.editMessage({
      messageId: messageId,
      swipeId: swipeId,
      content: content
    });

    // 更新本地状态
    const msg = chatStore.messages.find(m => m.id === messageId);
    if (msg) {
      const swipe = msg.swipes.find(s => s.id === swipeId);
      if (swipe) {
        swipe.content = content;
      }
    }
    message.success('消息已更新');
  } catch (e) {
    message.error('编辑失败');
  }
};

// 删除整个楼层
const handleDeleteMessage = async (id: number) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这条消息吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await chatClient.deleteMessage({ messageId: id });
        chatStore.removeMessage(id);
        message.success('消息已删除');
      } catch (e) {
        message.error('删除失败');
      }
    }
  });
};

// 删除单条消息（swipe）
const handleDeleteSwipe = async (messageId: number, swipeId: number) => {
  const msg = chatStore.messages.find(m => m.id === messageId);
  if (!msg || msg.swipes.length <= 1) {
    message.warning('无法删除最后一条回复');
    return;
  }

  // TODO: 后端实现删除swipe接口后启用
  // 本地删除
  const swipeIndex = msg.swipes.findIndex(s => s.id === swipeId);
  if (swipeIndex >= 0) {
    msg.swipes.splice(swipeIndex, 1);
    if (msg.activeSwipeIndex >= msg.swipes.length) {
      msg.activeSwipeIndex = msg.swipes.length - 1;
    }
  }
  message.success('回复已删除');
};

// 重新生成消息（添加新的 swipe）
const handleRegenerateMessage = async (id: number) => {
  const msg = chatStore.messages.find(m => m.id === id);
  if (!msg || msg.role !== 3) return; // 只能重新生成助手消息

  chatStore.startStreaming();
  const signal = chatStore.createAbortController();

  try {
    const stream = chatClient.regenerateMessage(
      { messageId: id },
      { signal: signal }
    );

    let newSwipeContent = '';

    for await (const response of stream) {
      if (response.newSwipe) {
        // 添加新的swipe
        msg.swipes.push(response.newSwipe);
        msg.activeSwipeIndex = msg.swipes.length - 1;
      }
      if (response.contentDelta) {
        newSwipeContent += response.contentDelta;
        // 更新最后一个swipe的内容
        const lastSwipe = msg.swipes[msg.swipes.length - 1];
        if (lastSwipe) {
          lastSwipe.content = newSwipeContent;
        }
      }
    }

    message.success('重新生成完成');
  } catch (error) {
    if ((error as Error).name === 'AbortError') {
      message.info('重新生成已取消');
    } else {
      message.error('重新生成失败');
    }
  } finally {
    chatStore.stopStreaming();
  }
};

// 切换 swipe
const handleSwipeChange = async (messageId: number, index: number) => {
  try {
    await chatClient.switchSwipe({
      messageId: messageId,
      swipeIndex: index
    });
    chatStore.localSwitchSwipe(messageId, index);
  } catch (e) {
    // 即使API失败，也更新本地状态
    chatStore.localSwitchSwipe(messageId, index);
  }
};

// 创建分支
const handleBranch = async (_messageId: number) => {
  // TODO: 后端实现分支接口后启用
  message.info('分支功能暂未实现');
};

// 复制当前消息为新版本（用户消息专用）
const handleDuplicateSwipe = async (messageId: number) => {
  const msg = chatStore.messages.find(m => m.id === messageId);
  if (!msg) return;

  const currentSwipe = msg.swipes[msg.activeSwipeIndex];
  if (!currentSwipe) return;

  // TODO: 后端实现复制swipe接口后启用
  // 本地复制
  const newSwipe = {
    ...currentSwipe,
    id: Date.now(),
    createdAt: BigInt(Date.now())
  };
  msg.swipes.push(newSwipe as any);
  msg.activeSwipeIndex = msg.swipes.length - 1;
  message.success('已创建副本');
};

// 处理人设切换
const handlePersonaChange = async (persona: { id: number; name: string; avatar: string }) => {
  try {
    await userStore.setActivePersonaId(persona.id);
  } catch (e) {
    message.error('切换人设失败');
  }
};
</script>

<style scoped>
.chat-view {
  display: flex;
  height: calc(100vh - 64px - 48px);
  background: var(--bg-primary);
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--border-color);
}

/* 会话列表 */
.chat-sessions {
  width: 300px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .chat-sessions {
    position: fixed;
    left: 0;
    top: 64px;
    bottom: 0;
    width: 300px;
    z-index: 100;
    transform: translateX(-100%);
    transition: transform var(--transition-normal);
  }

  .chat-sessions.show-mobile {
    transform: translateX(0);
    box-shadow: var(--shadow-xl);
  }
}

.sessions-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.sessions-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.sessions-list {
  flex: 1;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: all var(--transition-fast);
  border-bottom: 1px solid var(--border-color);
  transform: translateZ(0);
}

.session-item:hover {
  background: var(--bg-card-hover);
  transform: translateX(2px);
}

.session-item.active {
  background: rgba(0, 240, 255, 0.1);
  border-left: 3px solid var(--color-primary);
}

.session-avatar {
  flex-shrink: 0;
  background: var(--gradient-primary);
}

.session-info {
  flex: 1;
  min-width: 0;
}

.session-name {
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.session-preview {
  font-size: 13px;
  color: var(--text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 4px;
}

.session-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.sessions-empty {
  padding: 40px 20px;
}

/* 聊天主区域 */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  transition: transform var(--transition-normal);
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  transition: background var(--transition-fast);
}

.mobile-back-btn {
  margin-right: 4px;
}

.chat-character {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.character-avatar {
  background: var(--gradient-primary);
  box-shadow: var(--glow-soft);
}

.character-info {
  display: flex;
  flex-direction: column;
}

.character-name {
  font-weight: 600;
  color: var(--text-primary);
}

.character-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--color-success);
}

.status-dot {
  width: 8px;
  height: 8px;
  background: var(--color-success);
  border-radius: 50%;
  animation: pulse 2s ease-in-out infinite;
}

.chat-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 消息区域 */
.chat-messages {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

/* 消息项包装器 */
.message-item-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
  animation: fade-in 0.3s ease-out;
}

.welcome-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 60px 20px;
  color: var(--text-tertiary);
}

.welcome-icon {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border-radius: 20px;
  margin-bottom: 20px;
  color: var(--color-primary);
}

.welcome-message h2 {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.welcome-message p {
  margin: 0;
  font-size: 14px;
}

/* 输入指示器 */
.typing-indicator {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 0;
}

.typing-dots {
  display: flex;
  gap: 4px;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  background: var(--color-primary);
  border-radius: 50%;
  animation: typing-bounce 1.4s ease-in-out infinite;
}

.typing-dots span:nth-child(2) {
  animation-delay: 0.2s;
}

.typing-dots span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes typing-bounce {
  0%, 60%, 100% {
    transform: translateY(0);
    opacity: 0.4;
  }
  30% {
    transform: translateY(-8px);
    opacity: 1;
  }
}

.typing-text {
  font-size: 13px;
  color: var(--text-tertiary);
}

/* 输入区域 */
.chat-input-area {
  padding: 16px 20px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
}

/* 会话面板遮罩 */
.sessions-overlay {
  position: fixed;
  inset: 0;
  background: var(--bg-overlay);
  z-index: 99;
  backdrop-filter: blur(4px);
  transition: opacity var(--transition-fast);
}

/* 移动端优化 */
@media (max-width: 768px) {
  .chat-sessions {
    transition: transform var(--transition-mobile);
    will-change: transform;
  }
  
  .chat-sessions.show-mobile {
    transform: translateX(0);
  }
  
  .session-item:active {
    transform: scale(0.98);
    background: var(--bg-card-hover);
  }
  
  .chat-header {
    padding: 12px 16px;
  }
  
  .chat-input-area {
    padding: 12px 16px;
  }
  
  .message-item-wrapper {
    padding: 0 12px;
  }
}
</style>
