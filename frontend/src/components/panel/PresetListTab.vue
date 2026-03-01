<template>
  <div class="preset-list-tab">
    <!-- 当前使用的预设 -->
    <div class="current-preset">
      <div class="section-label">当前使用</div>
      <n-select
        :value="currentPresetId"
        :options="presetOptions"
        placeholder="选择预设"
        size="small"
        :loading="settingActive"
        @update:value="handleSetActive"
      />
    </div>

    <!-- 操作按钮 -->
    <div class="action-row">
      <n-button size="small" @click="handleImport">
        <template #icon>
          <n-icon><CloudUploadOutline /></n-icon>
        </template>
        导入
      </n-button>
      <n-button size="small" type="primary" @click="handleCreate">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        新建
      </n-button>
    </div>

    <!-- 预设列表 -->
    <n-spin :show="loading" size="small">
      <div class="list-container">
        <div
          v-for="preset in presets"
          :key="preset.id"
          class="list-item"
          :class="{ active: isEditing(preset.id) }"
          @click="openPresetEditor(preset.id)"
        >
          <div class="item-info">
            <div class="item-name">
              <span v-if="preset.id === currentPresetId" class="star-icon">★</span>
              {{ preset.name }}
            </div>
            <div class="item-desc">
              {{ preset.promptItems?.length ?? 0 }}个Prompt项{{ (preset.regexRules?.length ?? 0) > 0 ? `, ${preset.regexRules.length}个正则` : '' }}
            </div>
          </div>
          <n-dropdown
            trigger="click"
            :options="itemMenuOptions"
            @select="(key: string) => handleMenuSelect(key, preset)"
            @click.stop
          >
            <n-button quaternary circle size="tiny" @click.stop>
              <template #icon>
                <n-icon :size="16"><EllipsisHorizontal /></n-icon>
              </template>
            </n-button>
          </n-dropdown>
        </div>

        <!-- 空状态 -->
        <n-empty v-if="!loading && presets.length === 0" description="暂无预设" size="small" class="empty-state" />
      </div>
    </n-spin>

    <!-- 隐藏的文件选择器 -->
    <input
      ref="fileInputRef"
      type="file"
      accept=".json,.preset"
      style="display: none;"
      @change="handleFileSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { NSelect, NButton, NIcon, NDropdown, NEmpty, NSpin, useMessage, useDialog } from 'naive-ui';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';
import { usePresetStore } from '@/stores/preset';
import { useUserStore } from '@/stores/user';
import type { Preset } from '@/gen/muse/preset_pb';

const router = useRouter();
const route = useRoute();
const message = useMessage();
const dialog = useDialog();
const presetStore = usePresetStore();
const userStore = useUserStore();

// 状态
const loading = ref(false);
const settingActive = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

// 从 store 获取预设列表
const presets = computed(() => presetStore.presets);

// 当前活跃预设ID
const currentPresetId = computed(() => userStore.currentUser?.activePresetId ?? null);

// 下拉选项
const presetOptions = computed(() =>
  presets.value.map((p) => ({ label: p.name, value: p.id }))
);

// 右键菜单选项
const itemMenuOptions = [
  { label: '设为当前使用', key: 'use' },
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

// =====================
// 加载数据
// =====================

const loadPresets = async () => {
  loading.value = true;
  try {
    await presetStore.fetchAllPresets();
  } catch (error) {
    console.error('加载预设列表失败:', error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadPresets();
});

// =====================
// 操作方法
// =====================

const isEditing = (presetId: number): boolean => {
  return route.name === 'PresetEditor' && Number(route.params.id) === presetId;
};

const openPresetEditor = (presetId: number) => {
  router.push(`/preset/${presetId}`);
};

// 设置活跃预设
const handleSetActive = async (presetId: number) => {
  settingActive.value = true;
  try {
    await presetStore.setActivePreset(presetId);
    // 更新本地用户的活跃预设ID
    if (userStore.currentUser) {
      userStore.currentUser.activePresetId = presetId;
    }
    message.success('已切换当前预设');
  } catch (error) {
    console.error('设置活跃预设失败:', error);
  } finally {
    settingActive.value = false;
  }
};

// 导入预设
const handleImport = () => {
  fileInputRef.value?.click();
};

const handleFileSelected = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  try {
    const arrayBuffer = await file.arrayBuffer();
    const fileContent = new Uint8Array(arrayBuffer);
    const preset = await presetStore.importPreset(fileContent, file.name);
    if (preset) {
      message.success(`预设 "${preset.name}" 导入成功`);
    }
  } catch (error) {
    console.error('导入预设失败:', error);
  } finally {
    // 重置文件选择器
    input.value = '';
  }
};

// 新建预设
const handleCreate = async () => {
  try {
    const preset = await presetStore.createPreset({
      name: '新预设',
      temperature: 0.7,
      topP: 0.9,
      topK: 40,
      maxTokens: 2048,
      frequencyPenalty: 0,
      presencePenalty: 0,
    });
    if (preset) {
      message.success('预设创建成功');
      router.push(`/preset/${preset.id}`);
    }
  } catch (error) {
    console.error('创建预设失败:', error);
  }
};

// 菜单操作
const handleMenuSelect = (key: string, preset: Preset) => {
  switch (key) {
    case 'use':
      handleSetActive(preset.id);
      break;
    case 'copy':
      handleCopy(preset);
      break;
    case 'export':
      handleExport(preset);
      break;
    case 'delete':
      handleDelete(preset);
      break;
  }
};

// 复制预设
const handleCopy = async (preset: Preset) => {
  try {
    // 先获取完整预设数据（含 promptItems）
    const fullPreset = await presetStore.fetchPreset(preset.id);
    if (!fullPreset) return;

    const newPreset = await presetStore.createPreset({
      name: `${fullPreset.name} (副本)`,
      temperature: fullPreset.temperature,
      topP: fullPreset.topP,
      topK: fullPreset.topK,
      maxTokens: fullPreset.maxTokens,
      frequencyPenalty: fullPreset.frequencyPenalty,
      presencePenalty: fullPreset.presencePenalty,
      promptItems: fullPreset.promptItems?.map((item, index) => ({
        identifier: item.identifier,
        name: item.name,
        content: item.content,
        role: item.role,
        isEnabled: item.isEnabled,
        injectionPosition: item.injectionPosition,
        injectionDepth: item.injectionDepth,
        forbidOverrides: item.forbidOverrides,
        sortOrder: index,
      })),
    });
    if (newPreset) {
      message.success(`预设 "${newPreset.name}" 复制成功`);
    }
  } catch (error) {
    console.error('复制预设失败:', error);
  }
};

// 导出预设
const handleExport = async (preset: Preset) => {
  try {
    const result = await presetStore.exportPreset(preset.id);
    if (result.fileContent && result.fileName) {
      // 创建下载链接
      const blob = new Blob([new Uint8Array(result.fileContent)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = result.fileName;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
      message.success('导出成功');
    }
  } catch (error) {
    console.error('导出预设失败:', error);
  }
};

// 删除预设
const handleDelete = (preset: Preset) => {
  // 不允许删除当前活跃预设
  if (preset.id === currentPresetId.value) {
    message.warning('不能删除当前正在使用的预设');
    return;
  }

  dialog.warning({
    title: '确认删除',
    content: `确定要删除预设 "${preset.name}" 吗？此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await presetStore.deletePreset(preset.id);
        message.success('预设已删除');
      } catch (error) {
        console.error('删除预设失败:', error);
      }
    }
  });
};
</script>

<style scoped>
.preset-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.section-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.action-row {
  display: flex;
  gap: 8px;
}

.list-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background-color 150ms;
  border: 1px solid transparent;
}

.list-item:hover {
  background: var(--bg-hover);
}

.list-item.active {
  background: var(--bg-active);
  border-color: var(--color-primary);
}

.item-info {
  flex: 1;
  min-width: 0;
}

.item-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.star-icon {
  color: var(--color-primary);
  margin-right: 4px;
}

.item-desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.empty-state {
  padding: 32px 0;
}
</style>
