<template>
  <div class="preset-editor-view">
    <!-- 加载状态 -->
    <n-spin :show="pageLoading" description="加载中..." class="loading-spin">

    <!-- 编辑器顶部栏 -->
    <div class="editor-top-bar">
      <n-button quaternary @click="backToChat">
        <template #icon>
          <n-icon><ArrowBackOutline /></n-icon>
        </template>
        返回对话
      </n-button>
      <span class="editor-title">编辑预设: {{ presetName }}</span>
      <div class="editor-actions">
        <n-button size="small" :disabled="saving" @click="handleSaveAs">另存为</n-button>
        <n-button type="primary" size="small" :loading="saving" @click="handleSave">保存</n-button>
      </div>
    </div>

    <!-- 编辑器内容区域 -->
    <div class="editor-body">
      <div class="editor-content">
        <!-- 预设名称 -->
        <div class="field-group">
          <label class="field-label">预设名称</label>
          <n-input
            v-model:value="presetName"
            placeholder="输入预设名称"
            size="large"
          />
        </div>

        <!-- Prompt 排列 -->
        <div class="section-header">
          <span class="section-title">Prompt 排列</span>
          <span class="section-desc">拖拽调整顺序，展开编辑内容</span>
        </div>

        <draggable
          v-model="prompts"
          item-key="id"
          handle=".drag-handle"
          animation="200"
          ghost-class="prompt-ghost"
          class="prompt-list"
        >
          <template #item="{ element: prompt, index }">
          <div
            class="prompt-card"
            :class="{ expanded: prompt.expanded }"
          >
            <!-- 折叠态头部 -->
            <div class="prompt-header" @click="togglePrompt(index)">
              <div class="prompt-header-left">
                <n-icon class="drag-handle" :size="18" @click.stop>
                  <ReorderTwoOutline />
                </n-icon>
                <n-checkbox
                  :checked="prompt.isEnabled"
                  @update:checked="(val: boolean) => prompt.isEnabled = val"
                  @click.stop
                />
                <span class="prompt-name">{{ prompt.name || '未命名 Prompt' }}</span>
                <n-tag :type="roleTagType(prompt.role)" size="small" round>
                  {{ roleLabel(prompt.role) }}
                </n-tag>
                <n-tag v-if="prompt.identifier"
                       :color="{ borderColor: '#9B7ED9', textColor: '#9B7ED9' }"
                       size="small"
                       round>
                  系统内容引用
                </n-tag>
              </div>
              <div class="prompt-header-right">
                <n-button
                  quaternary
                  circle
                  size="tiny"
                  @click.stop="togglePrompt(index)"
                >
                  <template #icon>
                    <n-icon :size="16">
                      <ChevronDownOutline v-if="!prompt.expanded" />
                      <ChevronUpOutline v-else />
                    </n-icon>
                  </template>
                </n-button>
                <n-button
                  quaternary
                  circle
                  size="tiny"
                  @click.stop="removePrompt(index)"
                >
                  <template #icon>
                    <n-icon :size="16"><CloseOutline /></n-icon>
                  </template>
                </n-button>
              </div>
            </div>

            <!-- 展开态内容 -->
            <div v-if="prompt.expanded" class="prompt-body">
              <div class="prompt-fields">
                <div class="prompt-field-row">
                  <div class="prompt-field">
                    <label class="field-label-sm">名称</label>
                    <n-input
                      v-model:value="prompt.name"
                      placeholder="Prompt 名称"
                      size="small"
                    />
                  </div>
                  <div class="prompt-field">
                    <label class="field-label-sm">角色</label>
                    <n-select
                      v-model:value="prompt.role"
                      :options="roleOptions"
                      size="small"
                    />
                  </div>
                </div>
              </div>
              <div class="prompt-editor">
                <n-input
                  v-model:value="prompt.content"
                  type="textarea"
                  :placeholder="prompt.identifier ? '此提示词的内容是从其他地方提取的, 无法在此处进行编辑。\n来源: '
                  + PromptItemIdentifier[prompt.identifier]  : '输入 Prompt 内容，支持 {{char}}、{{user}} 等宏变量...'"
                  :autosize="{ minRows: 6, maxRows: 20 }"
                  :disabled="!!prompt.identifier"
                  class="prompt-textarea"
                />
              </div>
            </div>
          </div>
          </template>
        </draggable>

        <!-- 添加 Prompt 按钮 -->
        <n-button dashed block @click="addPrompt" class="add-prompt-btn">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          添加 Prompt 项
        </n-button>

        <!-- 预设内正则规则 -->
        <div class="section-header" style="margin-top: 32px;">
          <span class="section-title">预设内正则规则</span>
          <span class="section-desc">此预设关联的正则替换规则</span>
        </div>

        <div class="regex-list">
          <div
            v-for="(rule, index) in regexRules"
            :key="index"
            class="regex-item"
          >
            <n-checkbox
              :checked="rule.isEnabled"
              @update:checked="(val: boolean) => rule.isEnabled = val"
            />
            <span class="regex-name">{{ rule.name }}</span>
            <code class="regex-preview">{{ rule.findPattern }}</code>
            <span class="regex-arrow">→</span>
            <span class="regex-replace">{{ rule.replacePattern || '(空)' }}</span>
            <n-button
              quaternary
              circle
              size="tiny"
              @click="editRegex(index)"
            >
              <template #icon>
                <n-icon :size="14"><CreateOutline /></n-icon>
              </template>
            </n-button>
            <n-button
              quaternary
              circle
              size="tiny"
              @click="removeRegex(index)"
            >
              <template #icon>
                <n-icon :size="14"><CloseOutline /></n-icon>
              </template>
            </n-button>
          </div>

          <n-empty
            v-if="regexRules.length === 0"
            description="暂无关联正则规则"
            size="small"
            class="empty-regex"
          />
        </div>

        <n-button dashed block size="small" @click="addRegex" class="add-regex-btn">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
          添加正则
        </n-button>
      </div>
    </div>

    <!-- 正则规则编辑模态框 -->
    <n-modal
      v-model:show="showRegexModal"
      preset="card"
      :title="editingRegexIndex >= 0 ? '编辑正则规则' : '添加正则规则'"
      :style="{ width: '550px', maxWidth: '90vw' }"
    >
      <n-form ref="regexFormRef" :model="regexFormData" label-placement="top">
        <n-form-item label="规则名称" path="name">
          <n-input v-model:value="regexFormData.name" placeholder="输入规则名称" />
        </n-form-item>

        <n-form-item label="匹配模式" path="findPattern">
          <n-input v-model:value="regexFormData.findPattern" placeholder="输入正则表达式" style="font-family: monospace;" />
        </n-form-item>

        <n-form-item label="替换内容">
          <n-input v-model:value="regexFormData.replacePattern" placeholder="替换文本（留空表示删除匹配内容）" />
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
              <n-input-number v-model:value="regexFormData.sortOrder" :min="0" :max="1000" style="width: 100%;" />
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>

      <div class="modal-footer">
        <n-button @click="showRegexModal = false">取消</n-button>
        <n-button type="primary" @click="handleRegexSave">保存</n-button>
      </div>
    </n-modal>

    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import {
  NButton,
  NIcon,
  NInput,
  NSelect,
  NCheckbox,
  NTag,
  NEmpty,
  NModal,
  NForm,
  NFormItem,
  NInputNumber,
  NCheckboxGroup,
  NSpace,
  NGrid,
  NGi,
  NSpin,
  useMessage,
  useDialog
} from 'naive-ui';
import type { FormInst } from 'naive-ui';
import {
  ArrowBackOutline,
  AddOutline,
  CloseOutline,
  CreateOutline,
  ReorderTwoOutline,
  ChevronDownOutline,
  ChevronUpOutline
} from '@vicons/ionicons5';
import draggable from 'vuedraggable';
import { usePresetStore } from '@/stores/preset';
import {presetClient, regexRuleClient} from '@/api/client';
import { Role, InjectionPosition, PromptItemIdentifier } from '@/gen/muse/common_pb';
import type { PromptItem } from '@/gen/muse/preset_pb';
import type { RegexRule, RegexAffectFlags } from '@/gen/muse/regex_pb';

const route = useRoute();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const presetStore = usePresetStore();

const presetId = computed(() => Number(route.params.id));
const pageLoading = ref(true);
const saving = ref(false);

// 预设基础信息
const presetName = ref('');
const presetVersion = ref(0n);

// Prompt 角色选项
const roleOptions = [
  { label: '系统 (System)', value: Role.System },
  { label: '用户 (User)', value: Role.User },
  { label: '助手 (Assistant)', value: Role.Assistant }
];

// Prompt 项（带前端展开状态）
interface PromptItemUI extends PromptItem {
  expanded: boolean;
}

const prompts = ref<PromptItemUI[]>([]);
const regexRules = ref<RegexRule[]>([]);
// 记录加载时的原始正则规则ID，用于保存时对比差异
const originalRegexRuleIds = ref<Set<number>>(new Set());

// =====================
// 加载数据
// =====================

const loadPreset = async () => {
  pageLoading.value = true;
  try {
    const preset = await presetStore.fetchPreset(presetId.value);
    if (!preset || !preset.preset) {
      message.error('预设不存在');
      router.push('/');
      return;
    }
    presetName.value = preset.preset?.name ?? '';
    presetVersion.value = preset.preset?.version ?? 0n;
    prompts.value = (preset.promptItems || []).map(item => ({
      ...item,
      expanded: false
    }));

    // 通过后端接口获取预设关联的正则规则
    const regexResp = await regexRuleClient.listPresetRegexRules({ presetId: presetId.value });
    regexRules.value = regexResp.rules;
    originalRegexRuleIds.value = new Set(regexResp.rules.filter(r => r.id > 0).map(r => r.id));
  } catch {
    message.error('加载预设失败');
  } finally {
    pageLoading.value = false;
  }
};

onMounted(() => {
  loadPreset();
});

// 路由参数变化时重新加载
watch(presetId, (newId) => {
  if (newId) {
    loadPreset();
  }
});

// =====================
// Prompt 操作
// =====================

// 角色标签类型映射
const roleTagType = (role: Role) => {
  const map: Record<number, 'info' | 'success' | 'warning'> = {
    [Role.System]: 'info',
    [Role.User]: 'success',
    [Role.Assistant]: 'warning'
  };
  return map[role] || 'info';
};

// 角色标签文字
const roleLabel = (role: Role) => {
  const map: Record<number, string> = {
    [Role.System]: '系统',
    [Role.User]: '用户',
    [Role.Assistant]: '助手'
  };
  return map[role] || '未知';
};

// 展开/折叠 Prompt
const togglePrompt = (index: number) => {
  const item = prompts.value[index];
  if (item) {
    item.expanded = !item.expanded;
  }
};

// 添加 Prompt
const addPrompt = () => {
  prompts.value.push({
    $typeName: 'muse.PromptItem',
    id: 0,
    presetId: presetId.value,
    identifier: PromptItemIdentifier.PromptItemIdentifierUnspecified,
    name: '',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: false,
    createdAt: 0n,
    updatedAt: 0n,
    expanded: true
  } as PromptItemUI);
};

// 删除 Prompt
const removePrompt = (index: number) => {
  if (index >= prompts.value.length) {
    return
  }
  presetClient.deletePromptItem({ id: prompts.value[index]!.id });
  prompts.value.splice(index, 1);
};

// =====================
// 正则规则操作
// =====================

const showRegexModal = ref(false);
const editingRegexIndex = ref(-1);
const regexFormRef = ref<FormInst | null>(null);

const regexFormData = ref({
  name: '',
  findPattern: '',
  replacePattern: '',
  isEnabled: true,
  sortOrder: 100,
  affectFlags: {
    userInput: true,
    aiOutput: true,
    slashCommand: false,
    worldInfo: false,
    prompt: false
  }
});

const affectFlagsArray = computed({
  get: () => {
    const flags: string[] = [];
    if (regexFormData.value.affectFlags.userInput) flags.push('userInput');
    if (regexFormData.value.affectFlags.aiOutput) flags.push('aiOutput');
    if (regexFormData.value.affectFlags.slashCommand) flags.push('slashCommand');
    if (regexFormData.value.affectFlags.worldInfo) flags.push('worldInfo');
    if (regexFormData.value.affectFlags.prompt) flags.push('prompt');
    return flags;
  },
  set: (val: string[]) => {
    regexFormData.value.affectFlags = {
      userInput: val.includes('userInput'),
      aiOutput: val.includes('aiOutput'),
      slashCommand: val.includes('slashCommand'),
      worldInfo: val.includes('worldInfo'),
      prompt: val.includes('prompt')
    };
  }
});

const resetRegexForm = () => {
  regexFormData.value = {
    name: '',
    findPattern: '',
    replacePattern: '',
    isEnabled: true,
    sortOrder: 100,
    affectFlags: {
      userInput: true,
      aiOutput: true,
      slashCommand: false,
      worldInfo: false,
      prompt: false
    }
  };
};

const addRegex = () => {
  editingRegexIndex.value = -1;
  resetRegexForm();
  showRegexModal.value = true;
};

const editRegex = (index: number) => {
  editingRegexIndex.value = index;
  const rule = regexRules.value[index];
  if (!rule) return;
  regexFormData.value = {
    name: rule.name,
    findPattern: rule.findPattern,
    replacePattern: rule.replacePattern,
    isEnabled: rule.isEnabled,
    sortOrder: rule.sortOrder,
    affectFlags: {
      userInput: rule.affectFlags?.userInput ?? true,
      aiOutput: rule.affectFlags?.aiOutput ?? true,
      slashCommand: rule.affectFlags?.slashCommand ?? false,
      worldInfo: rule.affectFlags?.worldInfo ?? false,
      prompt: rule.affectFlags?.prompt ?? false
    }
  };
  showRegexModal.value = true;
};

const handleRegexSave = async () => {
  try {
    await regexFormRef.value?.validate();
  } catch {
    return;
  }

  // 验证正则表达式
  try {
    new RegExp(regexFormData.value.findPattern, 'g');
  } catch {
    message.error('正则表达式语法错误');
    return;
  }

  const affectFlags: RegexAffectFlags = {
    $typeName: 'muse.RegexAffectFlags',
    ...regexFormData.value.affectFlags
  };

  if (editingRegexIndex.value >= 0) {
    // 编辑已有规则
    const existing = regexRules.value[editingRegexIndex.value];
    regexRules.value[editingRegexIndex.value] = {
      ...existing,
      name: regexFormData.value.name,
      findPattern: regexFormData.value.findPattern,
      replacePattern: regexFormData.value.replacePattern,
      isEnabled: regexFormData.value.isEnabled,
      sortOrder: regexFormData.value.sortOrder,
      affectFlags
    } as RegexRule;
    message.success('规则已更新');
  } else {
    // 新增规则
    regexRules.value.push({
      $typeName: 'muse.RegexRule',
      id: 0,
      presetId: presetId.value,
      characterId: 0,
      name: regexFormData.value.name,
      findPattern: regexFormData.value.findPattern,
      replacePattern: regexFormData.value.replacePattern,
      isEnabled: regexFormData.value.isEnabled,
      runOnEdit: false,
      substituteRegex: true,
      minDepth: 0,
      maxDepth: 0,
      affectFlags,
      sortOrder: regexFormData.value.sortOrder,
      createdAt: 0n,
      updatedAt: 0n
    } as RegexRule);
    message.success('规则已添加');
  }

  showRegexModal.value = false;
};

// 删除正则
const removeRegex = (index: number) => {
  const rule = regexRules.value[index];
  if (!rule) return;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除规则"${rule.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      regexRules.value.splice(index, 1);
      message.success('规则已删除');
    }
  });
};

// =====================
// 保存操作
// =====================

const handleSave = async () => {
  if (!presetName.value.trim()) {
    message.warning('请输入预设名称');
    return;
  }

  saving.value = true;
  try {
    // 1. 更新预设基础信息
    await presetStore.updatePreset(presetId.value, {
      name: presetName.value,
      temperature: 0,
      topP: 0,
      topK: 0,
      maxTokens: 0,
      frequencyPenalty: 0,
      presencePenalty: 0,
      version: presetVersion.value
    });

    // 重新获取最新版本号
    const refreshedPreset = await presetStore.fetchPreset(presetId.value);
    if (refreshedPreset?.preset) {
      presetVersion.value = refreshedPreset.preset.version;
    }

    // 2. 更新prompt项排序
    const orderedIds = prompts.value.filter(p => p.id > 0).map(p => p.id);
    if (orderedIds.length > 1) {
      // 只需要传递第一个和最后一个ID
      await presetStore.updatePromptItemsOrder(presetId.value);
    }

    // 3. 处理prompt项（新增/更新/删除）
    const existingItems = await presetStore.fetchPromptItems(presetId.value);
    const existingIds = new Set(existingItems.map(item => item.id));
    const currentIds = new Set(prompts.value.filter(p => p.id > 0).map(p => p.id));

    // 删除已被移除的项
    for (const existingId of existingIds) {
      if (!currentIds.has(existingId)) {
        await presetStore.deletePromptItem(existingId);
      }
    }

    // 新增和更新
    for (let i = 0; i < prompts.value.length; i++) {
      const prompt = prompts.value[i];
      if (!prompt) continue;
      const promptData = {
        identifier: prompt.identifier || PromptItemIdentifier.PromptItemIdentifierUnspecified,
        name: prompt.name,
        content: prompt.content,
        role: prompt.role,
        isEnabled: prompt.isEnabled,
        injectionPosition: prompt.injectionPosition,
        injectionDepth: prompt.injectionDepth,
        forbidOverrides: prompt.forbidOverrides,
        sortOrder: i
      };

      if (prompt.id > 0) {
        // 更新已有项
        await presetStore.updatePromptItem(prompt.id, promptData);
      } else {
        // 新增项
        const newItem = await presetStore.addPromptItem(presetId.value, promptData);
        if (newItem) {
          prompts.value[i] = { ...prompt, ...newItem, expanded: prompt.expanded ?? false };
        }
      }
    }

    // 4. 处理正则规则（新增/更新/删除）
    const currentRegexIds = new Set(regexRules.value.filter(r => r.id > 0).map(r => r.id));

    // 删除已被移除的正则规则
    for (const originalId of originalRegexRuleIds.value) {
      if (!currentRegexIds.has(originalId)) {
        await regexRuleClient.deleteRegexRule({ id: originalId });
      }
    }

    // 新增和更新正则规则
    for (let i = 0; i < regexRules.value.length; i++) {
      const rule = regexRules.value[i];
      if (!rule) continue;

      const ruleData = {
        name: rule.name,
        findPattern: rule.findPattern,
        replacePattern: rule.replacePattern,
        isEnabled: rule.isEnabled,
        runOnEdit: rule.runOnEdit ?? false,
        substituteRegex: rule.substituteRegex ?? true,
        minDepth: rule.minDepth ?? 0,
        maxDepth: rule.maxDepth ?? 0,
        affectFlags: rule.affectFlags,
        sortOrder: i
      };

      if (rule.id > 0) {
        // 更新已有规则
        const resp = await regexRuleClient.updateRegexRule({ id: rule.id, ...ruleData });
        if (resp.rule) {
          regexRules.value[i] = resp.rule;
        }
      } else {
        // 新增规则
        const resp = await regexRuleClient.addRegexRule({ presetId: presetId.value, ...ruleData });
        if (resp.rule) {
          regexRules.value[i] = resp.rule;
        }
      }
    }

    // 更新原始ID记录
    originalRegexRuleIds.value = new Set(regexRules.value.filter(r => r.id > 0).map(r => r.id));

    message.success('预设保存成功');
  } catch {
    message.error('保存失败，请重试');
  } finally {
    saving.value = false;
  }
};

// 另存为
const handleSaveAs = () => {
  const newName = ref(presetName.value + ' (副本)');
  dialog.create({
    title: '另存为新预设',
    content: () =>
      h(NInput, {
        value: newName.value,
        'onUpdate:value': (v: string) => { newName.value = v; },
        placeholder: '输入新预设名称'
      }),
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      if (!newName.value.trim()) {
        message.warning('请输入预设名称');
        return false;
      }
      try {
        await presetStore.createPreset({
          name: newName.value,
          temperature: 0,
          topP: 0,
          topK: 0,
          maxTokens: 0,
          frequencyPenalty: 0,
          presencePenalty: 0,
          promptItems: prompts.value.map((p, i) => ({
            identifier: p.identifier || PromptItemIdentifier.PromptItemIdentifierUnspecified,
            name: p.name,
            content: p.content,
            role: p.role,
            isEnabled: p.isEnabled,
            injectionPosition: p.injectionPosition,
            injectionDepth: p.injectionDepth,
            forbidOverrides: p.forbidOverrides,
            sortOrder: i
          }))
        });

        // 重新加载列表获取新创建的预设
        await presetStore.fetchAllPresets();
        const newPreset = presetStore.presets.find(p => p.preset?.name === newName.value);
        
        if (newPreset && newPreset.preset) {
          // 为新预设创建关联的正则规则
          for (let i = 0; i < regexRules.value.length; i++) {
            const rule = regexRules.value[i];
            if (!rule) continue;
            await regexRuleClient.addRegexRule({
              presetId: newPreset.preset.id,
              name: rule.name,
              findPattern: rule.findPattern,
              replacePattern: rule.replacePattern,
              isEnabled: rule.isEnabled,
              runOnEdit: rule.runOnEdit ?? false,
              substituteRegex: rule.substituteRegex ?? true,
              minDepth: rule.minDepth ?? 0,
              maxDepth: rule.maxDepth ?? 0,
              affectFlags: rule.affectFlags,
              sortOrder: i
            });
          }

          message.success('预设已另存为: ' + newName.value);
          router.push(`/preset/${newPreset.preset.id}`);
        }
      } catch {
        message.error('另存为失败');
        return false;
      }
    }
  });
};

// 返回对话
const backToChat = () => {
  router.push('/');
};
</script>

<style scoped>
.preset-editor-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* n-spin 包裹了全部内容，需要让其内部容器也参与 flex 布局，否则高度会无限撑开 */
.loading-spin {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.loading-spin :deep(.n-spin-content) {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.editor-top-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}

.editor-title {
  flex: 1;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
}

.editor-actions {
  display: flex;
  gap: 8px;
}

.editor-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.editor-content {
  max-width: 960px;
  margin: 0 auto;
  width: 100%;
}

.field-group {
  margin-bottom: 24px;
}

.field-label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
}

.field-label-sm {
  display: block;
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.section-header {
  display: flex;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border-color);
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}

.section-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* Prompt 列表 */
.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 12px;
}

.prompt-card {
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  overflow: hidden;
  transition: border-color 150ms;
}

.prompt-card:hover {
  border-color: var(--color-primary-light, var(--border-color));
}

.prompt-card.expanded {
  border-color: var(--color-primary);
}

.prompt-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  min-height: 48px;
  transition: background-color 100ms;
}

.prompt-header:hover {
  background: var(--bg-hover);
}

.prompt-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.drag-handle {
  cursor: grab;
  color: var(--text-tertiary);
  flex-shrink: 0;
}

.drag-handle:active {
  cursor: grabbing;
}

.prompt-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.prompt-header-right {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

.prompt-body {
  padding: 0 16px 16px;
  border-top: 1px solid var(--border-color);
}

.prompt-fields {
  padding-top: 12px;
  margin-bottom: 12px;
}

.prompt-field-row {
  display: grid;
  grid-template-columns: 1fr 200px;
  gap: 12px;
}

.prompt-textarea :deep(textarea) {
  font-family: var(--font-mono);
  font-size: 13px;
  line-height: 1.6;
}

/* 拖拽幽灵样式 */
.prompt-ghost {
  opacity: .5;
  background: var(--color-primary-light, var(--bg-hover));
}

/* 添加 Prompt 按钮 */
.add-prompt-btn {
  margin-bottom: 8px;
}

/* 正则规则列表 */
.regex-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}

.regex-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
}

.regex-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
}

.regex-preview {
  font-family: var(--font-mono);
  font-size: 11px;
  background: var(--bg-tertiary);
  padding: 2px 6px;
  border-radius: 3px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 200px;
}

.regex-arrow {
  color: var(--text-tertiary);
  flex-shrink: 0;
  font-size: 12px;
}

.regex-replace {
  font-size: 12px;
  color: var(--text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-regex {
  padding: 16px 0;
}

.add-regex-btn {
  margin-bottom: 16px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

/* 响应式 */
@media (max-width: 767px) {
  .editor-top-bar {
    padding: 8px 16px;
  }

  .editor-body {
    padding: 16px 12px;
  }

  .prompt-field-row {
    grid-template-columns: 1fr;
  }

  .editor-title {
    font-size: 14px;
  }
}
</style>
