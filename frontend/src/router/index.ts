import { createRouter, createMemoryHistory, RouteRecordRaw } from 'vue-router'
import Home from '../views/Home.vue'
import About from '../views/About.vue'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'
import Upload from '../components/Upload.vue'
import Sky from '../components/Sky.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/about',
    name: 'About',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    // component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
    component: About
  },
  {
    path: '/',
    name: 'Login',
    component: Login
  },
  {
    path: '/',
    name: 'Register',
    component: Register
  },
  {
    path: '/',
    name: 'Upload',
    component: Upload
  },
  {
    path: '/',
    name: 'Sky',
    component: Sky
  },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes
})

export default router
