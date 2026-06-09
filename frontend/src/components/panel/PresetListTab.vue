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
      <n-infinite-scroll class="list-container" @load="loadMore" :distance="100">
        <div
          v-for="preset in presets"
          :key="preset.preset?.id"
          class="list-item"
          :class="{ active: isEditing(preset.preset?.id ?? 0) }"
          @click="openPresetEditor(preset.preset?.id ?? 0)"
        >
          <div class="item-info">
            <div class="item-name">
              <span v-if="preset.preset?.id === currentPresetId" class="star-icon">★</span>
              {{ preset.preset?.name }}
            </div>
            <div class="item-desc">
              {{ preset.promptLen ?? 0 }}个Prompt项
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

        <!-- 加载更多指示 -->
        <div v-if="loadingMore" class="loading-more">
          <n-spin size="small" />
        </div>

        <!-- 空状态 -->
        <n-empty v-if="!loading && presets.length === 0" description="暂无预设" size="small" class="empty-state" />
      </n-infinite-scroll>
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
import { NInfiniteScroll } from 'naive-ui/es/infinite-scroll';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';
import { usePresetStore } from '@/stores/preset';
import { useUserStore } from '@/stores/user';
import type { PresetWithPromptLen } from '@/gen/muse/preset_pb';

const router = useRouter();
const route = useRoute();
const message = useMessage();
const dialog = useDialog();
const presetStore = usePresetStore();
const userStore = useUserStore();

// 状态
const loading = computed(() => presetStore.loading)
const settingActive = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

const presets = computed(() => presetStore.presets || []);
const currentPresetId = computed(() => userStore.currentUser?.activePresetId || null);
const hasMore = computed(() => presetStore.hasMore);
const loadingMore = computed(() => presetStore.loadingMore);

const loadMore = () => {
  if (!hasMore.value || loadingMore.value) return Promise.resolve();
  return presetStore.loadMore();
};

// 下拉选项
const presetOptions = computed(() =>
  presets.value!.map((p) => ({ label: p.preset?.name ?? '', value: p.preset?.id ?? 0 }))
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
onMounted(async () => {
  await presetStore.loadPresets();
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
    if (preset && preset.preset) {
      message.success(`预设 "${preset.preset.name}" 导入成功`);
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
    await presetStore.createPreset({
      name: '新预设',
      temperature: 0.7,
      topP: 0.9,
      topK: 40,
      maxTokens: 2048,
      frequencyPenalty: 0,
      presencePenalty: 0,
    });
    message.success('预设创建成功');
    if (presets.value!.length > 0 && presets.value![0]?.preset) {
      router.push(`/preset/${presets.value![0].preset.id}`);
    }
  } catch (error) {
    console.error('创建预设失败:', error);
  }
};

// 菜单操作
const handleMenuSelect = (key: string, preset: PresetWithPromptLen) => {
  if (!preset.preset) return;
  
  switch (key) {
    case 'use':
      handleSetActive(preset.preset.id);
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
const handleCopy = async (preset: PresetWithPromptLen) => {
  if (!preset.preset) return;
  
  try {
    // 先获取完整预设数据（含 promptItems）
    const fullPreset = await presetStore.fetchPreset(preset.preset.id);
    if (!fullPreset || !fullPreset.preset) return;

    await presetStore.createPreset({
      name: `${fullPreset.preset.name} (副本)`,
      temperature: fullPreset.preset.temperature,
      topP: fullPreset.preset.topP,
      topK: fullPreset.preset.topK,
      maxTokens: fullPreset.preset.maxTokens,
      frequencyPenalty: fullPreset.preset.frequencyPenalty,
      presencePenalty: fullPreset.preset.presencePenalty,
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
    message.success(`预设 "${fullPreset.preset.name}" 复制成功`);
  } catch (error) {
    console.error('复制预设失败:', error);
  }
};

// 导出预设
const handleExport = async (preset: PresetWithPromptLen) => {
  if (!preset.preset) return;
  
  try {
    const result = await presetStore.exportPreset(preset.preset.id);
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
const handleDelete = (preset: PresetWithPromptLen) => {
  if (!preset.preset) return;
  
  // 不允许删除当前活跃预设
  if (preset.preset.id === currentPresetId.value) {
    message.warning('不能删除当前正在使用的预设');
    return;
  }

  dialog.warning({
    title: '确认删除',
    content: `确定要删除预设 "${preset.preset.name}" 吗？此操作不可撤销。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await presetStore.deletePreset(preset.preset!.id);
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
  height: 100%;
  min-height: 0;
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
  flex: 1;
  min-height: 0;
  overflow: auto;
}

/* 穿透 n-spin 内部结构，保证 flex 布局链条完整 */
.preset-list-tab :deep(.n-spin-container),
.preset-list-tab :deep(.n-spin-content) {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
}

.loading-more {
  display: flex;
  justify-content: center;
  padding: 12px 0;
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
