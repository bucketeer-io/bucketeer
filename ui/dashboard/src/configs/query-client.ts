import { QueryClient } from '@tanstack/react-query';

/**
 * Shared `QueryClient` for the entire dashboard.
 *
 * Freshness strategy:
 * - The axios response interceptor (see `@api/cache-invalidation-map.ts`)
 *   is the PRIMARY freshness mechanism: every mutation routed through
 *   `@api/axios-client` automatically invalidates the relevant queries, so
 *   the user's own writes are reflected immediately on the next render.
 * - `staleTime` is therefore set generously (5 min). It is no longer the
 *   safety net for write-driven freshness — it only controls how often
 *   "background" refetches (focus, reconnect, mount) actually hit the
 *   network. A 30-second value would cause 8–10 simultaneous refetches on
 *   every tab return for multi-query pages like feature-flag details.
 * - `refetchOnWindowFocus` / `refetchOnReconnect` stay enabled to cover
 *   changes the interceptor cannot see: edits by other admins, scheduled
 *   flag updates, auto-ops / progressive rollout side effects, and other
 *   server-driven state changes. Combined with the 5-minute `staleTime`,
 *   they only refetch queries that are actually stale, so the cost stays
 *   bounded.
 * - `gcTime` keeps queries in memory for a while so back/forward
 *   navigation stays snappy.
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      gcTime: 10 * 60 * 1000,
      refetchOnWindowFocus: true,
      refetchOnReconnect: true,
      refetchOnMount: true
    }
  }
});
