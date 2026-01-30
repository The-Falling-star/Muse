<template>
  <div class="login-container">
    <!-- 科幻背景效果 -->
    <div class="login-bg">
      <div class="bg-grid"></div>
      <div class="bg-glow"></div>
      <div class="bg-particles">
        <div
          v-for="i in 30"
          :key="i"
          class="particle"
          :style="{
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            animationDelay: `${Math.random() * 5}s`,
            animationDuration: `${3 + Math.random() * 4}s`
          }"
        ></div>
      </div>
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
          <div class="logo-glow"></div>
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

      <!-- 装饰边框 -->
      <div class="card-border card-border--top"></div>
      <div class="card-border card-border--bottom"></div>
      <div class="card-corner card-corner--tl"></div>
      <div class="card-corner card-corner--tr"></div>
      <div class="card-corner card-corner--bl"></div>
      <div class="card-corner card-corner--br"></div>
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

.bg-grid {
  position: absolute;
  width: 100%;
  height: 100%;
  background-image:
    linear-gradient(var(--grid-color) 1px, transparent 1px),
    linear-gradient(90deg, var(--grid-color) 1px, transparent 1px);
  background-size: 50px 50px;
  opacity: 0.5;
}

.bg-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 600px;
  height: 600px;
  background: radial-gradient(circle, rgba(0, 240, 255, 0.1) 0%, transparent 70%);
  animation: pulse 4s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.5;
    transform: translate(-50%, -50%) scale(1);
  }
  50% {
    opacity: 0.8;
    transform: translate(-50%, -50%) scale(1.1);
  }
}

.bg-particles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: var(--color-primary);
  border-radius: 50%;
  opacity: 0;
  animation: float 5s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    opacity: 0;
    transform: translateY(0);
  }
  50% {
    opacity: 0.6;
    transform: translateY(-100px);
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
  backdrop-filter: blur(20px);
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
    filter: drop-shadow(0 0 10px var(--color-primary));
  }
  50% {
    filter: drop-shadow(0 0 20px var(--color-primary));
  }
}

.logo-glow {
  position: absolute;
  width: 80px;
  height: 80px;
  background: radial-gradient(circle, rgba(0, 240, 255, 0.2) 0%, transparent 70%);
  border-radius: 50%;
}

.title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0 0 8px 0;
  letter-spacing: 2px;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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
  box-shadow: 0 0 0 2px rgba(0, 240, 255, 0.1);
}

.login-form :deep(.n-button--primary-type) {
  background: var(--gradient-primary);
  border: none;
  font-weight: 600;
  letter-spacing: 1px;
  transition: all var(--transition-normal);
}

.login-form :deep(.n-button--primary-type:hover) {
  transform: translateY(-2px);
  box-shadow: var(--glow-primary);
}

.login-form :deep(.n-button--primary-type:active) {
  transform: translateY(0);
}

/* 装饰边框 */
.card-border {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  width: 60%;
  height: 2px;
  background: var(--gradient-primary);
  opacity: 0.5;
}

.card-border--top {
  top: 0;
}

.card-border--bottom {
  bottom: 0;
}

/* 角落装饰 */
.card-corner {
  position: absolute;
  width: 20px;
  height: 20px;
  border: 2px solid var(--color-primary);
  opacity: 0.3;
}

.card-corner--tl {
  top: -1px;
  left: -1px;
  border-right: none;
  border-bottom: none;
  border-radius: 16px 0 0 0;
}

.card-corner--tr {
  top: -1px;
  right: -1px;
  border-left: none;
  border-bottom: none;
  border-radius: 0 16px 0 0;
}

.card-corner--bl {
  bottom: -1px;
  left: -1px;
  border-right: none;
  border-top: none;
  border-radius: 0 0 0 16px;
}

.card-corner--br {
  bottom: -1px;
  right: -1px;
  border-left: none;
  border-top: none;
  border-radius: 0 0 16px 0;
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
