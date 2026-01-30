<template>
  <div class="prompt-item-editor">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="名称" path="name">
            <n-input v-model:value="formData.name" placeholder="提示项名称" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="标识符" path="identifier">
            <n-input
              v-model:value="formData.identifier"
              placeholder="唯一标识符"
              :disabled="isSystemMarker"
            />
          </n-form-item>
        </n-gi>
      </n-grid>

      <n-form-item label="角色" path="role">
        <n-radio-group v-model:value="formData.role">
          <n-radio-button value="system">
            <n-icon><ServerOutline /></n-icon>
            系统
          </n-radio-button>
          <n-radio-button value="user">
            <n-icon><PersonOutline /></n-icon>
            用户
          </n-radio-button>
          <n-radio-button value="assistant">
            <n-icon><SparklesOutline /></n-icon>
            助手
          </n-radio-button>
        </n-radio-group>
      </n-form-item>

      <n-form-item label="内容" path="content">
        <n-input
          v-model:value="formData.content"
          type="textarea"
          :autosize="{ minRows: 4, maxRows: 12 }"
          placeholder="输入提示内容&#10;&#10;支持变量：&#10;{{char}} - 角色名&#10;{{user}} - 用户名&#10;{{scenario}} - 场景&#10;{{personality}} - 性格&#10;{{description}} - 描述"
        />
      </n-form-item>

      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="注入位置">
            <n-select
              v-model:value="formData.injectionPosition"
              :options="injectionOptions"
              placeholder="默认位置"
              clearable
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="注入深度">
            <n-input-number
              v-model:value="formData.injectionDepth"
              :min="0"
              :max="100"
              :disabled="!formData.injectionPosition"
              class="full-width"
            />
          </n-form-item>
        </n-gi>
      </n-grid>

      <n-form-item>
        <n-checkbox v-model:checked="formData.marker">
          <span class="checkbox-label">
            标记（占位符）
            <n-tooltip>
              <template #trigger>
                <n-icon size="14"><InformationCircleOutline /></n-icon>
              </template>
              标记项不包含实际内容，用于标记特定位置（如聊天历史、世界书等）
            </n-tooltip>
          </span>
        </n-checkbox>
      </n-form-item>
    </n-form>

    <div class="editor-footer">
      <n-button @click="$emit('cancel')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">
        {{ props.item ? '保存' : '添加' }}
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
  NInputNumber,
  NSelect,
  NRadioGroup,
  NRadioButton,
  NCheckbox,
  NGrid,
  NGi,
  NButton,
  NIcon,
  NTooltip,
  useMessage
} from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import {
  ServerOutline,
  PersonOutline,
  SparklesOutline,
  InformationCircleOutline
} from '@vicons/ionicons5';

// 提示项类型（组件本地使用）
interface PromptItem {
  id: number;
  identifier: string;
  name: string;
  role: 'system' | 'user' | 'assistant';
  content: string;
  enabled: boolean;
  marker?: boolean;
  injection?: {
    position: 'before' | 'after';
    depth: number;
  };
}

// 系统标记常量
const SYSTEM_MARKERS = {
  CHAT_HISTORY: 'chat_history',
  WORLD_INFO: 'world_info',
  PERSONA: 'persona',
  CHARACTER: 'character'
} as const;

interface FormData {
  name: string;
  identifier: string;
  role: 'system' | 'user' | 'assistant';
  content: string;
  marker: boolean;
  injectionPosition: 'before' | 'after' | null;
  injectionDepth: number;
}

const props = defineProps<{ item?: PromptItem | null }>();
const emit = defineEmits<{ save: [item: PromptItem]; cancel: [] }>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);

const formData = reactive<FormData>({
  name: '',
  identifier: '',
  role: 'system',
  content: '',
  marker: false,
  injectionPosition: null,
  injectionDepth: 0
});

const formRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  identifier: [
    { required: true, message: '请输入标识符', trigger: 'blur' },
    {
      pattern: /^[a-zA-Z_][a-zA-Z0-9_]*$/,
      message: '标识符只能包含字母、数字和下划线，且不能以数字开头',
      trigger: 'blur'
    }
  ]
};

const injectionOptions = [
  { label: '在聊天历史之前', value: 'before' },
  { label: '在聊天历史之后', value: 'after' }
];

// 检查是否为系统标记（不可编辑标识符）
const isSystemMarker = computed(() => {
  if (!props.item) return false;
  return Object.values(SYSTEM_MARKERS).includes(props.item.identifier as any);
});

// 监听 props 变化，初始化表单
watch(() => props.item, (item) => {
  if (item) {
    Object.assign(formData, {
      name: item.name,
      identifier: item.identifier,
      role: item.role,
      content: item.content,
      marker: item.marker,
      injectionPosition: item.injection?.position || null,
      injectionDepth: item.injection?.depth || 0
    });
  } else {
    Object.assign(formData, {
      name: '',
      identifier: '',
      role: 'system',
      content: '',
      marker: false,
      injectionPosition: null,
      injectionDepth: 0
    });
  }
}, { immediate: true });

// 自动生成标识符
watch(() => formData.name, (name) => {
  if (!props.item && name && !formData.identifier) {
    // 将名称转换为标识符格式
    formData.identifier = name
      .toLowerCase()
      .replace(/[^a-zA-Z0-9\u4e00-\u9fa5]/g, '_')
      .replace(/^[0-9]/, '_$&')
      .replace(/_+/g, '_')
      .replace(/^_|_$/g, '') || 'prompt';
  }
});

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();

    const item: PromptItem = {
      id: props.item?.id || Date.now(),
      identifier: formData.identifier,
      name: formData.name,
      role: formData.role,
      content: formData.content,
      enabled: props.item?.enabled ?? true,
      marker: formData.marker,
      ...(formData.injectionPosition ? {
        injection: {
          position: formData.injectionPosition,
          depth: formData.injectionDepth
        }
      } : {})
    };

    emit('save', item);
  } catch {
    message.error('请检查表单填写是否正确');
  }
};
</script>

<style scoped>
.prompt-item-editor {
  display: flex;
  flex-direction: column;
}

.full-width {
  width: 100%;
}

.checkbox-label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
  margin-top: 8px;
}

:deep(.n-radio-button) {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
</style>
