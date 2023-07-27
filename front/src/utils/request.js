/*
*
* axios 二次封装
*
* 全局配置
*
* 响应拦截
*
* request 请求的方法
*
* 封装成 对象调用的方式
*
* request.get('/api',{xxx:xxx})
*
*
* */

// 引入文件
import axios from "axios"
import { ElMessage } from "element-plus";

// 全局配置
const service = axios.create({
    baseURL: '/user', // 根路径
    timeout: 8000  // 请求超时时间
})

// 添加响应拦截器
service.interceptors.response.use(res => {
    const { code, message, data } = res.data
    if (code === 200) {
        //request success
        ElMessage.success(message)
        return { code, data }
    }
    else if (code === 400) {
        //request fail
        ElMessage.error(message)
    }
})

//request 请求的方法
function request(options) {
    options.method = options.method || 'get'
    if (options.method.toLowerCase() === 'get') {
        options.params = options.data
    }
    // console.log(options)
    return service(options)
}

//快捷方式
// // 发送GET请求
// request.get('/api', { xxx: xxx })
// // 发送POST请求
// request.post('/api', { xxx: xxx })
// // 发送PUT请求
// request.put('/api', { xxx: xxx })
// // 发送DELETE请求
// request.delete('/api', { xxx: xxx })
//当你调用 request.get('/api', { xxx: xxx }) 时，
//实际上是调用了 request({ url: '/api', data: { xxx: xxx }, method: 'get' })。

['get', 'post', 'put', 'delete'].forEach(item => {
    request[item] = (url, data) => {
        return request({
            url,
            data,
            method: item
        })
    }
})

// 导出
export default request