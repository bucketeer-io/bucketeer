import { QueryClient } from '@tanstack/react-query';
import { createRouter } from '@tanstack/react-router';
import { routeTree } from 'routeTree.gen';

export const queryClient = new QueryClient();

export const router = createRouter<typeof routeTree>({
  routeTree,
  notFoundMode: 'root',
  basepath: '/v3',
  context: {
    queryClient,
    auth: {
      isLogin: false,
      logout: () => {},

      consoleAccount: undefined,
      myOrganizations: [],

      syncSignIn: async () => {},
      onMeFetcher: async () => {},

      isInitialLoading: false,
      setIsInitialLoading: () => {},

      isGoogleAuthError: false,
      setIsGoogleAuthError: () => {}
    }
  }
});
