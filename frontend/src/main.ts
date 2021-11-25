import {createApp} from 'vue';
import App from './App.vue';
import router from './router';
import * as Wails from '@wailsapp/runtime';
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

Wails.Init(() => {
    createApp(App).use(router).use(ElementPlus).mount('#app');
});
