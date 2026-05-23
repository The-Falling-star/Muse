<template>
  <div class="character-detail">
    <!-- 头像区域 -->
    <div class="detail-header">
      <div class="avatar-section">
        <n-avatar
          :size="120"
          round
          :src="avatarUrl"
          class="detail-avatar"
        >
          {{ character.name.charAt(0) }}
        </n-avatar>
        <div class="avatar-glow"></div>
      </div>

      <h2 class="character-name">{{ character.name }}</h2>
    </div>

    <!-- 操作按钮 -->
    <div class="detail-actions">
      <n-button type="primary" @click="$emit('chat')">
        <template #icon>
          <n-icon><ChatbubbleOutline /></n-icon>
        </template>
        开始对话
      </n-button>
      <n-button @click="$emit('edit')">
        <template #icon>
          <n-icon><CreateOutline /></n-icon>
        </template>
        编辑
      </n-button>
    </div>

    <!-- 详细信息 -->
    <n-divider />

    <n-scrollbar class="detail-content">
      <n-collapse :default-expanded-names="['description', 'firstMessage']">
        <n-collapse-item title="描述" name="description">
          <p class="info-text">{{ character.description || '暂无描述' }}</p>
        </n-collapse-item>

        <n-collapse-item title="开场白" name="firstMessage">
          <p class="info-text">{{ character.firstMessage || '暂无开场白' }}</p>
        </n-collapse-item>

        <n-collapse-item title="示例对话" name="exampleDialogue">
          <div v-if="character.exampleDialogue && character.exampleDialogue.length > 0" class="example-dialogues">
            <div v-for="(dialogue, index) in character.exampleDialogue" :key="index" class="dialogue-item">
              <pre class="info-code">{{ dialogue }}</pre>
            </div>
          </div>
          <p v-else class="info-text">暂无示例对话</p>
        </n-collapse-item>

        <n-collapse-item title="创作者备注" name="creatorNotes">
          <p class="info-text">{{ character.creatorNotes || '暂无备注' }}</p>
        </n-collapse-item>
      </n-collapse>
    </n-scrollbar>

    <!-- 元信息 -->
    <div class="detail-meta">
      <span v-if="character.createdAt">
        创建于 {{ formatDate(character.createdAt) }}
      </span>
      <span v-if="character.updatedAt">
        更新于 {{ formatDate(character.updatedAt) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  NAvatar,
  NButton,
  NIcon,
  NDivider,
  NScrollbar,
  NCollapse,
  NCollapseItem
} from 'naive-ui';
import { ChatbubbleOutline, CreateOutline } from '@vicons/ionicons5';
import { toRef } from 'vue';

import type { Character } from '@/gen/muse/character_pb';
import { useAvatar } from '@/composables/useAvatar';

const props = defineProps<{
  character: Character;
}>();

defineEmits<{
  edit: [];
  chat: [];
}>();

// 使用头像加载hook
const { avatarUrl } = useAvatar(toRef(() => props.character.avatar));

const formatDate = (timestamp: bigint) => {
  return new Date(Number(timestamp)).toLocaleDateString('zh-CN');
};
</script>

<style scoped>
.character-detail {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 头部区域 */
.detail-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-bottom: 20px;
}

.avatar-section {
  position: relative;
  margin-bottom: 16px;
}

.detail-avatar {
  background: var(--gradient-primary);
  border: 4px solid var(--bg-card);
  box-shadow: var(--shadow-lg);
}

.avatar-glow {
  position: absolute;
  inset: -12px;
  background: var(--gradient-primary);
  border-radius: 50%;
  opacity: 0.3;
  filter: blur(24px);
  z-index: -1;
}

.character-name {
  font-size: 24px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 12px;
  text-align: center;
}

.character-tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
}

/* 操作按钮 */
.detail-actions {
  display: flex;
  gap: 12px;
  padding: 0 0 16px;
}

.detail-actions .n-button {
  flex: 1;
}

/* 内容区域 */
.detail-content {
  flex: 1;
  margin: 0 -24px;
  padding: 0 24px;
}

.info-text {
  font-size: 14px;
  line-height: 1.7;
  color: var(--text-secondary);
  margin: 0;
  white-space: pre-wrap;
}

.info-code {
  font-family: var(--font-mono, monospace);
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-secondary);
  background: var(--bg-tertiary);
  padding: 12px;
  border-radius: 8px;
  margin: 0 0 12px 0;
  overflow-x: auto;
  white-space: pre-wrap;
}

/* 示例对话列表样式 */
.example-dialogues {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.dialogue-item {
  position: relative;
}

.dialogue-item:not(:last-child)::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 0;
  right: 0;
  height: 1px;
  background: var(--border-color);
}

/* 元信息 */
.detail-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
  font-size: 12px;
  color: var(--text-tertiary);
}
</style>
