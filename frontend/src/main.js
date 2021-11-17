import 'core-js/stable';
import 'regenerator-runtime/runtime';
import Vue from 'vue';
import App from './App.vue';

import VueRouter from "vue-router";
//开启debug模式
Vue.config.debug = true;

Vue.config.productionTip = false;
Vue.config.devtools = true;

import * as Wails from '@wailsapp/runtime';
// 定义组件, 也可以像教程之前教的方法从别的文件引入
import login from './components/Login.vue'
import Element from 'element-ui'


// 创建一个路由器实例
// 并且配置路由规则
const router = new VueRouter({
    mode: 'history',
    base: __dirname,
    routes: [
        {
            path: '/login',
            component: login
        }
    ]
})
Vue.use(VueRouter)
Vue.use(Element)

Wails.Init(() => {
    new Vue({
        router: router,
        render: h => h(App),
    }).$mount('#app');
});
