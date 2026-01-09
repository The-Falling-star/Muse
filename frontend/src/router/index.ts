import { createRouter, createWebHistory } from 'vue-router';

const routes = [
  {
    path: '/',
    name: 'Chat',
    component: () => import('../views/ChatView.vue'),
    meta: { title: '聊天' }
  },
  {
    path: '/characters',
    name: 'Characters',
    component: () => import('../views/CharactersView.vue'),
    meta: { title: '角色' }
  },
  {
    path: '/worldinfo',
    name: 'WorldInfo',
    component: () => import('../views/WorldInfoView.vue'),
    meta: { title: '世界书' }
  },
  {
    path: '/presets',
    name: 'Presets',
    component: () => import('../views/PresetsView.vue'),
    meta: { title: '预设' }
  },
  {
    path: '/regex',
    name: 'Regex',
    component: () => import('../views/RegexView.vue'),
    meta: { title: '正则' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/SettingsView.vue'),
    meta: { title: '设置' }
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes
});

// 路由守卫 - 设置页面标题
router.beforeEach((to, _from, next) => {
  const title = to.meta.title as string;
  document.title = title ? `${title} - Muse` : 'Muse';
  next();
});

export default router;
