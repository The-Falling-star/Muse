<template>
  <div class="message-input-container">
    <!-- 附件预览区 -->
    <Transition name="slide-up">
      <div v-if="attachments.length > 0" class="attachments-preview">
        <div
          v-for="(file, index) in attachments"
          :key="index"
          class="attachment-item"
        >
          <n-icon size="16"><DocumentOutline /></n-icon>
          <span class="attachment-name">{{ file.name }}</span>
          <n-button quaternary circle size="tiny" @click="removeAttachment(index)">
            <template #icon>
              <n-icon size="12"><CloseOutline /></n-icon>
            </template>
          </n-button>
        </div>
      </div>
    </Transition>

    <!-- 输入区域 -->
    <div class="input-wrapper">
      <!-- Persona 选择器 -->
      <n-popover trigger="click" placement="top-start" :show-arrow="false">
        <template #trigger>
          <div class="persona-selector" :title="currentPersona.name">
            <n-avatar
              :size="36"
              round
              :src="currentPersona.avatar"
              class="persona-avatar"
            >
              {{ currentPersona.name.charAt(0) }}
            </n-avatar>
            <div class="persona-indicator">
              <n-icon size="10"><ChevronDownOutline /></n-icon>
            </div>
          </div>
        </template>

        <!-- Persona 列表弹出框 -->
        <div class="persona-list">
          <div class="persona-list-header">
            <span>选择人设</span>
            <n-button quaternary size="tiny" @click="openPersonaManager">
              <template #icon>
                <n-icon size="14"><SettingsOutline /></n-icon>
              </template>
            </n-button>
          </div>
          <n-scrollbar style="max-height: 240px">
            <div
              v-for="persona in personas"
              :key="persona.id"
              class="persona-option"
              :class="{ 'active': persona.id === currentPersona.id }"
              @click="selectPersona(persona)"
            >
              <n-avatar :size="32" round :src="persona.avatar" class="persona-option-avatar">
                {{ persona.name.charAt(0) }}
              </n-avatar>
              <span class="persona-option-name">{{ persona.name }}</span>
              <n-icon v-if="persona.id === currentPersona.id" size="16" class="persona-check">
                <CheckmarkOutline />
              </n-icon>
            </div>
          </n-scrollbar>
          <div class="persona-list-footer">
            <n-button quaternary block size="small" @click="openPersonaManager">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              新建人设
            </n-button>
          </div>
        </div>
      </n-popover>

      <!-- 附加功能按钮 -->
      <div class="input-actions-left">
        <n-dropdown :options="quickInsertOptions" trigger="click" @select="handleQuickInsert">
          <n-button quaternary circle size="small">
            <template #icon>
              <n-icon size="18"><FlashOutline /></n-icon>
            </template>
          </n-button>
        </n-dropdown>
      </div>

      <!-- 文本输入框 -->
      <n-input
        ref="inputRef"
        v-model:value="localValue"
        type="textarea"
        :autosize="{ minRows: 1, maxRows: 6 }"
        :placeholder="placeholder"
        :disabled="disabled"
        class="message-textarea"
        @keydown="handleKeydown"
      />

      <!-- 发送按钮 -->
      <div class="input-actions-right">
        <n-tooltip v-if="!isGenerating" trigger="hover">
          <template #trigger>
            <n-button
              type="primary"
              circle
              :disabled="!canSend"
              class="send-button"
              @click="handleSend"
            >
              <template #icon>
                <n-icon size="18"><SendOutline /></n-icon>
              </template>
            </n-button>
          </template>
          发送 (Enter)
        </n-tooltip>

        <n-tooltip v-else trigger="hover">
          <template #trigger>
            <n-button
              type="error"
              circle
              class="stop-button"
              @click="$emit('stop')"
            >
              <template #icon>
                <n-icon size="18"><StopOutline /></n-icon>
              </template>
            </n-button>
          </template>
          停止生成
        </n-tooltip>
      </div>
    </div>

    <!-- 提示信息 -->
    <div class="input-hint">
      <span class="char-count" :class="{ 'warning': charCount > 4000 }">
        {{ charCount }} / 4096
      </span>
    </div>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInputRef"
      type="file"
      multiple
      hidden
      @change="handleFileSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { NInput, NButton, NIcon, NTooltip, NDropdown, NAvatar, NPopover, NScrollbar } from 'naive-ui';
import {
  FlashOutline,
  SendOutline,
  StopOutline,
  CloseOutline,
  DocumentOutline,
  ChevronDownOutline,
  SettingsOutline,
  CheckmarkOutline,
  AddOutline
} from '@vicons/ionicons5';
import { useUserStore } from '@/stores/user';
import { userClient } from '@/api/client';
import type { Persona } from '@/gen/muse/muse_pb';

interface AttachmentFile {
  name: string;
  size: number;
  type: string;
  file: File;
}

const props = withDefaults(defineProps<{
  modelValue?: string;
  placeholder?: string;
  disabled?: boolean;
  isGenerating?: boolean;
}>(), {
  modelValue: '',
  placeholder: '输入消息...',
  disabled: false,
  isGenerating: false
});

const emit = defineEmits<{
  'update:modelValue': [value: string];
  send: [content: string];
  stop: [];
  'persona-change': [persona: Persona];
}>();

const router = useRouter();
const userStore = useUserStore();

const inputRef = ref<InstanceType<typeof NInput> | null>(null);
const fileInputRef = ref<HTMLInputElement | null>(null);
const attachments = ref<AttachmentFile[]>([]);

// 从 userStore 获取真实 personas 数据
const personas = computed(() => userStore.personas);

// 当前活跃人设（从 store 获取）
const currentPersona = computed(() => {
  return userStore.activePersona ?? { id: 0, name: '默认用户', avatar: '' } as Persona;
});

// 选择 Persona，调用后端 API 设置活跃人设
const selectPersona = async (persona: Persona) => {
  try {
    await userClient.setActivePersona({ personaId: persona.id });
    userStore.setActivePersonaId(persona.id);
    emit('persona-change', persona);
  } catch (e) {
    console.error('切换人设失败:', e);
  }
};

// 打开人设管理 - 导航到设置页
const openPersonaManager = () => {
  router.push('/settings');
};

// 本地值
const localValue = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
});

// 字符计数
const charCount = computed(() => localValue.value.length);

// 是否可以发送
const canSend = computed(() => localValue.value.trim().length > 0 && !props.disabled);

// 快速插入选项
const quickInsertOptions = [
  { label: '添加附件', key: 'attachment', icon: () => null },
  { type: 'divider', key: 'd0' },
  { label: '角色扮演开始', key: 'roleplay', icon: () => null },
  { label: '续写', key: 'continue', icon: () => null },
  { label: '重写', key: 'rewrite', icon: () => null },
  { type: 'divider', key: 'd1' },
  { label: '自定义指令...', key: 'custom', icon: () => null }
];

// 处理按键
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault();
    handleSend();
  }
};

// 发送消息
const handleSend = () => {
  if (!canSend.value) return;
  emit('send', localValue.value);
  localValue.value = '';
  attachments.value = [];
};

// 文件选择
const triggerFileInput = () => {
  fileInputRef.value?.click();
};

const handleFileSelect = (e: Event) => {
  const target = e.target as HTMLInputElement;
  const files = target.files;
  if (files) {
    for (const file of files) {
      attachments.value.push({
        name: file.name,
        size: file.size,
        type: file.type,
        file
      });
    }
  }
  target.value = '';
};

const removeAttachment = (index: number) => {
  attachments.value.splice(index, 1);
};

// 快速插入处理
const handleQuickInsert = (key: string) => {
  // 处理添加附件
  if (key === 'attachment') {
    triggerFileInput();
    return;
  }

  const inserts: Record<string, string> = {
    roleplay: '[开始角色扮演]\n',
    continue: '[续写上文]\n',
    rewrite: '[请重写以下内容]\n'
  };

  if (inserts[key]) {
    localValue.value += inserts[key];
    inputRef.value?.focus();
  }
};

// 聚焦输入框
const focus = () => {
  inputRef.value?.focus();
};

defineExpose({ focus });
</script>

<style scoped>
.message-input-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 附件预览 */
.attachments-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.attachment-name {
  max-width: 150px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 输入包装器 */
.input-wrapper {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-light);
  border-radius: 12px;
  padding: 8px 12px;
  transition: all var(--transition-fast);
  position: relative;
  overflow: hidden;
}

/* 输入框科幻光辉 */
.input-wrapper::before {
  content: '';
  position: absolute;
  inset: 0;
  background: var(--gradient-glow);
  opacity: .5;
  pointer-events: none;
  z-index: 0;
}

.input-wrapper > * {
  position: relative;
  z-index: 1;
}

.input-wrapper:focus-within {
  border-color: var(--color-primary);
  box-shadow: var(--glow-primary-sm);
}

/* Persona 选择器 */
.persona-selector {
  position: relative;
  cursor: pointer;
  flex-shrink: 0;
  transition: transform var(--transition-fast);
  /* GPU加速 */
  transform: translateZ(0);
}

.persona-selector:hover {
  transform: scale(1.05);
}

.persona-selector:active {
  transform: scale(0.98);
}

.persona-selector:hover .persona-avatar {
  box-shadow: 0 0 0 2px var(--color-primary);
}

/* 移动端优化 */
@media (max-width: 768px) {
  .persona-selector:hover {
    transform: scale(1.03);
  }

  .persona-selector:active {
    transform: scale(0.95);
  }
}

.persona-avatar {
  background: var(--gradient-primary);
  transition: box-shadow var(--transition-fast);
}

.persona-indicator {
  position: absolute;
  bottom: -2px;
  right: -2px;
  width: 16px;
  height: 16px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-tertiary);
}

/* Persona 列表弹出框 */
.persona-list {
  width: 220px;
  background: var(--bg-card);
  border-radius: 8px;
  overflow: hidden;
}

.persona-list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.persona-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.persona-option:hover {
  background: var(--bg-card-hover);
}

.persona-option.active {
  background: rgba(0, 240, 255, 0.1);
}

.persona-option-avatar {
  background: var(--gradient-primary);
  flex-shrink: 0;
}

.persona-option-name {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.persona-check {
  color: var(--color-primary);
  flex-shrink: 0;
}

.persona-list-footer {
  padding: 8px;
  border-top: 1px solid var(--border-color);
}

.input-actions-left,
.input-actions-right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
  height: 40px;
}

.input-actions-left .n-button {
  color: var(--text-tertiary);
  width: 36px;
  height: 36px;
}

.input-actions-left .n-button:hover {
  color: var(--color-primary);
}

/* 文本输入框 */
.message-textarea {
  flex: 1;
}

.message-textarea :deep(.n-input__textarea-el) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 6px 0 !important;
  resize: none;
}

.message-textarea :deep(.n-input-wrapper) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
}

/* 发送按钮 */
.send-button {
  width: 40px;
  height: 40px;
  transition: all var(--transition-fast);
  /* GPU加速 */
  transform: translateZ(0);
  will-change: transform, box-shadow;
}

.send-button:not(:disabled):hover {
  transform: scale(1.05);
  box-shadow: var(--glow-primary);
}

.send-button:not(:disabled):active {
  transform: scale(0.95);
}

.stop-button {
  width: 40px;
  height: 40px;
  animation: pulse 1.5s ease-in-out infinite;
}

/* 提示信息 */
.input-hint {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: var(--text-tertiary);
  padding: 0 4px;
}

.persona-hint {
  display: flex;
  align-items: center;
  gap: 4px;
}

.persona-hint strong {
  color: var(--color-primary);
  font-weight: 500;
}

.char-count {
  font-variant-numeric: tabular-nums;
}

.char-count.warning {
  color: var(--color-warning);
}

/* 动画 */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.2s var(--transition-fast);
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* 移动端优化 */
@media (max-width: 768px) {
  .input-wrapper {
    padding: 6px 10px;
    gap: 6px;
  }

  .input-actions-left .n-button,
  .send-button,
  .stop-button {
    width: 36px;
    height: 36px;
  }

  .message-textarea :deep(.n-input__textarea-el) {
    padding: 4px 0 !important;
  }

  .send-button:not(:disabled):hover {
    transform: scale(1.03);
  }

  .send-button:not(:disabled):active {
    transform: scale(0.98);
  }

  .persona-hint {
    display: none;
  }
}
</style>
