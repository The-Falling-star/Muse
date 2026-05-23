<template>
  <n-drawer
    :show="visible"
    :width="isMobile ? '100%' : 500"
    placement="right"
    :native-scrollbar="true"
    @update:show="$emit('update:visible', $event)"
  >
    <n-drawer-content title="设置" closable>
      <div class="settings-drawer">
        <!-- 通用设置 -->
        <div class="settings-group">
          <div class="group-title">通用设置</div>

          <div class="settings-item">
            <div class="item-label">语言</div>
            <n-select
              v-model:value="language"
              :options="languageOptions"
              size="small"
              class="item-control"
            />
          </div>

          <div class="settings-item">
            <div class="item-label">主题</div>
            <n-select
              v-model:value="currentTheme"
              :options="themeOptions"
              size="small"
              class="item-control"
              @update:value="handleThemeChange"
            />
          </div>

          <div class="settings-item">
            <div class="item-label">字体大小</div>
            <n-select
              v-model:value="fontSize"
              :options="fontSizeOptions"
              size="small"
              class="item-control"
            />
          </div>
        </div>

        <!-- 对话设置 -->
        <div class="settings-group">
          <div class="group-title">对话设置</div>

          <div class="settings-item">
            <div class="item-label">默认打招呼消息</div>
            <n-switch v-model:value="showGreeting" size="small" />
          </div>

          <div class="settings-item">
            <div class="item-label">历史消息条数</div>
            <n-input-number
              v-model:value="historyCount"
              :min="1"
              :max="100"
              size="small"
              class="item-control-sm"
            />
          </div>

          <div class="settings-item">
            <div class="item-label">流式输出</div>
            <n-switch v-model:value="streamOutput" size="small" />
          </div>
        </div>

        <!-- 账号管理 -->
        <div class="settings-group">
          <div class="group-title">账号管理</div>

          <div v-if="currentUser" class="user-info">
            <n-avatar :size="40" round>
              {{ currentUser.username?.charAt(0)?.toUpperCase() || 'U' }}
            </n-avatar>
            <div class="user-detail">
              <div class="user-name">{{ currentUser.username }}</div>
              <div class="user-role">{{ currentUser.username ? '用户' : '未知' }}</div>
            </div>
          </div>

          <div class="account-actions">
            <n-button size="small" block @click="handleChangePassword">
              修改密码
            </n-button>
            <n-button size="small" type="error" block ghost @click="handleLogout">
              退出登录
            </n-button>
          </div>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import {
  NDrawer,
  NDrawerContent,
  NSelect,
  NInputNumber,
  NSwitch,
  NButton,
  NAvatar,
  useDialog
} from 'naive-ui';
import { useMediaQuery } from '@vueuse/core';
import { useThemeStore } from '@/stores/theme';
import { useUserStore } from '@/stores/user';

const props = defineProps<{
  visible: boolean;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
}>();

const router = useRouter();
const dialog = useDialog();
const themeStore = useThemeStore();
const userStore = useUserStore();

const isMobile = useMediaQuery('(max-width: 767px)');

// ====== 通用设置 ======
const language = ref('zh-CN');
const languageOptions = [
  { label: '简体中文', value: 'zh-CN' },
  { label: 'English', value: 'en-US' }
];

const currentTheme = ref(themeStore.theme);
const themeOptions = [
  { label: '亮色', value: 'light' },
  { label: '暗色', value: 'dark' },
  { label: '跟随系统', value: 'system' }
];

const handleThemeChange = (value: string) => {
  if (value === 'light' || value === 'dark') {
    themeStore.setTheme(value);
  }
  // TODO: 支持跟随系统
};

const fontSize = ref('14');
const fontSizeOptions = [
  { label: '小 (12px)', value: '12' },
  { label: '默认 (14px)', value: '14' },
  { label: '大 (16px)', value: '16' },
  { label: '特大 (18px)', value: '18' }
];

// ====== 对话设置 ======
const showGreeting = ref(true);
const historyCount = ref(20);
const streamOutput = ref(true);

// ====== 账号管理 ======
const currentUser = computed(() => userStore.currentUser);

const handleChangePassword = () => {
  emit('update:visible', false);
  router.push('/settings');
};

const handleLogout = () => {
  dialog.warning({
    title: '确认退出',
    content: '确定要退出登录吗？',
    positiveText: '退出',
    negativeText: '取消',
    onPositiveClick: () => {
      userStore.logout();
      emit('update:visible', false);
      router.push('/login');
    }
  });
};
</script>

<style scoped>
.settings-drawer {
  display: flex;
  flex-direction: column;
  gap: 24px;
  padding-bottom: 24px;
}

.settings-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.group-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: .5px;
  padding-bottom: 4px;
  border-bottom: 1px solid var(--border-color);
}

.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 36px;
}

.item-label {
  font-size: 14px;
  color: var(--text-primary);
  white-space: nowrap;
}

.item-control {
  width: 160px;
  flex-shrink: 0;
}

.item-control-sm {
  width: 100px;
  flex-shrink: 0;
}


/* ====== 账号管理 ====== */
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  background: var(--bg-secondary);
}

.user-detail {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.user-role {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 2px;
}

.account-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
