<template>
  <div class="regex-editor-view">
    <!-- 编辑器顶部栏 -->
    <div class="editor-top-bar">
      <n-button quaternary @click="backToChat">
        <template #icon>
          <n-icon><ArrowBackOutline /></n-icon>
        </template>
        返回对话
      </n-button>
      <span class="editor-title">{{ isCreateMode ? '新建正则' : `编辑正则: ${ruleName}` }}</span>
      <div class="editor-actions">
        <n-button size="small" @click="scrollToTest">
          <template #icon>
            <n-icon><FlaskOutline /></n-icon>
          </template>
          测试
        </n-button>
        <n-button size="small" :disabled="saving" @click="handleSaveAs">另存为</n-button>
        <n-button type="primary" size="small" :loading="saving" :disabled="pageLoading" @click="handleSave">
          <template #icon>
            <n-icon><SaveOutline /></n-icon>
          </template>
          保存
        </n-button>
      </div>
    </div>

    <!-- 编辑器内容 -->
    <n-spin :show="pageLoading" description="加载中..." style="flex: 1; overflow: hidden;">
    <div class="editor-scroll-area">
      <div class="editor-content">
        <!-- 基本信息 -->
        <section class="editor-section">
          <div class="section-header">
            <h3 class="section-title">基本信息</h3>
          </div>
          <div class="form-row">
            <div class="form-item flex-1">
              <label class="form-label">规则名称</label>
              <n-input
                v-model:value="ruleName"
                placeholder="输入规则名称"
                @update:value="onFieldChange"
              />
            </div>
            <div class="form-item">
              <label class="form-label">启用状态</label>
              <n-switch
                v-model:value="ruleEnabled"
                @update:value="onFieldChange"
              >
                <template #checked>已启用</template>
                <template #unchecked>已禁用</template>
              </n-switch>
            </div>
          </div>
        </section>

        <!-- 匹配规则 -->
        <section class="editor-section">
          <div class="section-header">
            <h3 class="section-title">匹配规则</h3>
          </div>

          <!-- 查找模式 -->
          <div class="form-item">
            <label class="form-label">查找模式 (正则表达式)</label>
            <n-input
              v-model:value="findPattern"
              type="textarea"
              placeholder="输入正则表达式，例如：\(OOC:.*?\)"
              :autosize="{ minRows: 2, maxRows: 6 }"
              class="mono-input"
              @update:value="onPatternChange"
            />
          </div>

          <!-- 正则标志 -->
          <div class="form-item">
            <label class="form-label">正则标志</label>
            <div class="flags-row">
              <n-checkbox
                v-model:checked="flagG"
                @update:checked="onPatternChange"
              >
                <span class="flag-label">g</span>
                <span class="flag-desc">全局匹配</span>
              </n-checkbox>
              <n-checkbox
                v-model:checked="flagI"
                @update:checked="onPatternChange"
              >
                <span class="flag-label">i</span>
                <span class="flag-desc">忽略大小写</span>
              </n-checkbox>
              <n-checkbox
                v-model:checked="flagM"
                @update:checked="onPatternChange"
              >
                <span class="flag-label">m</span>
                <span class="flag-desc">多行模式</span>
              </n-checkbox>
              <n-checkbox
                v-model:checked="flagS"
                @update:checked="onPatternChange"
              >
                <span class="flag-label">s</span>
                <span class="flag-desc">单行模式</span>
              </n-checkbox>
              <n-checkbox
                v-model:checked="flagU"
                @update:checked="onPatternChange"
              >
                <span class="flag-label">u</span>
                <span class="flag-desc">Unicode</span>
              </n-checkbox>
            </div>
          </div>

          <!-- 替换为 -->
          <div class="form-item">
            <label class="form-label">替换为</label>
            <n-input
              v-model:value="replacePattern"
              type="textarea"
              placeholder="替换文本（留空表示删除匹配内容）"
              :autosize="{ minRows: 2, maxRows: 6 }"
              class="mono-input"
              @update:value="onPatternChange"
            />
          </div>
        </section>

        <!-- 作用域设置 -->
        <section class="editor-section">
          <div class="section-header">
            <h3 class="section-title">作用域设置</h3>
          </div>

          <!-- 作用于 -->
          <div class="form-item">
            <label class="form-label">作用于</label>
            <div class="scope-row">
              <n-checkbox
                v-model:checked="scopeAiOutput"
                @update:checked="onFieldChange"
              >
                AI 输出
              </n-checkbox>
              <n-checkbox
                v-model:checked="scopeUserInput"
                @update:checked="onFieldChange"
              >
                用户输入
              </n-checkbox>
              <n-checkbox
                v-model:checked="scopeSystemPrompt"
                @update:checked="onFieldChange"
              >
                系统提示
              </n-checkbox>
            </div>
          </div>

          <div class="form-row">
            <!-- 运行时机 -->
            <div class="form-item flex-1">
              <label class="form-label">运行时机</label>
              <n-select
                v-model:value="runTiming"
                :options="runTimingOptions"
                @update:value="onFieldChange"
              />
            </div>

            <!-- 最小目标 -->
            <div class="form-item">
              <label class="form-label">最小目标</label>
              <n-input-number
                v-model:value="minTarget"
                placeholder="可选"
                :min="0"
                size="medium"
                clearable
                style="width: 120px;"
                @update:value="onFieldChange"
              />
            </div>

            <!-- 最大目标 -->
            <div class="form-item">
              <label class="form-label">最大目标</label>
              <n-input-number
                v-model:value="maxTarget"
                placeholder="可选"
                :min="0"
                size="medium"
                clearable
                style="width: 120px;"
                @update:value="onFieldChange"
              />
            </div>
          </div>
        </section>

        <!-- 实时测试 -->
        <section ref="testSectionRef" class="editor-section test-section">
          <div class="section-header">
            <h3 class="section-title">实时测试</h3>
            <n-button size="tiny" quaternary @click="clearTestInput">清空</n-button>
          </div>

          <!-- 测试输入 -->
          <div class="form-item">
            <label class="form-label">测试输入</label>
            <n-input
              v-model:value="testInput"
              type="textarea"
              placeholder="在此粘贴待测试的文本..."
              :autosize="{ minRows: 4, maxRows: 10 }"
              class="mono-input"
            />
          </div>

          <!-- 匹配结果 -->
          <div class="form-item">
            <label class="form-label">匹配结果</label>
            <div class="test-result-box">
              <div v-if="testInput && findPattern" class="result-content mono-text">
                <!-- 无匹配时 -->
                <template v-if="matchCount === 0">
                  <span class="no-match-text">无匹配项</span>
                </template>
                <!-- 有匹配时，显示高亮结果 -->
                <template v-else>
                  <span v-html="highlightedResult"></span>
                </template>
              </div>
              <div v-else class="result-placeholder">
                输入查找模式和测试文本后自动显示结果
              </div>
            </div>
          </div>

          <!-- 替换结果 -->
          <div v-if="matchCount > 0" class="form-item">
            <label class="form-label">替换后结果</label>
            <div class="test-result-box replaced-box">
              <div class="result-content mono-text">{{ replacedResult }}</div>
            </div>
          </div>

          <!-- 统计信息 -->
          <div v-if="testInput && findPattern" class="test-stats">
            <span class="stat-item">
              <span class="stat-label">匹配数:</span>
              <span :class="['stat-value', matchCount > 0 ? 'has-match' : 'no-match']">
                {{ matchCount }}
              </span>
            </span>
            <span v-if="matchCount > 0" class="stat-item">
              <span class="stat-label">替换后长度变化:</span>
              <span class="stat-value">{{ lengthDelta >= 0 ? `+${lengthDelta}` : lengthDelta }}</span>
            </span>
            <span v-if="regexError" class="stat-item error-stat">
              <span class="stat-label">错误:</span>
              <span class="stat-value error-text">{{ regexError }}</span>
            </span>
          </div>
        </section>
      </div>
    </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { h, ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router';
import {
  NButton,
  NIcon,
  NInput,
  NInputNumber,
  NSwitch,
  NCheckbox,
  NSelect,
  NSpin,
  useMessage,
  useDialog
} from 'naive-ui';
import {
  ArrowBackOutline,
  SaveOutline,
  FlaskOutline
} from '@vicons/ionicons5';
import { useRegexRuleStore } from '@/stores/regexRule';
import { useCharacterStore } from '@/stores/character';
import { useUserStore } from '@/stores/user';
import { REGEX_RULE_TYPE } from '@/utils/constants';
import type { RegexRuleType } from '@/utils/constants';
import type { RegexRule } from '@/gen/muse/regex_pb';

const route = useRoute();
const router = useRouter();
const message = useMessage();
const dialog = useDialog();
const regexRuleStore = useRegexRuleStore();
const characterStore = useCharacterStore();
const userStore = useUserStore();

// ====== 模式判断 ======
const isCreateMode = computed(() => route.name === 'RegexCreator');

// ====== 路由参数 ======
const regexId = computed(() => {
  if (isCreateMode.value) return 0;
  return Number(route.params.id);
});
const pageLoading = ref(true);

// ====== 基本信息 ======
const ruleName = ref('');
const ruleEnabled = ref(true);

// ====== 匹配规则 ======
const findPattern = ref('');
const replacePattern = ref('');

// ====== 正则标志 ======
const flagG = ref(true);
const flagI = ref(false);
const flagM = ref(false);
const flagS = ref(false);
const flagU = ref(false);

// ====== 作用域设置 ======
const scopeAiOutput = ref(true);
const scopeUserInput = ref(false);
const scopeSystemPrompt = ref(false);
const runTiming = ref('after');
const minTarget = ref<number | null>(null);
const maxTarget = ref<number | null>(null);

// 原始规则数据（用于判断是否有变更）
let originalRule: RegexRule | null = null;

const runTimingOptions = [
  { label: '每次生成后', value: 'after' },
  { label: '每次生成前', value: 'before' },
  { label: '手动触发', value: 'manual' }
];

// ====== 实时测试 ======
const testInput = ref('这是一段测试文本 (OOC: 这是OOC内容) 继续正文 (OOC: 另一个OOC)');
const testSectionRef = ref<HTMLElement | null>(null);

// ====== 保存状态 ======
const saving = ref(false);
const isDirty = ref(false);

// ====== 正则标志字符串 ======
const flagsString = computed(() => {
  let flags = '';
  if (flagG.value) flags += 'g';
  if (flagI.value) flags += 'i';
  if (flagM.value) flags += 'm';
  if (flagS.value) flags += 's';
  if (flagU.value) flags += 'u';
  return flags;
});

// ====== 正则错误 ======
const regexError = ref<string | null>(null);

// ====== 构建正则对象 ======
const buildRegex = (): RegExp | null => {
  if (!findPattern.value) {
    regexError.value = null;
    return null;
  }
  try {
    const regex = new RegExp(findPattern.value, flagsString.value);
    regexError.value = null;
    return regex;
  } catch (e: unknown) {
    regexError.value = e instanceof Error ? e.message : '无效的正则表达式';
    return null;
  }
};

// ====== 匹配数量 ======
const matchCount = computed(() => {
  if (!testInput.value || !findPattern.value) return 0;
  const regex = buildRegex();
  if (!regex) return 0;

  // 为了统计匹配数，需要全局模式
  const globalRegex = new RegExp(regex.source, regex.flags.includes('g') ? regex.flags : regex.flags + 'g');
  const matches = testInput.value.match(globalRegex);
  return matches ? matches.length : 0;
});

// ====== 高亮匹配结果 ======
const highlightedResult = computed(() => {
  if (!testInput.value || !findPattern.value) return '';
  const regex = buildRegex();
  if (!regex) return escapeHtml(testInput.value);

  const globalRegex = new RegExp(regex.source, regex.flags.includes('g') ? regex.flags : regex.flags + 'g');
  return escapeHtml(testInput.value).replace(
    new RegExp(globalRegex.source, globalRegex.flags),
    (match) => `<mark class="highlight-match">${escapeHtml(match)}</mark>`
  );
});

// ====== 替换后结果 ======
const replacedResult = computed(() => {
  if (!testInput.value || !findPattern.value) return '';
  const regex = buildRegex();
  if (!regex) return testInput.value;

  const globalRegex = new RegExp(regex.source, regex.flags.includes('g') ? regex.flags : regex.flags + 'g');
  return testInput.value.replace(globalRegex, replacePattern.value);
});

// ====== 长度变化 ======
const lengthDelta = computed(() => {
  if (!testInput.value) return 0;
  return replacedResult.value.length - testInput.value.length;
});

// ====== 工具函数 ======
const escapeHtml = (str: string): string => {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
};

// ====== 事件处理 ======
const backToChat = () => {
  router.push('/');
};

const onFieldChange = () => {
  isDirty.value = true;
};

const onPatternChange = () => {
  isDirty.value = true;
  // 触发正则重新编译（通过 computed 自动处理）
};

const handleSave = async () => {
  if (!ruleName.value.trim()) {
    message.warning('请输入规则名称');
    return;
  }
  if (!findPattern.value.trim()) {
    message.warning('请输入查找模式');
    return;
  }

  // 验证正则表达式
  try {
    new RegExp(findPattern.value, flagsString.value);
  } catch {
    message.error('正则表达式语法错误，请检查查找模式');
    return;
  }

  saving.value = true;
  try {
    const affectFlags = {
      $typeName: 'muse.RegexAffectFlags' as const,
      aiOutput: scopeAiOutput.value,
      userInput: scopeUserInput.value,
      prompt: scopeSystemPrompt.value,
      slashCommand: false,
      worldInfo: false
    };

    if (isCreateMode.value) {
      // 创建模式：调用 addRule
      const type = route.query.tab as RegexRuleType;
      const presetId = type === REGEX_RULE_TYPE.PRESET ? (userStore.currentUser?.activePresetId ?? 0) : 0;
      const characterId = type === REGEX_RULE_TYPE.CHARACTER ? (characterStore.curCharId ?? 0) : 0;

      const newRule = await regexRuleStore.addRule({
        presetId,
        characterId,
        name: ruleName.value,
        findPattern: findPattern.value,
        replacePattern: replacePattern.value,
        isEnabled: ruleEnabled.value,
        runOnEdit: runTiming.value === 'after',
        substituteRegex: true,
        minDepth: minTarget.value ?? 0,
        maxDepth: maxTarget.value ?? 0,
        affectFlags,
        sortOrder: regexRuleStore.globalRules.length
      });

      if (newRule) {
        isDirty.value = false;
        message.success('规则创建成功');
        router.replace({ name: 'RegexEditor', params: { id: newRule.id } });
      }
    } else {
      // 编辑模式：调用 updateRule
      await regexRuleStore.updateRule(regexId.value, {
        name: ruleName.value,
        findPattern: findPattern.value,
        replacePattern: replacePattern.value,
        isEnabled: ruleEnabled.value,
        runOnEdit: runTiming.value === 'after',
        substituteRegex: true,
        minDepth: minTarget.value ?? 0,
        maxDepth: maxTarget.value ?? 0,
        affectFlags,
        sortOrder: originalRule?.sortOrder ?? 0
      });

      isDirty.value = false;
      message.success('规则保存成功');
    }
  } catch {
    message.error('保存失败，请重试');
  } finally {
    saving.value = false;
  }
};

const handleSaveAs = () => {
  const newName = ref(ruleName.value + ' (副本)');
  dialog.create({
    title: '另存为新规则',
    content: () =>
      h(NInput, {
        value: newName.value,
        'onUpdate:value': (v: string) => { newName.value = v; },
        placeholder: '输入新规则名称'
      }),
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      if (!newName.value.trim()) {
        message.warning('请输入规则名称');
        return false;
      }
      try {
        const affectFlags = {
          $typeName: 'muse.RegexAffectFlags' as const,
          aiOutput: scopeAiOutput.value,
          userInput: scopeUserInput.value,
          prompt: scopeSystemPrompt.value,
          slashCommand: false,
          worldInfo: false
        };

        const newRule = await regexRuleStore.addRule({
          presetId: originalRule?.presetId ?? 0,
          name: newName.value,
          findPattern: findPattern.value,
          replacePattern: replacePattern.value,
          isEnabled: ruleEnabled.value,
          runOnEdit: runTiming.value === 'after',
          substituteRegex: true,
          minDepth: minTarget.value ?? 0,
          maxDepth: maxTarget.value ?? 0,
          affectFlags,
          sortOrder: (originalRule?.sortOrder ?? 0) + 1
        });

        if (newRule) {
          message.success('规则已另存为: ' + newName.value);
          router.push(`/regex/${newRule.id}`);
        }
      } catch {
        message.error('另存为失败');
        return false;
      }
    }
  });
};

const scrollToTest = () => {
  testSectionRef.value?.scrollIntoView({ behavior: 'smooth' });
};

const clearTestInput = () => {
  testInput.value = '';
};

// ====== 加载数据 ======
const loadRule = async (id: number) => {
  // 创建模式：使用默认值，无需加载
  if (isCreateMode.value) {
    pageLoading.value = false;
    return;
  }

  if (!id || id <= 0) return;
  pageLoading.value = true;
  try {
    // 先确保已加载规则列表
    if (!regexRuleStore.loaded) {
      await regexRuleStore.loadRules();
    }
    const allRules = [
      ...regexRuleStore.globalRules,
      ...regexRuleStore.presetRules,
      ...regexRuleStore.characterRules
    ];
    const rule = allRules.find(r => r.id === id);
    if (!rule) {
      message.error('规则不存在');
      router.push('/');
      return;
    }
    originalRule = rule;
    ruleName.value = rule.name;
    ruleEnabled.value = rule.isEnabled;
    findPattern.value = rule.findPattern;
    replacePattern.value = rule.replacePattern;
    minTarget.value = rule.minDepth || null;
    maxTarget.value = rule.maxDepth || null;

    // 恢复作用域标志
    scopeAiOutput.value = rule.affectFlags?.aiOutput ?? true;
    scopeUserInput.value = rule.affectFlags?.userInput ?? false;
    scopeSystemPrompt.value = rule.affectFlags?.prompt ?? false;

    // 恢复运行时机
    runTiming.value = rule.runOnEdit ? 'after' : 'before';

    isDirty.value = false;
  } catch {
    message.error('加载规则失败');
  } finally {
    pageLoading.value = false;
  }
};

onMounted(() => {
  loadRule(regexId.value);
});

watch(regexId, (newId) => {
  if (newId && newId > 0) {
    loadRule(newId);
  }
});

// 离开页面时检查未保存的更改
onBeforeRouteLeave((_to, _from, next) => {
  if (!isDirty.value) {
    next();
    return;
  }
  dialog.warning({
    title: '未保存的更改',
    content: '你有未保存的更改，确定要离开吗？',
    positiveText: '离开',
    negativeText: '留下',
    onPositiveClick: () => {
      isDirty.value = false;
      next();
    },
    onNegativeClick: () => {
      next(false);
    }
  });
});
</script>

<style scoped>
.regex-editor-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 编辑器顶部栏 */
.editor-top-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 24px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  background: var(--bg-primary);
}

.editor-title {
  flex: 1;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.editor-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

/* 滚动区域 */
.editor-scroll-area {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.editor-content {
  max-width: 960px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 分区块 */
.editor-section {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

/* 表单 */
.form-row {
  display: flex;
  gap: 16px;
  align-items: flex-end;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.flex-1 {
  flex: 1;
  min-width: 0;
}

.form-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

/* 等宽字体输入框 */
.mono-input :deep(textarea),
.mono-input :deep(input) {
  font-family: 'JetBrains Mono', 'Fira Code', 'Source Code Pro', monospace !important;
  font-size: 13px;
  line-height: 1.6;
}

.mono-text {
  font-family: 'JetBrains Mono', 'Fira Code', 'Source Code Pro', monospace;
  font-size: 13px;
  line-height: 1.6;
}

/* 正则标志行 */
.flags-row {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.flag-label {
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-weight: 600;
  color: var(--color-primary);
  margin-right: 4px;
}

.flag-desc {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 作用域行 */
.scope-row {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

/* 测试区域 */
.test-section {
  background: var(--bg-secondary);
}

.test-result-box {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-primary);
  padding: 12px 16px;
  min-height: 80px;
}

.replaced-box {
  border-color: var(--color-primary);
  border-style: dashed;
}

.result-content {
  white-space: pre-wrap;
  word-break: break-all;
  color: var(--text-primary);
}

.result-placeholder {
  color: var(--text-tertiary);
  font-size: 13px;
  font-style: italic;
}

.no-match-text {
  color: var(--text-tertiary);
  font-style: italic;
}

/* 匹配高亮 */
.result-content :deep(.highlight-match) {
  background: #fef3cd;
  color: #856404;
  border-radius: 2px;
  padding: 1px 2px;
}

/* 暗色主题下的高亮 */
:root[data-theme='dark'] .result-content :deep(.highlight-match) {
  background: #5a4a00;
  color: #ffd60a;
}

/* 统计信息 */
.test-stats {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
  padding: 10px 0 0 0;
  border-top: 1px solid var(--border-color);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
}

.stat-label {
  color: var(--text-tertiary);
}

.stat-value {
  font-weight: 500;
  color: var(--text-primary);
}

.stat-value.has-match {
  color: var(--color-primary);
}

.stat-value.no-match {
  color: var(--text-tertiary);
}

.error-stat .stat-value.error-text {
  color: #e53935;
  font-weight: 400;
}

/* 响应式 */
@media (max-width: 767px) {
  .editor-top-bar {
    padding: 10px 16px;
  }

  .editor-scroll-area {
    padding: 16px 12px;
  }

  .editor-section {
    padding: 16px;
  }

  .form-row {
    flex-direction: column;
  }

  .flags-row {
    gap: 12px;
  }

  .scope-row {
    flex-direction: column;
    gap: 8px;
  }

  .test-stats {
    flex-direction: column;
    gap: 8px;
  }
}
</style>
