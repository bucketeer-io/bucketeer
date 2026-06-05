import { QueryClient } from '@tanstack/react-query';

/**
 * Shared `QueryClient` for the entire dashboard.
 *
 * - `staleTime` is intentionally short. Together with the axios response
 *   interceptor (see `@api/cache-invalidation-map.ts`), this ensures that
 *   data the user sees is at most a few seconds old without requiring every
 *   mutation site to remember which queries to invalidate.
 * - `gcTime` keeps queries in memory for a while so back/forward navigation
 *   stays snappy.
 * - `refetchOnWindowFocus` / `refetchOnReconnect` rely on `staleTime` to
 *   decide when to actually refetch; with the short stale window above they
 *   become useful again.
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 1000,
      gcTime: 10 * 60 * 1000,
      refetchOnWindowFocus: true,
      refetchOnReconnect: true,
      refetchOnMount: true
    }
  }
});
