<template>
  <div class="presets-view">
    <!-- 左侧：预设列表（移动端可折叠） -->
    <div class="presets-sidebar" :class="{ 'show-mobile': showSidebar }">
      <div class="sidebar-header">
        <div class="sidebar-header-content">
          <h3>预设列表</h3>
          <n-button 
            v-if="isMobile" 
            quaternary 
            circle 
            size="small" 
            class="close-sidebar-btn"
            @click="showSidebar = false"
          >
            <template #icon>
              <n-icon><CloseOutline /></n-icon>
            </template>
          </n-button>
        </div>
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
      <!-- 移动端：列表切换按钮 -->
      <div v-if="isMobile" class="mobile-list-toggle">
        <n-button 
          quaternary 
          @click="showSidebar = true"
          class="list-toggle-btn"
        >
          <template #icon>
            <n-icon><ListOutline /></n-icon>
          </template>
          预设列表
        </n-button>
      </div>
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
                      'editing': editingPromptId === item.id,
                      'forbid-overrides': item.forbidOverrides
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
                    <div class="prompt-info" @click="item.forbidOverrides ? null : toggleExpand(item.id)">
                      <div class="prompt-header-row">
                        <span class="prompt-name">{{ item.name }}</span>
                        <!-- 标签数组：同时显示标记和角色 -->
                        <div class="prompt-tags">
                          <n-tag v-if="item.forbidOverrides" type="warning" size="small" :bordered="false">
                            标记
                          </n-tag>
                          <n-tag :type="getRoleColor(item.role)" size="small" :bordered="false">
                            {{ getRoleName(item.role) }}
                          </n-tag>
                        </div>
                      </div>
                      <p v-if="item.content" class="prompt-preview">
                        {{ truncateContent(item.content) }}
                      </p>
                      <p v-if="item.forbidOverrides" class="prompt-marker-hint">
                        系统自动填充，名称和内容不可编辑
                      </p>
                    </div>

                    <!-- 操作按钮 -->
                    <div class="prompt-actions">
                      <n-button
                        quaternary
                        circle
                        size="tiny"
                        @click="editPromptItem(item)"
                      >
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

                    <!-- 展开的编辑区域（仅非标记项可展开） -->
                    <Transition name="expand">
                      <div v-if="expandedPromptId === item.id && !item.forbidOverrides" class="prompt-expanded">
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
import { ref, computed, onMounted, onUnmounted } from 'vue';
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
  ChatboxOutline,
  CloseOutline,
  ListOutline
} from '@vicons/ionicons5';
import draggable from 'vuedraggable';

import PromptItemEditor from '../components/preset/PromptItemEditor.vue';
import PresetRegexManager from '../components/preset/PresetRegexManager.vue';
import { presetClient, regexRuleClient } from '@/api/client';
import type {Preset, PromptItem} from '@/gen/muse/preset_pb';
import type {RegexRule} from '@/gen/muse/regex_pb';
import { Role, InjectionPosition, PromptItemIdentifier } from '@/gen/muse/common_pb';

// 设备检测
const isMobile = ref(false);
const showSidebar = ref(false);

// 检测是否为移动端
const detectMobile = () => {
  isMobile.value = window.innerWidth <= 768;
};

// 监听窗口大小变化
const handleResize = () => {
  detectMobile();
  // 在桌面端自动隐藏侧边栏
  if (!isMobile.value) {
    showSidebar.value = false;
  }
};

// 在组件挂载时检测
onMounted(() => {
  detectMobile();
  window.addEventListener('resize', handleResize);
});

// 组件卸载时清理事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
// 组件本地使用的提示项接口
interface LocalPromptItem {
  id: number;
  identifier: PromptItemIdentifier;
  name: string;
  role: Role;
  content: string;
  enabled: boolean;
  marker?: boolean;
  forbidOverrides?: boolean; // 标记项：名称和内容不可编辑
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
const regexRules = ref<RegexRule[]>([]);
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
    const response = await presetClient.getPreset({ id: presetId });
    if (response.preset) {
      // 更新提示项列表
      promptItems.value = response.preset.promptItems || [];
    }
};

// 加载预设正则规则
const loadPresetRegexRules = async (presetId: number) => {
    const response = await regexRuleClient.listPresetRegexRules({ presetId: presetId });
    // 直接使用 pb 类型
    regexRules.value = response.rules.map(r => ({
      ...r,
      affectFlags: r.affectFlags || {
        userInput: false,
        aiOutput: false,
        slashCommand: false,
        worldInfo: false,
        prompt: false,
        $typeName: 'muse.RegexAffectFlags'
      }
    }));
};

// 选择预设
const selectPreset = async (preset: Preset) => {
  selectedPreset.value = preset;
  expandedPromptId.value = null;
  
  // 移动端选择后自动关闭侧边栏
  if (isMobile.value) {
    showSidebar.value = false;
  }

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
        presencePenalty: 0,
        // 创建预设时自动添加默认提示项
        promptItems: defaultPromptItems.map(item => ({
          identifier: item.identifier as PromptItemIdentifier,
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
        await selectPreset(response.preset);
        message.success('预设已创建');
      }
    } finally {
      // 创建完成后自动处理
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
          promptItems: originalPreset.promptItems?.map((item, index) => ({
            identifier: item.identifier,
            name: item.name,
            content: item.content,
            role: item.role,
            isEnabled: item.isEnabled,
            injectionPosition: item.injectionPosition,
            injectionDepth: item.injectionDepth,
            forbidOverrides: item.forbidOverrides,
            sortOrder: index
          }))
        });
        if (response.preset) {
          presets.value.push(response.preset);
          message.success('预设已复制');
        }
      } finally {
        // 复制完成后自动处理
      }
      break;

    case 'export':
      try {
        const exportRes = await presetClient.exportPreset({ id: preset.id });
        // 下载文件
        const blob = new Blob([new Uint8Array(exportRes.fileContent).buffer], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = exportRes.fileName || `${preset.name}.json`;
        a.click();
        URL.revokeObjectURL(url);
        message.success('预设已导出');
      } finally {
        // 导出完成后自动处理
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
          } finally {
            // 删除完成后自动处理
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
  editingPromptItem.value = {
    id: item.id,
    identifier: item.identifier,
    name: item.name,
    role: item.role,
    content: item.content,
    enabled: item.isEnabled,
    forbidOverrides: item.forbidOverrides // 传递禁止覆盖标记
  };
  showPromptModal.value = true;
};

// 保存提示项
const handleSavePromptItem = async (itemData: Partial<LocalPromptItem>) => {
  const role = itemData.role ?? Role.System;
  if (!selectedPreset.value) return;

  try {
    if (editingPromptItem.value) {
      // 更新现有项 - 保留原有的 forbidOverrides 状态
      const originalItem = promptItems.value.find(p => p.id === editingPromptItem.value!.id);
      const response = await presetClient.updatePromptItem({
        id: editingPromptItem.value.id,
        identifier: itemData.identifier ?? PromptItemIdentifier.PromptItemIdentifierUnspecified,
        name: itemData.name || '',
        content: itemData.content,
        role: role,
        isEnabled: itemData.enabled ?? true,
        injectionPosition: InjectionPosition.Relative,
        injectionDepth: 0,
        forbidOverrides: itemData.forbidOverrides ?? originalItem?.forbidOverrides ?? false
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
        identifier: itemData.identifier ?? PromptItemIdentifier.PromptItemIdentifierUnspecified,
        name: itemData.name || '',
        content: itemData.content,
        role: role,
        isEnabled: itemData.enabled ?? true,
        injectionPosition: InjectionPosition.Relative,
        injectionDepth: 0,
        forbidOverrides: false
      });
      if (response.item) {
        promptItems.value.push(response.item);
        message.success('提示项已添加');
      }
    }
    showPromptModal.value = false;
    editingPromptItem.value = null;
  } finally {
    // 保存完成后自动处理
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
      } finally {
        // 删除完成后自动处理
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
      forbidOverrides: item.forbidOverrides
    });
    if (response.item) {
      const index = promptItems.value.findIndex(p => p.id === item.id);
      if (index >= 0) {
        promptItems.value[index] = response.item;
      }
    }
  } finally {
    // 切换完成后自动处理
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
  } finally {
    // 排序完成后自动处理
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
      forbidOverrides: item.forbidOverrides
    });
  } finally {
    // 内容更新完成后自动处理
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

// 默认提示项模板 - 参照 SillyTavern 的 chatCompletionDefaultPrompts
// marker: true 表示系统标记项，仅作为占位符使用，由系统自动填充内容
const defaultPromptItems = [
  {
    identifier: PromptItemIdentifier.Main,
    name: '主提示词',
    content: '在 {{char}} 和 {{user}} 的虚构聊天中，撰写 {{char}} 的下一条回复。',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: false,
    sortOrder: 0
  },
  {
    identifier: PromptItemIdentifier.WorldInfoBefore,
    name: '世界信息（前）',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true, // 系统标记项禁止覆盖
    sortOrder: 1
  },
  {
    identifier: PromptItemIdentifier.PersonaDescription,
    name: '人设描述',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 2
  },
  {
    identifier: PromptItemIdentifier.CharDescription,
    name: '角色描述',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 3
  },
  {
    identifier: PromptItemIdentifier.CharPersonality,
    name: '角色性格',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 4
  },
  {
    identifier: PromptItemIdentifier.Scenario,
    name: '场景',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 5
  },
  {
    identifier: PromptItemIdentifier.Nsfw,
    name: '辅助提示词',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: false,
    sortOrder: 6
  },
  {
    identifier: PromptItemIdentifier.WorldInfoAfter,
    name: '世界信息（后）',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 7
  },
  {
    identifier: PromptItemIdentifier.DialogueExamples,
    name: '对话示例',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 8
  },
  {
    identifier: PromptItemIdentifier.ChatHistory,
    name: '聊天记录',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: true,
    sortOrder: 9
  },
  {
    identifier: PromptItemIdentifier.Jailbreak,
    name: '越狱提示词',
    content: '',
    role: Role.System,
    isEnabled: true,
    injectionPosition: InjectionPosition.Relative,
    injectionDepth: 0,
    forbidOverrides: false,
    sortOrder: 10
  }
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
  position: relative;
}

/* 左侧预设列表 */
.presets-sidebar {
  width: 280px;
  min-width: 280px;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  transition: all var(--transition-normal);
}

/* 移动端侧边栏 */
@media (max-width: 768px) {
  .presets-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    bottom: 0;
    width: 300px;
    z-index: 100;
    transform: translateX(-100%);
    box-shadow: var(--shadow-xl);
  }
  
  .presets-sidebar.show-mobile {
    transform: translateX(0);
  }
  
  .presets-view {
    gap: 0;
  }
}

.sidebar-header {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}

.sidebar-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-header-content h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-sidebar-btn {
  margin-left: auto;
}

.search-input {
  margin: 0 16px;
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
  background: var(--color-primary);
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-secondary) 100%);
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
  position: relative;
}

/* 移动端列表切换按钮 */
.mobile-list-toggle {
  padding: 12px 20px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}

.list-toggle-btn {
  width: 100%;
  justify-content: flex-start;
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
  padding: 16px 20px 16px 20px;
}

/* 设置分区 */
.settings-section {
  margin-bottom: 24px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  overflow: hidden;
  /* 右侧间距调整：修改 margin-right 的值来控制设置区块右侧空白大小 */
  margin-right: 12px;
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

.prompt-item.forbid-overrides {
  background: var(--bg-tertiary);
  border-style: dashed;
}

.prompt-item.forbid-overrides .prompt-info {
  cursor: default;
}

.prompt-marker-hint {
  margin: 4px 0 0;
  font-size: 12px;
  color: var(--text-tertiary);
  font-style: italic;
}

.prompt-item.editing {
  border-color: var(--color-primary);
  box-shadow: var(--glow-soft);
}

.prompt-ghost {
  opacity: 0.5;
  background: var(--color-primary);
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

.prompt-tags {
  display: flex;
  gap: 4px;
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
  color: var(--color-primary);
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
    height: calc(100vh - 64px - 48px);
  }

  .presets-panel {
    flex: 1;
  }
  
  .panel-header {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
    padding: 12px 16px;
  }

  .header-title {
    flex-wrap: wrap;
  }

  .header-actions {
    justify-content: flex-end;
  }

  .preset-name-input {
    max-width: none;
    font-size: 16px;
  }

  .preset-name-input :deep(input) {
    font-size: 16px;
  }

  .param-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .param-control {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .param-input {
    width: 100%;
  }

  .prompt-item {
    padding: 10px 12px;
    gap: 8px;
  }

  .prompt-header-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .prompt-tags {
    flex-wrap: wrap;
  }

  .prompt-preview {
    font-size: 12px;
  }

  .prompt-actions {
    margin-left: auto;
  }

  .section-content {
    padding: 12px;
  }

  .settings-section {
    margin-bottom: 16px;
    margin-right: 0;
  }
  
  .panel-content {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .presets-view {
    margin: 0 -16px;
    border-radius: 0;
    border: none;
  }

  .sidebar-header,
  .panel-header,
  .panel-content {
    padding-left: 16px;
    padding-right: 16px;
  }

  .preset-item {
    padding: 8px 10px;
  }

  .prompt-item {
    padding: 8px 10px;
    gap: 6px;
  }

  .drag-handle {
    padding: 2px;
  }

  .prompt-index {
    width: 20px;
    height: 20px;
    font-size: 11px;
  }

  .prompt-name {
    font-size: 13px;
  }

  .prompt-preview {
    font-size: 11px;
    margin-top: 4px;
  }

  .param-label {
    font-size: 12px;
  }

  .section-title {
    font-size: 13px;
  }

  .variables-grid {
    gap: 6px;
  }

  .variable-tag {
    font-size: 11px;
    padding: 3px 8px;
  }
  
  .mobile-list-toggle {
    padding: 12px 16px;
  }
}

/* 测试按钮 */
.test-buttons {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
  background: var(--bg-secondary);
  padding: 12px;
  border-radius: 8px;
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow-lg);
}

@media (max-width: 768px) {
  .test-buttons {
    bottom: 10px;
    right: 10px;
    padding: 8px;
  }
  
  .test-buttons .n-button {
    font-size: 12px;
    padding: 8px 12px;
  }
}
</style>
