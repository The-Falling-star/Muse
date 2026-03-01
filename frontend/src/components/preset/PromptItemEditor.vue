<template>
  <div class="prompt-item-editor">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="名称" path="name">
            <n-input
              v-model:value="formData.name"
              placeholder="提示项名称"
              :disabled="isForbidOverrides"
            />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="标识符" path="identifier">
            <n-select
              v-model:value="formData.identifier"
              :options="identifierOptions"
              placeholder="选择标识符类型"
            />
          </n-form-item>
        </n-gi>
      </n-grid>

      <n-form-item label="角色" path="role">
        <n-radio-group v-model:value="formData.role">
          <n-radio-button :value="Role.System">
            <n-icon><ServerOutline /></n-icon>
            系统
          </n-radio-button>
          <n-radio-button :value="Role.User">
            <n-icon><PersonOutline /></n-icon>
            用户
          </n-radio-button>
          <n-radio-button :value="Role.Assistant">
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
          :disabled="isForbidOverrides"
        />
        <p v-if="isForbidOverrides" class="forbid-hint">
          此为系统标记项，名称和内容由系统自动填充，不可编辑
        </p>
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
import { PromptItemIdentifier, Role } from '@/gen/muse/common_pb';

interface PromptItem {
  id: number;
  identifier: PromptItemIdentifier;
  name: string;
  role: Role;
  content: string;
  enabled: boolean;
  marker?: boolean;
  forbidOverrides?: boolean; // 标记项：名称和内容不可编辑
  injection?: {
    position: 'before' | 'after';
    depth: number;
  };
}

// 标识符选项
const identifierOptions = [
  { label: '未指定', value: PromptItemIdentifier.PromptItemIdentifierUnspecified },
  { label: '主提示词', value: PromptItemIdentifier.Main },
  { label: '世界信息（前）', value: PromptItemIdentifier.WorldInfoBefore },
  { label: '人设描述', value: PromptItemIdentifier.PersonaDescription },
  { label: '角色描述', value: PromptItemIdentifier.CharDescription },
  { label: '角色性格', value: PromptItemIdentifier.CharPersonality },
  { label: '场景', value: PromptItemIdentifier.Scenario },
  { label: '辅助提示词', value: PromptItemIdentifier.Nsfw },
  { label: '世界信息（后）', value: PromptItemIdentifier.WorldInfoAfter },
  { label: '对话示例', value: PromptItemIdentifier.DialogueExamples },
  { label: '聊天记录', value: PromptItemIdentifier.ChatHistory },
  { label: '越狱提示词', value: PromptItemIdentifier.Jailbreak },
];

interface FormData {
  name: string;
  identifier: PromptItemIdentifier;
  role: Role;
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
  identifier: PromptItemIdentifier.PromptItemIdentifierUnspecified,
  role: Role.System,
  content: '',
  marker: false,
  injectionPosition: null,
  injectionDepth: 0
});

const formRules: FormRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }]
};

const injectionOptions = [
  { label: '在聊天历史之前', value: 'before' },
  { label: '在聊天历史之后', value: 'after' }
];

// 检查是否禁止覆盖（名称和内容不可编辑）
const isForbidOverrides = computed(() => {
  return props.item?.forbidOverrides === true;
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
      identifier: PromptItemIdentifier.PromptItemIdentifierUnspecified,
      role: Role.System,
      content: '',
      marker: false,
      injectionPosition: null,
      injectionDepth: 0
    });
  }
}, { immediate: true });

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
      forbidOverrides: props.item?.forbidOverrides ?? false, // 保留标记状态
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

.forbid-hint {
  margin: 8px 0 0;
  font-size: 12px;
  color: var(--text-tertiary);
  font-style: italic;
}
</style>
