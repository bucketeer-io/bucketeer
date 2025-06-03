import axios from 'axios';
import type { AxiosInstance } from 'axios';
import { getTokenStorage, setTokenStorage } from 'storage/token';
import { refreshTokenFetcher } from './auth';

const axiosClient: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_WEB_API_ENDPOINT
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
      return document.dispatchEvent(
        new CustomEvent('unauthenticated', {
          bubbles: true
        })
      );
    }
    if (
      authToken?.refreshToken &&
      error.response?.status === 401 &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true;
      refreshTokenFetcher(authToken?.refreshToken)
        .then(response => {
          const newAccessToken = response.token.accessToken;
          setTokenStorage(response.token);
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          document.dispatchEvent(
            new CustomEvent('tokenRefreshed', {
              bubbles: true
            })
          );
          return axiosClient(originalRequest);
        })
        .catch(err => {
          document.dispatchEvent(
            new CustomEvent('unauthenticated', {
              bubbles: true
            })
          );
          return Promise.reject(err);
        });
    }
    return Promise.reject(error);
  }
);

export default axiosClient;
