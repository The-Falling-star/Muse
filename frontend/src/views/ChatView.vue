<template>
  <div class="chat-view">
    <!-- 欢迎界面（无消息时显示） -->
    <div v-if="messages.length === 0" class="welcome-screen">
      <div class="welcome-content">
        <div class="welcome-logo">
          <n-icon size="56" color="var(--color-primary)">
            <SparklesOutline/>
          </n-icon>
        </div>
        <h1 class="welcome-title">Muse</h1>
        <p class="welcome-subtitle">开始你的创意之旅</p>
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
            <div v-if="item.id === -1" class="typing-indicator">
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
            <n-icon>
              <ChevronDownOutline/>
            </n-icon>
          </template>
        </n-button>
      </Transition>
    </div>

    <!-- 浮动输入区域 -->
    <div class="input-area-float">
      <div class="input-container">
        <MessageInput
            v-model="inputMessage"
            :disabled="inputDisabled"
            :placeholder="inputPlaceholder"
            @send="handleSendMessage"
            @stop="handleStopGeneration"
            @persona-change="(p: any) => handlePersonaChange(p)"
        />
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
import {computed, nextTick, onBeforeUnmount, onMounted, ref, watch} from 'vue';
import {useRoute} from 'vue-router';
import type {UploadFileInfo, VirtualListInst} from 'naive-ui';
import {NAvatar, NButton, NIcon, NUpload, NVirtualList, useDialog, useMessage} from 'naive-ui';
import {
  ChevronDownOutline,
  SparklesOutline
} from '@vicons/ionicons5';

import MessageItem from '@/components/chat/MessageItem.vue';
import MessageInput from '@/components/chat/MessageInput.vue';
import {useChatStore} from '@/stores/chat';
import {useUserStore} from '@/stores/user';
import {chatClient} from '@/api/client';
import type {Character, Persona} from '@/gen/muse/muse_pb';

// 本地Message类型适配
interface LocalMessage {
  id: number;
  role: 'user' | 'assistant' | 'system';
  swipes: Array<{
    id: number;
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

// 从 Store 获取数据
const isTyping = computed(() => chatStore.isStreaming);
const hasActiveSession = computed(() => !!chatStore.activeSessionId);
const inputDisabled = computed(() => !hasActiveSession.value || isTyping.value);
const inputPlaceholder = computed(() => {
  if (!hasActiveSession.value) return '请先选择或创建一个会话';
  return '输入消息...';
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
      id: -1,
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

// 优化PC端滚轮滚动：阻止Naive UI的n-scrollbar拦截wheel事件，让浏览器原生滚动接管
let wheelCleanup: (() => void) | null = null;

const setupNativeScroll = () => {
  const vvl = document.querySelector('.messages-area .v-vl') as HTMLElement | null;
  if (!vvl) return;

  // 在capture阶段拦截wheel事件，阻止Naive UI的scrollbar处理器接收到事件
  // 使用stopImmediatePropagation阻止同一元素上的其他监听器
  const handler = (e: WheelEvent) => {
    const canScrollUp = vvl.scrollTop > 0;
    const canScrollDown = vvl.scrollTop < vvl.scrollHeight - vvl.clientHeight - 1;
    const scrollingDown = e.deltaY > 0;
    const scrollingUp = e.deltaY < 0;

    if ((scrollingDown && canScrollDown) || (scrollingUp && canScrollUp)) {
      e.stopImmediatePropagation();
    }
  };

  // 在vvl的父元素(n-scrollbar-container)上capture阶段拦截
  const scrollbarEl = vvl.parentElement;
  if (scrollbarEl) {
    scrollbarEl.addEventListener('wheel', handler, { capture: true, passive: true });
    wheelCleanup = () => {
      scrollbarEl.removeEventListener('wheel', handler, { capture: true } as EventListenerOptions);
    };
  }
};

// 初始化
onMounted(async () => {
  // 如果URL中有sessionId参数，加载该会话
  const sessionId = route.params.sessionId;
  if (sessionId) {
    await loadSession(Number(sessionId));
  }

  // 延迟设置原生滚动优化，确保DOM已渲染
  nextTick(() => {
    setupNativeScroll();
  });
});

onBeforeUnmount(() => {
  wheelCleanup?.();
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
    virtualListRef.value.scrollTo({position: 'bottom', behavior: 'smooth'});
  }
  showScrollToBottom.value = false;
};

// 处理文件上传
const handleImportFile = (options: { file: UploadFileInfo }) => {
  const {file} = options;
  if (!file.file) return false;

  const reader = new FileReader();
  reader.onload = (e) => {
    const content = e.target?.result as string;
    JSON.parse(content);
    messageApi.info('导入功能暂未实现');
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
        {sessionId: chatStore.activeSessionId, content: content},
        {signal: signal}
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
};

// 删除整个楼层
const handleDeleteMessage = async (id: number) => {
  dialog.warning({
    title: '确认删除',
    content: '确定要删除这条消息吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await chatClient.deleteMessage({messageId: id});
      chatStore.removeMessage(id);
      messageApi.success('消息已删除');
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
        {messageId: id},
        {signal: signal}
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
  await userStore.setActivePersonaId(persona.id);
};
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 0;
  background: var(--bg-primary);
  position: relative;
  overflow: hidden;
}

/* ==================== 欢迎界面 ==================== */
.welcome-screen {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  padding-bottom: 100px;
  overflow-y: auto;
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
  margin: 0;
}

/* ==================== 消息区域 ==================== */
.messages-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  min-height: 0;
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
  min-height: 0;
  padding-bottom: 100px;
}

/* 优化虚拟列表滚动性能 */
.message-list :deep(.v-vl) {
  overscroll-behavior: contain;
  will-change: scroll-position;
  -webkit-overflow-scrolling: touch;
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
  bottom: 110px;
  right: 24px;
  z-index: 10;
  box-shadow: var(--glow-primary-sm);
  border: 1px solid var(--border-light);
}

/* ==================== 浮动输入区域 ==================== */
.input-area-float {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding: 12px 24px 20px;
  background: linear-gradient(to bottom, transparent, var(--bg-primary) 28%);
  pointer-events: none;
}

.input-area-float .input-container {
  pointer-events: auto;
}

.input-container {
  max-width: 768px;
  margin: 0 auto;
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

  .message-item-wrapper {
    padding: 0 16px;
  }

  .input-area-float {
    padding: 8px 16px 12px;
  }

  .character-banner {
    padding: 8px 16px;
  }
}
</style>
