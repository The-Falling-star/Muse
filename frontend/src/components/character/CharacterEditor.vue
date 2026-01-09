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

          <n-form-item label="性格" path="personality">
            <n-input
              v-model:value="formData.personality"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              placeholder="描述角色的性格特点"
            />
          </n-form-item>

          <n-form-item label="场景设定" path="scenario">
            <n-input
              v-model:value="formData.scenario"
              type="textarea"
              :autosize="{ minRows: 3, maxRows: 6 }"
              placeholder="角色所处的背景场景"
            />
          </n-form-item>
        </n-tab-pane>

        <n-tab-pane name="dialogue" tab="对话设定">
          <n-form-item label="开场白" path="firstMessage">
            <n-input
              v-model:value="formData.firstMessage"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 8 }"
              placeholder="角色的第一条消息"
            />
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

          <n-form-item label="系统提示词" path="systemPrompt">
            <n-input
              v-model:value="formData.systemPrompt"
              type="textarea"
              :autosize="{ minRows: 4, maxRows: 8 }"
              placeholder="自定义系统提示词（可选）"
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
import { ref, reactive, watch } from 'vue';
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
import { CameraOutline } from '@vicons/ionicons5';

import type { Character } from '../../types';

interface CharacterFormData {
  name: string;
  avatar: string;
  description: string;
  personality: string;
  scenario: string;
  firstMessage: string;
  exampleDialogue: string;
  creatorNotes: string;
  systemPrompt: string;
}

const props = defineProps<{
  character?: Character | null;
}>();

const emit = defineEmits<{
  save: [character: Character];
  cancel: [];
}>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);

// 表单数据
const formData = reactive<CharacterFormData>({
  name: '',
  avatar: '',
  description: '',
  personality: '',
  scenario: '',
  firstMessage: '',
  exampleDialogue: '',
  creatorNotes: '',
  systemPrompt: ''
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
      personality: char.personality || '',
      scenario: char.scenario || '',
      firstMessage: char.firstMessage || '',
      exampleDialogue: char.exampleDialogue || '',
      creatorNotes: char.creatorNotes || '',
      systemPrompt: ''
    });
  } else {
    // 重置表单
    Object.assign(formData, {
      name: '',
      avatar: '',
      description: '',
      personality: '',
      scenario: '',
      firstMessage: '',
      exampleDialogue: '',
      creatorNotes: '',
      systemPrompt: ''
    });
  }
}, { immediate: true });

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

    const character: Character = {
      id: props.character?.id || '',
      name: formData.name,
      avatar: formData.avatar,
      description: formData.description,
      personality: formData.personality,
      scenario: formData.scenario,
      firstMessage: formData.firstMessage,
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
</style>
