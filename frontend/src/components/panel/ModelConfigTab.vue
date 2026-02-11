<template>
  <div class="model-config-tab">
    <!-- 无活跃预设提示 -->
    <n-empty
      v-if="!loading && !activePreset"
      description="暂无活跃预设，请先在预设列表中设置一个活跃预设"
      style="margin-top: 40px;"
    />

    <template v-if="activePreset">
      <!-- 当前预设名称 -->
      <div class="active-preset-header">
        <span class="preset-label">当前预设：</span>
        <span class="preset-name">{{ activePreset.name }}</span>
      </div>

      <div class="config-section">
        <div class="section-label">Temperature</div>
        <div class="slider-row">
          <n-slider
            v-model:value="temperature"
            :min="0"
            :max="2"
            :step="0.01"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="temperature"
            :min="0"
            :max="2"
            :step="0.01"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-label">Top P</div>
        <div class="slider-row">
          <n-slider
            v-model:value="topP"
            :min="0"
            :max="1"
            :step="0.01"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="topP"
            :min="0"
            :max="1"
            :step="0.01"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-label">Top K</div>
        <div class="slider-row">
          <n-slider
            v-model:value="topK"
            :min="0"
            :max="100"
            :step="1"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="topK"
            :min="0"
            :max="100"
            :step="1"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-label">Max Tokens</div>
        <div class="slider-row">
          <n-slider
            v-model:value="maxTokens"
            :min="1"
            :max="32768"
            :step="1"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="maxTokens"
            :min="1"
            :max="32768"
            :step="1"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-label">频率惩罚</div>
        <div class="slider-row">
          <n-slider
            v-model:value="frequencyPenalty"
            :min="-2"
            :max="2"
            :step="0.01"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="frequencyPenalty"
            :min="-2"
            :max="2"
            :step="0.01"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <div class="config-section">
        <div class="section-label">存在惩罚</div>
        <div class="slider-row">
          <n-slider
            v-model:value="presencePenalty"
            :min="-2"
            :max="2"
            :step="0.01"
            :tooltip="true"
            @update:value="debouncedSave"
          />
          <n-input-number
            v-model:value="presencePenalty"
            :min="-2"
            :max="2"
            :step="0.01"
            size="small"
            class="slider-input"
            @update:value="debouncedSave"
          />
        </div>
      </div>

      <!-- 保存状态提示 -->
      <div v-if="saveStatus" class="save-status" :class="saveStatus">
        {{ saveStatus === 'saving' ? '保存中...' : saveStatus === 'saved' ? '已保存' : '保存失败' }}
      </div>
    </template>

    <!-- 加载遮罩 -->
    <n-spin v-if="loading" class="loading-spin" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue';
import {
  NSlider,
  NInputNumber,
  NSpin,
  NEmpty
} from 'naive-ui';
import { useUserStore } from '@/stores/user';
import { usePresetStore } from '@/stores/preset';
import type { Preset } from '@/gen/muse/muse_pb';

const userStore = useUserStore();
const presetStore = usePresetStore();

// 加载状态
const loading = ref(false);
// 保存状态：null | 'saving' | 'saved' | 'error'
const saveStatus = ref<string | null>(null);
// 当前活跃预设
const activePreset = ref<Preset | null>(null);

// 模型配置参数
const temperature = ref(0.7);
const topP = ref(0.9);
const topK = ref(40);
const maxTokens = ref(2048);
const frequencyPenalty = ref(0);
const presencePenalty = ref(0);

// 从预设数据同步到本地表单
const syncFromPreset = (preset: Preset) => {
  temperature.value = preset.temperature;
  topP.value = preset.topP;
  topK.value = preset.topK;
  maxTokens.value = preset.maxTokens;
  frequencyPenalty.value = preset.frequencyPenalty;
  presencePenalty.value = preset.presencePenalty;
};

// 加载活跃预设
const loadActivePreset = async () => {
  const activePresetId = userStore.currentUser?.activePresetId;
  if (!activePresetId) {
    activePreset.value = null;
    return;
  }

  loading.value = true;
  try {
    const preset = await presetStore.fetchPreset(activePresetId);
    if (!preset) {
      activePreset.value = null;
      return;
    }
    activePreset.value = preset;
    syncFromPreset(preset);
  } catch (error) {
    console.error('加载活跃预设失败:', error);
    activePreset.value = null;
  } finally {
    loading.value = false;
  }
};

// 保存模型配置到后端
const saveConfig = async () => {
  if (!activePreset.value) return;

  saveStatus.value = 'saving';
  try {
    const updated = await presetStore.updatePreset(activePreset.value.id, {
      name: activePreset.value.name,
      temperature: temperature.value,
      topP: topP.value,
      topK: topK.value,
      maxTokens: maxTokens.value,
      frequencyPenalty: frequencyPenalty.value,
      presencePenalty: presencePenalty.value,
      version: activePreset.value.version,
    });
    if (updated) {
      activePreset.value = updated;
    }
    saveStatus.value = 'saved';
    // 2秒后清除保存提示
    setTimeout(() => {
      if (saveStatus.value === 'saved') {
        saveStatus.value = null;
      }
    }, 2000);
  } catch (error) {
    console.error('保存模型配置失败:', error);
    saveStatus.value = 'error';
    setTimeout(() => {
      if (saveStatus.value === 'error') {
        saveStatus.value = null;
      }
    }, 3000);
  }
};

// 防抖保存（用户调整滑块时不会频繁请求）
let saveTimer: ReturnType<typeof setTimeout> | null = null;
const debouncedSave = () => {
  if (saveTimer) {
    clearTimeout(saveTimer);
  }
  saveTimer = setTimeout(() => {
    saveConfig();
  }, 800);
};

// 监听活跃预设ID变化，重新加载
watch(
  () => userStore.currentUser?.activePresetId,
  (newId) => {
    if (newId) {
      loadActivePreset();
    } else {
      activePreset.value = null;
    }
  }
);

onMounted(() => {
  loadActivePreset();
});
</script>

<style scoped>
.model-config-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: relative;
}

.active-preset-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--n-color-hover, rgba(255, 255, 255, .04));
  border-radius: 6px;
  font-size: 13px;
}

.preset-label {
  color: var(--text-secondary);
  flex-shrink: 0;
}

.preset-name {
  font-weight: 500;
  color: var(--n-text-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.config-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.slider-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.slider-input {
  width: 80px;
  flex-shrink: 0;
}

.save-status {
  text-align: center;
  font-size: 12px;
  padding: 4px 0;
  border-radius: 4px;
  transition: all .3s ease;
}

.save-status.saving {
  color: var(--n-text-color-3, #999);
}

.save-status.saved {
  color: var(--n-success-color, #18a058);
}

.save-status.error {
  color: var(--n-error-color, #d03050);
}

.loading-spin {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}
</style>
