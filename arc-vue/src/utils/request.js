/* * 封装axios库 */
import axios from 'axios';
import storageService from "@/service/storageService";

const service = axios.create({
    // * 配置选项
    baseURL: process.env["VUE_APP_API_BASE_URL"],
    timeout: 1000 * 5,
    headers: {
        'Authorization': `Bearer ${storageService.get(storageService.USER_TOKEN)}`,
        'Content-Type': 'application/json',
    },
});

// * 请求拦截器
service.interceptors.request.use(config => {
    Object.assign(config.headers, {
        'Authorization': `Bearer ${storageService.get(storageService.USER_TOKEN)}`,
    })
    return config;
}, error => {
    // * 统一处理错误
    return Promise.reject(error);
});

// * 响应拦截器
service.interceptors.response.use(response => {
    return response;
}, error => {
    const { status, data } = error.response;
    // * 统一处理错误
    return Promise.reject(error);
});

export default service;
