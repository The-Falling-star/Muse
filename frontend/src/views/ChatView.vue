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
            :src="session.avatar"
            class="session-avatar"
          >
            {{ session.characterName.charAt(0) }}
          </n-avatar>
          <div class="session-info">
            <div class="session-name">{{ session.characterName }}</div>
            <div class="session-preview">{{ session.lastMessage }}</div>
          </div>
          <div class="session-meta">
            <span class="session-time">{{ session.time }}</span>
            <n-badge v-if="session.unread" :value="session.unread" :max="99" />
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
        >
          <template #default="{ item }">
            <div :key="item.id" class="message-item-wrapper">
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
          @persona-change="handlePersonaChange"
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
import { ref, computed } from 'vue';
import {
  NScrollbar,
  NAvatar,
  NBadge,
  NButton,
  NIcon,
  NTooltip,
  NDropdown,
  NEmpty,
  NVirtualList,
  NUpload,
  useMessage
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
import type { Message, Character } from '../types';

// Persona 类型定义
interface Persona {
  id: string;
  name: string;
  avatar: string;
}

// 响应式状态
const showSessionsPanel = ref(false);
const activeSessionId = ref<string | null>('1');
const inputMessage = ref('');
const isTyping = ref(false);
const virtualListRef = ref<VirtualListInst | null>(null);
const uploadRef = ref<InstanceType<typeof NUpload> | null>(null);
const message = useMessage();

// 当前人设
const currentPersona = ref<Persona>({ id: '1', name: '默认用户', avatar: '' });

// 模拟会话数据
const sessions = ref([
  {
    id: '1',
    characterName: 'AI助手',
    avatar: '',
    lastMessage: '你好！有什么我可以帮助你的吗？',
    time: '刚刚',
    unread: 0
  },
  {
    id: '2',
    characterName: '小说角色',
    avatar: '',
    lastMessage: '这是一段很长的消息预览文本...',
    time: '5分钟前',
    unread: 2
  }
]);

// 当前角色
const currentCharacter = computed<Character | null>(() => {
  if (!activeSessionId.value) return null;
  const session = sessions.value.find(s => s.id === activeSessionId.value);
  if (!session) return null;
  return {
    id: session.id,
    name: session.characterName,
    avatar: session.avatar
  };
});

// 模拟消息数据（使用新的 swipes 结构）
const messages = ref<Message[]>([
  {
    id: '1',
    role: 'assistant',
    swipes: [
      {
        id: '1-1',
        content: '你好！我是Muse AI助手。我可以帮助你进行创意写作、角色扮演或者回答问题。有什么我可以帮你的吗？',
        timestamp: Date.now() - 60000
      }
    ],
    currentSwipeIndex: 0
  },
  {
    id: '2',
    role: 'user',
    swipes: [
      {
        id: '2-1',
        content: '你好，我想了解一下你有哪些功能？',
        timestamp: Date.now() - 30000
      }
    ],
    currentSwipeIndex: 0
  },
  {
    id: '3',
    role: 'assistant',
    swipes: [
      {
        id: '3-1',
        content: `当然！以下是我的主要功能：

1. **创意写作** - 帮助你创作故事、诗歌、剧本等
2. **角色扮演** - 可以扮演各种角色与你互动
3. **问题解答** - 回答各类问题
4. **文本编辑** - 帮助修改和润色文本

你对哪个功能最感兴趣呢？`,
        timestamp: Date.now() - 20000
      },
      {
        id: '3-2',
        content: `我有很多功能可以帮助你：

• **故事创作** - 我可以和你一起创作有趣的故事
• **角色扮演** - 扮演各种角色与你对话
• **写作润色** - 帮你修改和优化文本
• **知识问答** - 回答你的各种问题

想先试试哪个？`,
        timestamp: Date.now() - 10000
      },
      {
        id: '3-3',
        content: `很高兴你问这个！我的主要能力包括：

1. 📝 **创意写作** - 故事、诗歌、剧本
2. 🎭 **角色扮演** - 沉浸式互动体验
3. 💬 **智能对话** - 问答、讨论、建议
4. ✨ **文本优化** - 润色、改写、翻译

你最感兴趣的是哪个方向？`,
        timestamp: Date.now()
      }
    ],
    currentSwipeIndex: 2
  }
]);

// 计算属性：消息列表（包含输入指示器）
const messagesWithTyping = computed(() => {
  const items = [...messages.value];
  // 如果正在输入，添加一个特殊的输入指示器项
  if (isTyping.value) {
    items.push({
      id: 'typing-indicator',
      role: 'assistant',
      swipes: [],
      currentSwipeIndex: 0
    } as Message);
  }
  return items;
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
      const importedData = JSON.parse(content);

      // 验证导入的数据格式
      if (importedData.messages && Array.isArray(importedData.messages)) {
        // 合并或替换消息
        messages.value = importedData.messages;
        message.success('聊天记录导入成功');
        setTimeout(scrollToBottom, 50);
      } else if (Array.isArray(importedData)) {
        // 直接是消息数组
        messages.value = importedData;
        message.success('聊天记录导入成功');
        setTimeout(scrollToBottom, 50);
      } else {
        message.error('无效的聊天记录格式');
      }
    } catch (err) {
      message.error('解析文件失败，请确保是有效的JSON文件');
    }
  };
  reader.readAsText(file.file);
  return false; // 阻止默认上传行为
};

// 导出聊天记录
const handleExportChat = () => {
  const exportData = {
    sessionId: activeSessionId.value,
    character: currentCharacter.value,
    messages: messages.value,
    exportTime: new Date().toISOString()
  };

  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `chat-${currentCharacter.value?.name || 'export'}-${Date.now()}.json`;
  a.click();
  URL.revokeObjectURL(url);
  message.success('聊天记录已导出');
};

// 清空会话
const handleClearChat = () => {
  messages.value = [];
  message.success('会话已清空');
};

// 删除会话
const handleDeleteChat = () => {
  console.log('Delete chat');
};

// 方法
const createNewSession = () => {
  console.log('Create new session');
};

const selectSession = (id: string) => {
  activeSessionId.value = id;
  showSessionsPanel.value = false;
};

const handleSendMessage = (content: string) => {
  if (!content.trim()) return;

  const now = Date.now();

  // 添加用户消息（使用 swipes 结构）
  messages.value.push({
    id: now.toString(),
    role: 'user',
    swipes: [
      {
        id: `${now}-1`,
        content,
        timestamp: now
      }
    ],
    currentSwipeIndex: 0
  });

  inputMessage.value = '';
  isTyping.value = true;

  // 滚动到底部
  setTimeout(scrollToBottom, 50);

  // 模拟AI回复
  setTimeout(() => {
    const aiNow = Date.now();
    messages.value.push({
      id: aiNow.toString(),
      role: 'assistant',
      swipes: [
        {
          id: `${aiNow}-1`,
          content: '这是一条模拟的AI回复消息。在实际应用中，这里会调用API获取真实的AI响应。',
          timestamp: aiNow
        }
      ],
      currentSwipeIndex: 0
    });
    isTyping.value = false;
    // 滚动到底部
    setTimeout(scrollToBottom, 50);
  }, 2000);
};

const handleStopGeneration = () => {
  isTyping.value = false;
};

// 编辑消息内容
const handleEditMessage = (messageId: string, swipeId: string, content: string) => {
  const message = messages.value.find(m => m.id === messageId);
  if (!message) return;

  const swipe = message.swipes.find(s => s.id === swipeId);
  if (swipe) {
    swipe.content = content;
    swipe.timestamp = Date.now(); // 更新时间戳
  }
};

// 删除整个楼层
const handleDeleteMessage = (id: string) => {
  messages.value = messages.value.filter(m => m.id !== id);
};

// 删除单条消息（swipe）
const handleDeleteSwipe = (messageId: string, swipeId: string) => {
  const message = messages.value.find(m => m.id === messageId);
  if (!message || message.swipes.length <= 1) return;

  const swipeIndex = message.swipes.findIndex(s => s.id === swipeId);
  if (swipeIndex === -1) return;

  // 删除这条 swipe
  message.swipes.splice(swipeIndex, 1);

  // 调整当前索引
  if (message.currentSwipeIndex >= message.swipes.length) {
    message.currentSwipeIndex = message.swipes.length - 1;
  } else if (message.currentSwipeIndex > swipeIndex) {
    message.currentSwipeIndex--;
  }
};

// 重新生成消息（添加新的 swipe）
const handleRegenerateMessage = (id: string) => {
  const message = messages.value.find(m => m.id === id);
  if (!message || message.role !== 'assistant') return;

  isTyping.value = true;

  // 模拟生成新的回复
  setTimeout(() => {
    const now = Date.now();
    message.swipes.push({
      id: `${id}-${now}`,
      content: `这是重新生成的回复 #${message.swipes.length + 1}。每次点击重新生成都会添加一条新消息到这个楼层。`,
      timestamp: now
    });
    // 自动切换到新生成的消息
    message.currentSwipeIndex = message.swipes.length - 1;
    isTyping.value = false;
  }, 1500);
};

// 切换 swipe
const handleSwipeChange = (messageId: string, index: number) => {
  const message = messages.value.find(m => m.id === messageId);
  if (message) {
    message.currentSwipeIndex = index;
  }
};

// 创建分支（暂不实现业务逻辑，只打印日志）
const handleBranch = (messageId: string) => {
  console.log('Create branch from message:', messageId);
  // TODO: 实现分支逻辑
  // 1. 复制当前会话到新会话
  // 2. 新会话删除该消息之后的所有楼层
  // 3. 原会话保持不变
};

// 复制当前消息为新版本（用户消息专用）
const handleDuplicateSwipe = (messageId: string) => {
  const message = messages.value.find(m => m.id === messageId);
  if (!message) return;

  // 获取当前显示的 swipe 内容
  const currentSwipe = message.swipes[message.currentSwipeIndex];
  if (!currentSwipe) return;

  const now = Date.now();
  // 复制当前消息内容为新版本
  message.swipes.push({
    id: `${messageId}-${now}`,
    content: currentSwipe.content,
    timestamp: now
  });

  // 自动切换到新版本
  message.currentSwipeIndex = message.swipes.length - 1;
};

// 处理人设切换
const handlePersonaChange = (persona: Persona) => {
  currentPersona.value = persona;
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
  transition: background var(--transition-fast);
  border-bottom: 1px solid var(--border-color);
}

.session-item:hover {
  background: var(--bg-card-hover);
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
}

.chat-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
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
}

/* 消息项包装器 */
.message-item-wrapper {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 20px;
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
}

/* 消息列表动画 */
.message-list-enter-active,
.message-list-leave-active {
  transition: all 0.3s ease;
}

.message-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.message-list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
