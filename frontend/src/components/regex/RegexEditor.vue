<template>
  <div class="regex-editor">
    <n-form ref="formRef" :model="formData" :rules="formRules" label-placement="top">
      <n-form-item label="规则名称" path="name">
        <n-input v-model:value="formData.name" placeholder="输入规则名称" />
      </n-form-item>

      <n-form-item label="匹配模式" path="pattern">
        <n-input v-model:value="formData.pattern" placeholder="输入正则表达式" font-family="monospace" />
      </n-form-item>

      <n-form-item label="替换内容">
        <n-input v-model:value="formData.replacement" placeholder="替换文本（留空表示删除匹配内容）" />
      </n-form-item>

      <n-grid :cols="2" :x-gap="16">
        <n-gi>
          <n-form-item label="标志位">
            <n-checkbox-group v-model:value="flagsArray">
              <n-space>
                <n-checkbox value="g" label="全局(g)" />
                <n-checkbox value="i" label="忽略大小写(i)" />
                <n-checkbox value="m" label="多行(m)" />
                <n-checkbox value="s" label="dotAll(s)" />
              </n-space>
            </n-checkbox-group>
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="作用范围">
            <n-select v-model:value="formData.scope" :options="scopeOptions" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="优先级">
            <n-input-number v-model:value="formData.order" :min="0" :max="1000" class="full-width" />
          </n-form-item>
        </n-gi>
        <n-gi>
          <n-form-item label="启用">
            <n-switch v-model:value="formData.enabled" />
          </n-form-item>
        </n-gi>
      </n-grid>

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
import { NForm, NFormItem, NInput, NInputNumber, NSelect, NSwitch, NCheckboxGroup, NCheckbox, NSpace, NGrid, NGi, NDivider, NButton, useMessage } from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import type { RegexRule } from '../../types';

type RegexScope = 'input' | 'output' | 'both';

interface RegexFormData {
  name: string;
  pattern: string;
  replacement: string;
  flags: string;
  scope: RegexScope;
  order: number;
  enabled: boolean;
}

const props = defineProps<{ regex?: RegexRule | null }>();
const emit = defineEmits<{ save: [rule: RegexRule]; cancel: [] }>();

const message = useMessage();
const formRef = ref<FormInst | null>(null);
const testInput = ref('');

const formData = reactive<RegexFormData>({ name: '', pattern: '', replacement: '', flags: 'g', scope: 'output', order: 100, enabled: true });
const formRules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  pattern: [{ required: true, message: '请输入匹配模式', trigger: 'blur' }]
};

const scopeOptions = [
  { label: '仅输入', value: 'input' },
  { label: '仅输出', value: 'output' },
  { label: '双向', value: 'both' }
];

const flagsArray = computed({
  get: () => formData.flags.split(''),
  set: (val: string[]) => { formData.flags = val.join(''); }
});

const testResult = computed(() => {
  if (!testInput.value || !formData.pattern) return '';
  try {
    const regex = new RegExp(formData.pattern, formData.flags);
    return testInput.value.replace(regex, formData.replacement);
  } catch (e) {
    return `错误: ${(e as Error).message}`;
  }
});

watch(() => props.regex, (regex) => {
  if (regex) {
    Object.assign(formData, { name: regex.name, pattern: regex.pattern, replacement: regex.replacement, flags: regex.flags, scope: regex.scope, order: regex.order, enabled: regex.enabled });
  } else {
    Object.assign(formData, { name: '', pattern: '', replacement: '', flags: 'g', scope: 'output', order: 100, enabled: true });
  }
}, { immediate: true });

const handleSubmit = async () => {
  try {
    await formRef.value?.validate();
    new RegExp(formData.pattern, formData.flags);
    emit('save', { id: props.regex?.id || '', name: formData.name, pattern: formData.pattern, replacement: formData.replacement, flags: formData.flags, scope: formData.scope, order: formData.order, enabled: formData.enabled });
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
