import type { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'

export function enableAxiosLogger(instance: AxiosInstnace): AxiosInstance {
  instance.interceptors.request.use(request => {
    const x = showAxiosRequestConfig(request)
    console.log("Request:", x)
    return request
  })
  instance.interceptors.response.use(response => {
    const x = showAxiosRequestConfig(response.config)
    const y = showAxiosResponse(response)
    console.log("Response:", x, y)
    return response
  })
  return instance
}

function showAxiosRequestConfig(c: AxiosRequestConfig) {
  const method = c.method?.toUpperCase() || ''
  const baseURL = c.baseURL || ''
  const path = c.url || ''
  const params = c.params instanceof URLSearchParams
    ? c.params : new URLSearchParams(c.params || {})
  return method + ' ' + baseURL + path + '?' + params.toString()
}

function showAxiosResponse(r: AxiosResponse) {
  const status = r.status
  const statusText = r.statusText
  return status + ' ' + statusText
}
