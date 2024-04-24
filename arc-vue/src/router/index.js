import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import userRoutes from "@/router/module/user";
import store from "@/store";

const routes = [
  {
    path: '/',
    name: 'home',
    component: HomeView,
  },
  {
    path: '/about',
    name: 'about',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/AboutView.vue'),
  },
    ...userRoutes,
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach((to, from, next) => {
    // * 判断是否需要认证
    if(to.meta.auth){
        // * 判断是否登录
        if(store.state.userModule.token){
            // * 判断token有效性，如果token无效，请求token
            next();
        }else{
            // 跳转登录
            router.push({name: 'login'});
        }
    }else{
        next();
    }
});

export default router;
