import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router';
import App from './App.vue';
import { initGlobalMessage } from './composables/useGlobalMessage';

import './styles/global.css';
// 引入 highlight.js 代码高亮样式
import 'highlight.js/styles/atom-one-dark.css';
import {useUserStore} from "@/stores/user.ts";
import {useCommonStore} from "@/stores/common.ts";

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

app.mount('#app');

// 初始化全局消息实例（必须在 app.mount() 之后）
initGlobalMessage();

// 初始化公共配置
const commonStore = useCommonStore();
commonStore.initPublicConfig().catch(console.error);

const userStore = useUserStore();
userStore.initUserData().catch(console.error)