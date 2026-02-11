import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { userClient } from '@/api/client';

// 缓存公共配置
let publicConfigCache: { skipAuth: boolean } | null = null;

// 带超时的 Promise 包装
const withTimeout = <T>(promise: Promise<T>, ms: number): Promise<T> => {
  return Promise.race([
    promise,
    new Promise<T>((_, reject) =>
      setTimeout(() => reject(new Error('请求超时')), ms)
    )
  ]);
};

// 获取公共配置（带缓存 + 超时控制）
const getPublicConfig = async () => {
  if (publicConfigCache) {
    return publicConfigCache;
  }
  try {
    const res = await withTimeout(userClient.getPublicConfig({}), 2000);
    publicConfigCache = { skipAuth: res.skipAuth };
    return publicConfigCache;
  } catch (error) {
    console.error('获取公共配置失败:', error);
    return { skipAuth: false };
  }
};

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/LoginView.vue'),
    meta: { title: '登录', requiresAuth: false }
  },
  {
    // 共享三栏布局的父路由
    path: '/',
    component: () => import('../components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Chat',
        component: () => import('../views/ChatView.vue'),
        meta: { title: '对话' }
      },
      {
        path: 'preset/:id',
        name: 'PresetEditor',
        component: () => import('../views/PresetEditorView.vue'),
        meta: { title: '预设编辑' }
      },
      {
        path: 'worldinfo/:id',
        name: 'WorldInfoEditor',
        component: () => import('../views/WorldInfoEditorView.vue'),
        meta: { title: '世界书编辑' }
      },
      {
        path: 'regex/:id',
        name: 'RegexEditor',
        component: () => import('../views/RegexEditorView.vue'),
        meta: { title: '正则编辑' }
      }
    ]
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 路由守卫 - 认证检查和页面标题设置
router.beforeEach(async (to, _from, next) => {
  // 设置页面标题
  const title = to.meta.title as string;
  document.title = title ? `${title} - Muse` : 'Muse';

  // 获取公共配置
  const publicConfig = await getPublicConfig();

  // 如果配置了跳过认证，直接放行
  if (publicConfig.skipAuth) {
    // 如果访问的是登录页，跳转到首页
    if (to.name === 'Login') {
      next({ name: 'Chat' });
    } else {
      next();
    }
    return;
  }

  // 检查是否需要认证
  const requiresAuth = to.meta.requiresAuth !== false;
  const userStore = useUserStore();
  const token = userStore.token;

  if (requiresAuth && !token) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'Login', query: { redirect: to.fullPath } });
  } else if (to.name === 'Login' && token) {
    // 已登录但访问登录页，重定向到首页
    next({ name: 'Chat' });
  } else {
    next();
  }
});

export default router;
