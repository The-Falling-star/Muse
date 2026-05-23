import { ref, type Ref } from 'vue';

/**
 * 创建一个带有自动loading管理的请求封装
 * @param loadingRef - 可选的外部loading ref，如果不提供则使用内部ref
 */
export function useRequest(loadingRef?: Ref<boolean>) {
  const internalLoading = ref(false);
  const loading = loadingRef || internalLoading;

  /**
   * 包装一个异步函数，自动管理loading状态
   * @param fn - 要执行的异步函数
   * @returns 返回一个自动管理loading的新函数
   */
  const withLoading = <T, Args extends unknown[]>(
    fn: (...args: Args) => Promise<T>
  ) => {
    return async (...args: Args): Promise<T> => {
      loading.value = true;
      try {
        return await fn(...args);
      } finally {
        loading.value = false;
      }
    };
  };

  return {
    loading,
    withLoading
  };
}

/**
 * 全局loading状态管理
 * 用于显示全局加载指示器
 */
const globalLoading = ref(false);
const pendingRequests = ref(0);

export function useGlobalLoading() {
  const startLoading = () => {
    pendingRequests.value++;
    globalLoading.value = true;
  };

  const endLoading = () => {
    pendingRequests.value--;
    if (pendingRequests.value <= 0) {
      pendingRequests.value = 0;
      globalLoading.value = false;
    }
  };

  /**
   * 包装一个异步函数，自动管理全局loading状态
   */
  const withGlobalLoading = <T, Args extends unknown[]>(
    fn: (...args: Args) => Promise<T>
  ) => {
    return async (...args: Args): Promise<T> => {
      startLoading();
      try {
        return await fn(...args);
      } finally {
        endLoading();
      }
    };
  };

  return {
    globalLoading,
    pendingRequests,
    startLoading,
    endLoading,
    withGlobalLoading
  };
}
