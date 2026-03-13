import axios from 'axios';
import type { AxiosInstance } from 'axios';
import { urls } from 'configs';
import { getTokenStorage, setTokenStorage } from 'storage/token';
import { refreshTokenFetcher } from './auth';

let isRefreshing = false;

const axiosClient: AxiosInstance = axios.create({
  baseURL: urls.WEB_API_ENDPOINT
});

axiosClient.interceptors.request.use(
  config => {
    const authToken = getTokenStorage();
    if (authToken) {
      config.headers['Authorization'] = `Bearer ${authToken.accessToken}`;
    }
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

axiosClient.interceptors.response.use(
  response => response,
  async error => {
    const authToken = getTokenStorage();
    const originalRequest = error.config;
    if (!authToken && error.response?.status === 401) {
      document.dispatchEvent(
        new CustomEvent('unauthenticated', {
          bubbles: true
        })
      );
      return Promise.reject(error);
    }
    if (
      authToken?.refreshToken &&
      error.response?.status === 401 &&
      !isRefreshing &&
      !originalRequest._retry
    ) {
      isRefreshing = true;
      originalRequest._retry = true;
      return refreshTokenFetcher(authToken?.refreshToken)
        .then(response => {
          const newAccessToken = response.token.accessToken;
          setTokenStorage(response.token);
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          document.dispatchEvent(
            new CustomEvent('tokenRefreshed', {
              bubbles: true
            })
          );
          isRefreshing = false;
          return axiosClient(originalRequest);
        })
        .catch(err => {
          isRefreshing = false;
          document.dispatchEvent(
            new CustomEvent('unauthenticated', {
              bubbles: true
            })
          );
          return Promise.reject(err);
        });
    }
    if (error.response?.status === 401 && originalRequest._retry) {
      document.dispatchEvent(
        new CustomEvent('unauthenticated', {
          bubbles: true
        })
      );
    }
    return Promise.reject(error);
  }
);

export default axiosClient;
