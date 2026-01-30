<template>
  <div class="message-item" :class="[`message-${message.role}`, { 'streaming': message.isStreaming, 'editing': isEditing }]">
    <!-- 头像 -->
    <div class="message-avatar">
      <n-avatar
        v-if="message.role === 'assistant'"
        :size="36"
        round
        :src="character?.avatar"
        class="avatar"
      >
        {{ character?.name?.charAt(0) || 'A' }}
      </n-avatar>
      <n-avatar
        v-else-if="message.role === 'user'"
        :size="36"
        round
        class="avatar user-avatar"
      >
        <n-icon><PersonOutline /></n-icon>
      </n-avatar>
      <div v-else class="system-icon">
        <n-icon><InformationCircleOutline /></n-icon>
      </div>
    </div>

    <!-- 消息内容 -->
    <div class="message-content">
      <div class="message-header">
        <span class="message-author">
          {{ message.role === 'user' ? (persona?.name || '你') : message.role === 'assistant' ? character?.name : '系统' }}
        </span>
        <span class="message-time">{{ formatTime(currentTimestamp) }}</span>
      </div>

      <!-- 编辑模式 -->
      <div v-if="isEditing" class="message-edit-area">
        <n-input
          v-model:value="editContent"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 15 }"
          placeholder="输入消息内容..."
          class="edit-textarea"
        />
        <div class="edit-actions">
          <n-button size="small" @click="cancelEdit">取消</n-button>
          <n-button size="small" type="primary" @click="saveEdit">保存</n-button>
        </div>
      </div>

      <!-- 正常显示模式 -->
      <template v-else>
        <div
          class="message-body"
          :class="{
            'streaming-text': message.isStreaming,
            'swipe-animate': isAnimating,
            'swipe-left': swipeDirection === 'left',
            'swipe-right': swipeDirection === 'right'
          }"
        >
          <div class="message-text" v-html="formattedContent"></div>
        </div>

        <!-- Swipe 导航（当有多条消息时显示） -->
        <div v-if="message.swipes.length > 1" class="swipe-navigation">
          <n-button
            quaternary
            circle
            size="small"
            :disabled="message.currentSwipeIndex === 0"
            @click="handlePrevSwipe"
          >
            <template #icon>
              <n-icon><ChevronBackOutline /></n-icon>
            </template>
          </n-button>
          <span class="swipe-counter">
            {{ message.currentSwipeIndex + 1 }} / {{ message.swipes.length }}
          </span>
          <n-button
            quaternary
            circle
            size="small"
            :disabled="message.currentSwipeIndex === message.swipes.length - 1"
            @click="handleNextSwipe"
          >
            <template #icon>
              <n-icon><ChevronForwardOutline /></n-icon>
            </template>
          </n-button>
        </div>

        <!-- 操作按钮 -->
        <div class="message-actions" v-if="message.role !== 'system'">
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="tiny" @click="copyMessage">
                <template #icon>
                  <n-icon size="14"><CopyOutline /></n-icon>
                </template>
              </n-button>
            </template>
            复制
          </n-tooltip>

          <!-- 编辑按钮（用户和AI消息都可以编辑） -->
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="tiny" @click="startEdit">
                <template #icon>
                  <n-icon size="14"><PencilOutline /></n-icon>
                </template>
              </n-button>
            </template>
            编辑
          </n-tooltip>

          <!-- 用户消息：复制为新版本按钮 -->
          <n-tooltip v-if="message.role === 'user'" trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="tiny" @click="handleDuplicateSwipe">
                <template #icon>
                  <n-icon size="14"><DuplicateOutline /></n-icon>
                </template>
              </n-button>
            </template>
            复制为新版本
          </n-tooltip>

          <!-- AI消息：重新生成按钮 -->
          <n-tooltip v-if="message.role === 'assistant'" trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="tiny" @click="handleRegenerate">
                <template #icon>
                  <n-icon size="14"><RefreshOutline /></n-icon>
                </template>
              </n-button>
            </template>
            重新生成
          </n-tooltip>

          <!-- 分支按钮 -->
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button quaternary circle size="tiny" @click="handleBranch">
                <template #icon>
                  <n-icon size="14"><GitBranchOutline /></n-icon>
                </template>
              </n-button>
            </template>
            从此处创建分支
          </n-tooltip>

          <!-- 删除按钮：多条消息时显示下拉菜单 -->
          <n-dropdown
            v-if="message.swipes.length > 1"
            :options="deleteOptions"
            trigger="click"
            @select="handleDeleteSelect"
          >
            <n-tooltip trigger="hover">
              <template #trigger>
                <n-button quaternary circle size="tiny">
                  <template #icon>
                    <n-icon size="14"><TrashOutline /></n-icon>
                  </template>
                </n-button>
              </template>
              删除
            </n-tooltip>
          </n-dropdown>
          <!-- 只有一条消息时，直接删除整个楼层 -->
          <n-popconfirm
            v-else
            @positive-click="handleDeleteFloor"
          >
            <template #trigger>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-button quaternary circle size="tiny">
                    <template #icon>
                      <n-icon size="14"><TrashOutline /></n-icon>
                    </template>
                  </n-button>
                </template>
                删除
              </n-tooltip>
            </template>
            确定要删除这个楼层吗？
          </n-popconfirm>
        </div>
      </template>
    </div>

    <!-- 删除确认对话框 -->
    <n-modal v-model:show="showDeleteModal" preset="dialog" title="删除确认">
      <template #default>
        <p v-if="deleteType === 'swipe'">确定要删除这条消息吗？删除后将显示相邻的消息。</p>
        <p v-else>确定要删除整个楼层吗？该楼层的所有 {{ message.swipes.length }} 条消息都会被删除。</p>
      </template>
      <template #action>
        <n-button @click="showDeleteModal = false">取消</n-button>
        <n-button type="error" @click="confirmDelete">确认删除</n-button>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import {
  NAvatar,
  NButton,
  NIcon,
  NTooltip,
  NDropdown,
  NModal,
  NPopconfirm,
  NInput,
  useMessage
} from 'naive-ui';
import {
  PersonOutline,
  InformationCircleOutline,
  CopyOutline,
  PencilOutline,
  RefreshOutline,
  TrashOutline,
  ChevronBackOutline,
  ChevronForwardOutline,
  GitBranchOutline,
  DuplicateOutline
} from '@vicons/ionicons5';
import { marked } from 'marked';
import { markedHighlight } from 'marked-highlight';
// 只引入常用语言，减少打包体积
import hljs from 'highlight.js/lib/core';
import javascript from 'highlight.js/lib/languages/javascript';
import typescript from 'highlight.js/lib/languages/typescript';
import python from 'highlight.js/lib/languages/python';
import json from 'highlight.js/lib/languages/json';
import bash from 'highlight.js/lib/languages/bash';
import css from 'highlight.js/lib/languages/css';
import xml from 'highlight.js/lib/languages/xml';
import markdown from 'highlight.js/lib/languages/markdown';
import go from 'highlight.js/lib/languages/go';
import sql from 'highlight.js/lib/languages/sql';
import DOMPurify from 'dompurify';

import type { Character } from '@/gen/muse/muse_pb';

// 消息 Swipe 类型
interface MessageSwipe {
  id: number;
  content: string;
  timestamp: number;
}

// 消息类型（组件本地使用）
interface Message {
  id: number;
  role: 'user' | 'assistant' | 'system';
  swipes: MessageSwipe[];
  currentSwipeIndex: number;
  isStreaming?: boolean;
}

// Persona 类型定义
interface Persona {
  id: number;
  name: string;
  avatar: string;
}

// 注册常用语言
hljs.registerLanguage('javascript', javascript);
hljs.registerLanguage('js', javascript);
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('ts', typescript);
hljs.registerLanguage('python', python);
hljs.registerLanguage('py', python);
hljs.registerLanguage('json', json);
hljs.registerLanguage('bash', bash);
hljs.registerLanguage('sh', bash);
hljs.registerLanguage('css', css);
hljs.registerLanguage('html', xml);
hljs.registerLanguage('xml', xml);
hljs.registerLanguage('markdown', markdown);
hljs.registerLanguage('md', markdown);
hljs.registerLanguage('go', go);
hljs.registerLanguage('sql', sql);

// 配置 marked 使用 highlight.js 进行代码高亮
marked.use(markedHighlight({
  langPrefix: 'hljs language-',
  highlight(code: string, lang: string) {
    if (lang && hljs.getLanguage(lang)) {
      return hljs.highlight(code, { language: lang }).value;
    }
    // 不自动检测语言，避免误判
    return code;
  }
}));

marked.use({
  breaks: true,
  gfm: true
});

const props = defineProps<{
  message: Message;
  character?: Character | null;
  persona?: Persona | null;
}>();

const emit = defineEmits<{
  edit: [id: number, swipeId: number, content: string];  // 编辑消息内容
  delete: [id: number];                    // 删除整个楼层
  deleteSwipe: [id: number, swipeId: number];  // 删除单条消息
  regenerate: [id: number];                // 重新生成（会添加新的swipe）
  swipeChange: [id: number, index: number]; // 切换swipe
  branch: [id: number];                    // 创建分支
  duplicateSwipe: [id: number];            // 复制当前消息为新版本
}>();

const messageApi = useMessage();
const showDeleteModal = ref(false);
const deleteType = ref<'swipe' | 'floor'>('floor');

// 编辑模式状态
const isEditing = ref(false);
const editContent = ref('');

// 动画状态
const swipeDirection = ref<'left' | 'right' | null>(null);
const isAnimating = ref(false);

// 监听 swipe 变化，触发动画
watch(() => props.message.currentSwipeIndex, (newIndex, oldIndex) => {
  if (newIndex !== oldIndex && !isAnimating.value) {
    swipeDirection.value = newIndex > oldIndex ? 'left' : 'right';
    isAnimating.value = true;
    // 动画结束后重置状态
    setTimeout(() => {
      isAnimating.value = false;
      swipeDirection.value = null;
    }, 200);
  }
});

// 删除下拉选项
const deleteOptions = [
  { label: '删除这条消息', key: 'swipe' },
  { label: '删除整个楼层', key: 'floor' }
];

// 获取当前显示的内容
const currentContent = computed(() => {
  if (props.message.swipes.length === 0) return '';
  const index = Math.min(props.message.currentSwipeIndex, props.message.swipes.length - 1);
  return props.message.swipes[index]?.content || '';
});

// 获取当前 swipe
const currentSwipe = computed(() => {
  if (props.message.swipes.length === 0) return null;
  const index = Math.min(props.message.currentSwipeIndex, props.message.swipes.length - 1);
  return props.message.swipes[index] || null;
});

// 获取当前时间戳
const currentTimestamp = computed(() => {
  if (props.message.swipes.length === 0) return Date.now();
  const index = Math.min(props.message.currentSwipeIndex, props.message.swipes.length - 1);
  return props.message.swipes[index]?.timestamp || Date.now();
});

// 格式化时间
const formatTime = (timestamp: number) => {
  const date = new Date(timestamp);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`;

  return date.toLocaleString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
};

// 使用 marked + DOMPurify 安全渲染 Markdown
const formattedContent = computed(() => {
  const raw = currentContent.value;
  if (!raw) return '';
  // 使用 marked 解析 markdown，然后用 DOMPurify 清理 XSS
  const html = marked.parse(raw) as string;
  return DOMPurify.sanitize(html);
});

// 复制消息
const copyMessage = async () => {
  try {
    await navigator.clipboard.writeText(currentContent.value);
    messageApi.success('已复制到剪贴板');
  } catch {
    messageApi.error('复制失败');
  }
};

// 开始编辑
const startEdit = () => {
  editContent.value = currentContent.value;
  isEditing.value = true;
};

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false;
  editContent.value = '';
};

// 保存编辑
const saveEdit = () => {
  if (!editContent.value.trim()) {
    messageApi.warning('消息内容不能为空');
    return;
  }

  const swipe = currentSwipe.value;
  if (swipe) {
    emit('edit', props.message.id, swipe.id, editContent.value);
  }
  isEditing.value = false;
  editContent.value = '';
};

// 切换到上一条消息
const handlePrevSwipe = () => {
  if (props.message.currentSwipeIndex > 0) {
    emit('swipeChange', props.message.id, props.message.currentSwipeIndex - 1);
  }
};

// 切换到下一条消息
const handleNextSwipe = () => {
  if (props.message.currentSwipeIndex < props.message.swipes.length - 1) {
    emit('swipeChange', props.message.id, props.message.currentSwipeIndex + 1);
  }
};

// 重新生成（会添加新的swipe并切换到新消息）
const handleRegenerate = () => {
  emit('regenerate', props.message.id);
};

// 复制当前消息为新版本（用户消息专用）
const handleDuplicateSwipe = () => {
  emit('duplicateSwipe', props.message.id);
};

// 创建分支
const handleBranch = () => {
  emit('branch', props.message.id);
};

// 删除下拉选择处理
const handleDeleteSelect = (key: string) => {
  deleteType.value = key as 'swipe' | 'floor';
  showDeleteModal.value = true;
};

// 只有一条消息时直接删除楼层
const handleDeleteFloor = () => {
  emit('delete', props.message.id);
};

// 确认删除
const confirmDelete = () => {
  if (deleteType.value === 'swipe') {
    const swipe = currentSwipe.value;
    if (swipe) {
      emit('deleteSwipe', props.message.id, swipe.id);
    }
  } else {
    emit('delete', props.message.id);
  }
  showDeleteModal.value = false;
};
</script>

<style scoped>
.message-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  animation: fade-in 0.3s ease-out;
}

.message-item + .message-item {
  border-top: 1px solid var(--border-color);
}

/* 编辑模式样式 */
.message-item.editing {
  background: var(--bg-tertiary);
  margin: 0 -16px;
  padding: 16px;
  border-radius: 8px;
}

/* 头像 */
.message-avatar {
  flex-shrink: 0;
}

.avatar {
  background: var(--gradient-primary);
}

.user-avatar {
  background: var(--gradient-secondary);
}

.system-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  border-radius: 50%;
  color: var(--color-info);
}

/* 消息内容 */
.message-content {
  flex: 1;
  min-width: 0;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.message-author {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.message-assistant .message-author {
  color: var(--color-primary);
}

.message-time {
  font-size: 12px;
  color: var(--text-tertiary);
}

.message-body {
  position: relative;
}

.message-text {
  font-size: 15px;
  line-height: 1.7;
  color: var(--text-primary);
  word-wrap: break-word;
}

/* 编辑区域样式 */
.message-edit-area {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.edit-textarea {
  font-size: 15px;
  line-height: 1.7;
}

.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

/* 代码样式 - 适配 marked + highlight.js */
.message-text :deep(pre) {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 16px;
  margin: 12px 0;
  overflow-x: auto;
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: 13px;
  line-height: 1.5;
}

.message-text :deep(pre code) {
  background: transparent;
  padding: 0;
  font-family: inherit;
  font-size: inherit;
}

.message-text :deep(code) {
  background: var(--bg-tertiary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono, 'JetBrains Mono', monospace);
  font-size: 0.9em;
  color: var(--color-primary);
}

/* Markdown 列表样式 */
.message-text :deep(ul),
.message-text :deep(ol) {
  padding-left: 24px;
  margin: 8px 0;
}

.message-text :deep(li) {
  margin: 4px 0;
}

/* Markdown 链接样式 */
.message-text :deep(a) {
  color: var(--color-primary);
  text-decoration: none;
}

.message-text :deep(a:hover) {
  text-decoration: underline;
}

/* Markdown 引用样式 */
.message-text :deep(blockquote) {
  border-left: 4px solid var(--color-primary);
  margin: 12px 0;
  padding: 8px 16px;
  background: var(--bg-tertiary);
  border-radius: 0 8px 8px 0;
}

/* Markdown 表格样式 */
.message-text :deep(table) {
  border-collapse: collapse;
  width: 100%;
  margin: 12px 0;
}

.message-text :deep(th),
.message-text :deep(td) {
  border: 1px solid var(--border-color);
  padding: 8px 12px;
  text-align: left;
}

.message-text :deep(th) {
  background: var(--bg-tertiary);
  font-weight: 600;
}

/* Swipe 动画 - 使用 GPU 加速的 transform 和 opacity */
.message-body {
  will-change: transform, opacity;
}

.swipe-animate {
  animation-duration: 0.2s;
  animation-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  animation-fill-mode: forwards;
}

.swipe-animate.swipe-left {
  animation-name: swipe-in-left;
}

.swipe-animate.swipe-right {
  animation-name: swipe-in-right;
}

@keyframes swipe-in-left {
  0% {
    opacity: 0;
    transform: translateX(30px);
  }
  100% {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes swipe-in-right {
  0% {
    opacity: 0;
    transform: translateX(-30px);
  }
  100% {
    opacity: 1;
    transform: translateX(0);
  }
}

/* Swipe 导航 */
.swipe-navigation {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 12px;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  width: fit-content;
}

.swipe-counter {
  font-size: 13px;
  color: var(--text-secondary);
  min-width: 50px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}

/* 操作按钮 */
.message-actions {
  display: flex;
  gap: 4px;
  margin-top: 8px;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.message-item:hover .message-actions {
  opacity: 1;
}

.message-actions .n-button {
  color: var(--text-tertiary);
}

.message-actions .n-button:hover {
  color: var(--color-primary);
}

/* 系统消息样式 */
.message-system {
  background: var(--bg-tertiary);
  border-radius: 12px;
  padding: 12px 16px;
  margin: 8px 0;
}

.message-system .message-text {
  font-size: 13px;
  color: var(--text-secondary);
}

/* 流式输出动画 */
.streaming .message-body::after {
  content: '';
  display: inline-block;
  width: 8px;
  height: 16px;
  background: var(--color-primary);
  margin-left: 2px;
  animation: cursor-blink 0.8s step-end infinite;
  vertical-align: text-bottom;
}

@keyframes cursor-blink {
  0%, 50% { opacity: 1; }
  51%, 100% { opacity: 0; }
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
