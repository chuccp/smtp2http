import axios, { type AxiosRequestConfig, type InternalAxiosRequestConfig, type AxiosResponse } from 'axios'

// Override the AxiosInstance type to reflect that our response interceptor
// unwraps response.data, so callers receive the raw API payload, not AxiosResponse
declare module 'axios' {
  interface AxiosInstance {
    request<T = any>(config: AxiosRequestConfig): Promise<T>
    get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
    delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<T>
    post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
    put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
    patch<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T>
  }
}
import { ElMessage } from 'element-plus'
import i18n from '@/locales'

// The session is an HttpOnly cookie, so requests must carry credentials and stay
// same-origin (in dev the Vite proxy forwards /api to the backend).
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  withCredentials: true
})

// 处理未登录状态，跳转到登录页
let isRedirecting = false
const handleUnauthorized = () => {
  if (isRedirecting) return
  isRedirecting = true
  localStorage.removeItem('http2smtp-username')
  localStorage.removeItem('http2smtp-isAdmin')
  localStorage.removeItem('http2smtp-logged-in')
  ElMessage.error(i18n.global.t('auth.sessionExpired'))
  window.location.href = '/login'
}

// 请求拦截器：会话由 HttpOnly cookie 承载，无需手动附加凭据
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config
  },
  (error: any) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 如果是 blob 类型响应，直接返回
    if (response.config.responseType === 'blob') {
      return response.data
    }

    const res = response.data
    if (res.code === 0 || res.code === 200) {
      return res
    } else {
      ElMessage.error(res.msg || res.message || '请求失败')
      return Promise.reject(new Error(res.msg || res.message || '请求失败'))
    }
  },
  (error: any) => {
    // 处理 HTTP 错误状态码
    if (error.response) {
      const status = error.response.status
      const data = error.response.data

      // 尝试获取后端返回的错误消息
      const errorMsg = data?.msg || data?.message || `请求失败(${status})`

      // The auth filter answers an expired or missing session with 403 "Not logged in"
      // (the framework's NoLogin), never 401 — handling only 401 left the app on the
      // protected page showing the raw English message. Handler-level failures such as
      // "admin access required" come back as 500, so 403 here always means no session.
      if (status === 401 || status === 403) {
        handleUnauthorized()
      } else {
        ElMessage.error(errorMsg)
      }
    } else {
      ElMessage.error(error.message || '网络错误')
    }
    return Promise.reject(error)
  }
)

export default request
