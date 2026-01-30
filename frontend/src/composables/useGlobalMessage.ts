import { useMessage } from 'naive-ui';

// 全局消息实例引用
let messageInstance: ReturnType<typeof useMessage> | null = null;

// 消息队列（在消息实例未初始化时暂存消息）
const messageQueue: Array<{ type: 'error' | 'warning' | 'info' | 'success'; content: string }> = [];

// 初始化全局消息实例（必须在 NMessageProvider 的 setup 中调用）
export function initGlobalMessage() {
  console.log('[GlobalMessage] initGlobalMessage called, messageInstance:', messageInstance);
  if (!messageInstance) {
    try {
      messageInstance = useMessage();
      console.log('[GlobalMessage] messageInstance created:', messageInstance);

      // 处理积压的消息队列
      console.log('[GlobalMessage] Processing queue, size:', messageQueue.length);
      while (messageQueue.length > 0) {
        const msg = messageQueue.shift();
        if (msg) {
          console.log('[GlobalMessage] Showing queued message:', msg);
          showDirectMessage(msg.type, msg.content);
        }
      }
    } catch (error) {
      console.error('[GlobalMessage] Failed to create message instance:', error);
    }
  }
}

// 直接显示消息（在消息实例初始化后）
function showDirectMessage(type: 'error' | 'warning' | 'info' | 'success', content: string) {
  console.log('[GlobalMessage] showDirectMessage called, type:', type, 'content:', content);
  if (!messageInstance) {
    console.error('[GlobalMessage] messageInstance is null!');
    return;
  }

  console.log('[GlobalMessage] Showing message with instance');
  switch (type) {
    case 'error':
      messageInstance.error(content, { duration: 5000 });
      break;
    case 'warning':
      messageInstance.warning(content, { duration: 4000 });
      break;
    case 'info':
      messageInstance.info(content, { duration: 3000 });
      break;
    case 'success':
      messageInstance.success(content, { duration: 3000 });
      break;
  }
}

// 全局消息 API
export const globalMessage = {
  error: (content: string) => {
    console.log('[GlobalMessage] globalMessage.error called, content:', content);
    console.log('[GlobalMessage] messageInstance exists:', !!messageInstance);
    if (!messageInstance) {
      // 消息实例未初始化，加入队列
      console.log('[GlobalMessage] Adding to queue');
      messageQueue.push({ type: 'error', content });
    } else {
      showDirectMessage('error', content);
    }
  },
  warning: (content: string) => {
    if (!messageInstance) {
      messageQueue.push({ type: 'warning', content });
    } else {
      showDirectMessage('warning', content);
    }
  },
  info: (content: string) => {
    if (!messageInstance) {
      messageQueue.push({ type: 'info', content });
    } else {
      showDirectMessage('info', content);
    }
  },
  success: (content: string) => {
    if (!messageInstance) {
      messageQueue.push({ type: 'success', content });
    } else {
      showDirectMessage('success', content);
    }
  }
};
