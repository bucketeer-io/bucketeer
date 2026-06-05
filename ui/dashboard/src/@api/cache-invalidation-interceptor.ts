import type { AxiosInstance, InternalAxiosRequestConfig } from 'axios';
import { queryClient } from 'configs/query-client';
import { resolveInvalidationKeys } from './cache-invalidation-map';

/**
 * Invalidates React Query caches for a single response, based on its request
 * config. Mutating verbs (POST/PUT/PATCH/DELETE) trigger invalidation; GETs
 * are skipped. The set of keys to invalidate is resolved from the request
 * URL via `URL_TO_KEYS` in `cache-invalidation-map.ts`.
 */
export const invalidateCacheForResponse = (
  config: InternalAxiosRequestConfig
) => {
  const method = config.method?.toUpperCase();
  if (!method || method === 'GET') return;
  const url = config.url ?? '';
  const keys = resolveInvalidationKeys(url);
  if (keys.length === 0) return;
  keys.forEach(key => queryClient.invalidateQueries({ queryKey: [key] }));
};

/**
 * Installs the response interceptor that auto-invalidates React Query caches
 * for mutating requests. ANY axios instance created under `src/@api/` MUST
 * call this — otherwise its mutations will silently bypass cache invalidation
 * and leave stale data in the UI. The `cache-invalidation-map.test.ts`
 * guardrail enforces this for every `axios.create()` site in `src/@api/`.
 */
export const installCacheInvalidationInterceptor = (client: AxiosInstance) => {
  client.interceptors.response.use(response => {
    invalidateCacheForResponse(response.config);
    return response;
  });
};
