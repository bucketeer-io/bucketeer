import { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, redirect } from '@tanstack/react-router';
import { AuthContextType } from 'auth';
import { getTokenStorage } from 'storage/token';
import RootLayout from './-root-layout';

export interface RouterContext {
  queryClient: QueryClient;
  auth?: AuthContextType;
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
  beforeLoad: async context => {
    const { pathname } = context.location;
    const authToken = getTokenStorage();
    if (!authToken && !['/v3/', '/v3/auth/callback'].includes(pathname))
      return redirect({
        to: '/'
      });
  }
});
