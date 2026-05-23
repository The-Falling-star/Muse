<template>
  <div class="regex-editor">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
      <n-form-item label="规则名称" path="name">
        <n-input v-model:value="formData.name" placeholder="输入规则名称" />
      </n-form-item>

      <n-form-item label="匹配模式" path="findPattern">
        <n-input v-model:value="formData.findPattern" placeholder="输入正则表达式" font-family="monospace" />
      </n-form-item>

      <n-form-item label="替换内容">
        <n-input v-model:value="formData.replacePattern" placeholder="替换文本（留空表示删除匹配内容）" />
      </n-form-item>

      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="正则模式">
            <n-switch v-model:value="formData.substituteRegex" />
            <span style="margin-left: 8px; color: var(--text-tertiary)">{{ formData.substituteRegex ? '正则匹配' : '字面量匹配' }}</span>
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="编辑时运行">
            <n-switch v-model:value="formData.runOnEdit" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="优先级">
            <n-input-number v-model:value="formData.sortOrder" :min="0" :max="1000" class="full-width" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="启用">
            <n-switch v-model:value="formData.isEnabled" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="最小深度">
            <n-input-number v-model:value="formData.minDepth" :min="0" :max="100" class="full-width" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="最大深度">
            <n-input-number v-model:value="formData.maxDepth" :min="0" :max="100" class="full-width" />
          </n-form-item>
        </n-gi>
      </n-grid>

      <!-- 作用范围 -->
      <n-form-item label="作用范围">
        <n-checkbox-group v-model:value="affectFlagsArray">
          <n-space>
            <n-checkbox value="userInput" label="用户输入" />
            <n-checkbox value="aiOutput" label="AI输出" />
            <n-checkbox value="worldInfo" label="世界书" />
            <n-checkbox value="prompt" label="提示词" />
          </n-space>
        </n-checkbox-group>
      </n-form-item>

      <!-- 测试区域 -->
      <n-divider>测试</n-divider>
      <n-form-item label="测试文本">
        <n-input v-model:value="testInput" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" placeholder="输入测试文本" />
      </n-form-item>
      <n-form-item label="结果">
        <n-input :value="testResult" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" readonly placeholder="测试结果" />
      </n-form-item>
    </n-form>

    <div class="editor-footer">
      <n-button @click="$emit('cancel')">取消</n-button>
      <n-button type="primary" @click="handleSubmit">保存</n-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import { NForm, NFormItem, NInput, NInputNumber, NSwitch, NCheckboxGroup, NCheckbox, NSpace, NGrid, NGi, NDivider, NButton, useMessage } from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import type { RegexRule, RegexAffectFlags } from '@/gen/muse/regex_pb';

interface RegexFormData {
  name: string;
  findPattern: string;
  replacePattern: string;
  isEnabled: boolean;
  runOnEdit: boolean;
  substituteRegex: boolean;
  minDepth: number;
  maxDepth: number;
  sortOrder: number;
  affectFlags: RegexAffectFlags;
}

const props = defineProps<{ regex?: RegexRule | null }>();
const emit = defineEmits<{ save: [rule: Partial<RegexRule>]; cancel: [] }>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);
const testInput = ref('');

const formData = reactive<RegexFormData>({
  name: '',
  findPattern: '',
  replacePattern: '',
  isEnabled: true,
  runOnEdit: false,
  substituteRegex: true,
  minDepth: 0,
  maxDepth: 0,
  sortOrder: 100,
  affectFlags: {
    userInput: false,
    aiOutput: true,
    worldInfo: false,
    prompt: false
  } as RegexAffectFlags
});

const formRules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  findPattern: [{ required: true, message: '请输入匹配模式', trigger: 'blur' }]
};

// 作用范围数组（用于 checkbox-group）
const affectFlagsArray = computed({
  get: () => {
    const flags: string[] = [];
    if (formData.affectFlags?.userInput) flags.push('userInput');
    if (formData.affectFlags?.aiOutput) flags.push('aiOutput');
    if (formData.affectFlags?.worldInfo) flags.push('worldInfo');
    if (formData.affectFlags?.prompt) flags.push('prompt');
    return flags;
  },
  set: (val: string[]) => {
    formData.affectFlags = {
      userInput: val.includes('userInput'),
      aiOutput: val.includes('aiOutput'),
      worldInfo: val.includes('worldInfo'),
      prompt: val.includes('prompt')
    } as RegexAffectFlags;
  }
});

const testResult = computed(() => {
  if (!testInput.value || !formData.findPattern) return '';
  try {
    if (formData.substituteRegex) {
      const regex = new RegExp(formData.findPattern, 'g');
      return testInput.value.replace(regex, formData.replacePattern);
    } else {
      return testInput.value.split(formData.findPattern).join(formData.replacePattern);
    }
  } catch (e) {
    return `错误: ${(e as Error).message}`;
  }
});

watch(() => props.regex, (regex) => {
  if (regex) {
    Object.assign(formData, {
      name: regex.name || '',
      findPattern: regex.findPattern || '',
      replacePattern: regex.replacePattern || '',
      isEnabled: regex.isEnabled ?? true,
      runOnEdit: regex.runOnEdit ?? false,
      substituteRegex: regex.substituteRegex ?? true,
      minDepth: regex.minDepth ?? 0,
      maxDepth: regex.maxDepth ?? 0,
      sortOrder: regex.sortOrder ?? 100,
      affectFlags: regex.affectFlags || { userInput: false, aiOutput: true, worldInfo: false, prompt: false }
    });
  } else {
    Object.assign(formData, {
      name: '',
      findPattern: '',
      replacePattern: '',
      isEnabled: true,
      runOnEdit: false,
      substituteRegex: true,
      minDepth: 0,
      maxDepth: 0,
      sortOrder: 100,
      affectFlags: { userInput: false, aiOutput: true, worldInfo: false, prompt: false }
    });
  }
}, { immediate: true });

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    if (formData.substituteRegex) {
      new RegExp(formData.findPattern, 'g');
    }
    emit('save', {
      id: props.regex?.id || 0,
      name: formData.name,
      findPattern: formData.findPattern,
      replacePattern: formData.replacePattern,
      isEnabled: formData.isEnabled,
      runOnEdit: formData.runOnEdit,
      substituteRegex: formData.substituteRegex,
      minDepth: formData.minDepth,
      maxDepth: formData.maxDepth,
      sortOrder: formData.sortOrder,
      affectFlags: formData.affectFlags
    });
  } catch (e) {
    message.error((e as Error).message || '请检查表单填写是否正确');
  }
};
</script>

<style scoped>
.regex-editor { display: flex; flex-direction: column; }
.full-width { width: 100%; }
.editor-footer { display: flex; justify-content: flex-end; gap: 12px; padding-top: 16px; border-top: 1px solid var(--border-color); margin-top: 8px; }
</style>
