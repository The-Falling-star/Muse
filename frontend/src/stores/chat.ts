import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { ChatSession, Message } from '@/gen/muse/chat_pb';

/**
 * 聊天 Store
 * 管理当前会话状态、消息列表和流式生成状态
 * CRUD操作直接使用 @/api/chatApi 中的方法
 */
export const useChatStore = defineStore('chat', () => {
  // =====================
  // 全局共享状态
  // =====================

  // 会话列表缓存
  const sessions = ref<Map<number, ChatSession>>(new Map());
  // 当前激活的会话
  const activeSession = ref<ChatSession | null>(null);
  // 当前会话的消息列表
  const messages = ref<Message[]>([]);
  // 是否正在流式生成
  const isStreaming = ref(false);
  // 流式生成中止控制器
  let abortController: AbortController | null = null;

  // =====================
  // 计算属性
  // =====================

  // 当前会话ID
  const activeSessionId = computed(() => activeSession.value?.id ?? null);

  // 当前会话的角色信息
  const activeCharacter = computed(() => activeSession.value?.character ?? null);

  // =====================
  // 会话状态管理
  // =====================

  // 设置会话列表
  const setSessions = (list: ChatSession[], append = false) => {
    if (append) {
      list.forEach(s => sessions.value.set(s.id, s))
    } else {
      sessions.value = new Map(list.map(s => [s.id, s]));
    }
  };

  // 添加会话到列表头部
  const addSession = (session: ChatSession) => {
    sessions.value.set(session.id, session);
  };

  // 更新会话
  const updateSessionInList = (session: ChatSession) => {
    sessions.value.set(session.id, session);
    if (activeSession.value?.id === session.id) {
      activeSession.value = session;
      sessions.value.set(session.id, session);
    }
    if (activeSession.value?.id === session.id) {
      activeSession.value = session;
    }
  };

  // 移除会话
  const removeSession = (id: number) => {
    sessions.value.delete(id)
    if (activeSession.value?.id === id) {
      activeSession.value = null;
      messages.value = [];
    }
  };

  // 移除角色相关的会话
  const removeCharSession = (charID: number) => {
      let sessionsToRemove = Array.from(sessions.value.values()).filter(s => s.character?.id === charID);
      sessionsToRemove.forEach(s => removeSession(s.id))
  };

  // 设置当前激活会话
  const setActiveSession = (session: ChatSession | null) => {
    activeSession.value = session;
    if (session && session.messages) {
      messages.value = session.messages;
      const sessionInMap = sessions.value.get(session.id);
      if (!sessionInMap) {
          console.error("Session not found in map")
      } else {
          sessionInMap.messages = session.messages;
      }
      return
    }
    messages.value = [];
  };

  // =====================
  // 消息状态管理
  // =====================

  // 设置消息列表
  const setMessages = (list: Message[]) => {
    messages.value = list;
  };

  // 添加消息
  const addMessage = (message: Message) => {
    messages.value.push(message);
  };

  // 更新消息
  const updateMessage = (message: Message) => {
    const index = messages.value.findIndex(m => m.id === message.id);
    if (index >= 0) {
      messages.value[index] = message;
    }
  };

  // 移除消息
  const removeMessage = (id: number) => {
    messages.value = messages.value.filter(m => m.id !== id);
  };

  // 本地切换Swipe
  const localSwitchSwipe = (messageId: number, swipeIndex: number) => {
    const message = messages.value.find(m => m.id === messageId);
    if (message && swipeIndex >= 0 && swipeIndex < message.swipes.length) {
      message.activeSwipeIndex = swipeIndex;
    }
  };

  // =====================
  // 流式生成控制
  // =====================

  // 创建新的中止控制器
  const createAbortController = () => {
    abortController = new AbortController();
    return abortController.signal;
  };

  // 开始流式生成
  const startStreaming = () => {
    isStreaming.value = true;
  };

  // 结束流式生成
  const stopStreaming = () => {
    isStreaming.value = false;
    abortController = null;
  };

  // 中止生成
  const abortGeneration = () => {
    if (abortController) {
      abortController.abort();
    }
  };

  // =====================
  // 辅助方法
  // =====================

  // 获取消息的当前内容
  const getMessageContent = (message: Message): string => {
    if (!message.swipes || message.swipes.length === 0) return '';
    const index = Math.min(message.activeSwipeIndex, message.swipes.length - 1);
    return message.swipes[index]?.content ?? '';
  };

  // 创建临时用户消息（用于流式发送前显示）
  const createTempUserMessage = (content: string): Message => {
    const tempId = Date.now();
    return {
      id: tempId,
      sessionId: activeSession.value?.id ?? 0,
      role: 2, // User
      activeSwipeIndex: 0,
      sortOrder: messages.value.length,
      createdAt: BigInt(Date.now()),
      updatedAt: BigInt(Date.now()),
      swipes: [{
        id: tempId,
        messageId: tempId,
        content,
        sortOrder: 0,
        createdAt: BigInt(Date.now()),
        $typeName: 'muse.MessageSwipe'
      }],
      $typeName: 'muse.Message'
    } as Message;
  };

  // 创建临时AI消息占位（用于流式生成）
  const createTempAiMessage = (): Message => {
    const tempId = Date.now() + 1;
    return {
      id: tempId,
      sessionId: activeSession.value?.id ?? 0,
      role: 3, // Assistant
      activeSwipeIndex: 0,
      sortOrder: messages.value.length,
      createdAt: BigInt(Date.now()),
      updatedAt: BigInt(Date.now()),
      swipes: [],
      $typeName: 'muse.Message'
    } as Message;
  };

  // 重置状态
  const reset = () => {
    sessions.value = new Map();
    activeSession.value = null;
    messages.value = [];
    isStreaming.value = false;
    if (abortController) {
      abortController.abort();
      abortController = null;
    }
  };

  return {
    // 状态
    sessions,
    activeSession,
    activeSessionId,
    activeCharacter,
    messages,
    isStreaming,

    // 会话方法
    setSessions,
    addSession,
    updateSessionInList,
    removeSession,
    setActiveSession,
    removeCharSession,

    // 消息方法
    setMessages,
    addMessage,
    updateMessage,
    removeMessage,
    localSwitchSwipe,

    // 流式控制
    createAbortController,
    startStreaming,
    stopStreaming,
    abortGeneration,

    // 辅助方法
    getMessageContent,
    createTempUserMessage,
    createTempAiMessage,
    reset
  };
});
