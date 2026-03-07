import {defineStore} from 'pinia';
import {computed, ref} from 'vue';
import {fileClient} from '@/api/client';

/**
 * 文件 Store
 * 管理文件缓存、上传和下载
 * 使用 Map 存储文件路径到 Blob URL 的映射
 */
export const useFileStore = defineStore('file', () => {
  // =====================
  // 全局状态
  // =====================

  // 文件URL缓存（filePath -> blobUrl）
  const fileCache = ref<Map<string, string>>(new Map());

  // 正在加载中的文件（filePath -> Promise）
  const loadingPromises = ref<Map<string, Promise<string>>>(new Map());

  // =====================
  // 计算属性
  // =====================

  // 缓存大小
  const cacheSize = computed(() => fileCache.value.size);

  // =====================
  // 核心方法
  // =====================

  /**
   * 获取单个文件的URL
   * 如果缓存中存在则直接返回，否则下载并缓存
   */
  const getFileUrl = async (filePath: string | undefined | null): Promise<string> => {
    // 空路径处理
    if (!filePath) {
      return '';
    }

    // 如果是外部链接或 data URL，直接返回
    if (filePath.startsWith('http://') || filePath.startsWith('https://') || filePath.startsWith('data:')) {
      return filePath;
    }

    // 检查缓存
    const cachedUrl = fileCache.value.get(filePath);
    if (cachedUrl) {
      return cachedUrl;
    }

    // 检查是否正在加载，避免重复请求
    const existingPromise = loadingPromises.value.get(filePath);
    if (existingPromise) {
      return existingPromise;
    }

    // 发起下载请求
    const loadPromise = downloadAndCacheFile(filePath);
    loadingPromises.value.set(filePath, loadPromise);

    try {
        return await loadPromise;
    } finally {
      loadingPromises.value.delete(filePath);
    }
  };

  /**
   * 批量获取文件URL
   * 返回与输入数组顺序对应的URL数组（空路径对应空字符串）
   */
  const getFileUrls = async (filePaths: (string | undefined | null)[]): Promise<string[]> => {
    const promises = filePaths.map(path => getFileUrl(path));
    return Promise.all(promises);
  };

  /**
   * 上传文件
   * 返回服务器存储的文件路径
   */
  const uploadFile = async (
    file: File | Blob,
    fileName: string,
    fileType: string = 'general'
  ): Promise<string> => {
    const fileContent = file instanceof File
      ? await file.arrayBuffer()
      : await file.arrayBuffer();

    const response = await fileClient.uploadFile({
      fileContent: new Uint8Array(fileContent),
      fileName,
      fileType,
    });

    return response.filePath;
  };

  /**
   * 上传文件并获取URL
   * 便捷方法，上传后自动缓存并返回URL
   */
  const uploadAndGetUrl = async (
    file: File | Blob,
    fileName: string,
    fileType: string = 'general'
  ): Promise<string> => {
    const filePath = await uploadFile(file, fileName, fileType);
    return getFileUrl(filePath);
  };

  // =====================
  // 缓存管理
  // =====================

  /**
   * 同步获取缓存中的URL（不触发下载）
   */
  const getCachedUrl = (filePath: string | undefined | null): string | undefined => {
    if (!filePath) {
      return undefined;
    }
    if (filePath.startsWith('http://') || filePath.startsWith('https://') || filePath.startsWith('data:')) {
      return filePath;
    }
    return fileCache.value.get(filePath);
  };

  /**
   * 检查是否已缓存
   */
  const hasCache = (filePath: string): boolean => {
    return fileCache.value.has(filePath);
  };

  /**
   * 预加载文件到缓存
   */
  const preloadFile = async (filePath: string): Promise<void> => {
    await getFileUrl(filePath);
  };

  /**
   * 批量预加载文件
   */
  const preloadFiles = async (filePaths: string[]): Promise<void> => {
    await Promise.all(filePaths.map(path => preloadFile(path)));
  };

  /**
   * 清除指定文件的缓存
   */
  const removeCache = (filePath: string): void => {
    const blobUrl = fileCache.value.get(filePath);
    if (blobUrl) {
      URL.revokeObjectURL(blobUrl);
      fileCache.value.delete(filePath);
    }
  };

  /**
   * 清空所有缓存
   */
  const clearCache = (): void => {
    fileCache.value.forEach((url) => {
      URL.revokeObjectURL(url);
    });
    fileCache.value.clear();
    loadingPromises.value.clear();
  };

  // =====================
  // 内部方法
  // =====================

  /**
   * 下载文件并缓存
   */
  const downloadAndCacheFile = async (filePath: string): Promise<string> => {
    const response = await fileClient.downloadFile({
      filePath,
    });

    if (!response.fileContent || response.fileContent.length === 0) {
      return '';
    }

    // 创建 Blob URL
    const blob = new Blob([response.fileContent as BlobPart], {
      type: response.contentType || 'application/octet-stream',
    });
    const blobUrl = URL.createObjectURL(blob);

    // 缓存
    fileCache.value.set(filePath, blobUrl);

    return blobUrl;
  };

  return {
    // 状态
    fileCache,
    cacheSize,

    // 核心方法
    getFileUrl,
    getFileUrls,
    uploadFile,
    uploadAndGetUrl,

    // 缓存管理
    getCachedUrl,
    hasCache,
    preloadFile,
    preloadFiles,
    removeCache,
    clearCache,
  };
});
