import { defineStore } from 'pinia';
import { ref } from 'vue';
import { commonClient } from '@/api/client';
import {APIProvider, type CandidateModels} from '@/gen/muse/common_pb';

/**
 * Common Store
 * 管理公共配置等全局状态
 */
export const useCommonStore = defineStore('common', () => {
  // =====================
  // 全局状态
  // =====================

  // 是否跳过认证
  const skipAuth = ref(false);

  // 候选模型列表
  const candidateModels = ref<Map<APIProvider, string[]>>(new Map());

  // =====================
  // 初始化方法
  // =====================

  // 初始化公共配置
  const initPublicConfig = async () => {
    try {
      const res = await commonClient.getPublicConfig({});
      skipAuth.value = res.skipAuth;
      const newModels = new Map<APIProvider, string[]>();
      res.candidateModels.forEach((model: CandidateModels) => {
          newModels.set(model.provider, model.models);
      });
      candidateModels.value = newModels; // 整体替换引用
    } catch (error) {
      console.error('初始化公共配置失败:', error);
      throw error;
    }
  };

  return {
    // 状态
    skipAuth,
    candidateModels,

    // 初始化
    initPublicConfig
  };
});
