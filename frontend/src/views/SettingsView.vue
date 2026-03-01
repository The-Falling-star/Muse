<template>
  <div class="settings-view">
    <n-scrollbar class="settings-container">
      <div class="settings-content">
        <!-- API配置管理 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><CloudOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>API 配置</h2>
              <p>管理AI服务连接配置</p>
            </div>
            <n-button type="primary" size="small" @click="showApiConfigModal = true">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加配置
            </n-button>
          </div>

          <n-card class="settings-card">
            <n-spin :show="loading">
              <n-empty v-if="apiConfigs.length === 0" description="暂无API配置">
                <template #extra>
                  <n-button size="small" @click="showApiConfigModal = true">添加配置</n-button>
                </template>
              </n-empty>

              <n-list v-else>
                <n-list-item v-for="config in apiConfigs" :key="config.id">
                  <template #prefix>
                    <n-radio
                      :checked="config.isActive"
                      @click="handleSetActiveApiConfig(config.id)"
                    />
                  </template>
                  <n-thing :title="config.name" :description="getProviderLabel(config.provider)">
                    <template #header-extra>
                      <n-tag v-if="config.isActive" type="success" size="small">当前使用</n-tag>
                    </template>
                    <template #description>
                      <n-space :size="4">
                        <n-tag size="small">{{ config.model || '默认模型' }}</n-tag>
                        <span style="opacity: 0.7">{{ config.baseUrl || '默认地址' }}</span>
                      </n-space>
                    </template>
                  </n-thing>
                  <template #suffix>
                    <n-space>
                      <n-button size="small" quaternary @click="handleTestApiConfig(config.id)">
                        <template #icon>
                          <n-icon><FlashOutline /></n-icon>
                        </template>
                        测试
                      </n-button>
                      <n-button size="small" quaternary @click="handleEditApiConfig(config)">
                        <template #icon>
                          <n-icon><CreateOutline /></n-icon>
                        </template>
                      </n-button>
                      <n-button size="small" quaternary type="error" @click="handleDeleteApiConfig(config)">
                        <template #icon>
                          <n-icon><TrashOutline /></n-icon>
                        </template>
                      </n-button>
                    </n-space>
                  </template>
                </n-list-item>
              </n-list>
            </n-spin>
          </n-card>
        </section>

        <!-- 人设管理 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><PersonOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>人设管理</h2>
              <p>管理您的角色扮演人设</p>
            </div>
            <n-button type="primary" size="small" @click="showPersonaModal = true">
              <template #icon>
                <n-icon><AddOutline /></n-icon>
              </template>
              添加人设
            </n-button>
          </div>

          <n-card class="settings-card">
            <n-spin :show="loading">
              <n-empty v-if="personas.length === 0" description="暂无人设">
                <template #extra>
                  <n-button size="small" @click="showPersonaModal = true">添加人设</n-button>
                </template>
              </n-empty>

              <n-list v-else>
                <n-list-item v-for="persona in personas" :key="persona.id">
                  <template #prefix>
                    <n-radio
                      :checked="persona.id === userStore.currentUser?.activePersonaId"
                      @click="handleSetActivePersona(persona.id)"
                    />
                  </template>
                  <n-thing :title="persona.name" :description="persona.description || '暂无描述'">
                    <template #avatar>
                      <n-avatar :src="persona.avatar" round>
                        {{ persona.name.charAt(0) }}
                      </n-avatar>
                    </template>
                    <template #header-extra>
                      <n-tag v-if="persona.id === userStore.currentUser?.activePersonaId" type="success" size="small">当前使用</n-tag>
                    </template>
                  </n-thing>
                  <template #suffix>
                    <n-space>
                      <n-button size="small" quaternary @click="handleEditPersona(persona)">
                        <template #icon>
                          <n-icon><CreateOutline /></n-icon>
                        </template>
                      </n-button>
                      <n-button size="small" quaternary type="error" @click="handleDeletePersona(persona)">
                        <template #icon>
                          <n-icon><TrashOutline /></n-icon>
                        </template>
                      </n-button>
                    </n-space>
                  </template>
                </n-list-item>
              </n-list>
            </n-spin>
          </n-card>
        </section>

        <!-- 外观设置 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><ColorPaletteOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>外观设置</h2>
              <p>个性化界面风格</p>
            </div>
          </div>

          <n-card class="settings-card">
            <n-form label-placement="left" label-width="120">
              <n-form-item label="主题">
                <n-radio-group v-model:value="appearanceSettings.theme" @update:value="handleUpdateTheme">
                  <n-space>
                    <n-radio-button value="DARK">
                      <n-icon><MoonOutline /></n-icon>
                      暗黑
                    </n-radio-button>
                    <n-radio-button value="LIGHT">
                      <n-icon><SunnyOutline /></n-icon>
                      明亮
                    </n-radio-button>
                    <n-radio-button value="SYSTEM">
                      <n-icon><DesktopOutline /></n-icon>
                      跟随系统
                    </n-radio-button>
                  </n-space>
                </n-radio-group>
              </n-form-item>

              <n-form-item label="显示时间戳">
                <n-switch v-model:value="appearanceSettings.showTimestamps" @update:value="handleUpdateUserSetting" />
              </n-form-item>
            </n-form>
          </n-card>
        </section>

        <!-- 账户安全 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><ShieldOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>账户安全</h2>
              <p>管理账户密码</p>
            </div>
          </div>

          <n-card class="settings-card">
            <n-form label-placement="left" label-width="120">
              <n-form-item label="当前密码">
                <n-input
                  v-model:value="passwordForm.oldPassword"
                  type="password"
                  show-password-on="click"
                  placeholder="输入当前密码"
                />
              </n-form-item>

              <n-form-item label="新密码">
                <n-input
                  v-model:value="passwordForm.newPassword"
                  type="password"
                  show-password-on="click"
                  placeholder="输入新密码"
                />
              </n-form-item>

              <n-form-item label="确认密码">
                <n-input
                  v-model:value="passwordForm.confirmPassword"
                  type="password"
                  show-password-on="click"
                  placeholder="再次输入新密码"
                />
              </n-form-item>

              <n-form-item>
                <n-button type="primary" :loading="loading" @click="handleChangePassword">
                  修改密码
                </n-button>
              </n-form-item>
            </n-form>
          </n-card>
        </section>

        <!-- 关于 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><InformationCircleOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>关于</h2>
              <p>版本信息</p>
            </div>
          </div>

          <n-card class="settings-card about-card">
            <div class="about-logo">
              <div class="logo-icon">
                <span>M</span>
              </div>
              <div class="logo-info">
                <h3>Muse</h3>
                <p>v1.0.0</p>
              </div>
            </div>
            <p class="about-desc">
              Muse 是一个现代化的AI角色扮演平台，提供流畅的聊天体验和强大的角色管理功能。
            </p>
            <n-space justify="center">
              <n-button text type="primary">查看更新日志</n-button>
              <n-button text type="primary">GitHub</n-button>
              <n-button text type="primary">反馈问题</n-button>
            </n-space>
          </n-card>
        </section>
      </div>
    </n-scrollbar>

    <!-- API配置编辑模态框 -->
    <n-modal
      v-model:show="showApiConfigModal"
      preset="card"
      :title="editingApiConfig ? '编辑API配置' : '添加API配置'"
      :style="{ width: '500px', maxWidth: '90vw' }"
      :mask-closable="false"
    >
      <n-form ref="apiConfigFormRef" :model="apiConfigForm" :rules="apiConfigRules" label-placement="left" label-width="100">
        <n-form-item label="配置名称" path="name">
          <n-input v-model:value="apiConfigForm.name" placeholder="输入配置名称" />
        </n-form-item>

        <n-form-item label="服务类型" path="provider">
          <n-select
            v-model:value="apiConfigForm.provider"
            :options="providerOptions"
            placeholder="选择API类型"
          />
        </n-form-item>

        <n-form-item label="API地址" path="baseUrl">
          <n-input v-model:value="apiConfigForm.baseUrl" placeholder="https://api.openai.com/v1" />
        </n-form-item>

        <n-form-item label="API密钥" path="apiKey">
          <n-input
            v-model:value="apiConfigForm.apiKey"
            type="password"
            show-password-on="click"
            :placeholder="editingApiConfig ? '留空则不修改' : 'sk-...'"
          />
        </n-form-item>

        <n-form-item label="模型" path="model">
          <n-select
            v-model:value="apiConfigForm.model"
            :options="modelOptions"
            filterable
            tag
            placeholder="选择或输入模型名称"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showApiConfigModal = false">取消</n-button>
          <n-button type="primary" :loading="loading" @click="handleSaveApiConfig">
            {{ editingApiConfig ? '保存' : '添加' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 人设编辑模态框 -->
    <n-modal
      v-model:show="showPersonaModal"
      preset="card"
      :title="editingPersona ? '编辑人设' : '添加人设'"
      :style="{ width: '500px', maxWidth: '90vw' }"
      :mask-closable="false"
    >
      <n-form ref="personaFormRef" :model="personaForm" :rules="personaRules" label-placement="left" label-width="80">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="personaForm.name" placeholder="输入人设名称" />
        </n-form-item>

        <n-form-item label="头像" path="avatar">
          <n-input v-model:value="personaForm.avatar" placeholder="输入头像URL（可选）" />
        </n-form-item>

        <n-form-item label="描述" path="description">
          <n-input
            v-model:value="personaForm.description"
            type="textarea"
            :rows="4"
            placeholder="输入人设描述，将在聊天中作为您的角色设定"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showPersonaModal = false">取消</n-button>
          <n-button type="primary" :loading="loading" @click="handleSavePersona">
            {{ editingPersona ? '保存' : '添加' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import {
  NScrollbar,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSwitch,
  NRadioGroup,
  NRadioButton,
  NButton,
  NSpace,
  NIcon,
  NTag,
  NList,
  NListItem,
  NThing,
  NRadio,
  NAvatar,
  NModal,
  NEmpty,
  NSpin,
  useMessage,
  useDialog
} from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import {
  CloudOutline,
  FlashOutline,
  ColorPaletteOutline,
  MoonOutline,
  SunnyOutline,
  DesktopOutline,
  InformationCircleOutline,
  AddOutline,
  CreateOutline,
  TrashOutline,
  PersonOutline,
  ShieldOutline
} from '@vicons/ionicons5';

import { userClient } from '@/api/client';
import { useUserStore } from '@/stores/user';
import { useThemeStore } from '@/stores/theme';
import type { APIConfig, Persona } from '@/gen/muse/user_pb';
import { APIProvider, Theme } from '@/gen/muse/common_pb';

const message = useMessage();
const dialog = useDialog();
const userStore = useUserStore();
const themeStore = useThemeStore();

// 状态
const loading = ref(false);
const apiConfigs = ref<APIConfig[]>([]);
const personas = ref<Persona[]>([]);

// API配置相关
const showApiConfigModal = ref(false);
const editingApiConfig = ref<APIConfig | null>(null);
const apiConfigFormRef = ref<FormInst | null>(null);
const apiConfigForm = reactive({
  name: '',
  provider: APIProvider.OpenAI,
  baseUrl: '',
  apiKey: '',
  model: ''
});

const apiConfigRules: FormRules = {
  name: { required: true, message: '请输入配置名称', trigger: 'blur' },
  provider: { required: true, message: '请选择服务类型', trigger: 'change' }
};

// 人设相关
const showPersonaModal = ref(false);
const editingPersona = ref<Persona | null>(null);
const personaFormRef = ref<FormInst | null>(null);
const personaForm = reactive({
  name: '',
  avatar: '',
  description: ''
});

const personaRules: FormRules = {
  name: { required: true, message: '请输入人设名称', trigger: 'blur' }
};

// 外观设置
const appearanceSettings = reactive({
  theme: 'DARK',
  showTimestamps: true
});

// 密码修改
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

// 服务类型选项
const providerOptions = [
  { label: 'OpenAI', value: APIProvider.OpenAI },
  { label: 'Claude', value: APIProvider.Claude },
  { label: 'Google Gemini', value: APIProvider.Gemini }
];

// 模型选项
const modelOptions = [
  { label: 'GPT-4o', value: 'gpt-4o' },
  { label: 'GPT-4 Turbo', value: 'gpt-4-turbo' },
  { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo' },
  { label: 'Claude 3.5 Sonnet', value: 'claude-3-5-sonnet-20241022' },
  { label: 'Claude 3 Opus', value: 'claude-3-opus-20240229' },
  { label: 'Gemini 1.5 Pro', value: 'gemini-1.5-pro' },
  { label: 'Gemini 1.5 Flash', value: 'gemini-1.5-flash' }
];

// 获取服务类型标签
const getProviderLabel = (provider: APIProvider): string => {
  const option = providerOptions.find(o => o.value === provider);
  return option?.label || '未知';
};

// 加载数据
const loadData = async () => {
  loading.value = true;
  try {
    const [configsRes, personasRes, settingRes] = await Promise.all([
      userClient.listAPIConfigs({}),
      userClient.listPersonas({}),
      userClient.getUserSetting({})
    ]);

    apiConfigs.value = configsRes.configs;
    personas.value = personasRes.personas;

    // 更新Store
    userStore.setApiConfigs(configsRes.configs);
    userStore.setPersonas(personasRes.personas);

    // 更新外观设置
    if (settingRes.setting) {
      userStore.setUserSetting(settingRes.setting);
      appearanceSettings.theme = settingRes.setting.theme === Theme.Dark ? 'DARK' : settingRes.setting.theme === Theme.Light ? 'LIGHT' : 'SYSTEM';
      appearanceSettings.showTimestamps = settingRes.setting.showTimestamps;

      // 初始化主题
      handleUpdateTheme();
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  loadData();
});

// ==================== API配置相关方法 ====================

const handleEditApiConfig = (config: APIConfig) => {
  editingApiConfig.value = config;
  apiConfigForm.name = config.name;
  apiConfigForm.provider = config.provider;
  apiConfigForm.baseUrl = config.baseUrl || '';
  apiConfigForm.apiKey = '';
  apiConfigForm.model = config.model || '';
  showApiConfigModal.value = true;
};

const handleSaveApiConfig = async () => {
  await apiConfigFormRef.value?.validate();

  loading.value = true;
  try {
    if (editingApiConfig.value) {
      // 更新
      const res = await userClient.updateAPIConfig({
        id: editingApiConfig.value.id,
        name: apiConfigForm.name,
        provider: apiConfigForm.provider,
        baseUrl: apiConfigForm.baseUrl || undefined,
        apiKey: apiConfigForm.apiKey || undefined,
        model: apiConfigForm.model || undefined
      });
      if (res.config) {
        const index = apiConfigs.value.findIndex(c => c.id === editingApiConfig.value!.id);
        if (index >= 0) {
          apiConfigs.value[index] = res.config;
        }
        userStore.updateApiConfigInList(res.config);
        message.success('API配置已更新');
      }
    } else {
      // 创建
      const res = await userClient.createAPIConfig({
        name: apiConfigForm.name,
        provider: apiConfigForm.provider,
        apiKey: apiConfigForm.apiKey,
        baseUrl: apiConfigForm.baseUrl || undefined,
        model: apiConfigForm.model || undefined
      });
      if (res.config) {
        apiConfigs.value.push(res.config);
        userStore.addApiConfig(res.config);
        message.success('API配置已创建');
      }
    }
    showApiConfigModal.value = false;
    resetApiConfigForm();
  } finally {
    loading.value = false;
  }
};

const handleDeleteApiConfig = (config: APIConfig) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除API配置"${config.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      loading.value = true;
      try {
        await userClient.deleteAPIConfig({ id: config.id });
        apiConfigs.value = apiConfigs.value.filter(c => c.id !== config.id);
        userStore.removeApiConfig(config.id);
        message.success('API配置已删除');
      } finally {
        loading.value = false;
      }
    }
  });
};

const handleSetActiveApiConfig = async (configId: number) => {
  loading.value = true;
  try {
    await userClient.setActiveAPIConfig({ configId: configId });
    // 更新本地状态
    apiConfigs.value.forEach(c => {
      c.isActive = c.id === configId;
    });
    userStore.setActiveApiConfigId(configId);
    message.success('已切换API配置');
  } finally {
    loading.value = false;
  }
};

const handleTestApiConfig = async (configId: number) => {
  loading.value = true;
  try {
    const res = await userClient.testAPIConfig({ id: configId });
    if (res.success) {
      message.success('连接成功！');
    } else {
      message.error(`连接失败: ${res.errorMessage}`);
    }
  } finally {
    loading.value = false;
  }
};

const resetApiConfigForm = () => {
  editingApiConfig.value = null;
  apiConfigForm.name = '';
  apiConfigForm.provider = APIProvider.OpenAI;
  apiConfigForm.baseUrl = '';
  apiConfigForm.apiKey = '';
  apiConfigForm.model = '';
};

// ==================== 人设相关方法 ====================

const handleEditPersona = (persona: Persona) => {
  editingPersona.value = persona;
  personaForm.name = persona.name;
  personaForm.avatar = persona.avatar || '';
  personaForm.description = persona.description || '';
  showPersonaModal.value = true;
};

const handleSavePersona = async () => {
  await personaFormRef.value?.validate();

  loading.value = true;
  try {
    if (editingPersona.value) {
      // 更新
      const res = await userClient.updatePersona({
        id: editingPersona.value.id,
        name: personaForm.name,
        avatar: personaForm.avatar || undefined,
        description: personaForm.description || undefined
      });
      if (res.persona) {
        const index = personas.value.findIndex(p => p.id === editingPersona.value!.id);
        if (index >= 0) {
          personas.value[index] = res.persona;
        }
        userStore.updatePersonaInList(res.persona);
        message.success('人设已更新');
      }
    } else {
      // 创建
      const res = await userClient.createPersona({
        name: personaForm.name,
        avatar: personaForm.avatar || undefined,
        description: personaForm.description || undefined
      });
      if (res.persona) {
        personas.value.push(res.persona);
        userStore.addPersona(res.persona);
        message.success('人设已创建');
      }
    }
    showPersonaModal.value = false;
    resetPersonaForm();
  } finally {
    loading.value = false;
  }
};

const handleDeletePersona = (persona: Persona) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除人设"${persona.name}"吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      loading.value = true;
      try {
        await userClient.deletePersona({ id: persona.id });
        personas.value = personas.value.filter(p => p.id !== persona.id);
        userStore.removePersona(persona.id);
        message.success('人设已删除');
      } finally {
        loading.value = false;
      }
    }
  });
};

const handleSetActivePersona = async (personaId: number) => {
  loading.value = true;
  try {
    await userClient.setActivePersona({ personaId: personaId });
    userStore.setActivePersonaId(personaId);
    message.success('已切换人设');
  } catch (e) {
    console.error('设置活跃人设失败:', e);
    message.error('切换失败');
  } finally {
    loading.value = false;
  }
};

const resetPersonaForm = () => {
  editingPersona.value = null;
  personaForm.name = '';
  personaForm.avatar = '';
  personaForm.description = '';
};

// ==================== 用户设置相关方法 ====================

// 更新主题
const handleUpdateTheme = async () => {
  let themeMode: 'dark' | 'light';
  const themeValue = appearanceSettings.theme;

  if (themeValue === 'SYSTEM') {
    themeMode = window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
  } else {
    themeMode = themeValue === 'LIGHT' ? 'light' : 'dark';
  }

  themeStore.setTheme(themeMode);

  // 保存到后端
  const protobufTheme = themeValue === 'DARK' ? Theme.Dark : themeValue === 'LIGHT' ? Theme.Light : Theme.Auto;
  const res = await userClient.updateUserSetting({
    theme: protobufTheme,
    showTimestamps: appearanceSettings.showTimestamps
  });
  if (res.setting) {
    userStore.setUserSetting(res.setting);
  }
};

// 更新其他用户设置
const handleUpdateUserSetting = async () => {
  const themeValue = appearanceSettings.theme === 'DARK' ? Theme.Dark : appearanceSettings.theme === 'LIGHT' ? Theme.Light : Theme.Auto;
  const res = await userClient.updateUserSetting({
    theme: themeValue,
    showTimestamps: appearanceSettings.showTimestamps
  });
  if (res.setting) {
    userStore.setUserSetting(res.setting);
  }
};

// ==================== 密码修改 ====================

const handleChangePassword = async () => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword) {
    message.warning('请填写完整密码信息');
    return;
  }
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    message.warning('两次输入的新密码不一致');
    return;
  }
  if (passwordForm.newPassword.length < 6) {
    message.warning('新密码长度至少6位');
    return;
  }

  loading.value = true;
  try {
    await userClient.changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    });
    message.success('密码修改成功');
    passwordForm.oldPassword = '';
    passwordForm.newPassword = '';
    passwordForm.confirmPassword = '';
  } catch (e) {
    console.error('修改密码失败:', e);
    message.error('修改密码失败，请检查原密码是否正确');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.settings-view {
  height: calc(100vh - 64px - 48px);
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

.settings-container {
  height: 100%;
  /* 移动端滚动优化 */
  -webkit-overflow-scrolling: touch;
}

.settings-content {
  max-width: 800px;
  margin: 0 auto;
  padding-bottom: 40px;
}

/* 设置分区 */
.settings-section {
  margin-bottom: 32px;
  opacity: 1;
  transition: opacity var(--transition-normal), transform var(--transition-normal);
}

.settings-section:last-of-type {
  margin-bottom: 0;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .settings-section {
    margin-bottom: 24px;
  }
  
  .settings-content {
    padding: 0 12px;
  }
  
  .section-header {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    margin-bottom: 12px;
  }
  
  .section-title {
    text-align: center;
  }
}

.section-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 16px;
}

.section-icon {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-primary);
  border-radius: 12px;
  color: #000;
}

.section-title {
  flex: 1;
}

.section-title h2 {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.section-title p {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0;
}

.settings-card {
  border-radius: 12px;
}

/* 关于卡片 */
.about-card {
  text-align: center;
}

.about-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 16px;
}

.about-logo .logo-icon {
  width: 60px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-primary);
  border-radius: 16px;
  font-size: 28px;
  font-weight: 700;
  color: #000;
}

.about-logo .logo-info h3 {
  font-size: 24px;
  font-weight: 700;
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.about-logo .logo-info p {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0;
}

.about-desc {
  color: var(--text-secondary);
  margin-bottom: 16px;
}
</style>
