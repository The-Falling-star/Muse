import { ref, watch, type Ref } from 'vue';
import { getAvatarUrl, getAvatarUrlSync } from '@/utils/common';

/**
 * useAvatar 头像加载组合式函数
 * 用于异步加载和缓存头像URL
 */
export function useAvatar(avatarPath: Ref<string | undefined | null>) {
    const avatarUrl = ref(getAvatarUrlSync(avatarPath.value));
    const loading = ref(false);
    const error = ref<Error | null>(null);

    const loadAvatar = async () => {
        const path = avatarPath.value;
        if (!path || path.startsWith('http') || path.startsWith('data:')) {
            avatarUrl.value = path || '';
            return;
        }

        // 先尝试从缓存同步获取
        const cached = getAvatarUrlSync(path);
        if (cached) {
            avatarUrl.value = cached;
            return;
        }

        // 异步加载
        loading.value = true;
        error.value = null;

        try {
            const url = await getAvatarUrl(path);
            avatarUrl.value = url;
        } catch (err) {
            error.value = err instanceof Error ? err : new Error(String(err));
            console.error('加载头像失败:', err);
        } finally {
            loading.value = false;
        }
    };

    // 监听avatarPath变化
    watch(avatarPath, loadAvatar, { immediate: true });

    return {
        avatarUrl,
        loading,
        error,
        reload: loadAvatar,
    };
}

/**
 * useAvatars 批量头像加载组合式函数
 * 用于一次性加载多个头像
 */
export function useAvatars(avatarPaths: Ref<(string | undefined | null)[]>) {
    const avatarUrls = ref<Map<string, string>>(new Map());
    const loading = ref(false);

    const loadAvatars = async () => {
        loading.value = true;

        const promises = avatarPaths.value.map(async (path) => {
            if (!path) {
                return;
            }

            // 先检查缓存
            const cached = getAvatarUrlSync(path);
            if (cached) {
                avatarUrls.value.set(path, cached);
                return;
            }

            // 异步加载
            try {
                const url = await getAvatarUrl(path);
                avatarUrls.value.set(path, url);
            } catch (err) {
                console.error('加载头像失败:', path, err);
            }
        });

        await Promise.all(promises);
        loading.value = false;
    };

    const getUrl = (path: string | undefined | null): string => {
        if (!path) {
            return '';
        }
        return avatarUrls.value.get(path) || getAvatarUrlSync(path) || '';
    };

    watch(avatarPaths, loadAvatars, { immediate: true, deep: true });

    return {
        getUrl,
        loading,
        reload: loadAvatars,
    };
}
