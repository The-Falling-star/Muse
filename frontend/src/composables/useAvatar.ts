import { ref, type Ref, watch } from 'vue';
import { useFileStore } from '@/stores/file';

/**
 * useAvatar 头像加载组合式函数
 * 便捷封装，自动处理缓存和异步加载
 *
 * 如果需要 loading/error 状态，可直接使用 fileStore
 */
export function useAvatar(avatarPath: Ref<string | undefined | null>) {
  const fileStore = useFileStore();
  const avatarUrl = ref<string>('');

  const loadAvatar = async () => {
    const path = avatarPath.value;

    // 空路径
    if (!path) {
      avatarUrl.value = '';
      return;
    }

    // fileStore.getFileUrl 已处理缓存、外部链接等逻辑
    avatarUrl.value = await fileStore.getFileUrl(path);
  };

  // 监听路径变化，自动加载
  watch(avatarPath, loadAvatar, { immediate: true });

  return { avatarUrl };
}
