<script setup lang="ts">
import type {UploadFileInfo} from "naive-ui";
import {
  NButton,
  NEmpty,
  NIcon,
  NInput,
  NList,
  NListItem,
  NModal,
  NPopconfirm,
  NTabPane,
  NTabs,
  NUpload
} from "naive-ui";
import {Check, Edit, Trash} from "@vicons/tabler"
import {APIProvider, FileType} from "@/gen/muse/common_pb.ts";
import type {APIConfig, Persona} from "@/gen/muse/user_pb.ts";
import {ref, watch} from "vue";
import {useUserStore} from "@/stores/user.ts";
import {useFileStore} from "@/stores/file.ts";
import {AddOutline} from "@vicons/ionicons5";
import {userClient} from "@/api/client.ts";

const providerObj = Object.entries(APIProvider)
    .filter(([key]) => key !== 'APIProviderUnspecified' && isNaN(Number(key))) // 过滤掉反向映射
    .map(([key, value]) => ({name: key, value: value as APIProvider}));
const providers = ref<{ name: string, value: APIProvider }[]>(providerObj)

const userStore = useUserStore();
type APIConfigWithEdit = APIConfig & { edit: boolean, loading: boolean };
const apiConfigs = ref<APIConfigWithEdit[]>(userStore.apiConfigs.map(apiConfig => ({
  ...apiConfig,
  edit: apiConfig.id === 0,
  loading: false
})));

const handleSaveApiConfig = async (config: APIConfigWithEdit) => {
  config.loading = true
  try {
    if (config.id) {
      // 同步更新 Store 里的现有项
      const storeItem = userStore.apiConfigs.find(c => c.id === config.id)
      if (storeItem) {
        if (storeItem.apiKey === config.apiKey) {
          config.edit = false
          return
        }
        storeItem.apiKey = config.apiKey
      }
      await userClient.updateAPIConfig({
        id: config.id,
        apiKey: config.apiKey,
        provider: config.provider,
      })

      config.edit = false
      return
    }
    const rsp = await userClient.createAPIConfig({
      apiKey: config.apiKey,
      provider: config.provider,
    });
    config.id = rsp.configId
    // 将新项添加到 Store
    const {edit, loading, ...originalConfig} = config
    userStore.apiConfigs.push(originalConfig)
    config.edit = false
  } finally {
    config.loading = false
  }
}

const handleAddApiConfig = (_provider: APIProvider) => {
  apiConfigs.value.push({
    isActive: false,
    $typeName: 'muse.APIConfig',
    apiKey: '',
    createdAt: 0n,
    edit: true,
    loading: false,
    id: 0,
    provider: _provider,
    updatedAt: 0n,
    userId: 0
  })
}

const handleDeleteApiConfig = async (config: APIConfigWithEdit) => {
  if (config.id !== 0) {
    config.loading = true
    try {
      await userClient.deleteAPIConfig({id: config.id})
      // 从 Store 中同步移除
      const storeIndex = userStore.apiConfigs.findIndex(c => c.id === config.id)
      if (storeIndex > -1) {
        userStore.apiConfigs.splice(storeIndex, 1)
      }
    } finally {
      config.loading = false
    }
  }
  const index = apiConfigs.value.findIndex(c => c === config)
  if (index > -1) {
    apiConfigs.value.splice(index, 1)
  }
}

const proxyUrl = ref<string|undefined>(userStore.currentUser?.proxyUrl)
let proxyUrlUpdateTimer: ReturnType<typeof setTimeout> | null = null;
const isProxyUrlValid = ref<boolean>(true)
const updateProxyUrl = (proxyUrl: string) => {
  isProxyUrlValid.value = checkProxyUrl(proxyUrl)
  if (!isProxyUrlValid.value) {
    console.info("proxyUrl的检验结果: ", isProxyUrlValid.value)
    return
  }
  if (proxyUrlUpdateTimer) {
    clearTimeout(proxyUrlUpdateTimer)
  }
  proxyUrlUpdateTimer = setTimeout(async () => {
    const curUser = userStore.currentUser;
    if (!curUser) {
      console.warn("用户未登录, 请先登录")
      return
    }
    await userClient.updateUserInfo({
      proxyUrl: proxyUrl
    })
    console.info("更新代理Url成功")
  }, 1000)
}

const checkProxyUrl = (proxyUrl: string): boolean => {
  if (!proxyUrl) {
    return true
  }
  try {
    const url = new URL(proxyUrl);
    return url.protocol === 'http:' || url.protocol === 'https:'
  } catch (err) {
    console.warn("proxyUrl非法")
    return false
  }
}

// =====================
// 人设管理
// =====================

const fileStore = useFileStore();

type PersonaWithEdit = Persona & {
  edit: boolean;
  loading: boolean;
  _previewUrl?: string;
};

const personaList = ref<PersonaWithEdit[]>(userStore.personas.map(p => ({
  ...p,
  edit: false,
  loading: false,
})));

const avatarUrls = ref<Record<string, string>>({});

const previewVisible = ref(false);
const previewAvatarUrl = ref('');

const preloadPersonaAvatars = async () => {
  for (const p of personaList.value) {
    if (p.avatar && !avatarUrls.value[p.avatar]) {
      try {
        const url = await fileStore.getFileUrl(p.avatar);
        avatarUrls.value[p.avatar] = url;
      } catch (e) {
        console.warn('加载头像失败:', p.id, e);
      }
    }
  }
};

watch(
  () => personaList.value.map(p => p.avatar),
  () => { preloadPersonaAvatars(); },
  { immediate: true }
);

const getPersonaAvatarUrl = (persona: PersonaWithEdit): string => {
  if (!persona.avatar) {
    return '';
  }
  return avatarUrls.value[persona.avatar] || '';
};

const truncateDesc = (desc: string): string => {
  if (!desc) {
    return '';
  }
  return desc.length > 10 ? `${desc.slice(0, 10)}...` : desc;
};

const handleAddPersona = () => {
  personaList.value.push({
    id: 0,
    userId: userStore.currentUser?.id ?? 0,
    name: '',
    avatar: '',
    description: '',
    createdAt: 0n,
    updatedAt: 0n,
    $typeName: 'muse.Persona',
    edit: true,
    loading: false,
  });
};

const handleAvatarUpload = async (persona: PersonaWithEdit, options: { file: UploadFileInfo }) => {
  const rawFile = options.file.file;
  if (!rawFile) {
    return;
  }

  persona.loading = true;
  try {
    const file = rawFile as File;
    const filePath = await fileStore.uploadFile(file, file.name || 'avatar.png', FileType.PersonaAvatar);
    persona.avatar = filePath;

    const url = await fileStore.getFileUrl(filePath);
    avatarUrls.value[filePath] = url;

    if (persona._previewUrl) {
      URL.revokeObjectURL(persona._previewUrl);
      persona._previewUrl = undefined;
    }
  } catch (e) {
    console.warn('头像上传失败:', e);
  } finally {
    persona.loading = false;
  }
};

const handlePreviewAvatar = (persona: PersonaWithEdit) => {
  const url = getPersonaAvatarUrl(persona);
  if (!url) {
    return;
  }
  previewAvatarUrl.value = url;
  previewVisible.value = true;
};

const handleSavePersona = async (persona: PersonaWithEdit) => {
  persona.loading = true;
  try {
    if (persona.id > 0) {
      const rsp = await userClient.updatePersona({
        id: persona.id,
        name: persona.name,
        avatar: persona.avatar,
        description: persona.description,
      });
      userStore.updatePersonaInList(rsp.persona!);

      const idx = personaList.value.findIndex(p => p.id === persona.id);
      if (idx > -1) {
        personaList.value[idx] = { ...rsp.persona!, edit: false, loading: false };
      }
      return;
    }

    const rsp = await userClient.createPersona({
      name: persona.name,
      avatar: persona.avatar,
      description: persona.description,
    });
    userStore.addPersona(rsp.persona!);

    const idx = personaList.value.findIndex(p => p === persona);
    if (idx > -1) {
      personaList.value[idx] = { ...rsp.persona!, edit: false, loading: false };
    }

    if (userStore.personas.length === 1) {
      await userStore.setActivePersonaId(rsp.persona!.id);
    }
  } finally {
    persona.loading = false;
  }
};

const handleDeletePersona = async (persona: PersonaWithEdit) => {
  if (persona.id === 0) {
    cleanupPersonaPreview(persona);
    const idx = personaList.value.findIndex(p => p === persona);
    if (idx > -1) {
      personaList.value.splice(idx, 1);
    }
    return;
  }

  persona.loading = true;
  try {
    const isActive = persona.id === userStore.currentUser?.activePersonaId;
    await userClient.deletePersona({ id: persona.id });
    userStore.removePersona(persona.id);

    const idx = personaList.value.findIndex(p => p.id === persona.id);
    if (idx > -1) {
      personaList.value.splice(idx, 1);
    }

    if (isActive && personaList.value.length > 0) {
      const nextActive = personaList.value[0]!;
      await userStore.setActivePersonaId(nextActive.id);
    }
  } finally {
    persona.loading = false;
  }
};

const handleActivatePersona = async (persona: PersonaWithEdit) => {
  if (persona.id === 0 || persona.id === userStore.currentUser?.activePersonaId) {
    return;
  }
  await userStore.setActivePersonaId(persona.id);
};

const cleanupPersonaPreview = (persona: PersonaWithEdit) => {
  if (persona._previewUrl) {
    URL.revokeObjectURL(persona._previewUrl);
    persona._previewUrl = undefined;
  }
};

const handleCancelEditPersona = (persona: PersonaWithEdit) => {
  cleanupPersonaPreview(persona);
  if (persona.id === 0) {
    const idx = personaList.value.findIndex(p => p === persona);
    if (idx > -1) {
      personaList.value.splice(idx, 1);
    }
    return;
  }

  const stored = userStore.personas.find(p => p.id === persona.id);
  if (stored) {
    const idx = personaList.value.findIndex(p => p.id === persona.id);
    if (idx > -1) {
      personaList.value[idx] = { ...stored, edit: false, loading: false };
    }
    return;
  }

  persona.edit = false;
};


</script>

<template>
  <div class="user-config-tab">
    <div class="section-label">模型API</div>
    <div class="settings-item">
      <n-tabs type="line" animated>
        <n-tab-pane
            v-for="provider in providers"
            :key="provider.value"
            :name="provider.name"
            :tab="provider.name"
        >
          <div class="tab-content">
            <n-list :show-divider="false">
              <n-list-item
                  v-for="config in apiConfigs.filter(apiCfg => apiCfg.provider === provider.value)"
                  :key="config.id"
              >
                <div class="api-item">
                  <n-input
                      placeholder="请输入API Key"
                      v-model:value="config.apiKey"
                      :disabled="!config.edit"
                      :status="config.apiKey? 'success':'error'"
                      size="small"
                  />
                  <div class="item-actions">
                    <n-button
                        v-if="!config.edit"
                        quaternary
                        circle
                        size="small"
                        @click="config.edit = true"
                    >
                      <template #icon>
                        <n-icon>
                          <Edit/>
                        </n-icon>
                      </template>
                    </n-button>
                    <n-button
                        v-else
                        quaternary
                        circle
                        size="small"
                        type="primary"
                        :loading="config.loading"
                        @click="handleSaveApiConfig(config)"
                    >
                      <template #icon>
                        <n-icon>
                          <Check/>
                        </n-icon>
                      </template>
                    </n-button>

                    <n-popconfirm
                        v-if="config.id !== 0"
                        @positive-click="handleDeleteApiConfig(config)"
                    >
                      <template #trigger>
                        <n-button
                            quaternary
                            circle
                            size="small"
                            type="error"
                            :disabled="config.loading"
                        >
                          <template #icon>
                            <n-icon>
                              <Trash/>
                            </n-icon>
                          </template>
                        </n-button>
                      </template>
                      确定要删除此 API 配置吗？
                    </n-popconfirm>
                    <n-button
                        v-else
                        quaternary
                        circle
                        size="small"
                        type="error"
                        @click="handleDeleteApiConfig(config)"
                    >
                      <template #icon>
                        <n-icon>
                          <Trash/>
                        </n-icon>
                      </template>
                    </n-button>
                  </div>
                </div>
              </n-list-item>
            </n-list>
            <n-button
                size="small"
                dashed
                block
                class="add-btn"
                @click="handleAddApiConfig(provider.value)"
            >
              <template #icon>
                <n-icon>
                  <AddOutline/>
                </n-icon>
              </template>
              添加API配置
            </n-button>
          </div>
        </n-tab-pane>
      </n-tabs>
    </div>
    <div class="section-label">API 连接</div>
    <div class="settings-item" style="margin-top: 12px">
      <div class="item-label">代理地址</div>
      <n-input
          v-model:value="proxyUrl"
          placeholder="https://proxy.example.com"
          size="small"
          class="item-control"
          @input="updateProxyUrl"
          :status="isProxyUrlValid ? 'success':'error'"
      />
    </div>
    <div class="section-label">人设设置</div>
    <div class="settings-item persona-section">
      <n-empty
        v-if="personaList.length === 0"
        description="暂无人设"
      >
        <template #extra>
          <n-button size="small" @click="handleAddPersona">
            新增人设
          </n-button>
        </template>
      </n-empty>

      <div v-else class="persona-list">
        <div
          v-for="persona in personaList"
          :key="persona.id"
          class="persona-item"
          :class="{ 'persona-item--active': persona.id === userStore.currentUser?.activePersonaId }"
        >
          <!-- 展示模式 -->
          <template v-if="!persona.edit">
            <div
              class="persona-content"
              @click="handleActivatePersona(persona)"
            >
              <div
                class="persona-avatar"
                @click.stop="handlePreviewAvatar(persona)"
              >
                <img
                  v-if="getPersonaAvatarUrl(persona)"
                  :src="getPersonaAvatarUrl(persona)"
                  alt="头像"
                >
                <span v-else class="persona-avatar__placeholder">
                  {{ persona.name?.charAt(0) || '?' }}
                </span>
              </div>
              <div class="persona-info">
                <span class="persona-name">{{ persona.name || '未命名' }}</span>
                <span
                  v-if="persona.description"
                  class="persona-desc"
                >
                  {{ truncateDesc(persona.description) }}
                </span>
              </div>
            </div>
            <div class="persona-actions">
              <n-button
                quaternary
                circle
                size="small"
                @click="persona.edit = true"
              >
                <template #icon>
                  <n-icon>
                    <Edit />
                  </n-icon>
                </template>
              </n-button>
              <n-popconfirm @positive-click="handleDeletePersona(persona)">
                <template #trigger>
                  <n-button
                    quaternary
                    circle
                    size="small"
                    type="error"
                  >
                    <template #icon>
                      <n-icon>
                        <Trash />
                      </n-icon>
                    </template>
                  </n-button>
                </template>
                确定要删除此人设吗？
              </n-popconfirm>
            </div>
          </template>

          <!-- 编辑模式 -->
          <template v-else>
            <div class="persona-edit-row">
              <div class="persona-edit-avatar">
                <n-upload
                  :max="1"
                  accept="image/*"
                  :show-file-list="false"
                  @change="(opts) => handleAvatarUpload(persona, opts)"
                >
                  <div class="persona-avatar persona-avatar--upload persona-avatar--edit">
                    <img
                      v-if="persona._previewUrl || getPersonaAvatarUrl(persona)"
                      :src="persona._previewUrl || getPersonaAvatarUrl(persona)"
                      alt="头像"
                    >
                    <span v-else class="persona-avatar__placeholder">+</span>
                  </div>
                </n-upload>
              </div>
              <div class="persona-edit-fields">
                <n-input
                  v-model:value="persona.name"
                  placeholder="人设名称"
                  size="small"
                />
                <n-input
                  v-model:value="persona.description"
                  placeholder="人设描述"
                  size="small"
                  type="textarea"
                  :autosize="{ minRows: 2, maxRows: 4 }"
                />
              </div>
            </div>
            <div class="persona-actions">
              <n-button
                quaternary
                circle
                size="small"
                type="primary"
                :loading="persona.loading"
                @click="handleSavePersona(persona)"
              >
                <template #icon>
                  <n-icon>
                    <Check />
                  </n-icon>
                </template>
              </n-button>
              <n-button
                quaternary
                circle
                size="small"
                type="error"
                @click="handleCancelEditPersona(persona)"
              >
                <template #icon>
                  <n-icon>
                    <Trash />
                  </n-icon>
                </template>
              </n-button>
            </div>
          </template>
        </div>

        <n-button
          size="small"
          dashed
          block
          class="persona-add-btn"
          @click="handleAddPersona"
        >
          <template #icon>
            <n-icon>
              <AddOutline />
            </n-icon>
          </template>
          新增人设
        </n-button>
      </div>
    </div>
  </div>

  <!-- 头像预览弹窗 -->
  <n-modal
    v-model:show="previewVisible"
    preset="card"
    title="头像预览"
    :auto-focus="false"
  >
    <img
      :src="previewAvatarUrl"
      class="persona-preview-img"
      alt="头像预览"
    >
  </n-modal>
</template>

<style scoped>
.user-config-tab {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: .5px;
  padding-bottom: 4px;
  border-bottom: 1px solid var(--border-color);
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-top: 8px;
}

.api-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-actions {
  display: flex;
  gap: 4px;
}

.add-btn {
  margin-top: 4px;
}

.settings-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.item-label {
  font-size: 14px;
  color: var(--text-primary);
  white-space: nowrap;
}

.item-control {
  flex: 1;
}

.persona-section {
  flex-direction: column;
  align-items: stretch;
}

.persona-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.persona-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border-radius: 6px;
  border: 1px solid var(--border-color);
  transition: background-color .2s;
}

.persona-item--active {
  background-color: var(--primary-color-suppl, rgba(24, 160, 88, .08));
  border-color: var(--primary-color, #18a058);
}

.persona-content {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  min-width: 0;
  cursor: pointer;
}

.persona-item--active .persona-content {
  cursor: default;
}

.persona-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
}

.persona-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.persona-avatar--upload {
  cursor: pointer;
  transition: opacity .2s;
}

.persona-avatar--upload:hover {
  opacity: .8;
}

.persona-avatar--edit {
  width: 60px;
  height: 60px;
}

.persona-avatar__placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background-color: var(--border-color);
  color: var(--text-secondary);
  font-size: 14px;
  font-weight: 600;
  user-select: none;
}

.persona-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.persona-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.persona-desc {
  font-size: 12px;
  color: var(--text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.persona-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.persona-edit-row {
  display: flex;
  align-items: stretch;
  gap: 8px;
  flex: 1;
  min-width: 0;
}

.persona-edit-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.persona-edit-fields {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.persona-add-btn {
  margin-top: 4px;
}

.persona-preview-img {
  max-width: 100%;
  max-height: 70vh;
  border-radius: 8px;
  object-fit: contain;
}


</style>