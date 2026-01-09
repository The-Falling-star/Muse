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
              <span class="preset-item-count">{{ preset.prompts.length }} 项</span>
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

        <n-empty v-if="filteredPresets.length === 0" description="暂无预设" size="small" />
      </n-scrollbar>
    </div>

    <!-- 右侧：预设编辑面板 -->
    <div class="presets-panel">
      <template v-if="selectedPreset">
        <!-- 预设头部信息 -->
        <div class="panel-header">
          <div class="header-title">
            <n-input
              v-model:value="selectedPreset.name"
              :bordered="false"
              placeholder="预设名称"
              class="preset-name-input"
            />
          </div>
          <div class="header-actions">
            <n-button type="primary" size="small" @click="savePreset">
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
                      v-model:value="selectedPreset.temperature"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.temperature"
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
                      v-model:value="selectedPreset.topP"
                      :min="0"
                      :max="1"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.topP"
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
                      v-model:value="selectedPreset.topK"
                      :min="0"
                      :max="500"
                      :step="1"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.topK"
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
                      v-model:value="selectedPreset.maxTokens"
                      :min="1"
                      :max="9999999"
                      :step="10000"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.maxTokens"
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
                      v-model:value="selectedPreset.frequencyPenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.frequencyPenalty"
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
                      v-model:value="selectedPreset.presencePenalty"
                      :min="0"
                      :max="2"
                      :step="0.05"
                      :tooltip="false"
                      class="param-slider"
                    />
                    <n-input-number
                      v-model:value="selectedPreset.presencePenalty"
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
              <PresetRegexManager v-model:regexRules="regexRules" />
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
                v-model="orderedPromptItems"
                item-key="identifier"
                handle=".drag-handle"
                animation="200"
                ghost-class="prompt-ghost"
                class="prompt-list"
              >
                <template #item="{ element: item, index }">
                  <div
                    class="prompt-item"
                    :class="{
                      'disabled': !getPromptEnabled(item.identifier),
                      'marker': item.marker,
                      'editing': editingPromptId === item.id
                    }"
                  >
                    <!-- 拖拽手柄 -->
                    <div class="drag-handle">
                      <n-icon size="16"><ReorderTwoOutline /></n-icon>
                    </div>

                    <!-- 启用开关 -->
                    <n-switch
                      :value="getPromptEnabled(item.identifier)"
                      size="small"
                      @update:value="(val: boolean) => togglePromptEnabled(item.identifier, val)"
                    />

                    <!-- 序号 -->
                    <span class="prompt-index">{{ index + 1 }}</span>

                    <!-- 提示项信息 -->
                    <div class="prompt-info" @click="!item.marker && toggleExpand(item.id)">
                      <div class="prompt-header-row">
                        <span class="prompt-name">{{ item.name }}</span>
                        <!-- 系统标记项只显示"系统占位"标签，不显示角色标签 -->
                        <n-tag v-if="item.marker" type="warning" size="small" :bordered="false">
                          系统占位
                        </n-tag>
                        <!-- 用户自定义提示项显示角色标签 -->
                        <n-tag v-else :type="getRoleColor(item.role)" size="small" :bordered="false">
                          {{ getRoleName(item.role) }}
                        </n-tag>
                      </div>
                      <p v-if="!item.marker && item.content" class="prompt-preview">
                        {{ truncateContent(item.content) }}
                      </p>
                      <p v-if="item.marker" class="prompt-marker-hint">
                        {{ getMarkerHint(item.identifier) }}
                      </p>
                    </div>

                    <!-- 操作按钮 - 系统标记项不显示编辑和删除按钮 -->
                    <div v-if="!item.marker" class="prompt-actions">
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
                      <div v-if="expandedPromptId === item.id && !item.marker" class="prompt-expanded">
                        <n-input
                          v-model:value="item.content"
                          type="textarea"
                          :autosize="{ minRows: 3, maxRows: 10 }"
                          placeholder="输入提示内容，支持 {{char}} {{user}} 等变量"
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
import { ref, computed, watch } from 'vue';
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
import type { Preset, PromptItem, RegexRule } from '../types';
import { DEFAULT_PROMPT_ITEMS, DEFAULT_PROMPT_ORDER } from '../types';

const message = useMessage();
const dialog = useDialog();

// 状态
const searchQuery = ref('');
const selectedPreset = ref<Preset | null>(null);
const showPromptModal = ref(false);
const editingPromptItem = ref<PromptItem | null>(null);
const editingPromptId = ref<string | null>(null);
const expandedPromptId = ref<string | null>(null);
const regexRules = ref<RegexRule[]>([]);

// 预设菜单选项
const presetMenuOptions = [
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

// 模拟预设数据
const presets = ref<Preset[]>([
  {
    id: '1',
    name: 'Default',
    description: '默认预设',
    prompts: JSON.parse(JSON.stringify(DEFAULT_PROMPT_ITEMS)),
    promptOrder: JSON.parse(JSON.stringify(DEFAULT_PROMPT_ORDER)),
    temperature: 1,
    topP: 1,
    topK: 0,
    maxTokens: 300,
    frequencyPenalty: 0,
    presencePenalty: 0
  },
  {
    id: '2',
    name: '角色扮演增强',
    description: '适合角色扮演场景',
    prompts: [
      ...JSON.parse(JSON.stringify(DEFAULT_PROMPT_ITEMS)),
      {
        id: 'rp_enhance',
        identifier: 'rp_enhance',
        name: 'RP增强提示',
        role: 'system' as const,
        content: '请深入扮演{{char}}，保持角色一致性，使用生动的描写...',
        enabled: true,
        marker: false
      }
    ],
    promptOrder: [
      ...JSON.parse(JSON.stringify(DEFAULT_PROMPT_ORDER)),
      { identifier: 'rp_enhance', enabled: true }
    ],
    temperature: 1.2,
    topP: 0.95,
    topK: 0,
    maxTokens: 500,
    frequencyPenalty: 0,
    presencePenalty: 0
  }
]);

// 计算属性：过滤后的预设列表
const filteredPresets = computed(() => {
  if (!searchQuery.value) return presets.value;
  const query = searchQuery.value.toLowerCase();
  return presets.value.filter(p =>
    p.name.toLowerCase().includes(query) ||
    p.description?.toLowerCase().includes(query)
  );
});

// 计算属性：按排序顺序排列的提示项
const orderedPromptItems = computed({
  get() {
    if (!selectedPreset.value) return [];
    const orderMap = new Map(
      selectedPreset.value.promptOrder.map((entry, index) => [entry.identifier, index])
    );
    return [...selectedPreset.value.prompts].sort((a, b) => {
      const indexA = orderMap.get(a.identifier) ?? 999;
      const indexB = orderMap.get(b.identifier) ?? 999;
      return indexA - indexB;
    });
  },
  set(newOrder: PromptItem[]) {
    if (!selectedPreset.value) return;
    // 更新排序
    selectedPreset.value.promptOrder = newOrder.map(item => ({
      identifier: item.identifier,
      enabled: getPromptEnabled(item.identifier)
    }));
  }
});

// 方法：获取提示项是否启用
const getPromptEnabled = (identifier: string): boolean => {
  if (!selectedPreset.value) return false;
  const entry = selectedPreset.value.promptOrder.find(e => e.identifier === identifier);
  return entry?.enabled ?? true;
};

// 方法：切换提示项启用状态
const togglePromptEnabled = (identifier: string, enabled: boolean) => {
  if (!selectedPreset.value) return;
  const entry = selectedPreset.value.promptOrder.find(e => e.identifier === identifier);
  if (entry) {
    entry.enabled = enabled;
  }
};

// 方法：选择预设
const selectPreset = (preset: Preset) => {
  selectedPreset.value = preset;
  expandedPromptId.value = null;
  // 同步正则规则
  regexRules.value = preset.regexRules ? [...preset.regexRules] : [];
};

// 方法：创建预设
const createPreset = () => {
  const newPreset: Preset = {
    id: Date.now().toString(),
    name: '新预设',
    description: '',
    prompts: JSON.parse(JSON.stringify(DEFAULT_PROMPT_ITEMS)),
    promptOrder: JSON.parse(JSON.stringify(DEFAULT_PROMPT_ORDER)),
    temperature: 1,
    topP: 1,
    topK: 0,
    maxTokens: 300,
    frequencyPenalty: 0,
    presencePenalty: 0
  };
  presets.value.push(newPreset);
  selectedPreset.value = newPreset;
  regexRules.value = [];
  message.success('预设已创建');
};

// 方法：处理预设菜单
const handlePresetMenu = (key: string, preset: Preset) => {
  switch (key) {
    case 'copy':
      const copied: Preset = {
        ...JSON.parse(JSON.stringify(preset)),
        id: Date.now().toString(),
        name: `${preset.name} (副本)`
      };
      presets.value.push(copied);
      message.success('预设已复制');
      break;
    case 'export':
      // TODO: 导出功能
      message.info('导出功能开发中');
      break;
    case 'delete':
      dialog.warning({
        title: '确认删除',
        content: `确定要删除预设"${preset.name}"吗？`,
        positiveText: '删除',
        negativeText: '取消',
        onPositiveClick: () => {
          presets.value = presets.value.filter(p => p.id !== preset.id);
          if (selectedPreset.value?.id === preset.id) {
            selectedPreset.value = presets.value[0] ?? null;
          }
          message.success('预设已删除');
        }
      });
      break;
  }
};

// 方法：保存预设
const savePreset = () => {
  if (selectedPreset.value) {
    // 同步正则规则到预设
    selectedPreset.value.regexRules = [...regexRules.value];
  }
  // TODO: 调用API保存
  message.success('预设已保存');
};

// 方法：添加提示项
const addPromptItem = () => {
  editingPromptItem.value = null;
  showPromptModal.value = true;
};

// 方法：编辑提示项
const editPromptItem = (item: PromptItem) => {
  editingPromptItem.value = item;
  showPromptModal.value = true;
};

// 方法：保存提示项
const handleSavePromptItem = (item: PromptItem) => {
  if (!selectedPreset.value) return;

  if (editingPromptItem.value) {
    // 更新现有项
    const index = selectedPreset.value.prompts.findIndex(
      p => p.id === editingPromptItem.value!.id
    );
    if (index >= 0) {
      selectedPreset.value.prompts[index] = item;
    }
    message.success('提示项已更新');
  } else {
    // 添加新项
    selectedPreset.value.prompts.push(item);
    selectedPreset.value.promptOrder.push({
      identifier: item.identifier,
      enabled: true
    });
    message.success('提示项已添加');
  }
  showPromptModal.value = false;
  editingPromptItem.value = null;
};

// 方法：删除提示项
const deletePromptItem = (item: PromptItem) => {
  if (!selectedPreset.value || item.marker) return;

  dialog.warning({
    title: '确认删除',
    content: `确定要删除提示项"${item.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      selectedPreset.value!.prompts = selectedPreset.value!.prompts.filter(
        p => p.id !== item.id
      );
      selectedPreset.value!.promptOrder = selectedPreset.value!.promptOrder.filter(
        e => e.identifier !== item.identifier
      );
      message.success('提示项已删除');
    }
  });
};

// 方法：切换展开/收起
const toggleExpand = (itemId: string) => {
  expandedPromptId.value = expandedPromptId.value === itemId ? null : itemId;
};

// 辅助方法
const getRoleName = (role: string) => {
  const names: Record<string, string> = {
    system: '系统',
    user: '用户',
    assistant: '助手'
  };
  return names[role] || role;
};

const getRoleColor = (role: string): 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error' => {
  const colors: Record<string, 'default' | 'primary' | 'info' | 'success' | 'warning' | 'error'> = {
    system: 'primary',
    user: 'success',
    assistant: 'info'
  };
  return colors[role] || 'default';
};

const truncateContent = (content: string, maxLength = 100) => {
  if (content.length <= maxLength) return content;
  return content.slice(0, maxLength) + '...';
};

// 系统标记项的说明文字
const getMarkerHint = (identifier: string): string => {
  const hints: Record<string, string> = {
    chatHistory: '此处将自动插入聊天历史记录',
    worldInfoBefore: '此处将自动插入世界书内容（前置）',
    worldInfoAfter: '此处将自动插入世界书内容（后置）',
    charDescription: '此处将自动插入角色描述',
    charPersonality: '此处将自动插入角色性格',
    scenario: '此处将自动插入场景设定',
    dialogueExamples: '此处将自动插入对话示例'
  };
  return hints[identifier] || '系统自动填充内容';
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

// 初始化选择第一个预设
watch(presets, (val) => {
  if (val.length > 0 && !selectedPreset.value) {
    selectedPreset.value = val[0] ?? null;
    if (selectedPreset.value) {
      regexRules.value = selectedPreset.value.regexRules ? [...selectedPreset.value.regexRules] : [];
    }
  }
}, { immediate: true });
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
