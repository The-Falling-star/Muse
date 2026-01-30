<template>
  <div class="presets-view">
    <!-- 左侧：预设列表 -->
    <div class="presets-sidebar">
      <div class="sidebar-header">
        <h3>预设列表</h3>
        <n-button type="primary" size="small" @click="createPreset">
          <template #icon>
            <n-icon><AddOutline /></n-icon>
          </template>
        </n-button>
      </div>

      <n-input
        v-model:value="searchQuery"
        placeholder="搜索预设..."
        clearable
        size="small"
        class="search-input"
      >
        <template #prefix>
          <n-icon><SearchOutline /></n-icon>
        </template>
      </n-input>

      <n-scrollbar class="preset-list-container">
        <n-spin :show="loading">
          <div class="preset-list">
            <div
              v-for="preset in filteredPresets"
              :key="preset.id"
              class="preset-item"
              :class="{ 'active': selectedPreset?.id === preset.id }"
              @click="selectPreset(preset)"
            >
              <div class="preset-item-info">
                <span class="preset-item-name">{{ preset.name }}</span>
                <span class="preset-item-count">{{ preset.promptItems?.length || 0 }} 项</span>
              </div>
              <n-dropdown :options="presetMenuOptions" trigger="click" @select="(key: string) => handlePresetMenu(key, preset)">
                <n-button quaternary circle size="tiny" @click.stop>
                  <template #icon>
                    <n-icon><EllipsisVerticalOutline /></n-icon>
                  </template>
                </n-button>
              </n-dropdown>
            </div>
          </div>

          <n-empty v-if="filteredPresets.length === 0 && !loading" description="暂无预设" size="small" />
        </n-spin>
      </n-scrollbar>
    </div>

    <!-- 右侧：预设编辑面板 -->
    <div class="presets-panel">
      <template v-if="selectedPreset">
        <!-- 预设头部信息 -->
        <div class="panel-header">
          <div class="header-title">
            <n-input
              v-model:value="editForm.name"
              :bordered="false"
              placeholder="预设名称"
              class="preset-name-input"
            />
          </div>
          <div class="header-actions">
            <n-button type="primary" size="small" :loading="saving" @click="savePreset">
              <template #icon>
                <n-icon><SaveOutline /></n-icon>
              </template>
              保存
            </n-button>
          </div>
        </div>

        <!-- 主内容区：上下结构 -->
        <n-scrollbar class="panel-content">
          <!-- 第一部分：AI生成参数 -->
          <section class="settings-section">
            <div class="section-header">
              <n-icon size="20"><OptionsOutline /></n-icon>
              <span class="section-title">AI生成参数</span>
            </div>
            <div class="section-content">
              <div class="param-grid">
                <!-- Temperature -->
                <div class="param-item">
                  <span class="param-label">Temperature</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.temperature"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.temperature"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
                <!-- Top P -->
                <div class="param-item">
                  <span class="param-label">Top P</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.topP"
                      :min="0"
                      :max="1"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.topP"
                      :min="0"
                      :max="1"
                      :step="0.05"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
                <!-- Top K -->
                <div class="param-item">
                  <span class="param-label">Top K</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.topK"
                      :min="0"
                      :max="500"
                      :step="1"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.topK"
                      :min="0"
                      :max="500"
                      :step="1"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
                <!-- 最大 Tokens -->
                <div class="param-item">
                  <span class="param-label">最大 Tokens</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.maxTokens"
                      :min="1"
                      :max="9999999"
                      :step="10000"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.maxTokens"
                      :min="1"
                      :max="9999999"
                      :step="100"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
                <!-- 频率惩罚 -->
                <div class="param-item">
                  <span class="param-label">频率惩罚</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.frequencyPenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.frequencyPenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
                <!-- 存在惩罚 -->
                <div class="param-item">
                  <span class="param-label">存在惩罚</span>
                  <div class="param-control">
                    <n-slider
                      v-model:value="editForm.presencePenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="editForm.presencePenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      size="small"
                      :show-button="true"
                      button-placement="both"
                      class="param-input"
                    />
                  </div>
                </div>
              </div>
            </div>
          </section>

          <!-- 第二部分：预设正则 -->
          <section class="settings-section">
            <div class="section-header">
              <n-icon size="20"><CodeSlashOutline /></n-icon>
              <span class="section-title">预设正则</span>
            </div>
            <div class="section-content">
              <PresetRegexManager
                v-model:regexRules="regexRules"
                :preset-id="selectedPreset.id"
              />
            </div>
          </section>

          <!-- 第三部分：预设提示词 -->
          <section class="settings-section">
            <div class="section-header">
              <n-icon size="20"><ChatboxOutline /></n-icon>
              <span class="section-title">预设提示词</span>
              <n-button size="tiny" @click="addPromptItem" class="section-action">
                <template #icon>
                  <n-icon><AddOutline /></n-icon>
                </template>
                添加
              </n-button>
            </div>
            <div class="section-content">
              <!-- 提示项列表 - 可拖拽排序 -->
              <draggable
                v-model="promptItems"
                item-key="id"
                handle=".drag-handle"
                animation="200"
                ghost-class="prompt-ghost"
                class="prompt-list"
                @end="handlePromptOrderChange"
              >
                <template #item="{ element: item, index }">
                  <div
                    class="prompt-item"
                    :class="{
                      'disabled': !item.isEnabled,
                      'editing': editingPromptId === item.id
                    }"
                  >
                    <!-- 拖拽手柄 -->
                    <div class="drag-handle">
                      <n-icon size="16"><ReorderTwoOutline /></n-icon>
                    </div>

                    <!-- 启用开关 -->
                    <n-switch
                      :value="item.isEnabled"
                      size="small"
                      @update:value="(val: boolean) => togglePromptEnabled(item, val)"
                    />

                    <!-- 序号 -->
                    <span class="prompt-index">{{ index + 1 }}</span>

                    <!-- 提示项信息 -->
                    <div class="prompt-info" @click="toggleExpand(item.id)">
                      <div class="prompt-header-row">
                        <span class="prompt-name">{{ item.name }}</span>
                        <n-tag :type="getRoleColor(item.role)" size="small" :bordered="false">
                          {{ getRoleName(item.role) }}
                        </n-tag>
                      </div>
                      <p v-if="item.content" class="prompt-preview">
                        {{ truncateContent(item.content) }}
                      </p>
                    </div>

                    <!-- 操作按钮 -->
                    <div class="prompt-actions">
                      <n-button quaternary circle size="tiny" @click="editPromptItem(item)">
                        <template #icon>
                          <n-icon><CreateOutline /></n-icon>
                        </template>
                      </n-button>
                      <n-button
                        quaternary
                        circle
                        size="tiny"
                        @click="deletePromptItem(item)"
                      >
                        <template #icon>
                          <n-icon><TrashOutline /></n-icon>
                        </template>
                      </n-button>
                    </div>

                    <!-- 展开的编辑区域 -->
                    <Transition name="expand">
                      <div v-if="expandedPromptId === item.id" class="prompt-expanded">
                        <n-input
                          v-model:value="item.content"
                          type="textarea"
                          :autosize="{ minRows: 3, maxRows: 10 }"
                          placeholder="输入提示内容，支持 {{char}} {{user}} 等变量"
                          @blur="handlePromptContentChange(item)"
                        />
                      </div>
                    </Transition>
                  </div>
                </template>
              </draggable>

              <!-- 提示词变量说明 -->
              <n-collapse class="variables-collapse">
                <n-collapse-item title="可用变量" name="variables">
                  <div class="variables-grid">
                    <span v-for="v in variablesList" :key="v" class="variable-tag">{{ v }}</span>
                  </div>
                </n-collapse-item>
              </n-collapse>
            </div>
          </section>
        </n-scrollbar>
      </template>

      <!-- 未选择预设时的占位 -->
      <div v-else class="empty-panel">
        <n-empty description="选择一个预设进行编辑">
          <template #extra>
            <n-button type="primary" @click="createPreset">创建预设</n-button>
          </template>
        </n-empty>
      </div>
    </div>

    <!-- 提示项编辑模态框 -->
    <n-modal
      v-model:show="showPromptModal"
      preset="card"
      :title="editingPromptItem ? '编辑提示项' : '新建提示项'"
      :style="{ width: '600px', maxWidth: '90vw' }"
    >
      <PromptItemEditor
        :item="editingPromptItem"
        @save="handleSavePromptItem"
        @cancel="showPromptModal = false"
      />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import {
  NInput,
  NInputNumber,
  NSlider,
  NButton,
  NIcon,
  NScrollbar,
  NSwitch,
  NTag,
  NDropdown,
  NEmpty,
  NModal,
  NCollapse,
  NCollapseItem,
  NSpin,
  useMessage,
  useDialog
} from 'naive-ui';
import {
  SearchOutline,
  AddOutline,
  EllipsisVerticalOutline,
  SaveOutline,
  ReorderTwoOutline,
  CreateOutline,
  TrashOutline,
  OptionsOutline,
  CodeSlashOutline,
  ChatboxOutline
} from '@vicons/ionicons5';
import draggable from 'vuedraggable';

import PromptItemEditor from '../components/preset/PromptItemEditor.vue';
import PresetRegexManager from '../components/preset/PresetRegexManager.vue';
import { presetClient, regexRuleClient } from '@/api/client';
import type { Preset, PromptItem } from '@/gen/muse/muse_pb';
import { Role, InjectionPosition } from '@/gen/muse/muse_pb';

// 组件本地使用的正则规则接口
interface LocalRegexRule {
  id: number;
  name: string;
  pattern: string;
  replacement: string;
  flags: string;
  scope: 'input' | 'output' | 'both';
  order: number;
  enabled: boolean;
}

// 组件本地使用的提示项接口
interface LocalPromptItem {
  id: number;
  identifier: string;
  name: string;
  role: 'system' | 'user' | 'assistant';
  content: string;
  enabled: boolean;
  marker?: boolean;
  injection?: {
    position: 'before' | 'after';
    depth: number;
  };
}

const message = useMessage();
const dialog = useDialog();

// 状态
const searchQuery = ref('');
const loading = ref(false);
const saving = ref(false);
const presets = ref<Preset[]>([]);
const selectedPreset = ref<Preset | null>(null);
const promptItems = ref<PromptItem[]>([]);
const regexRules = ref<LocalRegexRule[]>([]);
const showPromptModal = ref(false);
const editingPromptItem = ref<LocalPromptItem | null>(null);
const editingPromptId = ref<number | null>(null);
const expandedPromptId = ref<number | null>(null);

// 编辑表单
const editForm = ref({
  name: '',
  temperature: 1,
  topP: 1,
  topK: 0,
  maxTokens: 300,
  frequencyPenalty: 0,
  presencePenalty: 0
});

// 预设菜单选项
const presetMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

// 计算属性：过滤后的预设列表
const filteredPresets = computed(() => {
  if (!searchQuery.value) return presets.value;
  const query = searchQuery.value.toLowerCase();
  return presets.value.filter(p =>
    p.name.toLowerCase().includes(query)
  );
});

// 加载预设列表
const loadPresets = async () => {
  loading.value = true;
  try {
    const response = await presetClient.listPresets({});
    presets.value = response.presets;
    // 如果有预设且没有选中的，选择第一个
    if (presets.value.length > 0 && !selectedPreset.value) {
      const firstPreset = presets.value[0];
      if (firstPreset) {
        await selectPreset(firstPreset);
      }
    }
  } finally {
    loading.value = false;
  }
};

// 加载预设详情（包括提示项）
const loadPresetDetail = async (presetId: number) => {
  try {
    const response = await presetClient.getPreset({ id: presetId });
    if (response.preset) {
      // 更新提示项列表
      promptItems.value = response.preset.promptItems || [];
    }
  } catch (e) {
    console.error('加载预设详情失败:', e);
  }
};

// 加载预设正则规则
const loadPresetRegexRules = async (presetId: number) => {
  try {
    const response = await regexRuleClient.listRegexRules({ presetId: presetId });
    // 转换为组件本地类型
    regexRules.value = response.rules.map(r => ({
      id: r.id,
      name: r.name,
      pattern: r.findPattern,
      replacement: r.replacePattern,
      flags: 'gi',
      scope: 'both' as const,
      order: r.sortOrder,
      enabled: r.isEnabled
    }));
  } catch (e) {
    console.error('加载正则规则失败:', e);
  }
};

// 选择预设
const selectPreset = async (preset: Preset) => {
  selectedPreset.value = preset;
  expandedPromptId.value = null;

  // 同步编辑表单
  editForm.value = {
    name: preset.name,
    temperature: preset.temperature,
    topP: preset.topP,
    topK: preset.topK,
    maxTokens: preset.maxTokens,
    frequencyPenalty: preset.frequencyPenalty,
    presencePenalty: preset.presencePenalty
  };

  // 加载详情和正则规则
  await Promise.all([
    loadPresetDetail(preset.id),
    loadPresetRegexRules(preset.id)
  ]);
};

// 创建预设
const createPreset = async () => {
  try {
    const response = await presetClient.createPreset({
      name: '新预设',
      temperature: 1,
      topP: 1,
      topK: 0,
      maxTokens: 300,
      frequencyPenalty: 0,
      presencePenalty: 0
    });
    if (response.preset) {
      presets.value.push(response.preset);
      await selectPreset(response.preset);
      message.success('预设已创建');
    }
  } catch (e) {
    console.error('创建预设失败:', e);
    message.error('创建预设失败');
  }
};

// 处理预设菜单
const handlePresetMenu = async (key: string, preset: Preset) => {
  switch (key) {
    case 'copy':
      try {
        // 获取原预设的详情
        const detailRes = await presetClient.getPreset({ id: preset.id });
        const originalPreset = detailRes.preset;
        if (!originalPreset) return;

        // 创建副本
        const response = await presetClient.createPreset({
          name: `${preset.name} (副本)`,
          temperature: originalPreset.temperature,
          topP: originalPreset.topP,
          topK: originalPreset.topK,
          maxTokens: originalPreset.maxTokens,
          frequencyPenalty: originalPreset.frequencyPenalty,
          presencePenalty: originalPreset.presencePenalty,
          promptItems: originalPreset.promptItems?.map(item => ({
            identifier: item.identifier,
            name: item.name,
            content: item.content,
            role: item.role,
            isEnabled: item.isEnabled,
            injectionPosition: item.injectionPosition,
            injectionDepth: item.injectionDepth,
            forbidOverrides: item.forbidOverrides,
            sortOrder: item.sortOrder
          }))
        });
        if (response.preset) {
          presets.value.push(response.preset);
          message.success('预设已复制');
        }
      } catch (e) {
        console.error('复制预设失败:', e);
        message.error('复制预设失败');
      }
      break;

    case 'export':
      try {
        const response = await presetClient.exportPreset({ id: preset.id });
        // 下载文件
        const blob = new Blob([new Uint8Array(response.fileContent).buffer], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = response.fileName || `${preset.name}.json`;
        a.click();
        URL.revokeObjectURL(url);
        message.success('预设已导出');
      } catch (e) {
        console.error('导出预设失败:', e);
        message.error('导出预设失败');
      }
      break;

    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除预设"${preset.name}"吗？`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: async () => {
          try {
            await presetClient.deletePreset({ id: preset.id });
            presets.value = presets.value.filter(p => p.id !== preset.id);
            if (selectedPreset.value?.id === preset.id) {
              selectedPreset.value = presets.value[0] ?? null;
              if (selectedPreset.value) {
                await selectPreset(selectedPreset.value);
              } else {
                promptItems.value = [];
                regexRules.value = [];
              }
            }
            message.success('预设已删除');
          } catch (e) {
            console.error('删除预设失败:', e);
            message.error('删除预设失败');
          }
        }
      });
      break;
  }
};

// 保存预设
const savePreset = async () => {
  if (!selectedPreset.value) return;

  saving.value = true;
  try {
    const response = await presetClient.updatePreset({
      id: selectedPreset.value.id,
      name: editForm.value.name,
      temperature: editForm.value.temperature,
      topP: editForm.value.topP,
      topK: editForm.value.topK,
      maxTokens: editForm.value.maxTokens,
      frequencyPenalty: editForm.value.frequencyPenalty,
      presencePenalty: editForm.value.presencePenalty,
      version: selectedPreset.value.version
    });
    if (response.preset) {
      // 更新列表中的预设
      const index = presets.value.findIndex(p => p.id === selectedPreset.value!.id);
      if (index >= 0) {
        presets.value[index] = response.preset;
      }
      selectedPreset.value = response.preset;
      message.success('预设已保存');
    }
  } catch (e) {
    console.error('保存预设失败:', e);
    message.error('保存预设失败');
  } finally {
    saving.value = false;
  }
};

// 添加提示项
const addPromptItem = () => {
  editingPromptItem.value = null;
  showPromptModal.value = true;
};

// 编辑提示项
const editPromptItem = (item: PromptItem) => {
  // 转换为组件本地类型
  const roleMap: Record<number, 'system' | 'user' | 'assistant'> = {
    [Role.System]: 'system',
    [Role.User]: 'user',
    [Role.Assistant]: 'assistant'
  };
  editingPromptItem.value = {
    id: item.id,
    identifier: item.identifier,
    name: item.name,
    role: roleMap[item.role] || 'system',
    content: item.content,
    enabled: item.isEnabled
  };
  showPromptModal.value = true;
};

// 保存提示项
const handleSavePromptItem = async (itemData: Partial<LocalPromptItem>) => {
  // 角色映射
  const roleMap: Record<string, Role> = {
    'system': Role.System,
    'user': Role.User,
    'assistant': Role.Assistant
  };
  const role = itemData.role ? roleMap[itemData.role] : Role.System;
  if (!selectedPreset.value) return;

  try {
    if (editingPromptItem.value) {
      // 更新现有项
      const response = await presetClient.updatePromptItem({
        id: editingPromptItem.value.id,
        identifier: itemData.identifier || '',
        name: itemData.name || '',
        content: itemData.content,
        role: role,
        isEnabled: itemData.enabled ?? true,
        injectionPosition: InjectionPosition.Relative,
        injectionDepth: 0,
        forbidOverrides: false,
        sortOrder: 0
      });
      if (response.item) {
        const index = promptItems.value.findIndex(p => p.id === editingPromptItem.value!.id);
        if (index >= 0) {
          promptItems.value[index] = response.item;
        }
        message.success('提示项已更新');
      }
    } else {
      // 添加新项
      const response = await presetClient.addPromptItem({
        presetId: selectedPreset.value.id,
        identifier: itemData.identifier || `prompt_${Date.now()}`,
        name: itemData.name || '',
        content: itemData.content,
        role: role,
        isEnabled: itemData.enabled ?? true,
        injectionPosition: InjectionPosition.Relative,
        injectionDepth: 0,
        forbidOverrides: false,
        sortOrder: promptItems.value.length
      });
      if (response.item) {
        promptItems.value.push(response.item);
        message.success('提示项已添加');
      }
    }
    showPromptModal.value = false;
    editingPromptItem.value = null;
  } catch (e) {
    console.error('保存提示项失败:', e);
    message.error('保存提示项失败');
  }
};

// 删除提示项
const deletePromptItem = (item: PromptItem) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除提示项"${item.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await presetClient.deletePromptItem({ id: item.id });
        promptItems.value = promptItems.value.filter(p => p.id !== item.id);
        message.success('提示项已删除');
      } catch (e) {
        console.error('删除提示项失败:', e);
        message.error('删除提示项失败');
      }
    }
  });
};

// 切换提示项启用状态
const togglePromptEnabled = async (item: PromptItem, enabled: boolean) => {
  try {
    const response = await presetClient.updatePromptItem({
      id: item.id,
      identifier: item.identifier,
      name: item.name,
      content: item.content,
      role: item.role,
      isEnabled: enabled,
      injectionPosition: item.injectionPosition,
      injectionDepth: item.injectionDepth,
      forbidOverrides: item.forbidOverrides,
      sortOrder: item.sortOrder
    });
    if (response.item) {
      const index = promptItems.value.findIndex(p => p.id === item.id);
      if (index >= 0) {
        promptItems.value[index] = response.item;
      }
    }
  } catch (e) {
    console.error('更新提示项失败:', e);
  }
};

// 处理提示项排序变化
const handlePromptOrderChange = async () => {
  if (!selectedPreset.value) return;

  try {
    const itemIds = promptItems.value.map(item => item.id);
    await presetClient.updatePromptItemsOrder({
      presetId: selectedPreset.value.id,
      itemIds: itemIds
    });
  } catch (e) {
    console.error('更新排序失败:', e);
  }
};

// 处理提示项内容变化
const handlePromptContentChange = async (item: PromptItem) => {
  try {
    await presetClient.updatePromptItem({
      id: item.id,
      identifier: item.identifier,
      name: item.name,
      content: item.content,
      role: item.role,
      isEnabled: item.isEnabled,
      injectionPosition: item.injectionPosition,
      injectionDepth: item.injectionDepth,
      forbidOverrides: item.forbidOverrides,
      sortOrder: item.sortOrder
    });
  } catch (e) {
    console.error('更新提示项内容失败:', e);
  }
};

// 切换展开/收起
const toggleExpand = (itemId: number) => {
  expandedPromptId.value = expandedPromptId.value === itemId ? null : itemId;
};

// 辅助方法
const getRoleName = (role: Role) => {
  const names: Record<number, string> = {
    [Role.System]: '系统',
    [Role.User]: '用户',
    [Role.Assistant]: '助手'
  };
  return names[role] || '未知';
};

const getRoleColor = (role: Role): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' => {
  const colors: Record<number, 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'> = {
    [Role.System]: 'primary',
    [Role.User]: 'success',
    [Role.Assistant]: 'info'
  };
  return colors[role] || 'default';
};

const truncateContent = (content: string, maxLength = 100) => {
  if (content.length <= maxLength) return content;
  return content.slice(0, maxLength) + '...';
};

// 变量列表（用于显示）
const variablesList = [
  '{{char}}',
  '{{user}}',
  '{{scenario}}',
  '{{personality}}',
  '{{description}}',
  '{{persona}}'
];

// 初始化
onMounted(() => {
  loadPresets();
});
</script>

<style scoped>
.presets-view {
  display: flex;
  height: calc(100vh - 64px - 48px);
  gap: 16px;
}

/* 左侧预设列表 */
.presets-sidebar {
  width: 280px;
  min-width: 280px;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.search-input {
  margin: 12px 16px;
  width: calc(100% - 32px);
  box-sizing: border-box;
}

.preset-list-container {
  flex: 1;
}

.preset-list {
  padding: 0 8px 8px;
}

.preset-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  margin-bottom: 4px;
  border-radius: 8px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.preset-item:hover {
  background: var(--bg-tertiary);
}

.preset-item.active {
  background: var(--primary-color);
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
}

.preset-item.active .preset-item-name,
.preset-item.active .preset-item-count {
  color: #fff;
}

.preset-item-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.preset-item-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.preset-item-count {
  font-size: 12px;
  color: var(--text-tertiary);
}

/* 右侧预设面板 */
.presets-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.header-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.preset-name-input {
  font-size: 18px;
  font-weight: 600;
  max-width: 300px;
}

.preset-name-input :deep(input) {
  font-size: 18px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 8px;
}

/* 主内容区 */
.panel-content {
  flex: 1;
  padding: 16px 20px;
}

/* 设置分区 */
.settings-section {
  margin-bottom: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-color);
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  flex: 1;
}

.section-action {
  margin-left: auto;
}

.section-content {
  padding: 16px;
}

/* 参数网格布局 */
.param-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.param-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.param-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.param-control {
  display: flex;
  align-items: center;
  gap: 12px;
}

.param-slider {
  flex: 1;
}

.param-input {
  width: 140px;
  flex-shrink: 0;
}

.param-input :deep(input) {
  text-align: center;
}

@media (max-width: 900px) {
  .param-grid {
    grid-template-columns: 1fr;
  }
}

/* 提示项列表 */
.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.prompt-item {
  display: flex;
  flex-wrap: wrap;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 10px;
  transition: all var(--transition-fast);
}

.prompt-item:hover {
  border-color: var(--border-glow);
}

.prompt-item.disabled {
  opacity: 0.5;
}

.prompt-item.marker {
  background: var(--bg-tertiary);
  border-style: dashed;
}

.prompt-item.marker .prompt-info {
  cursor: default;
}

.prompt-marker-hint {
  margin: 4px 0 0;
  font-size: 12px;
  color: var(--text-tertiary);
  font-style: italic;
}

.prompt-item.editing {
  border-color: var(--primary-color);
  box-shadow: var(--glow-soft);
}

.prompt-ghost {
  opacity: 0.5;
  background: var(--primary-color);
}

.drag-handle {
  cursor: grab;
  color: var(--text-tertiary);
  padding: 4px;
  display: flex;
  align-items: center;
}

.drag-handle:hover {
  color: var(--text-secondary);
}

.drag-handle:active {
  cursor: grabbing;
}

.prompt-index {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-tertiary);
  background: var(--bg-tertiary);
  border-radius: 6px;
}

.prompt-info {
  flex: 1;
  min-width: 0;
  cursor: pointer;
}

.prompt-header-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.prompt-name {
  font-weight: 500;
  color: var(--text-primary);
}

.prompt-preview {
  margin: 6px 0 0;
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.prompt-actions {
  display: flex;
  gap: 4px;
}

/* 展开编辑区域 */
.prompt-expanded {
  width: 100%;
  margin-top: 8px;
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}

.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
}

/* 变量提示 */
.variables-collapse {
  margin-top: 16px;
}

.variables-collapse :deep(.n-collapse-item__header) {
  font-size: 13px;
}

.variables-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.variable-tag {
  padding: 4px 10px;
  font-size: 12px;
  font-family: 'Fira Code', monospace;
  color: var(--primary-color);
  background: var(--bg-tertiary);
  border-radius: 4px;
}

/* 空面板 */
.empty-panel {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 响应式 */
@media (max-width: 768px) {
  .presets-view {
    flex-direction: column;
  }

  .presets-sidebar {
    width: 100%;
    min-width: unset;
    max-height: 300px;
    border-right: none;
    border-bottom: 1px solid var(--border-color);
  }

  .preset-list-container {
    max-height: 220px;
  }

  .panel-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }

  .header-title {
    flex-wrap: wrap;
  }

  .header-actions {
    justify-content: flex-end;
  }
}
</style>
