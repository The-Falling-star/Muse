<template>
  <div class="chat-view">
    <!-- 欢迎界面（无消息时显示） -->
    <div v-if="messages.length === 0" class="welcome-screen">
      <div class="welcome-content">
        <div class="welcome-logo">
          <n-icon size="56" color="var(--color-primary)">
            <SparklesOutline />
          </n-icon>
        </div>
        <h1 class="welcome-title">Muse</h1>
        <p class="welcome-subtitle">开始你的创意之旅</p>

        <div class="quick-prompts">
          <div
            v-for="(prompt, index) in quickPrompts"
            :key="index"
            class="quick-prompt-card"
            @click="handleQuickPrompt(prompt.text)"
          >
            <n-icon size="20" class="prompt-icon">
              <component :is="prompt.icon" />
            </n-icon>
            <span class="prompt-text">{{ prompt.text }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 消息列表 -->
    <div v-else class="messages-area">
      <!-- 角色信息条（当有角色时显示） -->
      <div v-if="currentCharacter" class="character-banner">
        <n-avatar :size="28" round :src="currentCharacter.avatar" class="character-avatar">
          {{ currentCharacter.name?.charAt(0) }}
        </n-avatar>
        <span class="character-name">{{ currentCharacter.name }}</span>
      </div>

      <!-- 虚拟滚动消息列表 -->
      <n-virtual-list
        ref="virtualListRef"
        class="message-list"
        :items="messagesWithTyping"
        :item-size="120"
        item-resizable
        key-field="id"
        :padding-top="20"
        :padding-bottom="20"
        :intersection-observer-options="{ rootMargin: '100px 0px 100px 0px' }"
      >
        <template #default="{ item }">
          <div :key="item.id" class="message-item-wrapper">
            <!-- 正在输入指示器 -->
            <div v-if="item.id === 'typing-indicator'" class="typing-indicator">
              <div class="typing-avatar">
                <n-avatar :size="32" round :src="currentCharacter?.avatar">
                  {{ currentCharacter?.name?.charAt(0) || '?' }}
                </n-avatar>
              </div>
              <div class="typing-bubble">
                <div class="typing-dots">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
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

      <!-- 跳到最新按钮 -->
      <Transition name="fade">
        <n-button
          v-if="showScrollToBottom"
          class="scroll-to-bottom-btn"
          circle
          secondary
          size="small"
          @click="scrollToBottom"
        >
          <template #icon>
            <n-icon><ChevronDownOutline /></n-icon>
          </template>
        </n-button>
      </Transition>
    </div>

    <!-- 输入区域 -->
    <div class="input-area">
      <div class="input-container">
        <MessageInput
          v-model="inputMessage"
          :disabled="isTyping"
          @send="handleSendMessage"
          @stop="handleStopGeneration"
          @persona-change="(p: any) => handlePersonaChange(p)"
        />
        <div class="input-footer">
          <span class="model-info">{{ currentModelName }}</span>
        </div>
      </div>
    </div>

    <!-- 隐藏的文件上传组件 -->
    <n-upload
      ref="uploadRef"
      :show-file-list="false"
      accept=".json"
      :custom-request="handleImportFile"
      style="display: none;"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick } from 'vue';
import { useRoute } from 'vue-router';
import {
  NAvatar,
  NButton,
  NIcon,
  NVirtualList,
  NUpload,
  useMessage,
  useDialog
} from 'naive-ui';
import type { UploadFileInfo } from 'naive-ui';
import type { VirtualListInst } from 'naive-ui';
import {
  SparklesOutline,
  ChevronDownOutline,
  ChatbubbleEllipsesOutline,
  BookOutline,
  ColorPaletteOutline,
  CodeSlashOutline
} from '@vicons/ionicons5';

import MessageItem from '@/components/chat/MessageItem.vue';
import MessageInput from '@/components/chat/MessageInput.vue';
import { useChatStore } from '@/stores/chat';
import { useUserStore } from '@/stores/user';
import { chatClient } from '@/api/client';
import type { Character, Persona } from '@/gen/muse/muse_pb';

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
const messageApi = useMessage();
const dialog = useDialog();
const chatStore = useChatStore();
const userStore = useUserStore();

// 响应式状态
const inputMessage = ref('');
const virtualListRef = ref<VirtualListInst | null>(null);
const loading = ref(false);
const showScrollToBottom = ref(false);

// 快捷提示
const quickPrompts = [
  { text: '帮我写一段创意故事', icon: ChatbubbleEllipsesOutline },
  { text: '解释一个复杂的概念', icon: BookOutline },
  { text: '帮我头脑风暴创意', icon: ColorPaletteOutline },
  { text: '写一段代码片段', icon: CodeSlashOutline }
];

// 从 Store 获取数据
const isTyping = computed(() => chatStore.isStreaming);

// 当前模型名称
const currentModelName = computed(() => {
  // TODO: 从模型配置中获取
  return 'Muse AI';
});

// 当前角色
const currentCharacter = computed<Character | null>(() => {
  return chatStore.activeCharacter ?? null;
});

// 当前人设
const currentPersona = computed<Persona | null>(() => {
  return userStore.activePersona ?? null;
});

// 将PbMessage转换为本地格式
const convertToLocalMessage = (msg: any): LocalMessage => {
  const roleMap: Record<number, 'user' | 'assistant' | 'system'> = {
    1: 'system',
    2: 'user',
    3: 'assistant'
  };
  return {
    id: msg.id,
    role: roleMap[msg.role] || 'assistant',
    swipes: msg.swipes.map((s: any) => ({
      id: s.id,
      content: s.content,
      timestamp: Number(s.createdAt)
    })),
    currentSwipeIndex: msg.activeSwipeIndex
  };
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

// 加载会话详情（包含消息）
const loadSession = async (sessionId: number) => {
  loading.value = true;
  try {
    const response = await chatClient.getChatSession({
      id: sessionId,
      includeMessages: true
    });
    if (response.session) {
      chatStore.setActiveSession(response.session);
    }
  } finally {
    loading.value = false;
  }
};

// 初始化
onMounted(async () => {
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
  showScrollToBottom.value = false;
};

// 快捷提示点击
const handleQuickPrompt = (text: string) => {
  inputMessage.value = text;
};

// 处理文件上传
const handleImportFile = (options: { file: UploadFileInfo }) => {
  const { file } = options;
  if (!file.file) return false;

  const reader = new FileReader();
  reader.onload = (e) => {
    try {
      const content = e.target?.result as string;
      JSON.parse(content);
      messageApi.info('导入功能暂未实现');
    } catch {
      messageApi.error('解析文件失败，请确保是有效的JSON文件');
    }
  };
  reader.readAsText(file.file);
  return false;
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
  await nextTick();
  scrollToBottom();

  // 开始流式生成
  chatStore.startStreaming();
  const signal = chatStore.createAbortController();

  try {
    const stream = chatClient.sendMessage(
      { sessionId: chatStore.activeSessionId, content: content },
      { signal: signal }
    );

    const swipeContents: Map<number, string> = new Map();

    for await (const response of stream) {
      const index = response.index;
      const currentContent = swipeContents.get(index) || '';
      const newContent = currentContent + response.content;
      swipeContents.set(index, newContent);

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

    await loadSession(chatStore.activeSessionId);
  } catch (error) {
    if ((error as Error).name === 'AbortError') {
      messageApi.info('消息生成已取消');
    } else {
      messageApi.error('发送消息失败');
      console.error('发送消息失败:', error);
    }
  } finally {
    chatStore.stopStreaming();
    await nextTick();
    scrollToBottom();
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

    const msg = chatStore.messages.find(m => m.id === messageId);
    if (msg) {
      const swipe = msg.swipes.find(s => s.id === swipeId);
      if (swipe) {
        swipe.content = content;
      }
    }
    messageApi.success('消息已更新');
  } catch {
    messageApi.error('编辑失败');
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
        messageApi.success('消息已删除');
      } catch {
        messageApi.error('删除失败');
      }
    }
  });
};

// 删除单条消息（swipe）
const handleDeleteSwipe = async (messageId: number, swipeId: number) => {
  const msg = chatStore.messages.find(m => m.id === messageId);
  if (!msg || msg.swipes.length <= 1) {
    messageApi.warning('无法删除最后一条回复');
    return;
  }

  const swipeIndex = msg.swipes.findIndex(s => s.id === swipeId);
  if (swipeIndex >= 0) {
    msg.swipes.splice(swipeIndex, 1);
    if (msg.activeSwipeIndex >= msg.swipes.length) {
      msg.activeSwipeIndex = msg.swipes.length - 1;
    }
  }
  messageApi.success('回复已删除');
};

// 重新生成消息
const handleRegenerateMessage = async (id: number) => {
  const msg = chatStore.messages.find(m => m.id === id);
  if (!msg || msg.role !== 3) return;

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
        msg.swipes.push(response.newSwipe);
        msg.activeSwipeIndex = msg.swipes.length - 1;
      }
      if (response.contentDelta) {
        newSwipeContent += response.contentDelta;
        const lastSwipe = msg.swipes[msg.swipes.length - 1];
        if (lastSwipe) {
          lastSwipe.content = newSwipeContent;
        }
      }
    }

    messageApi.success('重新生成完成');
  } catch (error) {
    if ((error as Error).name === 'AbortError') {
      messageApi.info('重新生成已取消');
    } else {
      messageApi.error('重新生成失败');
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
  } catch {
    chatStore.localSwitchSwipe(messageId, index);
  }
};

// 创建分支
const handleBranch = async (_messageId: number) => {
  messageApi.info('分支功能暂未实现');
};

// 复制当前消息为新版本
const handleDuplicateSwipe = async (messageId: number) => {
  const msg = chatStore.messages.find(m => m.id === messageId);
  if (!msg) return;

  const currentSwipe = msg.swipes[msg.activeSwipeIndex];
  if (!currentSwipe) return;

  const newSwipe = {
    ...currentSwipe,
    id: Date.now(),
    createdAt: BigInt(Date.now())
  };
  msg.swipes.push(newSwipe as any);
  msg.activeSwipeIndex = msg.swipes.length - 1;
  messageApi.success('已创建副本');
};

// 处理人设切换
const handlePersonaChange = async (persona: { id: number; name: string; avatar: string }) => {
  try {
    await userStore.setActivePersonaId(persona.id);
  } catch {
    messageApi.error('切换人设失败');
  }
};
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-primary);
  position: relative;
}

/* ==================== 欢迎界面 ==================== */
.welcome-screen {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
}

.welcome-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 640px;
  width: 100%;
}

.welcome-logo {
  width: 88px;
  height: 88px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-glow);
  border: 1px solid var(--border-light);
  border-radius: 24px;
  margin-bottom: 20px;
  box-shadow: var(--glow-primary);
  animation: logo-breathe 4s ease-in-out infinite;
}

@keyframes logo-breathe {
  0%, 100% {
    box-shadow: var(--glow-primary-sm);
  }
  50% {
    box-shadow: var(--glow-primary);
  }
}

.welcome-title {
  font-size: 32px;
  font-weight: 600;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 8px;
  letter-spacing: -.5px;
}

.welcome-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
  margin: 0 0 40px;
}

.quick-prompts {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  width: 100%;
  max-width: 480px;
}

.quick-prompt-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  background: var(--gradient-glow);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  transition: all 250ms ease;
  text-align: left;
  position: relative;
  overflow: hidden;
}

.quick-prompt-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--gradient-primary);
  opacity: 0;
  transition: opacity 250ms ease;
}

.quick-prompt-card:hover {
  border-color: var(--color-primary);
  box-shadow: var(--glow-primary-sm);
  transform: translateY(-1px);
}

.quick-prompt-card:hover::before {
  opacity: .06;
}

.prompt-icon {
  color: var(--color-primary);
  flex-shrink: 0;
}

.prompt-text {
  font-size: 13px;
  color: var(--text-primary);
  line-height: 1.4;
}

/* ==================== 消息区域 ==================== */
.messages-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.character-banner {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 24px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.character-avatar {
  flex-shrink: 0;
}

.character-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.message-list {
  flex: 1;
  height: 100%;
}

.message-item-wrapper {
  max-width: 768px;
  margin: 0 auto;
  padding: 0 40px;
}

/* 输入指示器 */
.typing-indicator {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px 0;
}

.typing-bubble {
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-radius: 18px 18px 18px 4px;
}

.typing-dots {
  display: flex;
  gap: 5px;
}

.typing-dots span {
  width: 8px;
  height: 8px;
  background: var(--color-primary);
  border-radius: 50%;
  animation: typing-bounce 1.4s ease-in-out infinite;
  box-shadow: 0 0 6px rgba(77, 168, 255, .4);
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
    opacity: .4;
  }
  30% {
    transform: translateY(-6px);
    opacity: 1;
  }
}

/* 跳到最新按钮 */
.scroll-to-bottom-btn {
  position: absolute;
  bottom: 16px;
  right: 24px;
  z-index: 10;
  box-shadow: var(--glow-primary-sm);
  border: 1px solid var(--border-light);
}

/* ==================== 输入区域 ==================== */
.input-area {
  padding: 0 24px 20px;
  background: var(--bg-primary);
}

.input-container {
  max-width: 768px;
  margin: 0 auto;
}

.input-footer {
  display: flex;
  justify-content: center;
  padding-top: 6px;
}

.model-info {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* ==================== 响应式 ==================== */
@media (max-width: 1023px) {
  .message-item-wrapper {
    padding: 0 24px;
  }
}

@media (max-width: 767px) {
  .welcome-title {
    font-size: 24px;
  }

  .welcome-subtitle {
    font-size: 14px;
    margin-bottom: 28px;
  }

  .quick-prompts {
    grid-template-columns: 1fr;
    max-width: 100%;
  }

  .message-item-wrapper {
    padding: 0 16px;
  }

  .input-area {
    padding: 0 16px 12px;
  }

  .character-banner {
    padding: 8px 16px;
  }
}
</style>
