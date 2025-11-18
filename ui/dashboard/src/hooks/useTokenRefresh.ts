import { useEffect, useRef } from 'react';
import { refreshTokenFetcher } from '@api/auth';
import { getTokenStorage, setTokenStorage } from 'storage/token';

const REFRESH_BEFORE_EXPIRY_MS = 60 * 1000; // 1 minute before expiry (for tokens >= 10 min)
const MIN_TTL_FOR_FIXED_BUFFER_MS = 10 * 60 * 1000; // 10 minutes
const MIN_TTL_FOR_PROACTIVE_REFRESH_MS = 60 * 1000; // 1 minute - minimum TTL for proactive refresh

/**
 * Calculate when to refresh the token based on its TTL
 *
 * Strategy:
 * - For tokens with TTL <= 1 minute: SKIP proactive refresh (return null)
 *   - Rely on reactive refresh (401 interceptor) only
 *   - Prevents excessive server load (e.g., refresh every 5s for 10s TTL)
 *
 * - For tokens with TTL >= 10 minutes: refresh 1 minute before expiry
 *   - Optimal for production use cases
 *
 * - For tokens with 1 min < TTL < 10 minutes: refresh at 75% of lifetime
 *   - Balances UX with server load
 *   - Provides security buffer if refresh fails
 *
 * @param ttlMs Token time-to-live in milliseconds
 * @returns Refresh buffer in milliseconds, or null to skip proactive refresh
 */
const calculateRefreshBuffer = (ttlMs: number): number | null => {
  // Skip proactive refresh for very short TTLs (< 1 minute)
  // Reactive refresh (via 401 interceptor) is sufficient and avoids excessive load
  if (ttlMs < MIN_TTL_FOR_PROACTIVE_REFRESH_MS) {
    return null;
  }

  if (ttlMs >= MIN_TTL_FOR_FIXED_BUFFER_MS) {
    // For normal TTLs (>= 10 min), refresh 1 minute before expiry
    return REFRESH_BEFORE_EXPIRY_MS;
  }

  // For short TTLs (1 min < TTL < 10 min), refresh at 75% of lifetime
  // Calculation: Buffer = TTL - (TTL * 0.75) = TTL * 0.25
  return ttlMs * 0.25;
};

interface UseTokenRefreshOptions {
  onRefreshSuccess?: () => void;
  onRefreshError?: () => void;
}

export const useTokenRefresh = (options?: UseTokenRefreshOptions) => {
  const timerRef = useRef<NodeJS.Timeout | null>(null);
  const isRefreshingRef = useRef(false);

  const clearTimer = () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  };

  const scheduleTokenRefresh = () => {
    clearTimer();

    const authToken = getTokenStorage();
    if (!authToken || !authToken.refreshToken || !authToken.expiry) {
      return;
    }

    // Calculate token TTL and appropriate refresh buffer
    const now = Date.now();
    const expiryTime = authToken.expiry * 1000; // Convert Unix timestamp to milliseconds
    const ttl = expiryTime - now;

    // Use dynamic refresh buffer based on TTL
    const refreshBuffer = calculateRefreshBuffer(ttl);

    // Skip proactive refresh for very short TTLs (< 1 minute)
    // Rely on reactive refresh (401 interceptor) instead
    if (refreshBuffer === null) {
      return;
    }

    const refreshTime = expiryTime - refreshBuffer;
    const timeUntilRefresh = refreshTime - now;

    // If token is already expired or should be refreshed now, refresh immediately
    if (timeUntilRefresh <= 0) {
      performTokenRefresh();
      return;
    }

    // Schedule refresh
    timerRef.current = setTimeout(() => {
      performTokenRefresh();
    }, timeUntilRefresh);
  };

  const performTokenRefresh = async () => {
    if (isRefreshingRef.current) {
      return;
    }

    const authToken = getTokenStorage();
    if (!authToken || !authToken.refreshToken) {
      return;
    }

    isRefreshingRef.current = true;

    try {
      const response = await refreshTokenFetcher(authToken.refreshToken);

      if (response.token) {
        setTokenStorage(response.token);

        // Dispatch event to notify other components
        document.dispatchEvent(
          new CustomEvent('tokenRefreshed', {
            bubbles: true
          })
        );

        // Schedule next refresh (will be rescheduled by event handler, but that's ok - it cancels first)
        scheduleTokenRefresh();

        options?.onRefreshSuccess?.();
      }
    } catch (error: unknown) {
      // Only log out on authentication/authorization errors (401, 403)
      const isAxiosError =
        error && typeof error === 'object' && 'response' in error;
      const status = isAxiosError
        ? (error as { response?: { status?: number } }).response?.status
        : undefined;
      const isAuthError = status === 401 || status === 403;

      if (isAuthError) {
        // Account disabled, token invalid, or permission denied
        document.dispatchEvent(
          new CustomEvent('unauthenticated', {
            bubbles: true
          })
        );

        options?.onRefreshError?.();
      } else {
        // Network error or server error - retry on next scheduled refresh
        // The token will still expire naturally and trigger reactive refresh
        scheduleTokenRefresh();
      }
    } finally {
      isRefreshingRef.current = false;
    }
  };

  const startTokenRefresh = () => {
    scheduleTokenRefresh();
  };

  const stopTokenRefresh = () => {
    clearTimer();
  };

  useEffect(() => {
    // Start proactive token refresh on mount
    startTokenRefresh();

    // Listen for token storage changes (from other tabs or login events)
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'auth_token' && e.newValue) {
        scheduleTokenRefresh();
      }
    };

    // Listen for token refresh events (from axios interceptor or manual refresh)
    const handleTokenRefreshed = () => {
      console.log(
        '[useTokenRefresh] Token refreshed, rescheduling next refresh'
      );
      // Always reschedule when token is refreshed (from any source)
      // scheduleTokenRefresh() calls clearTimer() first, so no duplicate timers
      scheduleTokenRefresh();
    };

    window.addEventListener('storage', handleStorageChange);
    window.addEventListener('tokenRefreshed', handleTokenRefreshed);

    return () => {
      stopTokenRefresh();
      window.removeEventListener('storage', handleStorageChange);
      window.removeEventListener('tokenRefreshed', handleTokenRefreshed);
    };
  }, []);

  return {
    scheduleTokenRefresh,
    stopTokenRefresh
  };
};
