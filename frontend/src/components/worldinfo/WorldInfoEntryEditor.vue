<template>
  <div class="entry-editor">
    <n-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-placement="top"
    >
      <n-form-item label="关键词">
        <n-dynamic-tags v-model:value="formData.keys" />
        <template #feedback>
          多个关键词用回车分隔，当消息包含关键词时触发此条目
        </template>
      </n-form-item>

      <n-form-item label="次要关键词">
        <n-dynamic-tags v-model:value="formData.secondaryKeys" />
        <template #feedback>
          可选，用于更精确的匹配
        </template>
      </n-form-item>

      <n-form-item label="内容" path="content">
        <n-input
          v-model:value="formData.content"
          type="textarea"
          :autosize="{ minRows: 4, maxRows: 10 }"
          placeholder="条目内容将被插入到上下文中"
        />
      </n-form-item>

      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="优先级">
            <n-input-number
              v-model:value="formData.order"
              :min="0"
              :max="1000"
              class="full-width"
            />
          </n-form-item>
        </n-gi>

        <n-gi>
          <n-form-item label="深度">
            <n-input-number
              v-model:value="formData.depth"
              :min="0"
              :max="100"
              class="full-width"
            />
          </n-form-item>
        </n-gi>

        <n-gi>
          <n-form-item label="触发概率 (%)">
            <n-input-number
              v-model:value="formData.probability"
              :min="0"
              :max="100"
              class="full-width"
            />
          </n-form-item>
        </n-gi>

        <n-gi>
          <n-form-item label="选择逻辑">
            <n-select
              v-model:value="formData.selectiveLogic"
              :options="logicOptions"
            />
          </n-form-item>
        </n-gi>
      </n-grid>

      <n-form-item label="备注">
        <n-input
          v-model:value="formData.comment"
          placeholder="可选的备注信息"
        />
      </n-form-item>

      <n-form-item label="启用">
        <n-switch v-model:value="formData.enabled" />
      </n-form-item>
    </n-form>

    <div class="editor-footer">
      <n-button @click="$emit('cancel')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">保存</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import {
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NSelect,
  NSwitch,
  NDynamicTags,
  NGrid,
  NGi,
  NButton,
  useMessage
} from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';

// 世界书条目类型（组件本地使用）
interface WorldInfoEntry {
  id: number;
  worldId: number;
  keys: string[];
  secondaryKeys?: string[];
  content: string;
  comment?: string;
  enabled: boolean;
  order: number;
  probability: number;
  depth: number;
  selectiveLogic: 'and' | 'or' | 'not';
}

interface EntryFormData {
  keys: string[];
  secondaryKeys: string[];
  content: string;
  comment: string;
  enabled: boolean;
  order: number;
  probability: number;
  depth: number;
  selectiveLogic: 'and' | 'or' | 'not';
}

const props = defineProps<{
  entry?: WorldInfoEntry | null;
}>();

const emit = defineEmits<{
  save: [entry: WorldInfoEntry];
  cancel: [];
}>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);

const formData = reactive<EntryFormData>({
  keys: [],
  secondaryKeys: [],
  content: '',
  comment: '',
  enabled: true,
  order: 100,
  probability: 100,
  depth: 4,
  selectiveLogic: 'or'
});

const formRules: FormRules = {
  content: [
    { required: true, message: '请输入条目内容', trigger: 'blur' }
  ]
};

const logicOptions = [
  { label: '任意匹配 (OR)', value: 'or' },
  { label: '全部匹配 (AND)', value: 'and' },
  { label: '排除 (NOT)', value: 'not' }
];

watch(() => props.entry, (entry) => {
  if (entry) {
    Object.assign(formData, {
      keys: entry.keys || [],
      secondaryKeys: entry.secondaryKeys || [],
      content: entry.content || '',
      comment: entry.comment || '',
      enabled: entry.enabled ?? true,
      order: entry.order ?? 100,
      probability: entry.probability ?? 100,
      depth: entry.depth ?? 4,
      selectiveLogic: entry.selectiveLogic || 'or'
    });
  } else {
    Object.assign(formData, {
      keys: [],
      secondaryKeys: [],
      content: '',
      comment: '',
      enabled: true,
      order: 100,
      probability: 100,
      depth: 4,
      selectiveLogic: 'or'
    });
  }
}, { immediate: true });

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();

    const entry: WorldInfoEntry = {
      id: props.entry?.id || 0,
      worldId: props.entry?.worldId || 0,
      keys: formData.keys,
      secondaryKeys: formData.secondaryKeys,
      content: formData.content,
      comment: formData.comment || undefined,
      enabled: formData.enabled,
      order: formData.order,
      probability: formData.probability,
      depth: formData.depth,
      selectiveLogic: formData.selectiveLogic
    };

    emit('save', entry);
  } catch {
    message.error('请检查表单填写是否正确');
  }
};
</script>

<style scoped>
.entry-editor {
  display: flex;
  flex-direction: column;
}

.full-width {
  width: 100%;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
  margin-top: 8px;
}
</style>
