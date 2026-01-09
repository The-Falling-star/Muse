import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router';
import App from './App.vue';

import './styles/global.css';
// 引入 highlight.js 代码高亮样式
import 'highlight.js/styles/atom-one-dark.css';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);
app.use(router);

app.mount('#app');
