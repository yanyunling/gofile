import axios, { AxiosRequestConfig } from "axios";
import { ElMessage } from "element-plus";
import { host, context, basicTokenKey } from "@/config";
import { useTokenStore } from "@/store/token";
import { useCommonStore } from "@/store/common";
import sha256 from "crypto-js/sha256";

interface ResponseResult<T> {
  code: number;
  message: string;
  data: T;
}

// 刷新token中
let refreshing = false;

// 数据接口axios实例
const dataInstance = axios.create({
  baseURL: host + context + "/data",
  timeout: 10000,
});

// 授权接口axios实例
const authInstance = axios.create({
  baseURL: host + context + "/token",
  timeout: 10000,
});

// 开放接口axios实例
const openInstance = axios.create({
  baseURL: host + context + "/open",
  timeout: 10000,
});

/**
 * 数据接口request拦截器
 */
dataInstance.interceptors.request.use(
  (config) => {
    const tokenStore = useTokenStore();
    // 设置请求头认证信息
    if (!config.headers.Authorization) {
      config.headers.Authorization = "Bearer " + tokenStore.accessToken;
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  },
);

/**
 * 授权接口request拦截器
 */
authInstance.interceptors.request.use(
  (config) => {
    // 设置请求头认证信息
    if (!config.headers.Authorization) {
      let token = sha256(basicTokenKey + Math.floor(new Date().getTime() / 600000)).toString();
      config.headers.Authorization = "Basic " + token;
    }
    return config;
  },
  (error) => {
    Promise.reject(error);
  },
);

/**
 * 数据接口respone拦截器
 */
dataInstance.interceptors.response.use(
  async (response) => {
    const tokenStore = useTokenStore();
    // 请求结果码
    if (response.status === 200) {
      // 检验blob类型的返回结果
      let blobResult = await blobCheck(response);
      if (!blobResult.isBlob) {
        response.data = blobResult.data;
      }
      // 业务结果码，数据请求成功，或返回了有效的blob数据
      if (response.data.code === 200 || blobResult.isBlob) {
        return response;
      } else if (response.data.code === 401) {
        // 认证失败
        return new Promise((resolve, reject) => {
          if (tokenStore.refreshToken) {
            if (!refreshing) {
              // 未刷新token时才进行刷新
              refreshing = true;
              authInstance({ method: "post", url: "/refresh", data: { refreshToken: tokenStore.refreshToken } })
                .then((tokenResult) => {
                  refreshing = false;
                  // 刷新token成功，保存新的token
                  tokenStore.setToken(tokenResult.data.data);
                  // 重设头信息中的token
                  response.config.headers.Authorization = "Bearer " + tokenStore.accessToken;
                  // 重新进行上一次的请求
                  dataInstance(response.config)
                    .then((dataResult) => {
                      resolve(dataResult);
                    })
                    .catch((dataError) => {
                      reject(dataError);
                    });
                })
                .catch(() => {
                  refreshing = false;
                  // 刷新token失败，清空token信息，弹出登录窗口
                  tokenStore.removeToken();
                  useCommonStore().loginVisible = true;
                  reject({});
                });
            } else {
              // 如已经在刷新token，等待token刷新完毕，重新进行上一次的请求
              const interval = setInterval(() => {
                if (!refreshing) {
                  clearInterval(interval);
                  // 重设头信息中的token
                  response.config.headers.Authorization = "Bearer " + tokenStore.accessToken;
                  // 重新进行上一次的请求
                  dataInstance(response.config)
                    .then((dataResult) => {
                      resolve(dataResult);
                    })
                    .catch((dataError) => {
                      reject(dataError);
                    });
                }
              }, 500);
            }
          } else {
            // token不存在直接弹出登录窗口
            tokenStore.removeToken();
            useCommonStore().loginVisible = true;
            reject({});
          }
        });
      }
      ElMessage.error(response.data.message ? response.data.message : "服务器错误");
    } else {
      ElMessage.error(response.statusText ? response.statusText : "连接超时");
    }
    return Promise.reject(response);
  },
  (error) => {
    ElMessage.error("连接超时");
    return Promise.reject(error);
  },
);

/**
 * 授权接口respone拦截器
 */
authInstance.interceptors.response.use(
  (response) => {
    // 请求结果码
    if (response.status === 200) {
      // 业务结果码
      if (response.data.code === 200) {
        return response;
      }
      ElMessage.error(response.data.message ? response.data.message : "服务器错误");
    } else {
      ElMessage.error(response.statusText ? response.statusText : "连接超时");
    }
    return Promise.reject(response);
  },
  (error) => {
    ElMessage.error("连接超时");
    return Promise.reject(error);
  },
);

/**
 * 开放接口respone拦截器
 */
openInstance.interceptors.response.use(
  async (response) => {
    // 请求结果码
    if (response.status === 200) {
      // 检验blob类型的返回结果
      let blobResult = await blobCheck(response);
      if (!blobResult.isBlob) {
        response.data = blobResult.data;
      }
      // 业务结果码，数据请求成功，或返回了有效的blob数据
      if (response.data.code === 200 || blobResult.isBlob) {
        return response;
      }
      ElMessage.error(response.data.message ? response.data.message : "服务器错误");
    } else {
      ElMessage.error(response.statusText ? response.statusText : "连接超时");
    }
    return Promise.reject(response);
  },
  (error) => {
    ElMessage.error("连接超时");
    return Promise.reject(error);
  },
);

/**
 * 调用数据接口
 * @param config
 * @returns
 */
export default <T>(config: AxiosRequestConfig) => {
  return new Promise<ResponseResult<T>>((resolve, reject) => {
    dataInstance
      .request<ResponseResult<T>>(config)
      .then((res) => {
        resolve(res.data);
      })
      .catch((err) => {
        reject(err.data ? err.data : err);
      });
  });
};

/**
 * 调用授权接口
 * @param config
 * @returns
 */
export const authRequest = <T>(config: AxiosRequestConfig) => {
  return new Promise<ResponseResult<T>>((resolve, reject) => {
    authInstance
      .request<ResponseResult<T>>(config)
      .then((res) => {
        resolve(res.data);
      })
      .catch((err) => {
        reject(err.data ? err.data : err);
      });
  });
};

/**
 * 调用开放接口
 * @param config
 * @returns
 */
export const openRequest = <T>(config: AxiosRequestConfig) => {
  return new Promise<ResponseResult<T>>((resolve, reject) => {
    openInstance
      .request<ResponseResult<T>>(config)
      .then((res) => {
        resolve(res.data);
      })
      .catch((err) => {
        reject(err.data ? err.data : err);
      });
  });
};

/**
 * 检验blob类型的返回结果
 * @param {*} response
 */
const blobCheck = (response): Promise<{ isBlob: boolean; data: Object }> => {
  return new Promise((resolve, reject) => {
    // blob类型
    if (response.config.responseType && response.config.responseType.toLowerCase() === "blob") {
      // blob类型的json信息
      if (response.data.type && response.data.type.indexOf("application/json") != -1) {
        let reader = new FileReader();
        reader.readAsText(response.data);
        reader.onload = (res) => {
          resolve({ isBlob: false, data: JSON.parse(String(res.target.result)) });
        };
        reader.onerror = () => {
          resolve({ isBlob: false, data: { code: 500, message: "文件解析失败" } });
        };
      } else {
        // 正常blob类型
        resolve({ isBlob: true, data: response.data });
      }
    } else {
      // 非blob类型
      resolve({ isBlob: false, data: response.data });
    }
  });
};
