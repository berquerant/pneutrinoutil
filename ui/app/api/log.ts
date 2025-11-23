import { AxiosRequestConfig, InternalAxiosRequestConfig, AxiosInstance, AxiosResponse } from "axios";

export function enableAxiosLogger(instance: AxiosInstance): AxiosInstance {
  instance.interceptors.request.use((request: InternalAxiosRequestConfig) => {
    const x = showAxiosRequestConfig(request);
    console.log("Request :", x);
    return request;
  });
  instance.interceptors.response.use((response: AxiosResponse) => {
    const x = showAxiosRequestConfig(response.config);
    const y = showAxiosResponse(response);
    console.log("Response:", x, y);
    return response;
  });
  return instance;
}

function showAxiosRequestConfig(c: AxiosRequestConfig) {
  const method = c.method?.toUpperCase() || "";
  const baseURL = c.baseURL || "";
  const path = c.url || "";
  const params = c.params instanceof URLSearchParams
    ? c.params
    : new URLSearchParams(c.params || {});
  const paramsString = params.toString();
  const msg = method + " " + baseURL + path;
  if (paramsString != "") {
    return msg + "?" + paramsString;
  }
  return msg;
}

function showAxiosResponse(r: AxiosResponse) {
  const status = r.status;
  const statusText = r.statusText;
  return status + " " + statusText;
}
