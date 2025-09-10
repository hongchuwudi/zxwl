import {createRouter, createWebHistory} from 'vue-router'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('../views/Home.vue'),
        },
        {
            path: '/test',
            name: 'test',
            component: () => import('../views/logView.vue'),
        },
        {
            path: '/login',
            name: 'Login',
            meta: {title: '登录系统'},
            component: () => import('../views/Login.vue'),
        },
        {
            path: '/policy',
            name: 'policy',
            meta: {title: '登录系统'},
            component: () => import('../views/policy.vue'),
        },
        {
            path: '/Register',
            name: 'Register',
            meta: {title: '注册系统'},
            component: () => import('../views/Register.vue'),
        },
        {
            path: '/Reset',
            name: 'Reset',
            meta: {title: '重置密码系统'},
            component: () => import('../views/Reset.vue'),
        },
        {
            path: '/simulate',
            name: 'simulate',
            meta: {title: '志愿模拟模块'},
            component: () => import('../views/simulate.vue'),
        },
        {
            path: '/AiSmartsel',
            name: 'AiSmartsel',
            meta: {title: 'ai对话模块'},
            component: () => import('../views/AiSmartsel.vue'),
        },
        {
            path: '/Zxwl',
            name: 'Zxwl',
            meta: {title: '智选未来首页'},
            component: () => import('../views/Zxwl-s.vue'),
        },
        {
            path: '/profile',
            name: 'profile',
            meta: {title: '个人中心'},
            component: () => import('../views/profile.vue'),
        },
        {
            path: '/chat/:id',
            name: 'ChatRoom',
            component: () => import('@/views/ChatRoom.vue')
        },

        {
            path: '/professional',
            name: 'professional',
            component: () => import('../views/professional.vue'),
            props: true
        },
        {
            path: '/recommend',
            name: 'recommend',
            component: () => import('../views/recommend.vue'),
            props: true
        },

        {
            path: '/admin',
            children: [
                {path: 'users', component: () => import('@/views/AdminUsers.vue')},
                {path: 'policies/create', component: () => import('@/views/AdminCreatePolicy.vue')},
                {path: 'logs', component: () => import('@/views/SystemLogs.vue')}
            ]
        },

        // new insert
        {
            path: '/specialDetail',
            name: 'specialDetail',
            component: () => import('../views/new_zxwl/specialDetail.vue'),
            props: (route) => ({id: route.query.id})
        },
        {
            path: '/recommends',
            name: 'recommends',
            component: () => import('../views/new_zxwl/recommends.vue'),
            props: true
        },
        {
            path: '/schoolDetail',
            name: 'schoolDetail',
            component: () => import('../views/new_zxwl/schoolDetail.vue'),
            props: (route) => ({id: route.query.id})
        },
        {
            path: '/allSchool',
            name: 'allSchool',
            component: () => import('../views/new_zxwl/allSchool.vue'),
            props: true
        },
        {
            path: '/searchUAS', // universities and Specials
            name: 'searchUAS',
            component: () => import('../views/new_zxwl/components/searchSchAndSpe.vue'),
            props: true
        },
        {
            path: '/news',  // news是资讯的意思
            name: 'allNews',
            component: () => import('../views/new_zxwl/allNews.vue'),
            props: true
        },
        {
            path: '/newsDetail',
            name: 'newsDetail',
            component: () => import('../views/new_zxwl/newsDetail.vue')
            // props: true
        },
        {
          path: '/videoCall',
          name: 'videoCall',
          component: () => import('../views/new_zxwl/videoCall.vue'),
        },
        // 404 页面 - 必须放在最后
        {
            path: '/:pathMatch(.*)*',
            name: 'NotFound',
            component: () => import('../views/NotFound.vue'),
            meta: { title: '页面未找到' }
        }


    ],
})

export default router
