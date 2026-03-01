<template>
  <div class="worldinfo-list-tab">
    <!-- 操作按钮 -->
    <div class="action-row">
      <n-button size="small" :loading="importing" @click="handleImport">
        <template #icon>
          <n-icon><CloudUploadOutline /></n-icon>
        </template>
        导入
      </n-button>
      <n-button size="small" type="primary" :loading="creating" @click="handleCreate">
        <template #icon>
          <n-icon><AddOutline /></n-icon>
        </template>
        新建
      </n-button>
    </div>

    <!-- 加载状态 -->
    <n-spin :show="loading" description="加载中..." style="min-height: 60px;">
      <!-- 世界书列表 -->
      <div class="list-container">
        <div
          v-for="wi in worldInfoStore.worldInfos"
          :key="wi.id"
          class="list-item"
          :class="{ active: isEditing(wi.id) }"
          @click="openWorldInfoEditor(wi.id)"
        >
          <n-checkbox
            :checked="wi.isGlobal"
            @update:checked="(val: boolean) => toggleGlobal(wi, val)"
            @click.stop
          />
          <div class="item-info">
            <div class="item-name">{{ wi.name }}</div>
            <div class="item-desc">
              {{ wi.entries.length }}个词条
              <n-tag v-if="wi.isGlobal" size="tiny" type="info" :bordered="false" style="margin-left: 4px;">全局</n-tag>
            </div>
          </div>
          <n-dropdown
            trigger="click"
            :options="itemMenuOptions"
            @select="(key: string) => handleMenuSelect(key, wi)"
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
        <n-empty v-if="worldInfoStore.worldInfos.length === 0 && !loading" description="暂无世界书" size="small" class="empty-state" />
      </div>
    </n-spin>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInputRef"
      type="file"
      accept=".json"
      style="display: none;"
      @change="onFileSelected"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { NButton, NIcon, NCheckbox, NDropdown, NEmpty, NTag, NSpin, useDialog } from 'naive-ui';
import { CloudUploadOutline, AddOutline, EllipsisHorizontal } from '@vicons/ionicons5';
import { worldInfoClient } from '@/api/client';
import { useWorldInfoStore } from '@/stores/worldInfo';
import { globalMessage } from '@/composables/useGlobalMessage';
import type { WorldInfo } from '@/gen/muse/muse_pb';

const router = useRouter();
const route = useRoute();
const dialog = useDialog();
const worldInfoStore = useWorldInfoStore();

const loading = ref(false);
const importing = ref(false);
const creating = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

const itemMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

// 加载世界书列表
const loadWorldInfos = async () => {
  loading.value = true;
  try {
    const resp = await worldInfoClient.listWorldInfos({ pageSize: 100 });
    worldInfoStore.setWorldInfos(resp.worldInfos);
  } catch (e) {
    console.error('加载世界书列表失败:', e);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadWorldInfos();
});

// 判断是否正在编辑
const isEditing = (wiId: number): boolean => {
  return route.name === 'WorldInfoEditor' && Number(route.params.id) === wiId;
};

// 打开世界书编辑器
const openWorldInfoEditor = (wiId: number) => {
  router.push(`/worldinfo/${wiId}`);
};

// 切换全局状态
const toggleGlobal = async (wi: WorldInfo, isGlobal: boolean) => {
  try {
    const resp = await worldInfoClient.updateWorldInfo({
      id: wi.id,
      name: wi.name,
      description: wi.description,
      isGlobal
    });
    if (resp.worldInfo) {
      worldInfoStore.updateWorldInfoInList(resp.worldInfo);
    }
    globalMessage.success(isGlobal ? '已设为全局世界书' : '已取消全局');
  } catch (e) {
    console.error('更新世界书失败:', e);
  }
};

// 导入世界书
const handleImport = () => {
  fileInputRef.value?.click();
};

const onFileSelected = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  importing.value = true;
  try {
    const arrayBuffer = await file.arrayBuffer();
    const fileContent = new Uint8Array(arrayBuffer);
    const resp = await worldInfoClient.importWorldInfo({
      fileContent,
      fileName: file.name
    });
    if (resp.worldInfo) {
      worldInfoStore.addWorldInfo(resp.worldInfo);
      globalMessage.success(`导入世界书"${resp.worldInfo.name}"成功`);
    }
  } catch (e) {
    console.error('导入世界书失败:', e);
  } finally {
    importing.value = false;
    // 重置文件输入
    if (input) input.value = '';
  }
};

// 创建新世界书
const handleCreate = async () => {
  creating.value = true;
  try {
    const resp = await worldInfoClient.createWorldInfo({
      name: '新世界书',
      isGlobal: false
    });
    if (resp.worldInfo) {
      worldInfoStore.addWorldInfo(resp.worldInfo);
      globalMessage.success('创建成功');
      // 跳转到编辑页
      router.push(`/worldinfo/${resp.worldInfo.id}`);
    }
  } catch (e) {
    console.error('创建世界书失败:', e);
  } finally {
    creating.value = false;
  }
};

// 菜单选择处理
const handleMenuSelect = (key: string, wi: WorldInfo) => {
  switch (key) {
    case 'copy':
      handleCopy(wi);
      break;
    case 'export':
      handleExport(wi);
      break;
    case 'delete':
      handleDelete(wi);
      break;
  }
};

// 复制世界书
const handleCopy = async (wi: WorldInfo) => {
  try {
    const resp = await worldInfoClient.createWorldInfo({
      name: `${wi.name} (副本)`,
      description: wi.description,
      isGlobal: wi.isGlobal
    });
    if (resp.worldInfo) {
      worldInfoStore.addWorldInfo(resp.worldInfo);
      globalMessage.success('复制成功');
    }
  } catch (e) {
    console.error('复制世界书失败:', e);
  }
};

// 导出世界书
const handleExport = async (wi: WorldInfo) => {
  try {
    const resp = await worldInfoClient.exportWorldInfo({ id: wi.id });
    // 触发下载
    const blob = new Blob([new Uint8Array(resp.fileContent)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = resp.fileName || `${wi.name}.json`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
    globalMessage.success('导出成功');
  } catch (e) {
    console.error('导出世界书失败:', e);
  }
};

// 删除世界书
const handleDelete = (wi: WorldInfo) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除世界书"${wi.name}"吗？此操作不可恢复。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await worldInfoClient.deleteWorldInfo({ id: wi.id });
        worldInfoStore.removeWorldInfo(wi.id);
        globalMessage.success('删除成功');
      } catch (e) {
        console.error('删除世界书失败:', e);
      }
    }
  });
};
</script>

<style scoped>
.worldinfo-list-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
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
  gap: 10px;
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

.item-desc {
  font-size: 11px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.empty-state {
  padding: 32px 0;
}
</style>
