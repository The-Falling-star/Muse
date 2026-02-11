<template>
  <div class="model-config-tab">
    <div class="config-section">
      <div class="section-label">Temperature</div>
      <div class="slider-row">
        <n-slider
          v-model:value="temperature"
          :min="0"
          :max="2"
          :step="0.01"
          :tooltip="true"
        />
        <n-input-number
          v-model:value="temperature"
          :min="0"
          :max="2"
          :step="0.01"
          size="small"
          class="slider-input"
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
        />
        <n-input-number
          v-model:value="topP"
          :min="0"
          :max="1"
          :step="0.01"
          size="small"
          class="slider-input"
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
        />
        <n-input-number
          v-model:value="topK"
          :min="0"
          :max="100"
          :step="1"
          size="small"
          class="slider-input"
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
        />
        <n-input-number
          v-model:value="maxTokens"
          :min="1"
          :max="32768"
          :step="1"
          size="small"
          class="slider-input"
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
        />
        <n-input-number
          v-model:value="frequencyPenalty"
          :min="-2"
          :max="2"
          :step="0.01"
          size="small"
          class="slider-input"
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
        />
        <n-input-number
          v-model:value="presencePenalty"
          :min="-2"
          :max="2"
          :step="0.01"
          size="small"
          class="slider-input"
        />
      </div>
    </div>

    <!-- 高级设置 -->
    <n-collapse>
      <n-collapse-item title="高级设置" name="advanced">
        <div class="config-section">
          <div class="section-label">停止序列</div>
          <n-dynamic-tags v-model:value="stopSequences" />
        </div>
        <div class="config-section">
          <div class="section-label">代理地址</div>
          <n-input
            v-model:value="proxyUrl"
            placeholder="https://..."
            size="small"
          />
        </div>
      </n-collapse-item>
    </n-collapse>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import {
  NSlider,
  NInputNumber,
  NInput,
  NCollapse,
  NCollapseItem,
  NDynamicTags
} from 'naive-ui';

// TODO: 后续对接真实的模型配置 store
const temperature = ref(0.7);
const topP = ref(0.9);
const topK = ref(40);
const maxTokens = ref(2048);
const frequencyPenalty = ref(0);
const presencePenalty = ref(0);
const stopSequences = ref<string[]>([]);
const proxyUrl = ref('');
</script>

<style scoped>
.model-config-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 20px;
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
</style>
