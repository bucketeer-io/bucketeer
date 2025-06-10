import { createRoute, redirect } from '@tanstack/react-router';
import { Route as EnvRoute } from './__env';

export const Route = createRoute({
  id: 'envIndex',
  beforeLoad: ({ params }) => {
    throw redirect({
      to: '/$env/features',
      params
    });
  },
  getParentRoute: () => EnvRoute
});
