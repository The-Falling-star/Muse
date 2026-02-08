<template>
  <div class="character-editor">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="top"
    >
      <!-- 基本信息 -->
      <n-grid :cols="24" :x-gap="16">
        <n-gi :span="24">
          <div class="avatar-upload">
            <n-upload
              :max="1"
              accept="image/*"
              :show-file-list="false"
              @change="handleAvatarChange"
            >
              <n-avatar
                :size="100"
                round
                :src="formData.avatar"
                class="upload-avatar"
              >
                <template #placeholder>
                  <n-icon size="32"><CameraOutline /></n-icon>
                </template>
                {{ formData.name?.charAt(0) || '?' }}
              </n-avatar>
            </n-upload>
            <span class="upload-hint">点击上传头像</span>
          </div>
        </n-gi>

        <n-gi :span="24">
          <n-form-item label="角色名称" path="name">
            <n-input v-model:value="formData.name" placeholder="输入角色名称" />
          </n-form-item>
        </n-gi>

      </n-grid>

      <!-- 详细信息标签页 -->
      <n-tabs type="line" animated>
        <n-tab-pane name="basic" tab="基本信息">
          <n-form-item label="描述" path="description">
            <n-input
              v-model:value="formData.description"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              placeholder="角色的简短描述"
            />
          </n-form-item>
        </n-tab-pane>

        <n-tab-pane name="dialogue" tab="对话设定">
          <n-form-item label="开场白">
            <div class="first-message-carousel">
              <div class="carousel-header">
                <span class="message-indicator">
                  {{ currentMessageIndex + 1 }} / {{ formData.firstMessages.length + (isAddingNew ? 1 : 0) }}
                </span>
              </div>
              <div class="carousel-content">
                <n-button
                  text
                  class="carousel-arrow"
                  :disabled="!isAddingNew && currentMessageIndex === 0"
                  @click="prevMessage"
                >
                  <n-icon size="24"><ChevronBackOutline /></n-icon>
                </n-button>
                <div class="message-input-wrapper">
                  <n-input
                    v-model:value="currentMessageContent"
                    type="textarea"
                    :autosize="{ minRows: 4, maxRows: 8 }"
                    :placeholder="isAddingNew ? '输入新的开场白...' : `开场白 ${currentMessageIndex + 1}`"
                  />
                  <div v-if="!isAddingNew && formData.firstMessages.length > 1" class="message-actions">
                    <n-button
                      text
                      type="error"
                      size="small"
                      @click="deleteCurrentMessage"
                    >
                      <template #icon>
                        <n-icon><TrashOutline /></n-icon>
                      </template>
                      删除此开场白
                    </n-button>
                  </div>
                </div>
                <n-button
                  text
                  class="carousel-arrow"
                  :disabled="!canGoNext"
                  @click="nextMessage"
                >
                  <n-icon size="24"><ChevronForwardOutline /></n-icon>
                </n-button>
              </div>
            </div>
          </n-form-item>

          <n-form-item label="示例对话" path="exampleDialogue">
            <n-input
              v-model:value="formData.exampleDialogue"
              type="textarea"
              :autosize="{ minRows: 6, maxRows: 12 }"
              placeholder="示例对话格式：&#10;{{user}}: 用户消息&#10;{{char}}: 角色回复"
            />
          </n-form-item>
        </n-tab-pane>

        <n-tab-pane name="advanced" tab="高级设置">
          <n-form-item label="创作者备注" path="creatorNotes">
            <n-input
              v-model:value="formData.creatorNotes"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              placeholder="给其他用户的备注信息"
            />
          </n-form-item>
        </n-tab-pane>
      </n-tabs>
    </n-form>

    <!-- 底部按钮 -->
    <div class="editor-footer">
      <n-button @click="$emit('cancel')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">
        {{ character ? '保存修改' : '创建角色' }}
      </n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, computed } from 'vue';
import {
  NForm,
  NFormItem,
  NInput,
  NGrid,
  NGi,
  NAvatar,
  NUpload,
  NTabs,
  NTabPane,
  NButton,
  NIcon,
  useMessage
} from 'naive-ui';
import type { FormInst, FormRules, UploadFileInfo } from 'naive-ui';
import { CameraOutline, ChevronBackOutline, ChevronForwardOutline, TrashOutline } from '@vicons/ionicons5';

import type { Character } from '@/gen/muse/muse_pb';

interface CharacterFormData {
  name: string;
  avatar: string;
  description: string;
  firstMessages: string[];
  exampleDialogue: string;
  creatorNotes: string;
}

const props = defineProps<{
  character?: Character | null;
}>();

const emit = defineEmits<{
  save: [character: Partial<Character>];
  cancel: [];
}>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);

// 开场白轮播相关
const currentMessageIndex = ref(0);
const isAddingNew = ref(false);
const newMessageContent = ref('');

// 表单数据
const formData = reactive<CharacterFormData>({
  name: '',
  avatar: '',
  description: '',
  firstMessages: [''],
  exampleDialogue: '',
  creatorNotes: ''
});

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入角色名称', trigger: 'blur' },
    { min: 1, max: 50, message: '名称长度在1-50个字符之间', trigger: 'blur' }
  ]
};

// 监听character变化，初始化表单
watch(() => props.character, (char) => {
  if (char) {
    Object.assign(formData, {
      name: char.name || '',
      avatar: char.avatar || '',
      description: char.description || '',
      firstMessages: (char.firstMessage && char.firstMessage.length > 0) 
        ? [...char.firstMessage] 
        : [''],
      exampleDialogue: char.exampleDialogue || '',
      creatorNotes: char.creatorNotes || ''
    });
  } else {
    // 重置表单
    Object.assign(formData, {
      name: '',
      avatar: '',
      description: '',
      firstMessages: [''],
      exampleDialogue: '',
      creatorNotes: ''
    });
  }
  // 重置轮播状态
  currentMessageIndex.value = 0;
  isAddingNew.value = false;
  newMessageContent.value = '';
}, { immediate: true });

// 当前显示的消息内容（双向绑定）
const currentMessageContent = computed({
  get() {
    if (isAddingNew.value) {
      return newMessageContent.value;
    }
    return formData.firstMessages[currentMessageIndex.value] || '';
  },
  set(value: string) {
    if (isAddingNew.value) {
      newMessageContent.value = value;
    } else {
      formData.firstMessages[currentMessageIndex.value] = value;
    }
  }
});

// 是否可以向右切换
const canGoNext = computed(() => {
  // 如果正在新增，不能继续向右
  if (isAddingNew.value) {
    return false;
  }
  // 如果不是最后一个，可以向右
  if (currentMessageIndex.value < formData.firstMessages.length - 1) {
    return true;
  }
  // 如果是最后一个，可以向右进入新增模式
  return true;
});

// 上一条消息
const prevMessage = () => {
  if (isAddingNew.value) {
    // 从新增模式返回
    if (newMessageContent.value.trim() === '') {
      // 如果新增的内容为空，直接返回，不添加
      isAddingNew.value = false;
      newMessageContent.value = '';
    } else {
      // 如果有内容，保存到数组
      formData.firstMessages.push(newMessageContent.value);
      currentMessageIndex.value = formData.firstMessages.length - 1;
      isAddingNew.value = false;
      newMessageContent.value = '';
    }
  } else if (currentMessageIndex.value > 0) {
    currentMessageIndex.value--;
  }
};

// 下一条消息
const nextMessage = () => {
  if (isAddingNew.value) {
    return;
  }
  
  if (currentMessageIndex.value < formData.firstMessages.length - 1) {
    // 切换到下一条
    currentMessageIndex.value++;
  } else {
    // 切换到新增模式
    isAddingNew.value = true;
    newMessageContent.value = '';
  }
};

// 删除当前开场白
const deleteCurrentMessage = () => {
  // 至少保留一条开场白
  if (formData.firstMessages.length <= 1) {
    message.warning('至少需要保留一条开场白');
    return;
  }

  // 删除当前索引的开场白
  formData.firstMessages.splice(currentMessageIndex.value, 1);

  // 调整当前索引
  if (currentMessageIndex.value >= formData.firstMessages.length) {
    currentMessageIndex.value = formData.firstMessages.length - 1;
  }

  message.success('已删除开场白');
};

// 处理头像上传
const handleAvatarChange = (options: { fileList: UploadFileInfo[] }) => {
  const file = options.fileList[0];
  if (file?.file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      formData.avatar = e.target?.result as string;
    };
    reader.readAsDataURL(file.file);
  }
};

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value?.validate();

    // 如果当前在新增模式，需要检查是否保存
    if (isAddingNew.value) {
      if (newMessageContent.value.trim() !== '') {
        formData.firstMessages.push(newMessageContent.value);
      }
      isAddingNew.value = false;
      newMessageContent.value = '';
    }

    // 过滤掉空的开场白
    const filteredFirstMessages = formData.firstMessages.filter(msg => msg.trim() !== '');

    const character: Partial<Character> = {
      id: props.character?.id || 0,
      name: formData.name,
      avatar: formData.avatar,
      description: formData.description,
      firstMessage: filteredFirstMessages.length > 0 ? filteredFirstMessages : [''],
      exampleDialogue: formData.exampleDialogue,
      creatorNotes: formData.creatorNotes
    };

    emit('save', character);
  } catch {
    message.error('请检查表单填写是否正确');
  }
};
</script>

<style scoped>
.character-editor {
  display: flex;
  flex-direction: column;
}

/* 头像上传 */
.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
}

.upload-avatar {
  cursor: pointer;
  background: var(--gradient-primary);
  transition: all var(--transition-fast);
}

.upload-avatar:hover {
  transform: scale(1.05);
  box-shadow: var(--glow-primary);
}

.upload-hint {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 底部按钮 */
.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 24px;
  margin-top: 16px;
  border-top: 1px solid var(--border-color);
}

/* 开场白轮播 */
.first-message-carousel {
  width: 100%;
}

.carousel-header {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 8px;
}

.message-indicator {
  font-size: 12px;
  color: var(--text-tertiary);
}

.carousel-content {
  display: flex;
  align-items: flex-start;
  gap: 8px;
}

.message-input-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.message-actions {
  display: flex;
  justify-content: flex-end;
}

.carousel-arrow {
  flex-shrink: 0;
  margin-top: 8px;
}
</style>
