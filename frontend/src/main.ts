import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router';
import App from './App.vue';
import { initGlobalMessage } from './composables/useGlobalMessage';

import './styles/global.css';
// 引入 highlight.js 代码高亮样式
import 'highlight.js/styles/atom-one-dark.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

console.log('[Main] App mounting...');
app.mount('#app');
console.log('[Main] App mounted!');

// 初始化全局消息实例（必须在 app.mount() 之后）
console.log('[Main] Calling initGlobalMessage...');
initGlobalMessage();
console.log('[Main] initGlobalMessage returned');
