import { createApp } from 'vue';
import App from './App.vue';
import axios from 'axios';
import VueAxios from "vue-axios";
import router from './router';
import store from './store';
import Toast from "vue-toastification";
import "vue-toastification/dist/index.css";

import '@/assets/css/main.css';
import '@/assets/css/output.css';

const app = createApp(App);
app.provide('axios', axios);
app.use(Toast, {
    // 这里可以配置全局的toast选项
    position: 'top-right',
    timeout: 5000,
    closeOnClick: true,
    pauseOnFocusLoss: true,
    pauseOnHover: true,
    draggable: true,
    draggablePercent: 0.6,
    showCloseButtonOnHover: false,
    hideProgressBar: true,
    closeButton: 'button',
    icon: true,
    rtl: false
});
app.use(store).use(router).mount('#app');
