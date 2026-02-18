import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '@/stores/user';
import { userClient } from '@/api/client';

// 缓存公共配置
let publicConfigCache: { skipAuth: boolean } | null = null;

// 获取公共配置（带缓存）
const getPublicConfig = async () => {
  if (publicConfigCache) {
    return publicConfigCache;
  }
  try {
    const res = await userClient.getPublicConfig({});
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
    path: '/chat/:sessionId',
    name: 'ChatSession',
    component: () => import('../views/ChatView.vue'),
    meta: { title: '聊天会话', requiresAuth: true }
  },
  {
    path: '/chat',
    redirect: '/'
  },
  {
    path: '/',
    name: 'Chat',
    component: () => import('../views/ChatView.vue'),
    meta: { title: '聊天', requiresAuth: true }
  },
  {
    path: '/characters',
    name: 'Characters',
    component: () => import('../views/CharactersView.vue'),
    meta: { title: '角色', requiresAuth: true }
  },
  {
    path: '/worldinfo',
    name: 'WorldInfo',
    component: () => import('../views/WorldInfoView.vue'),
    meta: { title: '世界书', requiresAuth: true }
  },
  {
    path: '/presets',
    name: 'Presets',
    component: () => import('../views/PresetsView.vue'),
    meta: { title: '预设', requiresAuth: true }
  },
  {
    path: '/regex',
    name: 'Regex',
    component: () => import('../views/RegexView.vue'),
    meta: { title: '正则', requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/SettingsView.vue'),
    meta: { title: '设置', requiresAuth: true }
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
