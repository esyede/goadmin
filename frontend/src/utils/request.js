import store from '@/store'
import axios from 'axios'
import { Message, MessageBox } from 'element-ui'
// import router from '@/router'
import router from '@/router'
import { getToken } from '@/utils/auth'

// create an axios instance
const service = axios.create({
  baseURL: process.env.NODE_ENV === 'production' ? process.env.VUE_APP_BASE_API : '/', // api 的 base_url
  // withCredentials: true, // send cookies when cross-domain requests
  timeout: 5000 // request timeout
})

// request interceptor
service.interceptors.request.use(
  config => {
    // do something before request is sent
    if (store.getters.token) {
      // let each request carry token
      // ['X-Token'] is a custom headers key
      // please modify it according to the actual situation
      config.headers['Authorization'] = 'Bearer ' + getToken()
      // config.headers['Content-Type'] = 'application/json'
    }
    return config
  },
  error => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * If you want to get http information such as headers or status
   * Please return  response => response
  */

  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  response => {
    const res = response.data
    return res
  },
  error => {
    if (error.response.status === 401) {
      let err = error.response.data.message.toLowerCase()
      if (err.includes('jwt') && err.includes('failed')) {
        MessageBox.confirm(
          'Session expired, refresh or stay on the current page?',
          'Session Expired',
          {
            confirmButtonText: 'Refresh',
            cancelButtonText: 'Stay',
            type: 'warning'
          }
        ).then(() => {
          // store.dispatch('user/resetToken').then(() => {
          //   location.reload()
          // })
          store.dispatch('user/logout').then(() => {
            location.reload()
          })
        })
      } else {
        Message({
          showClose: true,
          message: error.response.data.message,
          type: 'error',
          duration: 5 * 1000
        })
        return Promise.reject(error)
      }
    } else if (error.response.status === 403) {
      router.push({ path: '/401' })
    } else {
      Message({
        showClose: true,
        message: error.response.data.message || error.message,
        type: 'error',
        duration: 5 * 1000
      })
      return Promise.reject(error)
    }
  }
)

export default service
