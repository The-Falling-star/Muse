import { fileClient } from '@/api/client';

// 缓存相关常量
const AVATAR_CACHE_PREFIX = 'avatar_cache_';
const CACHE_EXPIRE_KEY = 'avatar_cache_expire_min';

// 缓存项结构
interface CacheItem {
    data: string; // base64 data URL
    hash: string;
    expireAt: number; // 过期时间戳
}

// getAvatarUrl 获取头像URL，优先使用缓存
const getAvatarUrl = async (avatar: string | undefined | null): Promise<string> => {
    if (!avatar) {
        return '';
    }

    // 如果已经是完整的URL或data URI，直接返回
    if (avatar.startsWith('http') || avatar.startsWith('data:')) {
        return avatar;
    }

    // 计算缓存key
    const cacheKey = AVATAR_CACHE_PREFIX + avatar.replace(/[\/\\]/g, '_');

    // 检查缓存
    const cached = getFromCache(cacheKey);
    if (cached) {
        return cached;
    }

    // 从服务器下载
    try {
        const response = await fileClient.downloadFile({ filePath: avatar })

        // 转换为data URL
        const base64 = uint8ArrayToBase64(response.fileContent);
        const dataUrl = `data:${response.contentType};base64,${base64}`;

        // 存入缓存
        saveToCache(cacheKey, dataUrl, avatar);

        return dataUrl;
    } catch (error) {
        console.error('下载头像失败:', error);
        return '';
    }
};

// getAvatarUrlSync 同步获取头像URL（从缓存），用于需要同步渲染的场景
const getAvatarUrlSync = (avatar: string | undefined | null): string => {
    if (!avatar) {
        return '';
    }

    // 如果已经是完整的URL或data URI，直接返回
    if (avatar.startsWith('http') || avatar.startsWith('data:')) {
        return avatar;
    }

    // 计算缓存key
    const cacheKey = AVATAR_CACHE_PREFIX + avatar.replace(/[\/\\]/g, '_');

    // 检查缓存
    const cached = getFromCache(cacheKey);
    return cached || '';
};

// getFromCache 从缓存获取数据
const getFromCache = (key: string): string | null => {
    try {
        const itemStr = localStorage.getItem(key);
        if (!itemStr) {
            return null;
        }

        const item: CacheItem = JSON.parse(itemStr);

        // 检查是否过期
        if (Date.now() > item.expireAt) {
            localStorage.removeItem(key);
            return null;
        }

        return item.data;
    } catch {
        return null;
    }
};

// saveToCache 保存数据到缓存
const saveToCache = (key: string, data: string, hash: string): void => {
    try {
        const expireMin = getCacheExpireMinutes();
        const item: CacheItem = {
            data,
            hash,
            expireAt: Date.now() + expireMin * 60 * 1000,
        };
        localStorage.setItem(key, JSON.stringify(item));
    } catch (error) {
        // localStorage已满，清理过期缓存
        console.warn('缓存保存失败，尝试清理过期缓存:', error);
        cleanExpiredCache();
        try {
            const expireMin = getCacheExpireMinutes();
            const item: CacheItem = {
                data,
                hash,
                expireAt: Date.now() + expireMin * 60 * 1000,
            };
            localStorage.setItem(key, JSON.stringify(item));
        } catch {
            console.error('缓存保存仍然失败');
        }
    }
};

// getCacheExpireMinutes 获取缓存过期时间（分钟）
const getCacheExpireMinutes = (): number => {
    const stored = localStorage.getItem(CACHE_EXPIRE_KEY);
    if (stored) {
        const min = parseInt(stored, 10);
        if (!isNaN(min) && min > 0) {
            return min;
        }
    }
    return 1440; // 默认1天
};

// setCacheExpireMinutes 设置缓存过期时间（分钟）
const setCacheExpireMinutes = (minutes: number): void => {
    if (minutes > 0) {
        localStorage.setItem(CACHE_EXPIRE_KEY, minutes.toString());
    }
};

// cleanExpiredCache 清理过期缓存
const cleanExpiredCache = (): void => {
    const now = Date.now();
    const keysToRemove: string[] = [];

    for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i);
        if (key && key.startsWith(AVATAR_CACHE_PREFIX)) {
            try {
                const itemStr = localStorage.getItem(key);
                if (itemStr) {
                    const item: CacheItem = JSON.parse(itemStr);
                    if (now > item.expireAt) {
                        keysToRemove.push(key);
                    }
                }
            } catch {
                keysToRemove.push(key);
            }
        }
    }

    keysToRemove.forEach(key => localStorage.removeItem(key));
};

// uint8ArrayToBase64 将Uint8Array转换为base64字符串
const uint8ArrayToBase64 = (bytes: Uint8Array): string => {
    let binary = '';
    for (let i = 0; i < bytes.length; i++) {
        binary += String.fromCharCode(bytes[i]!);
    }
    return btoa(binary);
};

// preloadAvatars 预加载多个头像
const preloadAvatars = async (avatars: (string | undefined | null)[]): Promise<void> => {
    const promises = avatars
        .filter(a => a && !a.startsWith('http') && !a.startsWith('data:'))
        .map(a => getAvatarUrl(a));

    await Promise.all(promises);
};

// invalidateAvatarCache 使头像缓存失效
const invalidateAvatarCache = (avatar: string): void => {
    if (!avatar || avatar.startsWith('http') || avatar.startsWith('data:')) {
        return;
    }
    const cacheKey = AVATAR_CACHE_PREFIX + avatar.replace(/[\/\\]/g, '_');
    localStorage.removeItem(cacheKey);
};

export {
    generateSessionName,
    formatDateTime,
    getAvatarUrl,
    getAvatarUrlSync,
    preloadAvatars,
    invalidateAvatarCache,
    setCacheExpireMinutes,
    getCacheExpireMinutes,
    cleanExpiredCache,
};

const generateSessionName = (charName: string) => {
    return `${charName}_${formatDateTime()}`;
};

const formatDateTime = (date: Date = new Date()): string => {
    return new Intl.DateTimeFormat('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false
    }).format(date).replace(/\//g, '-');
};
