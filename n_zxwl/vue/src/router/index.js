// router/index.js
import { createRouter, createWebHistory } from 'vue-router'

// 预加载状态追踪
const preloadStatus = {
    success: new Set(),
    failed: new Set(),
    pending: new Set()
}

// 预加载策略配置
const preloadStrategies = {
    home: [
        { loader: () => import('../views/open_base/Login.vue'), name: 'Login' },
        { loader: () => import('../views/open_base/Register.vue'), name: 'Register' },
        { loader: () => import('../views/open_base/Reset.vue'), name: 'Reset' },
        { loader: () => import('../views/Zxwl-s.vue'), name: 'zxwl' },
    ],
    Login: [
        { loader: () => import('../views/open_base/Home.vue'), name: 'home' },
        { loader: () => import('../views/Zxwl-s.vue'), name: 'zxwl' },
        { loader: () => import('../views/open_base/Register.vue'), name: 'Register' },
    ],
    Register: [
        { loader: () => import('../views/open_base/Home.vue'), name: 'home' },
        { loader: () => import('../views/open_base/Login.vue'), name: 'Login' },
        { loader: () => import('../views/Zxwl-s.vue'), name: 'zxwl' }
    ],
    zxwl: [
        { loader: () => import('../views/level_1_service/userFriend.vue'), name: 'user-friend' },
        { loader: () => import('../views/level_1_service/simulate.vue'), name: 'simulate' },
        { loader: () => import('../views/level_1_service/policy.vue'), name: 'policy' },
        { loader: () => import('../views/level_1_service/profile.vue'), name: 'profile' },
        { loader: () => import('../views/level_1_service/AiSmartsel.vue'), name: 'AiSmartsel' },
        { loader: () => import('../views/level_1_service/allNews.vue'), name: 'allNews' },
        { loader: () => import('../components/searchSchAndSpe.vue'), name: 'searchUAS' },
        { loader: () => import('../views/level_1_service/allSchool.vue'), name: 'allSchool' },
        { loader: () => import('../views/level_1_service/recommends.vue'), name: 'recommends' },
        { loader: () => import('../views/level_1_service/professional.vue'), name: 'professional' }
    ],
    default: [
        { loader: () => import('../views/open_base/Home.vue'), name: 'home' },
        { loader: () => import('../views/open_base/Login.vue'), name: 'Login' },
        { loader: () => import('../views/open_base/Register.vue'), name: 'Register' }
    ],
    professional:[
        { loader: () => import('../views/level_2_service/specialDetail.vue'),name:'specialDetail'}
    ],
    allSchool:[
        { loader: () => import('../views/level_2_service/schoolDetail.vue'), name: 'schoolDetail' }
    ],
    allNews:[
        { loader:() => import('../views/level_2_service/newsDetail.vue'), name: 'newsDetail' }
    ],
}

// 增强的预加载函数
const preloadBasedOnRoute = async (routeName) => {
    const strategies = preloadStrategies[routeName] || preloadStrategies.default

    console.log(`🚀 开始预加载 ${routeName} 的相关路由`)

    const promises = strategies.map(({ loader, name }) => {
        // 如果已经加载过，跳过
        if (preloadStatus.success.has(name) || preloadStatus.pending.has(name)) {
            return Promise.resolve({ name, status: 'skipped' })
        }

        preloadStatus.pending.add(name)

        return loader()
            .then(module => {
                preloadStatus.success.add(name)
                preloadStatus.pending.delete(name)
                console.log(`✅ 预加载成功: ${name}`)
                return { name, status: 'success' }
            })
            .catch(error => {
                preloadStatus.failed.add(name)
                preloadStatus.pending.delete(name)
                console.warn(`❌ 预加载失败: ${name}`, error)
                return { name, status: 'failed', error }
            })
    })

    const results = await Promise.allSettled(promises)

    // 打印预加载报告
    setTimeout(() => {
        printPreloadReport(routeName)
    }, 100)
}

// 打印预加载报告
const printPreloadReport = (triggerRoute) => {
    console.group(`📊 预加载报告 - 由 ${triggerRoute} 触发`)
    console.log('✅ 成功预加载:', Array.from(preloadStatus.success))
    console.log('❌ 预加载失败:', Array.from(preloadStatus.failed))
    console.log('⏳ 正在加载:', Array.from(preloadStatus.pending))
    console.groupEnd()
}


const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('../views/open_base/Home.vue'),
        },
        {
            path: '/test',
            name: 'test',
            component: () => import('../views/admin/logView.vue'),
        },
        {
            path: '/login',
            name: 'Login',
            meta: {title: '登录系统'},
            component: () => import('../views/open_base/Login.vue'),
        },
        {
            path: '/policy',
            name: 'policy',
            meta: {title: '政策系统'},
            component: () => import('../views/level_1_service/policy.vue'),
        },
        {
            path: '/Register',
            name: 'Register',
            meta: {title: '注册系统'},
            component: () => import('../views/open_base/Register.vue'),
        },
        {
            path: '/Reset',
            name: 'Reset',
            meta: {title: '重置密码系统'},
            component: () => import('../views/open_base/Reset.vue'),
        },
        {
            path: '/simulate',
            name: 'simulate',
            meta: {title: '志愿模拟模块'},
            component: () => import('../views/level_1_service/simulate.vue'),
        },
        {
            path: '/AiSmartsel',
            name: 'AiSmartsel',
            meta: {title: 'ai对话模块'},
            component: () => import('../views/level_1_service/AiSmartsel.vue'),
        },
        {
            path: '/zxwl',
            name: 'zxwl',
            meta: {title: '智选未来首页'},
            component: () => import('../views/Zxwl-s.vue'),
        },
        {
            path: '/profile',
            name: 'profile',
            meta: {title: '个人中心'},
            component: () => import('../views/level_1_service/profile.vue'),
        },
        {
            path: '/professional',
            name: 'professional',
            meta: {title: '专业详情'},
            component: () => import('../views/level_1_service/professional.vue'),
            props: true
        },
        {
            path: '/admin',
            children: [
                {path: 'users', component: () => import('@/views/admin/AdminUsers.vue')},
                {path: 'policies/create', component: () => import('@/views/admin/AdminCreatePolicy.vue')},
                {path: 'logs', component: () => import('@/views/admin/SystemLogs.vue')}
            ]
        },

        // new insert
        {
            path: '/specialDetail',
            name: 'specialDetail',
            meta: {title: '专业详情'},
            component: () => import('../views/level_2_service/specialDetail.vue'),
            props: (route) => ({id: route.query.id})
        },
        {
            path: '/recommends',
            name: 'recommends',
            meta: {title: '推荐'},
            component: () => import('../views/level_1_service/recommends.vue'),
            props: true
        },
        {
            path: '/schoolDetail',
            name: 'schoolDetail',
            meta: {title: '学校详情'},
            component: () => import('../views/level_2_service/schoolDetail.vue'),
            props: (route) => ({id: route.query.id})
        },
        {
            path: '/allSchool',
            name: 'allSchool',
            meta: {title: '学校列表'},
            component: () => import('../views/level_1_service/allSchool.vue'),
            props: true
        },
        {
            path: '/searchUAS', // universities and Specials
            name: 'searchUAS',
            meta: {title: '搜索学校列表'},
            component: () => import('../components/searchSchAndSpe.vue'),
            props: true
        },
        {
            path: '/news',  // news是资讯的意思
            name: 'allNews',
            meta: {title: '新闻列表'},
            component: () => import('../views/level_1_service/allNews.vue'),
            props: true
        },
        {
            path: '/newsDetail',
            name: 'newsDetail',
            meta: {title: '新闻详情'},
            component: () => import('../views/level_2_service/newsDetail.vue')
            // props: true
        },
        {
          path: '/videoCall',
          name: 'videoCall',
            meta: {title: '视频通话'},
          component: () => import('../views/level_2_service/videoCallMain.vue'),
        },
        {
            path: '/user-friend',
            name: 'user-friend',
            meta: {title: '用户好友'},
            component: () => import('../views/level_1_service/userFriend.vue'),
        },
        {
            path: '/score-section',
            name: 'score-section',
            component: () => import('../views/level_1_service/score.vue'),
        },
        // 404 页面 - 必须放在最后
        {
            path: '/:pathMatch(.*)*',
            name: 'NotFound',
            redirect: '/',
            component: () => import('../views/open_base/NotFound.vue'),
            meta: { title: '页面未找到' }
        }


    ],
})

// 记录已经触发预加载的路由
const preloadedRoutes = new Set()

router.afterEach((to) => {
    if (to.name && !preloadedRoutes.has(to.name)) {
        preloadedRoutes.add(to.name)

        // 延迟预加载，避免影响当前页面渲染
        setTimeout(() => {
            preloadBasedOnRoute(to.name)
        }, 300)
    }
})

// 添加全局方法以便在控制台检查状态
window.__preloadStatus = preloadStatus
window.__checkPreload = () => {
    printPreloadReport('manual check')
}

// 导出预加载状态供其他组件使用
export { preloadStatus }

export default router
