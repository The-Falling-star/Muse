<template>
  <div class="regex-list-tab">
    <!-- 操作按钮 -->
    <div class="action-row">
      <n-button size="small" @click="handleImport">
        <template #icon>
          <n-icon>
            <CloudUploadOutline/>
          </n-icon>
        </template>
        导入
      </n-button>
      <n-button size="small" type="primary" @click="handleCreate">
        <template #icon>
          <n-icon>
            <AddOutline/>
          </n-icon>
        </template>
        新建规则
      </n-button>
    </div>

    <!-- 分类 Tab -->
    <n-tabs v-model:value="activeTab" animated class="rules-tabs" type="segment">
      <n-tab-pane :name="REGEX_RULE_TYPE.GLOBAL" :tab="`全局 (${sortedGlobal.length})`">
        <RuleListSection
            :loading="loading"
            :rules="sortedGlobal"
            empty-text="暂无全局正则规则"
            @edit="openRegexEditor"
            @toggle="toggleEnabled"
            @menu-action="handleMenuSelect"
        />
      </n-tab-pane>
      <n-tab-pane :name="REGEX_RULE_TYPE.PRESET" :tab="`预设 (${sortedPreset.length})`">
        <RuleListSection
            :loading="loading"
            :rules="sortedPreset"
            empty-text="暂无预设正则规则"
            @edit="openRegexEditor"
            @toggle="toggleEnabled"
            @menu-action="handleMenuSelect"
        />
      </n-tab-pane>
      <n-tab-pane :name="REGEX_RULE_TYPE.CHARACTER" :tab="`角色 (${sortedCharacter.length})`">
        <RuleListSection
            :loading="loading"
            :rules="sortedCharacter"
            empty-text="暂无角色正则规则"
            @edit="openRegexEditor"
            @toggle="toggleEnabled"
            @menu-action="handleMenuSelect"
        />
      </n-tab-pane>
    </n-tabs>
  </div>
</template>

<script lang="ts" setup>
import {computed, ref, watch} from 'vue';
import {useRoute, useRouter} from 'vue-router';
import {NButton, NIcon, NTabPane, NTabs, useDialog, useMessage} from 'naive-ui';
import {AddOutline, CloudUploadOutline} from '@vicons/ionicons5';
import type {RegexRuleType} from '@/utils/constants';
import {REGEX_RULE_TYPE} from '@/utils/constants';
import {useRegexRuleStore} from '@/stores/regexRule';
import {useCharacterStore} from '@/stores/character';
import {useUserStore} from '@/stores/user';
import RuleListSection from './RuleListSection.vue';

const router = useRouter();
const route = useRoute();
const dialog = useDialog();
const message = useMessage();
const regexStore = useRegexRuleStore();
const charStore = useCharacterStore();
const userStore = useUserStore();

const activeTab = ref<RegexRuleType>(REGEX_RULE_TYPE.GLOBAL);
const loading = computed(() => regexStore.loading);

// 排序后的三类规则
const sortedGlobal = computed(() =>
    [...regexStore.globalRegex].sort((a, b) => a.sortOrder - b.sortOrder)
);
const sortedPreset = computed(() => {
      const presetId = userStore.currentUser?.activePresetId;
      const regexs = regexStore.presetRegex.get(presetId ?? 0);
      return regexs ?? [];
    }
);
const sortedCharacter = computed(() => {
      const charId = charStore.curCharId;
      const regexs = regexStore.charRegex.get(charId ?? 0);
      return regexs ?? [];
    }
);

// 监听角色/预设变化，自动加载对应正则
watch(
    () => [charStore.curCharId, userStore.currentUser?.activePresetId] as const,
    async ([newCharId, newPresetId]) => {
      await regexStore.loadRegexs(newPresetId ?? 0, newCharId ?? 0);
    },
    {immediate: true}
);

// 判断是否正在编辑
const isEditing = (ruleId: number): boolean =>
    route.name === 'RegexEditor' && Number(route.params.id) === ruleId;

// 打开正则编辑器
const openRegexEditor = (ruleId: number) => {
  router.push(`/regex/${ruleId}`);
};

// 切换启用/禁用
const toggleEnabled = async (ruleId: number) => {
  try {
    await regexStore.toggleRegexEnabled(ruleId);
  } catch (e) {
    console.error('切换规则状态失败:', e);
  }
};

// 菜单操作
const handleMenuSelect = (key: string, ruleId: number) => {
  switch (key) {
    case 'copy':
      handleCopy(ruleId);
      break;
    case 'export':
      handleExport();
      break;
    case 'delete':
      handleDelete(ruleId);
      break;
  }
};

// 复制规则
const handleCopy = async (ruleId: number) => {
  const allRules = [
    ...regexStore.globalRegex,
    ...[...regexStore.presetRegex.values()].flat(),
    ...[...regexStore.charRegex.values()].flat()
  ];
  const rule = allRules.find(r => r.id === ruleId);
  if (!rule) return;
  try {
    const newRule = await regexStore.addRegex({
      presetId: rule.presetId,
      characterId: rule.characterId,
      name: `${rule.name} (副本)`,
      findPattern: rule.findPattern,
      replacePattern: rule.replacePattern,
      isEnabled: rule.isEnabled,
      runOnEdit: rule.runOnEdit,
      substituteRegex: rule.substituteRegex,
      minDepth: rule.minDepth,
      maxDepth: rule.maxDepth,
      affectFlags: rule.affectFlags,
      sortOrder: allRules.length
    });
    if (newRule) {
      message.success('规则已复制');
    }
  } catch (e) {
    console.error('复制规则失败:', e);
  }
};

// 导出规则
const handleExport = async () => {
  try {
    const result = await regexStore.exportRegex(0);
    const blob = new Blob([new Uint8Array(result.fileContent)], {type: 'application/json'});
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = result.fileName || 'regex-rules.json';
    a.click();
    URL.revokeObjectURL(url);
    message.success('导出成功');
  } catch (e) {
    console.error('导出正则规则失败:', e);
  }
};

// 删除规则
const handleDelete = (ruleId: number) => {
  const allRules = [
    ...regexStore.globalRegex,
    ...[...regexStore.presetRegex.values()].flat(),
    ...[...regexStore.charRegex.values()].flat()
  ];
  const rule = allRules.find(r => r.id === ruleId);
  if (!rule) return;
  dialog.warning({
    title: '确认删除',
    content: `确定要删除正则规则「${rule.name}」吗？此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await regexStore.deleteRegex(ruleId);
        message.success('规则已删除');
        if (isEditing(ruleId)) {
          router.push('/');
        }
      } catch (e) {
        console.error('删除规则失败:', e);
      }
    }
  });
};

// 导入正则规则
const handleImport = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = '.json';
  input.onchange = async (event: Event) => {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (!file) return;
    try {
      const buffer = await file.arrayBuffer();
      await regexStore.importRegex(new Uint8Array(buffer), file.name);
      message.success('导入成功');
    } catch (e) {
      console.error('导入正则规则失败:', e);
    }
  };
  input.click();
};

// 新建规则
const handleCreate = () => {
  router.push({name: 'RegexCreator', query: {tab: activeTab.value}});
};
</script>

<style scoped>
.regex-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  height: 100%;
  min-height: 0;
}

.action-row {
  display: flex;
  gap: 8px;
}

.rules-tabs {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.rules-tabs :deep(.n-tabs-pane-wrapper) {
  flex: 1;
  min-height: 0;
}

.rules-tabs :deep(.n-tab-pane) {
  height: 100%;
  min-height: 0;
}
</style>
