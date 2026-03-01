<template>
  <div class="login-container">
    <!-- 简洁背景 -->
    <div class="login-bg">
      <div class="bg-glow"></div>
    </div>

    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- Logo和标题 -->
      <div class="login-header">
        <div class="logo">
          <div class="logo-icon">
            <n-icon size="48">
              <SparklesOutline />
            </n-icon>
          </div>
        </div>
        <h1 class="title">Muse</h1>
        <p class="subtitle">AI 对话助手</p>
      </div>

      <!-- 登录表单 -->
      <n-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        size="large"
        class="login-form"
      >
        <n-form-item path="username" label="用户名">
          <n-input
            v-model:value="formData.username"
            placeholder="请输入用户名"
            :input-props="{ autocomplete: 'username' }"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <n-icon :component="PersonOutline" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password" label="密码">
          <n-input
            v-model:value="formData.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
            :input-props="{ autocomplete: 'current-password' }"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <n-icon :component="LockClosedOutline" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            block
            strong
            :loading="loading"
            @click="handleLogin"
          >
            <template #icon>
              <n-icon :component="LogInOutline" />
            </template>
            登录
          </n-button>
        </n-form-item>
      </n-form>
    </div>

    <!-- 版权信息 -->
    <div class="login-footer">
      <span>Powered by Muse</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage, type FormInst, type FormRules } from 'naive-ui';
import {
  NForm,
  NFormItem,
  NInput,
  NButton,
  NIcon
} from 'naive-ui';
import {
  SparklesOutline,
  PersonOutline,
  LockClosedOutline,
  LogInOutline
} from '@vicons/ionicons5';
import { useUserStore } from '@/stores/user';
import { userClient } from '@/api/client';
import { ConnectError } from '@connectrpc/connect';

const router = useRouter();
const message = useMessage();
const userStore = useUserStore();

const formRef = ref<FormInst | null>(null);
const loading = ref(false);

const formData = reactive({
  username: '',
  password: ''
});

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: ['blur', 'input'] }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: ['blur', 'input'] }
  ]
};

const handleLogin = async () => {
  // 表单验证
  const valid = await formRef.value?.validate().catch(() => false);
  if (!valid) return;

  loading.value = true;
  try {
    const response = await userClient.login({
      username: formData.username,
      password: formData.password
    });

    // 保存 token 和用户信息
    userStore.setToken(response.token);
    userStore.setCurrentUser(response.user ?? null);

    // 初始化用户数据（人设、设置、API配置等）
    await userStore.initUserData();

    message.success('登录成功');

    // 跳转到首页
    router.push('/');
  } catch (error) {
    if (error instanceof ConnectError) {
      message.error(error.message || '登录失败');
    } else {
      message.error('网络错误，请检查网络连接');
    }
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  position: relative;
  overflow: hidden;
  background: var(--bg-primary);
}

/* 背景效果 */
.login-bg {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 0;
}

.bg-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(138, 180, 248, .12) 0%, transparent 70%);
  animation: pulse 4s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: .5;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: .8;
    transform: translate(-50%, -50%) scale(1.1);
  }
}

/* 登录卡片 */
.login-card {
  position: relative;
  width: 100%;
  max-width: 400px;
  padding: 48px 40px;
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: var(--shadow-xl);
  z-index: 1;
}

/* 登录头部 */
.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.logo-icon {
  color: var(--color-primary);
  animation: glow 2s ease-in-out infinite;
}

@keyframes glow {
  0%, 100% {
    filter: drop-shadow(0 0 8px rgba(138, 180, 248, .4));
  }
  50% {
    filter: drop-shadow(0 0 16px rgba(138, 180, 248, .6));
  }
}

.title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: 2px;
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

/* 表单样式 */
.login-form {
  margin-top: 24px;
}

.login-form :deep(.n-form-item-label) {
  color: var(--text-secondary);
  font-size: 13px;
}

.login-form :deep(.n-input) {
  background: var(--bg-tertiary);
  border-color: var(--border-color);
}

.login-form :deep(.n-input:hover) {
  border-color: var(--color-primary);
}

.login-form :deep(.n-input:focus-within) {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px rgba(138, 180, 248, .15);
}

.login-form :deep(.n-button--primary-type) {
  background: var(--color-primary);
  border: none;
  font-weight: 600;
  letter-spacing: 1px;
  transition: all var(--transition-normal);
}

.login-form :deep(.n-button--primary-type:hover) {
  background: var(--color-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(138, 180, 248, .3);
}

.login-form :deep(.n-button--primary-type:active) {
  transform: translateY(0);
}

/* 页脚 */
.login-footer {
  position: absolute;
  bottom: 24px;
  color: var(--text-tertiary);
  font-size: 12px;
}

/* 响应式 */
@media (max-width: 480px) {
  .login-card {
    padding: 32px 24px;
  }

  .title {
    font-size: 28px;
  }
}
</style>
