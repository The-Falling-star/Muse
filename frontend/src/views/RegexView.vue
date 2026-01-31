<template>
  <div class="regex-view">
    <!-- 工具栏 -->
    <div class="toolbar">
      <n-input v-model:value="searchQuery" placeholder="搜索正则..." clearable class="search-input">
        <template #prefix><n-icon><SearchOutline /></n-icon></template>
      </n-input>
      <n-button-group>
        <n-button @click="handleImport">
          <template #icon><n-icon><CloudUploadOutline /></n-icon></template>
          导入
        </n-button>
        <n-button @click="handleExport">
          <template #icon><n-icon><CloudDownloadOutline /></n-icon></template>
          导出
        </n-button>
      </n-button-group>
      <n-button type="primary" @click="createRegex">
        <template #icon><n-icon><AddOutline /></n-icon></template>
        新建规则
      </n-button>
    </div>

    <!-- 规则列表 -->
    <n-spin :show="loading">
      <n-scrollbar class="regex-container">
        <div class="regex-list">
          <TransitionGroup name="regex-list">
            <div
              v-for="rule in filteredRules"
              :key="rule.id"
              class="regex-card"
              :class="{ 'disabled': !rule.isEnabled }"
            >
              <div class="regex-header">
                <n-switch
                  :value="rule.isEnabled"
                  size="small"
                  @update:value="(val: boolean) => handleToggleEnabled(rule, val)"
                />
                <h3 class="regex-name">{{ rule.name }}</h3>
                <div class="regex-tags">
                  <n-tag v-if="rule.affectFlags?.userInput" type="info" size="small">输入</n-tag>
                  <n-tag v-if="rule.affectFlags?.aiOutput" type="success" size="small">输出</n-tag>
                  <n-tag v-if="rule.affectFlags?.worldInfo" type="warning" size="small">世界书</n-tag>
                  <n-tag v-if="rule.affectFlags?.prompt" type="error" size="small">提示词</n-tag>
                </div>
                <div class="regex-actions">
                  <n-button quaternary circle size="tiny" @click="editRegex(rule)">
                    <template #icon><n-icon size="14"><CreateOutline /></n-icon></template>
                  </n-button>
                  <n-button quaternary circle size="tiny" @click="deleteRegex(rule)">
                    <template #icon><n-icon size="14"><TrashOutline /></n-icon></template>
                  </n-button>
                </div>
              </div>

              <div class="regex-body">
                <div class="regex-pattern">
                  <span class="label">匹配:</span>
                  <code>{{ rule.findPattern }}</code>
                  <n-tag v-if="rule.substituteRegex" size="tiny" :bordered="false">正则</n-tag>
                  <n-tag v-else size="tiny" :bordered="false">字面量</n-tag>
                </div>
                <div class="regex-replacement">
                  <span class="label">替换:</span>
                  <code>{{ rule.replacePattern || '(删除匹配内容)' }}</code>
                </div>
              </div>

              <div class="regex-footer">
                <span class="regex-order">
                  <n-icon><SwapVerticalOutline /></n-icon>
                  优先级: {{ rule.sortOrder }}
                </span>
                <span v-if="rule.minDepth > 0 || rule.maxDepth > 0" class="regex-depth">
                  深度: {{ rule.minDepth }} - {{ rule.maxDepth || '∞' }}
                </span>
                <span v-if="rule.runOnEdit" class="regex-run-on-edit">
                  <n-tag size="tiny" type="info">编辑时运行</n-tag>
                </span>
              </div>
            </div>
          </TransitionGroup>
        </div>

        <n-empty v-if="filteredRules.length === 0 && !loading" description="暂无正则规则" class="empty-state">
          <template #extra><n-button type="primary" @click="createRegex">创建规则</n-button></template>
        </n-empty>
      </n-scrollbar>
    </n-spin>

    <!-- 编辑模态框 -->
    <n-modal
      v-model:show="showEditModal"
      preset="card"
      :title="editingRegex ? '编辑规则' : '新建规则'"
      :style="{ width: '600px', maxWidth: '90vw' }"
    >
      <RegexEditor
        :regex="editingRegex"
        @save="handleSaveRegex"
        @cancel="showEditModal = false"
      />
    </n-modal>

    <!-- 导入模态框 -->
    <n-modal
      v-model:show="showImportModal"
      preset="card"
      title="导入正则规则"
      :style="{ width: '500px', maxWidth: '90vw' }"
    >
      <n-upload
        :max="1"
        accept=".json"
        :default-upload="false"
        @change="handleImportFile"
      >
        <n-upload-dragger>
          <div class="upload-content">
            <n-icon size="48" class="upload-icon"><CloudUploadOutline /></n-icon>
            <p class="upload-text">点击或拖拽文件到此处</p>
            <p class="upload-hint">支持 JSON 格式的正则规则文件</p>
          </div>
        </n-upload-dragger>
      </n-upload>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
  NInput,
  NButton,
  NButtonGroup,
  NIcon,
  NScrollbar,
  NSwitch,
  NTag,
  NEmpty,
  NModal,
  NSpin,
  NUpload,
  NUploadDragger,
  useMessage,
  useDialog
} from 'naive-ui';
import type { UploadFileInfo } from 'naive-ui';
import {
  SearchOutline,
  AddOutline,
  CreateOutline,
  TrashOutline,
  SwapVerticalOutline,
  CloudUploadOutline,
  CloudDownloadOutline
} from '@vicons/ionicons5';
import RegexEditor from '../components/regex/RegexEditor.vue';
import { regexRuleClient } from '@/api/client';
import type { RegexRule } from '@/gen/muse/muse_pb';

const message = useMessage();
const dialog = useDialog();

// 状态
const searchQuery = ref('');
const showEditModal = ref(false);
const showImportModal = ref(false);
const editingRegex = ref<RegexRule | null>(null);
const loading = ref(false);
const rules = ref<RegexRule[]>([]);

// 当前预设ID（0表示全局规则）
const currentPresetId = ref(0);

// 过滤后的规则列表
const filteredRules = computed(() => {
  if (!searchQuery.value) return rules.value;
  const query = searchQuery.value.toLowerCase();
  return rules.value.filter(r =>
    r.name.toLowerCase().includes(query) ||
    r.findPattern.toLowerCase().includes(query)
  );
});

// 加载规则列表
const loadRules = async () => {
  loading.value = true;
  try {
    const response = await regexRuleClient.listRegexRules({});
    rules.value = response.rules;
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadRules();
});

// 创建规则
const createRegex = () => {
  editingRegex.value = null;
  showEditModal.value = true;
};

// 编辑规则
const editRegex = (rule: RegexRule) => {
  editingRegex.value = rule;
  showEditModal.value = true;
};

// 删除规则
const deleteRegex = (rule: RegexRule) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除规则"${rule.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      await regexRuleClient.deleteRegexRule({ id: rule.id });
      rules.value = rules.value.filter(r => r.id !== rule.id);
      message.success('规则已删除');
    }
  });
};

// 切换启用状态
const handleToggleEnabled = async (rule: RegexRule, enabled: boolean) => {
  const response = await regexRuleClient.updateRegexRule({
    id: rule.id,
    name: rule.name,
    findPattern: rule.findPattern,
    replacePattern: rule.replacePattern,
    isEnabled: enabled,
    runOnEdit: rule.runOnEdit,
    substituteRegex: rule.substituteRegex,
    minDepth: rule.minDepth,
    maxDepth: rule.maxDepth,
    affectFlags: rule.affectFlags,
    sortOrder: rule.sortOrder
  });
  if (response.rule) {
    const index = rules.value.findIndex(r => r.id === rule.id);
    if (index >= 0) {
      rules.value[index] = response.rule;
    }
  }
};

// 保存规则
const handleSaveRegex = async (ruleData: {
  name?: string;
  findPattern?: string;
  replacePattern?: string;
  isEnabled?: boolean;
  runOnEdit?: boolean;
  substituteRegex?: boolean;
  minDepth?: number;
  maxDepth?: number;
  affectFlags?: { userInput?: boolean; aiOutput?: boolean; slashCommands?: boolean; worldInfo?: boolean };
  sortOrder?: number;
}) => {
  try {
    if (editingRegex.value) {
      // 编辑模式
      const response = await regexRuleClient.updateRegexRule({
        id: editingRegex.value.id,
        name: ruleData.name || '',
        findPattern: ruleData.findPattern || '',
        replacePattern: ruleData.replacePattern || '',
        isEnabled: ruleData.isEnabled ?? true,
        runOnEdit: ruleData.runOnEdit ?? false,
        substituteRegex: ruleData.substituteRegex ?? true,
        minDepth: ruleData.minDepth ?? 0,
        maxDepth: ruleData.maxDepth ?? 0,
        affectFlags: ruleData.affectFlags,
        sortOrder: ruleData.sortOrder ?? 0
      });
      if (response.rule) {
        const index = rules.value.findIndex(r => r.id === editingRegex.value!.id);
        if (index >= 0) {
          rules.value[index] = response.rule;
        }
      }
      message.success('规则已更新');
    } else {
      // 创建模式
      const response = await regexRuleClient.addRegexRule({
        presetId: currentPresetId.value,
        name: ruleData.name || '',
        findPattern: ruleData.findPattern || '',
        replacePattern: ruleData.replacePattern || '',
        isEnabled: ruleData.isEnabled ?? true,
        runOnEdit: ruleData.runOnEdit ?? false,
        substituteRegex: ruleData.substituteRegex ?? true,
        minDepth: ruleData.minDepth ?? 0,
        maxDepth: ruleData.maxDepth ?? 0,
        affectFlags: ruleData.affectFlags,
        sortOrder: ruleData.sortOrder ?? rules.value.length
      });
      if (response.rule) {
        rules.value.push(response.rule);
      }
      message.success('规则已创建');
    }
    showEditModal.value = false;
    editingRegex.value = null;
  } finally {
    // loading状态在modal关闭后自动处理
  }
};

// 导入
const handleImport = () => {
  showImportModal.value = true;
};

const handleImportFile = async (options: { file: UploadFileInfo; fileList: UploadFileInfo[]; event?: Event }) => {
  const file = options.file.file;
  if (!file) return;

  try {
    const arrayBuffer = await file.arrayBuffer();
    const fileContent = new Uint8Array(arrayBuffer);
    await regexRuleClient.importRegexRules({
      fileContent: fileContent,
      fileName: file.name
    });
    message.success('导入成功');
    showImportModal.value = false;
    await loadRules();
  } finally {
    // Modal关闭后自动处理
  }
};

// 导出
const handleExport = async () => {
  try {
    const response = await regexRuleClient.exportRegexRules({
      presetId: currentPresetId.value
    });
    // 下载文件
    const blob = new Blob([new Uint8Array(response.fileContent).buffer], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = response.fileName || 'regex-rules.json';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    message.success('导出成功');
  } finally {
    // 下载完成后自动处理
  }
};
</script>

<style scoped>
.regex-view { display: flex; flex-direction: column; height: calc(100vh - 64px - 48px); }
.toolbar { display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; }
.search-input { flex: 1; max-width: 300px; min-width: 200px; }
.regex-container { flex: 1; }
.regex-list { display: flex; flex-direction: column; gap: 12px; padding-bottom: 20px; }
.regex-card { background: var(--bg-card); border: 1px solid var(--border-color); border-radius: 12px; padding: 16px; transition: all var(--transition-normal); }
.regex-card:hover { border-color: var(--border-glow); box-shadow: var(--glow-soft); }
.regex-card.disabled { opacity: 0.6; }
.regex-header { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; }
.regex-name { flex: 1; font-size: 15px; font-weight: 600; margin: 0; color: var(--text-primary); }
.regex-tags { display: flex; gap: 4px; }
.regex-actions { display: flex; gap: 4px; opacity: 0; transition: opacity var(--transition-fast); }
.regex-card:hover .regex-actions { opacity: 1; }
.regex-body { display: flex; flex-direction: column; gap: 8px; margin-bottom: 12px; }
.regex-pattern, .regex-replacement { display: flex; align-items: center; gap: 8px; font-size: 13px; }
.regex-pattern .label, .regex-replacement .label { color: var(--text-tertiary); min-width: 40px; }
.regex-pattern code, .regex-replacement code { background: var(--bg-tertiary); padding: 4px 8px; border-radius: 4px; font-family: var(--font-mono, monospace); color: var(--color-primary); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.regex-footer { display: flex; align-items: center; gap: 16px; font-size: 12px; color: var(--text-tertiary); }
.regex-order, .regex-depth { display: flex; align-items: center; gap: 4px; }
.regex-run-on-edit { margin-left: auto; }
.empty-state { padding: 60px 20px; }
.upload-content { display: flex; flex-direction: column; align-items: center; padding: 40px 20px; }
.upload-icon { color: var(--color-primary); margin-bottom: 16px; }
.upload-text { font-size: 16px; color: var(--text-primary); margin: 0 0 8px; }
.upload-hint { font-size: 13px; color: var(--text-tertiary); margin: 0; }
.regex-list-enter-active, .regex-list-leave-active { transition: all 0.3s ease; }
.regex-list-enter-from, .regex-list-leave-to { opacity: 0; transform: translateX(-20px); }
</style>
