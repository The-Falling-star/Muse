import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { userClient } from '@/api/client';
import type {
  SysUser,
  Persona,
  APIConfig
} from '@/gen/muse/user_pb';

/**
 * 用户 Store
 * 管理认证状态、用户信息、人设、API配置等全局状态
 * 这些都是需要全局共享的核心数据
 */
export const useUserStore = defineStore('user', () => {
  // =====================
  // 全局状态
  // =====================

  // 认证状态
  const token = ref<string | null>(localStorage.getItem('token'));
  const currentUser = ref<SysUser | null>(null);

  // 人设
  const personas = ref<Persona[]>([]);

  // API配置
  const apiConfigs = ref<APIConfig[]>([]);

  // =====================
  // 计算属性
  // =====================

  // 是否已认证
  const isAuthenticated = computed(() => !!token.value && !!currentUser.value);

  // 活跃人设
  const activePersona = computed(() => {
    if (!currentUser.value) return null;
    return personas.value.find(p => p.id === currentUser.value!.activePersonaId) ?? null;
  });

  // 活跃API配置
  const activeApiConfig = computed(() => {
    return apiConfigs.value.find(c => c.isActive) ?? null;
  });

  // =====================
  // 认证方法
  // =====================

  // 设置Token
  const setToken = (newToken: string | null) => {
    token.value = newToken;
    if (newToken) {
      localStorage.setItem('token', newToken);
      return
    }
    localStorage.removeItem('token');
  };

  // 设置当前用户
  const setCurrentUser = (user: SysUser | null) => {
      if (!user) {
          return
      }
      currentUser.value = user;
  };

  // 登录
  const login = async (username: string, password: string) => {
    const res = await userClient.login({ username, password });
    if (res.token) {
      setToken(res.token);
    }
    if (res.user) {
      currentUser.value = res.user;
    }
    return res;
  };

  // 登出
  const logout = () => {
    setToken(null);
    currentUser.value = null;
    personas.value = [];
    apiConfigs.value = [];
  };

  // =====================
  // 人设状态管理
  // =====================

  // 设置人设列表
  const setPersonas = (list: Persona[]) => {
    personas.value = list;
  };

  // 添加人设
  const addPersona = (persona: Persona) => {
    personas.value.push(persona);
  };

  // 更新人设
  const updatePersonaInList = (persona: Persona) => {
    const index = personas.value.findIndex(p => p.id === persona.id);
    if (index >= 0) {
      personas.value[index] = persona;
    }
  };

  // 移除人设
  const removePersona = (id: number) => {
    personas.value = personas.value.filter(p => p.id !== id);
  };

  // 设置活跃人设ID
  const setActivePersonaId = (personaId: number) => {
    if (currentUser.value) {
      currentUser.value.activePersonaId = personaId;
    }
  };

  // =====================
  // API配置状态管理
  // =====================

  // 设置API配置列表
  const setApiConfigs = (configs: APIConfig[]) => {
    apiConfigs.value = configs;
  };

  // 添加API配置
  const addApiConfig = (config: APIConfig) => {
    apiConfigs.value.push(config);
  };

  // 更新API配置
  const updateApiConfigInList = (config: APIConfig) => {
    const index = apiConfigs.value.findIndex(c => c.id === config.id);
    if (index >= 0) {
      apiConfigs.value[index] = config;
    }
  };

  // 移除API配置
  const removeApiConfig = (id: number) => {
    apiConfigs.value = apiConfigs.value.filter(c => c.id !== id);
  };

  // 设置活跃API配置
  const setActiveApiConfigId = (configId: number) => {
    apiConfigs.value.forEach(c => {
      c.isActive = c.id === configId;
    });
  };

  // =====================
  // 初始化方法
  // =====================

  // 初始化用户数据
  const initUserData = async () => {
    try {
      const [userRes, personaRes, configsRes] = await Promise.all([
        userClient.getCurrentUser({}),
        userClient.listPersonas({}),
        userClient.listAPIConfigs({})
      ]);

      currentUser.value = userRes.user ?? null;
      personas.value = personaRes.personas;
      apiConfigs.value = configsRes.configs;
    } catch (error) {
      // 如果获取用户数据失败（如token过期），则登出
      console.error('初始化用户数据失败:', error);
      logout();
      throw error;
    }
  };

  return {
    // 状态
    token,
    currentUser,
    isAuthenticated,
    personas,
    activePersona,
    apiConfigs,
    activeApiConfig,

    // 认证方法
    setToken,
    setCurrentUser,
    login,
    logout,

    // 人设方法
    setPersonas,
    addPersona,
    updatePersonaInList,
    removePersona,
    setActivePersonaId,

    // API配置方法
    setApiConfigs,
    addApiConfig,
    updateApiConfigInList,
    removeApiConfig,
    setActiveApiConfigId,

    // 初始化
    initUserData
  };
});
