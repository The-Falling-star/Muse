<script setup lang="ts">
import {NButton, NEmpty, NIcon, NInput, NList, NListItem, NPopconfirm, NTabPane, NTabs} from "naive-ui";
import {Check, Edit, Trash} from "@vicons/tabler"
import {APIProvider} from "@/gen/muse/common_pb.ts";
import type {APIConfig} from "@/gen/muse/user_pb.ts";
import {ref} from "vue";
import {useUserStore} from "@/stores/user.ts";
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
  try {
    const url = new URL(proxyUrl);
    return url.protocol === 'http:' || url.protocol === 'https:'
  } catch (err) {
    console.warn("proxyUrl非法")
    return false
  }
}


</script>

<template>
  <div class="user-config-tab">
    <div class="section-label">模型API</div>
    <n-tabs type="line" animated>
      <n-tab-pane
          v-for="provider in providers"
          :key="provider.value"
          :name="provider.name"
          :tab="provider.name"
      >
        <div class="tab-content">
          <n-empty
              v-if="apiConfigs.filter(apiCfg => apiCfg.provider === provider.value).length === 0"
              description="暂无API配置"
              size="small"
              class="empty-state"
          />
          <n-list v-else :show-divider="false">
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
    <br>
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
  </div>
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

.empty-state {
  padding: 32px 0;
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

</style>