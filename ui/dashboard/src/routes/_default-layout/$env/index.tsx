import { createRoute, redirect } from '@tanstack/react-router';
import { Route as EnvRoute } from './__env';

export const Route = createRoute({
  path: '/',
  beforeLoad: ({ params }) => {
    throw redirect({
      to: '/$envUrlCode/features',
      params
    });
  },
  getParentRoute: () => EnvRoute
});
