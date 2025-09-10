import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import'@/assets/global.css'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import store from '@/vuex/userStorage';
import Antd from 'ant-design-vue'
import 'ant-design-vue/dist/reset.css'

const app = createApp(App)
// 引入代码高亮样式
import 'highlight.js/styles/atom-one-dark.css'
import 'highlight.js/styles/atom-one-light.css'
import { BookOpenIcon, PencilIcon, DocumentTextIcon } from '@heroicons/vue/24/outline'

app.component('BookOpenIcon', BookOpenIcon)
app.component('PencilIcon', PencilIcon)
app.component('DocumentTextIcon', DocumentTextIcon)

// 设置等宽字体（在index.html的<head>中添加）
app.config.devtools = false

app.use(Antd)
app.use(store);
app.use(router)
app.use(ElementPlus,{
    locale:zhCn,
})

for(const [key,component] of Object.entries(ElementPlusIconsVue)){
    app.component(key,component)
}

app.mount('#app')
