import { QueryClient } from '@tanstack/react-query';
import { createRouter } from '@tanstack/react-router';
import { routeTree } from 'routeTree.gen';
import NotFoundPage from 'pages/not-found';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30 * 60 * 1000 // Set the global stale time to 30 minutes
    }
  }
});

export const router = createRouter<typeof routeTree>({
  routeTree,
  notFoundMode: 'root',
  basepath: '/v3',
  context: {
    queryClient,
    auth: undefined
  },
  defaultNotFoundComponent: NotFoundPage
});
