<template>
  <div class="character-card" @click="$emit('click')">
    <!-- 卡片头部 - 头像背景 -->
    <div class="card-header">
      <div class="avatar-container">
        <n-avatar
          :size="80"
          round
          :src="character.avatar"
          class="character-avatar"
        >
          {{ character.name.charAt(0) }}
        </n-avatar>
        <div class="avatar-glow"></div>
      </div>

      <!-- 操作按钮 -->
      <div class="card-actions">
        <n-dropdown :options="actionOptions" trigger="click" @select="handleAction">
          <n-button quaternary circle size="small" @click.stop>
            <template #icon>
              <n-icon><EllipsisVerticalOutline /></n-icon>
            </template>
          </n-button>
        </n-dropdown>
      </div>
    </div>

    <!-- 卡片内容 -->
    <div class="card-content">
      <h3 class="character-name">{{ character.name }}</h3>

      <p class="character-description">
        {{ character.description || '暂无描述' }}
      </p>
    </div>

    <!-- 卡片底部 -->
    <div class="card-footer">
      <n-button type="primary" size="small" @click.stop="$emit('chat', character)">
        <template #icon>
          <n-icon><ChatbubbleOutline /></n-icon>
        </template>
        开始对话
      </n-button>
    </div>

    <!-- 装饰性边框 -->
    <div class="card-border-glow"></div>
  </div>
</template>

<script setup lang="ts">
import { NAvatar, NButton, NIcon, NDropdown } from 'naive-ui';
import {
  EllipsisVerticalOutline,
  ChatbubbleOutline
} from '@vicons/ionicons5';

import type { Character } from '../../types';

const props = defineProps<{
  character: Character;
}>();

const emit = defineEmits<{
  click: [];
  edit: [character: Character];
  delete: [character: Character];
  chat: [character: Character];
}>();

// 操作菜单选项
const actionOptions = [
  { label: '编辑', key: 'edit' },
  { label: '复制', key: 'copy' },
  { label: '导出', key: 'export' },
  { type: 'divider', key: 'd1' },
  { label: '删除', key: 'delete' }
];

const handleAction = (key: string) => {
  switch (key) {
    case 'edit':
      emit('edit', props.character);
      break;
    case 'delete':
      emit('delete', props.character);
      break;
  }
};
</script>

<style scoped>
.character-card {
  position: relative;
  background: var(--gradient-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  transition: all var(--transition-normal);
}

.character-card:hover {
  transform: translateY(-4px);
  border-color: var(--border-glow);
  box-shadow: var(--shadow-lg), var(--glow-soft);
}

.character-card:hover .card-border-glow {
  opacity: 1;
}

.character-card:hover .avatar-glow {
  opacity: 0.6;
}

/* 边框发光效果 */
.card-border-glow {
  position: absolute;
  inset: -1px;
  background: var(--gradient-primary);
  border-radius: 16px;
  opacity: 0;
  transition: opacity var(--transition-normal);
  z-index: -1;
  filter: blur(8px);
}

/* 卡片头部 */
.card-header {
  position: relative;
  display: flex;
  justify-content: center;
  padding: 24px 16px 16px;
  background: linear-gradient(180deg, var(--bg-tertiary) 0%, transparent 100%);
}

.avatar-container {
  position: relative;
}

.character-avatar {
  background: var(--gradient-primary);
  border: 3px solid var(--bg-card);
  box-shadow: var(--shadow-md);
}

.avatar-glow {
  position: absolute;
  inset: -8px;
  background: var(--gradient-primary);
  border-radius: 50%;
  opacity: 0;
  filter: blur(20px);
  transition: opacity var(--transition-normal);
  z-index: -1;
}

.card-actions {
  position: absolute;
  top: 12px;
  right: 12px;
}

.card-actions .n-button {
  color: var(--text-tertiary);
  background: var(--bg-card);
}

.card-actions .n-button:hover {
  color: var(--color-primary);
}

/* 卡片内容 */
.card-content {
  padding: 0 20px 16px;
  text-align: center;
}

.character-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.character-description {
  font-size: 13px;
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 卡片底部 */
.card-footer {
  padding: 12px 20px 20px;
  display: flex;
  justify-content: center;
}

.card-footer .n-button {
  width: 100%;
}
</style>
