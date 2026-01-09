<template>
  <div class="settings-view">
    <n-scrollbar class="settings-container">
      <div class="settings-content">
        <!-- API设置 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><CloudOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>API 设置</h2>
              <p>配置AI服务连接</p>
            </div>
          </div>

          <n-card class="settings-card">
            <n-form label-placement="left" label-width="120">
              <n-form-item label="服务类型">
                <n-select
                  v-model:value="apiSettings.type"
                  :options="apiTypeOptions"
                  placeholder="选择API类型"
                />
              </n-form-item>

              <n-form-item label="API地址">
                <n-input
                  v-model:value="apiSettings.baseUrl"
                  placeholder="https://api.openai.com/v1"
                />
              </n-form-item>

              <n-form-item label="API密钥">
                <n-input
                  v-model:value="apiSettings.apiKey"
                  type="password"
                  show-password-on="click"
                  placeholder="sk-..."
                />
              </n-form-item>

              <n-form-item label="模型">
                <n-select
                  v-model:value="apiSettings.model"
                  :options="modelOptions"
                  filterable
                  tag
                  placeholder="选择或输入模型名称"
                />
              </n-form-item>

              <n-form-item>
                <n-space>
                  <n-button type="primary" @click="testConnection">
                    <template #icon>
                      <n-icon><FlashOutline /></n-icon>
                    </template>
                    测试连接
                  </n-button>
                  <n-button @click="saveApiSettings">保存设置</n-button>
                </n-space>
              </n-form-item>
            </n-form>
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
                <n-radio-group v-model:value="appearance.theme">
                  <n-space>
                    <n-radio-button value="dark">
                      <n-icon><MoonOutline /></n-icon>
                      暗黑
                    </n-radio-button>
                    <n-radio-button value="light">
                      <n-icon><SunnyOutline /></n-icon>
                      明亮
                    </n-radio-button>
                    <n-radio-button value="auto">
                      <n-icon><DesktopOutline /></n-icon>
                      跟随系统
                    </n-radio-button>
                  </n-space>
                </n-radio-group>
              </n-form-item>

              <n-form-item label="字体大小">
                <n-slider
                  v-model:value="appearance.fontSize"
                  :min="12"
                  :max="20"
                  :step="1"
                  :marks="{ 12: '小', 14: '默认', 16: '中', 18: '大', 20: '特大' }"
                />
              </n-form-item>

              <n-form-item label="消息气泡">
                <n-switch v-model:value="appearance.showBubbles" />
              </n-form-item>

              <n-form-item label="显示时间戳">
                <n-switch v-model:value="appearance.showTimestamp" />
              </n-form-item>
            </n-form>
          </n-card>
        </section>

        <!-- 快捷键 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><KeypadOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>快捷键</h2>
              <p>自定义键盘快捷方式</p>
            </div>
          </div>

          <n-card class="settings-card">
            <n-form label-placement="left" label-width="140">
              <n-form-item label="发送消息">
                <n-tag>Enter</n-tag>
              </n-form-item>
              <n-form-item label="换行">
                <n-tag>Shift + Enter</n-tag>
              </n-form-item>
              <n-form-item label="新建会话">
                <n-tag>Ctrl + N</n-tag>
              </n-form-item>
              <n-form-item label="切换侧边栏">
                <n-tag>Ctrl + B</n-tag>
              </n-form-item>
            </n-form>
          </n-card>
        </section>

        <!-- 数据管理 -->
        <section class="settings-section">
          <div class="section-header">
            <div class="section-icon">
              <n-icon size="24"><ServerOutline /></n-icon>
            </div>
            <div class="section-title">
              <h2>数据管理</h2>
              <p>备份与恢复数据</p>
            </div>
          </div>

          <n-card class="settings-card">
            <n-space vertical :size="16">
              <div class="data-action">
                <div class="data-action-info">
                  <h4>导出数据</h4>
                  <p>导出所有角色、会话和设置</p>
                </div>
                <n-button>
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  导出
                </n-button>
              </div>

              <n-divider />

              <div class="data-action">
                <div class="data-action-info">
                  <h4>导入数据</h4>
                  <p>从备份文件恢复数据</p>
                </div>
                <n-button>
                  <template #icon>
                    <n-icon><CloudUploadOutline /></n-icon>
                  </template>
                  导入
                </n-button>
              </div>

              <n-divider />

              <div class="data-action danger">
                <div class="data-action-info">
                  <h4>清除数据</h4>
                  <p>删除所有本地数据，此操作不可撤销</p>
                </div>
                <n-button type="error">
                  <template #icon>
                    <n-icon><TrashOutline /></n-icon>
                  </template>
                  清除
                </n-button>
              </div>
            </n-space>
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
            <n-space>
              <n-button text type="primary">查看更新日志</n-button>
              <n-button text type="primary">GitHub</n-button>
              <n-button text type="primary">反馈问题</n-button>
            </n-space>
          </n-card>
        </section>
      </div>
    </n-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import {
  NScrollbar,
  NCard,
  NForm,
  NFormItem,
  NInput,
  NSelect,
  NSlider,
  NSwitch,
  NRadioGroup,
  NRadioButton,
  NButton,
  NSpace,
  NIcon,
  NDivider,
  NTag,
  useMessage
} from 'naive-ui';
import {
  CloudOutline,
  FlashOutline,
  ColorPaletteOutline,
  MoonOutline,
  SunnyOutline,
  DesktopOutline,
  KeypadOutline,
  ServerOutline,
  DownloadOutline,
  CloudUploadOutline,
  TrashOutline,
  InformationCircleOutline
} from '@vicons/ionicons5';

const message = useMessage();

// API设置
const apiSettings = reactive({
  type: 'openai',
  baseUrl: '',
  apiKey: '',
  model: 'gpt-4'
});

// API类型选项
const apiTypeOptions = [
  { label: 'OpenAI', value: 'openai' },
  { label: 'Claude', value: 'claude' },
  { label: 'Google Gemini', value: 'gemini' },
  { label: '自定义', value: 'custom' }
];

// 模型选项
const modelOptions = [
  { label: 'GPT-4', value: 'gpt-4' },
  { label: 'GPT-4 Turbo', value: 'gpt-4-turbo' },
  { label: 'GPT-3.5 Turbo', value: 'gpt-3.5-turbo' },
  { label: 'Claude 3 Opus', value: 'claude-3-opus' },
  { label: 'Claude 3 Sonnet', value: 'claude-3-sonnet' },
  { label: 'Gemini Pro', value: 'gemini-pro' }
];

// 外观设置
const appearance = reactive({
  theme: 'dark',
  fontSize: 14,
  showBubbles: true,
  showTimestamp: true
});

// 方法
const testConnection = () => {
  message.info('正在测试连接...');
  setTimeout(() => {
    message.success('连接成功！');
  }, 1500);
};

const saveApiSettings = () => {
  message.success('设置已保存');
};
</script>

<style scoped>
.settings-view {
  height: calc(100vh - 64px - 48px);
}

.settings-container {
  height: 100%;
}

.settings-content {
  max-width: 800px;
  margin: 0 auto;
  padding-bottom: 40px;
}

/* 设置分区 */
.settings-section {
  margin-bottom: 32px;
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

/* 数据操作 */
.data-action {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.data-action-info h4 {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 4px;
}

.data-action-info p {
  font-size: 13px;
  color: var(--text-tertiary);
  margin: 0;
}

.data-action.danger .data-action-info h4 {
  color: var(--color-error);
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
