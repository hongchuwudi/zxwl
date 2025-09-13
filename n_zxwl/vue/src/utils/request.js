import axios from "axios";
import {ElMessage} from "element-plus";
import router from "@/router/index.js";

const request = axios.create({
    baseURL:'/gapi',
    timeout:30000
})

request.interceptors.request.use(config=>{
    config.headers['Content-Type'] = 'application/json;charset=utf-8';
    return config
},error=>{
    return Promise.reject(error)
});

request.interceptors.request.use(response=>{
    let res = response.data;
    if(typeof res === 'string'){
        res = res ? JSON.parse(res):res
    }
    return res;
},error=>{
    if(error.response.status === 404){
        ElMessage.error('未找到接口')
    }else {
        console.error(error.message)
    }
    return Promise.reject(error)
});

export default request