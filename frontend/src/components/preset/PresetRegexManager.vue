<template>
  <div class="preset-regex-manager">
    <!-- 工具栏 -->
    <div class="manager-header">
      <span class="header-title">预设正则规则</span>
      <n-button size="small" @click="addRegex">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        添加规则
      </n-button>
    </div>

    <!-- 正则列表 -->
    <div class="regex-list" v-if="rules.length > 0">
      <div
        v-for="(rule, index) in rules"
        :key="rule.id"
        class="regex-item"
        :class="{ 'disabled': !rule.isEnabled }"
      >
        <div class="regex-item-left">
          <n-switch v-model:value="rule.isEnabled" size="small" />
          <span class="regex-index">{{ index + 1 }}</span>
          <div class="regex-info">
            <div class="regex-name">{{ rule.name }}</div>
            <div class="regex-pattern">
              <code>{{ rule.findPattern }}</code>
              <n-tag size="tiny">{{ getAffectFlagsDisplay(rule.affectFlags) }}</n-tag>
            </div>
          </div>
        </div>
        <div class="regex-item-actions">
          <n-button quaternary circle size="tiny" @click="editRegex(index)">
            <template #icon><n-icon size="14"><CreateOutline /></n-icon></template>
          </n-button>
          <n-button quaternary circle size="tiny" @click="deleteRegex(index)">
            <template #icon><n-icon size="14"><TrashOutline /></n-icon></template>
          </n-button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <n-empty v-else description="暂无预设正则规则" size="small" class="empty-state">
      <template #extra>
        <n-button size="small" @click="addRegex">添加规则</n-button>
      </template>
    </n-empty>

    <!-- 编辑模态框 -->
    <n-modal
      v-model:show="showEditModal"
      preset="card"
      :title="editingIndex >= 0 ? '编辑正则规则' : '添加正则规则'"
      :style="{ width: '550px', maxWidth: '90vw' }"
    >
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
            <n-form-item label="作用范围">
              <n-checkbox-group v-model:value="affectFlagsArray">
                <n-space vertical>
                  <n-checkbox value="userInput" label="用户输入" />
                  <n-checkbox value="aiOutput" label="AI输出" />
                  <n-checkbox value="slashCommand" label="斜杠命令" />
                  <n-checkbox value="worldInfo" label="世界书" />
                  <n-checkbox value="prompt" label="提示词" />
                </n-space>
              </n-checkbox-group>
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
        </n-grid>

        <!-- 测试区域 -->
        <n-divider>测试</n-divider>
        <n-form-item label="测试文本">
          <n-input
            v-model:value="testInput"
            type="textarea"
            :autosize="{ minRows: 2, maxRows: 4 }"
            placeholder="输入测试文本"
          />
        </n-form-item>
        <n-form-item label="结果">
          <n-input
            :value="testResult"
            type="textarea"
            :autosize="{ minRows: 2, maxRows: 4 }"
            readonly
            placeholder="测试结果"
          />
        </n-form-item>
      </n-form>

      <div class="modal-footer">
        <n-button @click="showEditModal = false">取消</n-button>
        <n-button type="primary" @click="handleSave">保存</n-button>
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue';
import {
  NButton,
  NIcon,
  NSwitch,
  NTag,
  NEmpty,
  NModal,
  NForm,
  NFormItem,
  NInput,
  NInputNumber,
  NCheckboxGroup,
  NCheckbox,
  NSpace,
  NGrid,
  NGi,
  NDivider,
  useMessage,
  useDialog
} from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import { AddOutline, CreateOutline, TrashOutline } from '@vicons/ionicons5';
import type { RegexRule } from '@/gen/muse/muse_pb';

const props = defineProps<{
  regexRules: RegexRule[];
}>();

const emit = defineEmits<{
  'update:regexRules': [rules: RegexRule[]];
}>();

const message = useMessage();
const dialog = useDialog();

// 内部数据
const rules = computed({
  get: () => props.regexRules || [],
  set: (val) => emit('update:regexRules', val)
});

// 编辑状态
const showEditModal = ref(false);
const editingIndex = ref(-1);
const formRef = ref<FormInst | null>(null);
const testInput = ref('');

// 创建默认的 RegexRule 对象
const createDefaultRegexRule = (): RegexRule => ({
  $typeName: 'muse.RegexRule',
  id: 0,
  presetId: 0,
  characterId: 0,
  name: '',
  findPattern: '',
  replacePattern: '',
  isEnabled: true,
  runOnEdit: false,
  substituteRegex: true,
  minDepth: 0,
  maxDepth: 0,
  affectFlags: {
    userInput: true,
    aiOutput: true,
    slashCommand: false,
    worldInfo: false,
    prompt: false,
    $typeName: 'muse.RegexAffectFlags'
  },
  sortOrder: 100,
  createdAt: 0n,
  updatedAt: 0n
});

const formData = reactive<RegexRule>(createDefaultRegexRule());

const affectFlagsArray = computed({
  get: () => {
    const flags: string[] = [];
    if (formData.affectFlags?.userInput) flags.push('userInput');
    if (formData.affectFlags?.aiOutput) flags.push('aiOutput');
    if (formData.affectFlags?.slashCommand) flags.push('slashCommand');
    if (formData.affectFlags?.worldInfo) flags.push('worldInfo');
    if (formData.affectFlags?.prompt) flags.push('prompt');
    return flags;
  },
  set: (val: string[]) => {
    formData.affectFlags = {
      userInput: val.includes('userInput'),
      aiOutput: val.includes('aiOutput'),
      slashCommand: val.includes('slashCommand'),
      worldInfo: val.includes('worldInfo'),
      prompt: val.includes('prompt'),
      $typeName: 'muse.RegexAffectFlags'
    };
  }
});

const formRules: FormRules = {
  name: { required: true, message: '请输入规则名称', trigger: 'blur' },
  pattern: { required: true, message: '请输入匹配模式', trigger: 'blur' }
};

// 测试结果
const testResult = computed(() => {
  if (!testInput.value || !formData.findPattern) return '';
  try {
    const regex = new RegExp(formData.findPattern, 'g');
    return testInput.value.replace(regex, formData.replacePattern);
  } catch {
    return '正则表达式无效';
  }
});

// 方法
const getAffectFlagsDisplay = (flags: { userInput?: boolean; aiOutput?: boolean; slashCommand?: boolean; worldInfo?: boolean; prompt?: boolean; $typeName?: string } | undefined) => {
  const names: { [key: string]: string } = {
    userInput: '输入',
    aiOutput: '输出',
    slashCommand: '命令',
    worldInfo: '世界书',
    prompt: '提示词'
  };
  const activeFlags = Object.entries(flags || {})
    .filter(([key, enabled]) => key !== '$typeName' && enabled)
    .map(([key, _]) => names[key] || key);
  return activeFlags.join(', ') || '-';
};

const resetForm = () => {
  Object.assign(formData, createDefaultRegexRule());
  testInput.value = '';
};

const addRegex = () => {
  editingIndex.value = -1;
  resetForm();
  showEditModal.value = true;
};

const editRegex = (index: number) => {
  editingIndex.value = index;
  const rule = rules.value[index];
  if (!rule) return;
  Object.assign(formData, rule);
  showEditModal.value = true;
};

const deleteRegex = (index: number) => {
  const rule = rules.value[index];
  if (!rule) return;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除规则"${rule.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      const newRules = [...rules.value];
      newRules.splice(index, 1);
      rules.value = newRules;
      message.success('规则已删除');
    }
  });
};

const handleSave = async () => {
  try {
    await formRef.value?.validate();
  } catch {
    return;
  }

  // 验证正则表达式
  try {
    new RegExp(formData.findPattern, 'g');
  } catch {
    message.error('正则表达式语法错误');
    return;
  }

  const newRules = [...rules.value];
  const existingRule = editingIndex.value >= 0 ? rules.value[editingIndex.value] : null;
  const newRule: RegexRule = {
    ...existingRule,
    ...formData
  };

  if (editingIndex.value >= 0) {
    newRules[editingIndex.value] = newRule;
    message.success('规则已更新');
  } else {
    newRules.push(newRule);
    message.success('规则已添加');
  }

  rules.value = newRules;
  showEditModal.value = false;
};
</script>

<style scoped>
.preset-regex-manager {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.regex-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.regex-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  transition: all var(--transition-fast);
}

.regex-item:hover {
  background: var(--bg-secondary);
}

.regex-item.disabled {
  opacity: 0.5;
}

.regex-item-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.regex-index {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-tertiary);
  background: var(--bg-card);
  border-radius: 4px;
}

.regex-info {
  flex: 1;
  min-width: 0;
}

.regex-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.regex-pattern {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.regex-pattern code {
  background: var(--bg-card);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: var(--font-mono, monospace);
  color: var(--color-primary);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.regex-item-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.regex-item:hover .regex-item-actions {
  opacity: 1;
}

.empty-state {
  padding: 24px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.full-width {
  width: 100%;
}
</style>
