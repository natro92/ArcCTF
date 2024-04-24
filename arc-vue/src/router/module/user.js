const userRoutes = [
    {
        path: '/register',
        name: 'register',
        component: () => import(/* webpackChunkName: "register" */ '@/views/register/RegisterForm.vue'),
    },
    {
        path: '/login',
        name: 'login',
        component: () => import(/* webpackChunkName: "login" */ '@/views/login/LoginForm.vue'),
    },
    {
        path: '/profile',
        name: 'profile',
        meta: {
          //  * 需要认证
          auth: true
        },
        component: () => import(/* webpackChunkName: "login" */ '@/views/profile/UserProfile.vue'),
    },
]

export default userRoutes;
